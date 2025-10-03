package golden

import (
	"io/fs"
	"os"

	"github.com/mitchellh/go-testing-interface"

	devpfs "github.com/mws-cloud-platform/util-toolset/pkg/internal/os/fs"
)

type FS struct {
	t   testing.T
	dir *Dir
}

var _ devpfs.FS = &FS{}

func NewCodegenFS(t testing.T, d *Dir) devpfs.FS {
	return &FS{t: t, dir: d}
}

func (fs *FS) MkdirAll(path string, mode fs.FileMode) error {
	return os.MkdirAll(fs.dir.fileName(path), mode)
}

func (fs *FS) Rename(src, dst string) error {
	return os.Rename(fs.dir.fileName(src), fs.dir.fileName(dst))
}

func (fs *FS) WriteFile(name string, data []byte, _ fs.FileMode) error {
	fs.dir.Bytes(fs.t, name, data)
	return nil
}

func (fs *FS) Open(name string) (fs.File, error) {
	return os.Open(fs.dir.fileName(name))
}

func (fs *FS) OpenFile(name string, flag int, perm os.FileMode) (devpfs.WritableFile, error) {
	return os.OpenFile(fs.dir.fileName(name), flag, perm)
}

func (fs *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(fs.dir.fileName(name))
}

func (fs *FS) Remove(name string) error {
	return os.Remove(fs.dir.fileName(name))
}

func (fs *FS) RemoveAll(path string) error {
	return os.RemoveAll(fs.dir.fileName(path))
}
