package s3

import "io"

type ReaderConverter struct {
	reader io.Reader
}

func (r ReaderConverter) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r ReaderConverter) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
