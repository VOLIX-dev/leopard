package files

import (
	"errors"
	"io"
	"io/fs"
)

var registry = make(map[string]func(config any) (Driver, error))

type Driver interface {
	// Create a new file.
	Create(path string) (File, error)

	// Open a file.
	Open(path string) (File, error)

	// Remove a file.
	Remove(path string) error

	// Rename a file.
	Rename(oldPath, newPath string) error

	// Stat a file.
	Stat(path string) (fs.FileInfo, error)

	// MakeDir make a directory.
	MakeDir(path string) error

	// RemoveDir remove a directory.
	RemoveDir(path string) error

	// RenameDir rename a directory.
	RenameDir(oldPath, newPath string) error

	// StatDir stat a directory.
	StatDir(path string) (fs.FileInfo, error)

	// ListDir list the contents of a directory.
	ListDir(path string) ([]fs.DirEntry, error)

	// ReadFile read the contents of a file.
	ReadFile(path string) ([]byte, error)

	// WriteFile write the contents of a file.
	WriteFile(path string, data []byte) error

	// ReadFileStream read the contents of a file.
	ReadFileStream(path string) (io.ReadCloser, error)

	// WriteFileStream write the contents of a file.
	WriteFileStream(path string, data io.Reader) error

	// RemoveFile remove a file.
	RemoveFile(path string) error

	// FileSize get the file's size.
	FileSize(path string) (int64, error)

	// FileModTime get the file's last modification time.
	FileModTime(path string) (int64, error)
}

// Register register a driver.
func Register(name string, driverCreator func(config any) (Driver, error)) {
	registry[name] = driverCreator
}

// Get a driver.
func Get(name string, config any) (Driver, error) {
	driverCreator, ok := registry[name]
	if !ok {
		return nil, errors.New("driver not found")
	}
	return driverCreator(config)
}
