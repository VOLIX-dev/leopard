package s3

import (
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/volix-dev/leopard/files"
	"io"
	"io/fs"
	"io/ioutil"
	"strings"
	"time"
)

type S3Driver struct {
	client *s3.S3
	bucket string
}

type S3Config struct {
	Bucket     string
	Region     *string
	AccessKey  string
	SecretKey  string
	Endpoint   *string
	DisableSSL *bool
}

func init() {
	files.Register("s3", func(config any) (files.Driver, error) {
		if conf, isConfig := config.(S3Config); isConfig {
			client := s3.New(nil, &aws.Config{
				Region:      conf.Region,
				Credentials: credentials.NewStaticCredentials(conf.AccessKey, conf.SecretKey, ""),
				Endpoint:    conf.Endpoint,
				DisableSSL:  conf.DisableSSL,
			})

			return &S3Driver{
				client: client,
				bucket: conf.Bucket,
			}, nil

		}

		return nil, errors.New("invalid config")
	})
}

func (s *S3Driver) Create(path string) (files.File, error) {
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	return &S3File{
		path:   path,
		client: s.client,
		bucket: s.bucket,
	}, err
}

func (s *S3Driver) Open(path string) (files.File, error) {
	return &S3File{
		path:   path,
		client: s.client,
		bucket: s.bucket,
	}, nil
}

func (s *S3Driver) Remove(path string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	return err
}

func (s *S3Driver) Rename(oldPath, newPath string) error {
	_, err := s.client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(s.bucket),
		Key:        aws.String(newPath),
		CopySource: aws.String(s.bucket + "/" + oldPath),
	})

	if err == nil {
		_, err = s.client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(oldPath),
		})
	}

	return err
}

func (s *S3Driver) Stat(path string) (fs.FileInfo, error) {
	resp, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return nil, err
	}

	parts := strings.Split(path, "/")
	fileName := parts[len(parts)-1]

	return &S3FileInfo{
		object: &s3.Object{
			Key:          &fileName,
			LastModified: resp.LastModified,
			Size:         resp.ContentLength,
		},
	}, nil
}

func (s *S3Driver) MakeDir(path string) error {
	return nil
}

func (s *S3Driver) RemoveDir(path string) error {
	return nil
}

func (s *S3Driver) RenameDir(oldPath, newPath string) error {
	return nil
}

func (s *S3Driver) StatDir(path string) (fs.FileInfo, error) {
	return nil, nil
}

func (s *S3Driver) ListDir(path string) ([]fs.DirEntry, error) {
	resp, err := s.client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(path),
	})

	if err != nil {
		return nil, err
	}

	entries := make([]fs.DirEntry, len(resp.Contents))

	for i, object := range resp.Contents {
		parts := strings.Split(aws.StringValue(object.Key), "/")
		fileName := parts[len(parts)-1]

		entries[i] = &S3FileInfo{
			object: &s3.Object{
				Key:          &fileName,
				LastModified: object.LastModified,
				Size:         object.Size,
			},
		}
	}

	return entries, nil
}

func (s *S3Driver) ReadFile(path string) ([]byte, error) {
	resp, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func (s *S3Driver) WriteFile(path string, data []byte) error {
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	})

	return err
}

func (s *S3Driver) ReadFileStream(path string) (io.ReadCloser, error) {
	resp, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (s *S3Driver) WriteFileStream(path string, data io.Reader) error {
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   ReaderConverter{reader: data},
	})

	return err
}

func (s *S3Driver) RemoveFile(path string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	return err
}

func (s *S3Driver) FileSize(path string) (int64, error) {
	resp, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return 0, err
	}

	return getPointingNum(resp.ContentLength), nil
}

func (s *S3Driver) FileModTime(path string) (int64, error) {
	resp, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return 0, err
	}

	return getPointingStruct(resp.LastModified).Unix(), nil
}

type S3File struct {
	client         *s3.S3
	bucket         string
	path           string
	incoming       *s3.GetObjectOutput
	outgoingBuffer []byte
}

func (s *S3File) Read(p []byte) (n int, err error) {
	if s.incoming == nil {
		s.incoming, err = s.client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(s.path),
		})
		if err != nil {
			return 0, err
		}
	}
	read, err := s.incoming.Body.Read(p)

	return read, err
}

func (s *S3File) Write(p []byte) (n int, err error) {
	s.outgoingBuffer = append(s.outgoingBuffer, p...)
	return len(p), nil
}

func (s *S3File) Close() error {
	if len(s.outgoingBuffer) > 0 {
		_, err := s.client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(s.path),
			Body:   bytes.NewReader(s.outgoingBuffer),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *S3File) Stat() (fs.FileInfo, error) {
	object, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.path),
	})

	if err != nil {
		return nil, err
	}

	parts := strings.Split(s.path, "/")
	fileName := parts[len(parts)-1]

	return &S3FileInfo{
		object: &s3.Object{
			Key:          &fileName,
			Size:         object.ContentLength,
			LastModified: object.LastModified,
		},
	}, nil
}

type S3FileInfo struct {
	object *s3.Object
}

func (s *S3FileInfo) Name() string {
	return getPointingString(s.object.Key)
}

func (s *S3FileInfo) Size() int64 {
	return getPointingNum(s.object.Size)
}

func (s *S3FileInfo) Mode() fs.FileMode {
	return fs.ModeSticky
}

func (s *S3FileInfo) ModTime() time.Time {
	return getPointingStruct(s.object.LastModified)
}

func (s *S3FileInfo) IsDir() bool {
	return false
}

func (s *S3FileInfo) Sys() any {
	return nil
}

func (s *S3FileInfo) Type() fs.FileMode {
	return fs.ModeSticky
}

func (s *S3FileInfo) Info() (fs.FileInfo, error) {
	return s, nil
}

func getPointingString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func getPointingNum[N int | int8 | int16 | int32 | int64 | float32 | float64](i *N) N {
	if i == nil {
		return 0
	}

	return *i
}

func getPointingStruct[T any](s *T) T {
	return *s
}
