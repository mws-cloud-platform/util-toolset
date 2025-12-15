package golden

import (
	"io/fs"
	"os"

	"github.com/mitchellh/go-testing-interface"

	devpfs "go.mws.cloud/util-toolset/pkg/internal/os/fs"
)

// FS is a filesystem implementation for code generation tests.
type FS struct {
	t   testing.T
	dir *Dir
}

var _ devpfs.FS = (*FS)(nil)

// NewCodegenFS creates a new filesystem instance for code generation tests.
func NewCodegenFS(t testing.T, d *Dir) devpfs.FS {
	return &FS{t: t, dir: d}
}

// MkdirAll creates a directory named path, along with any necessary parents, and returns nil.
func (fs *FS) MkdirAll(path string, mode fs.FileMode) error {
	return os.MkdirAll(fs.dir.fileName(path), mode)
}

// Rename renames (moves) file src to dst.
func (fs *FS) Rename(src, dst string) error {
	return os.Rename(fs.dir.fileName(src), fs.dir.fileName(dst))
}

// WriteFile writes data content to the named file.
func (fs *FS) WriteFile(name string, data []byte, _ fs.FileMode) error {
	fs.dir.Bytes(fs.t, name, data)
	return nil
}

// Open opens the named file for reading.
func (fs *FS) Open(name string) (fs.File, error) {
	return os.Open(fs.dir.fileName(name))
}

// OpenFile is the generalized open call that supports both read and write modes.
func (fs *FS) OpenFile(name string, flag int, perm os.FileMode) (devpfs.WritableFile, error) {
	return os.OpenFile(fs.dir.fileName(name), flag, perm)
}

// ReadDir reads the named directory and returns a list of directory entries.
func (fs *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(fs.dir.fileName(name))
}

// Remove removes the named file or (empty) directory.
func (fs *FS) Remove(name string) error {
	return os.Remove(fs.dir.fileName(name))
}

// RemoveAll removes path and any children it contains.
func (fs *FS) RemoveAll(path string) error {
	return os.RemoveAll(fs.dir.fileName(path))
}
