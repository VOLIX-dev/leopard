package osFs

import (
	"errors"
	"github.com/volix-dev/leopard/files"
	"io"
	"io/fs"
	"os"
	path2 "path"
)

func init() {
	files.Register("os", func(config any) (files.Driver, error) {
		if root, isString := config.(string); isString {
			err := os.Mkdir(root, 0777)

			if err != nil && !os.IsExist(err) {
				return nil, err
			}
			return &OsFs{root: root}, nil
		}
		return nil, errors.New("invalid config for osFs")
	})
}

type OsFs struct {
	root string
}

func (o OsFs) addRoot(path string) string {
	return path2.Join(o.root, path)
}

func (o *OsFs) Create(path string) (files.File, error) {
	return os.Create(o.addRoot(path))
}

func (o *OsFs) Open(path string) (files.File, error) {
	return os.Open(o.addRoot(path))
}

func (o *OsFs) Remove(path string) error {
	return os.Remove(o.addRoot(path))
}

func (o *OsFs) Rename(oldPath, newPath string) error {
	return os.Rename(o.addRoot(oldPath), o.addRoot(newPath))
}

func (o *OsFs) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(o.addRoot(path))
}

func (o *OsFs) MakeDir(path string) error {
	return os.Mkdir(o.addRoot(path), 0755)
}

func (o *OsFs) RemoveDir(path string) error {
	return os.Remove(o.addRoot(path))
}

func (o *OsFs) RenameDir(oldPath, newPath string) error {
	return os.Rename(o.addRoot(oldPath), o.addRoot(newPath))
}

func (o *OsFs) StatDir(path string) (fs.FileInfo, error) {
	return os.Stat(o.addRoot(path))
}

func (o *OsFs) ListDir(path string) ([]fs.DirEntry, error) {
	return os.ReadDir(o.addRoot(path))
}

func (o *OsFs) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(o.addRoot(path))
}

func (o *OsFs) WriteFile(path string, data []byte) error {
	return os.WriteFile(o.addRoot(path), data, 0644)
}

func (o *OsFs) ReadFileStream(path string) (io.ReadCloser, error) {
	return os.Open(o.addRoot(path))
}

func (o *OsFs) WriteFileStream(path string, data io.Reader) error {
	create, err := os.Create(o.addRoot(path))

	if err != nil {
		return err
	}

	defer func(create *os.File) {
		_ = create.Close()
	}(create)

	_, err = io.Copy(create, data)

	return err
}

func (o *OsFs) RemoveFile(path string) error {
	return os.Remove(o.addRoot(path))
}

func (o *OsFs) FileSize(path string) (int64, error) {
	info, err := os.Stat(o.addRoot(path))

	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

func (o *OsFs) FileModTime(path string) (int64, error) {
	info, err := os.Stat(o.addRoot(path))

	if err != nil {
		return 0, err
	}

	return info.ModTime().Unix(), nil
}
