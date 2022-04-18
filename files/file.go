package files

import (
	"io"
	"io/fs"
)

type File interface {
	io.ReadWriteCloser

	// Stat returns the FileInfo structure describing file.
	Stat() (fs.FileInfo, error)
}
