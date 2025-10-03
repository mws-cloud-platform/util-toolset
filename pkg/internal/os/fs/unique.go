package fs

import (
	"io/fs"
	"sync"
)

// NonUniqueError is returned when trying to create a file or directory
// with a name that already exists in the unique wrapper.
type NonUniqueError struct {
	Name string
}

func (e *NonUniqueError) Error() string {
	return "trying to write " + e.Name + " more than once"
}

// unique doesn't allow to create more than one file with same name
type unique struct {
	FS

	mu    sync.Mutex
	files map[string]struct{}
}

// WithUnique is an option for [NewFS] that wraps the [FS] so that MkdirAll,
// WriteFile, and Rename return an error if the target path was already used by
// this [FS] instance.
func WithUnique() Option {
	return func(fs FS) FS {
		return &unique{
			FS: fs,

			files: make(map[string]struct{}),
		}
	}
}

func (u *unique) MkdirAll(path string, perm fs.FileMode) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, ok := u.files[path]; ok {
		return &NonUniqueError{Name: path}
	}
	if err := u.FS.MkdirAll(path, perm); err != nil {
		return err
	}

	u.files[path] = struct{}{}
	return nil
}

func (u *unique) WriteFile(name string, data []byte, perm fs.FileMode) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, ok := u.files[name]; ok {
		return &NonUniqueError{Name: name}
	}
	if err := u.FS.WriteFile(name, data, perm); err != nil {
		return err
	}

	u.files[name] = struct{}{}
	return nil
}

func (u *unique) Rename(src, dst string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, ok := u.files[dst]; ok {
		return &NonUniqueError{Name: dst}
	}
	if err := u.FS.Rename(src, dst); err != nil {
		return err
	}

	delete(u.files, src)
	u.files[dst] = struct{}{}
	return nil
}

func (u *unique) Remove(name string) error {
	if err := u.FS.Remove(name); err != nil {
		return err
	}

	u.mu.Lock()
	delete(u.files, name)
	u.mu.Unlock()
	return nil
}
