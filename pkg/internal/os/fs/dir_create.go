package fs

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type dirCreate struct {
	FS

	mu      sync.Mutex
	dirs    map[string]struct{}
	dirMode fs.FileMode
}

// WithDirCreate is an option for [NewFS] that wraps the [FS] so that WriteFile,
// Rename and OpenFile all try to create parent directories with the specified
// mode if they do not exist.
func WithDirCreate(dirMode fs.FileMode) Option {
	return func(fs FS) FS {
		return &dirCreate{
			FS:      fs,
			dirs:    map[string]struct{}{},
			dirMode: dirMode,
		}
	}
}

func (d *dirCreate) MkdirAll(path string, perm fs.FileMode) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.dirs[path]; ok {
		return nil
	}

	if err := d.FS.MkdirAll(path, perm); err != nil {
		return err
	}

	d.dirs[path] = struct{}{}
	return nil
}

func (d *dirCreate) RemoveAll(name string) error {
	if err := d.FS.RemoveAll(name); err != nil {
		return err
	}
	d.removeDir(name)
	return nil
}

func (d *dirCreate) Remove(name string) error {
	if err := d.FS.Remove(name); err != nil {
		return err
	}
	d.removeDir(name)
	return nil
}

func (d *dirCreate) OpenFile(name string, flag int, perm os.FileMode) (WritableFile, error) {
	if flag&os.O_CREATE != 0 {
		if err := d.MkdirAll(path.Dir(name), d.dirMode); err != nil {
			return nil, err
		}
	}

	return d.FS.OpenFile(name, flag, perm)
}

func (d *dirCreate) Rename(src, dst string) error {
	if err := d.MkdirAll(path.Dir(dst), d.dirMode); err != nil {
		return err
	}

	d.mu.Lock()
	if _, ok := d.dirs[src]; ok {
		d.dirs[dst] = struct{}{}
		delete(d.dirs, src)
	}
	d.mu.Unlock()

	return d.FS.Rename(src, dst)
}

func (d *dirCreate) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if err := d.MkdirAll(path.Dir(name), d.dirMode); err != nil {
		return err
	}

	return d.FS.WriteFile(name, data, perm)
}

func (d *dirCreate) removeDir(name string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for k := range d.dirs {
		if isSubPath(name, k) {
			delete(d.dirs, k)
		}
	}
}

func isSubPath(basepath, targpath string) bool {
	result, err := filepath.Rel(basepath, targpath)
	return err == nil && !strings.Contains(result, "..")
}
