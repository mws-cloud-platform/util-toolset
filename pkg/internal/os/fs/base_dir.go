package fs

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// baseDir adds prefix dir to all names
type baseDir struct {
	FS
	dir string
}

// WithBaseDir is an option for [NewFS] that wraps the [FS] so that all paths are
// relative to the specified base directory.
func WithBaseDir(dir string) Option {
	return func(fs FS) FS {
		return &baseDir{
			FS:  fs,
			dir: dir,
		}
	}
}

func (b *baseDir) Open(name string) (_ fs.File, err error) {
	if name, err = b.path(name); err != nil {
		return nil, &os.PathError{Op: "open", Path: name, Err: err}
	}

	return b.FS.Open(name)
}

func (b *baseDir) OpenFile(name string, flag int, mode os.FileMode) (f WritableFile, err error) {
	if name, err = b.path(name); err != nil {
		return nil, &os.PathError{Op: "openfile", Path: name, Err: err}
	}
	return b.FS.OpenFile(name, flag, mode)
}

func (b *baseDir) ReadDir(name string) (entries []fs.DirEntry, err error) {
	if name, err = b.path(name); err != nil {
		return nil, &os.PathError{Op: "readdir", Path: name, Err: err}
	}
	return b.FS.ReadDir(name)
}

func (b *baseDir) MkdirAll(name string, mode os.FileMode) (err error) {
	if name, err = b.path(name); err != nil {
		return &os.PathError{Op: "mkdir", Path: name, Err: err}
	}
	return b.FS.MkdirAll(name, mode)
}

func (b *baseDir) WriteFile(name string, data []byte, perm fs.FileMode) (err error) {
	if name, err = b.path(name); err != nil {
		return &os.PathError{Op: "write_file", Path: name, Err: err}
	}
	return b.FS.WriteFile(name, data, perm)
}

func (b *baseDir) Rename(src, dst string) (err error) {
	if src, err = b.path(src); err != nil {
		return &os.PathError{Op: "rename", Path: src, Err: err}
	}
	if dst, err = b.path(dst); err != nil {
		return &os.PathError{Op: "rename", Path: dst, Err: err}
	}

	return b.FS.Rename(src, dst)
}

func (b *baseDir) Remove(name string) (err error) {
	if name, err = b.path(name); err != nil {
		return &os.PathError{Op: "remove", Path: name, Err: err}
	}
	return b.FS.Remove(name)
}

func (b *baseDir) RemoveAll(name string) (err error) {
	if name, err = b.path(name); err != nil {
		return &os.PathError{Op: "remove_all", Path: name, Err: err}
	}
	return b.FS.RemoveAll(name)
}

func (b *baseDir) path(name string) (path string, err error) {
	if b.dir == "." || b.dir == "" {
		return name, nil
	}

	bpath := filepath.Clean(b.dir)
	path = filepath.Clean(filepath.Join(bpath, name))
	if !strings.HasPrefix(path, bpath) {
		return name, os.ErrNotExist
	}

	return path, nil
}
