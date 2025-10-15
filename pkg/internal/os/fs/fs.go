// Package fs provides an abstraction over the file system with read and write capabilities.
package fs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/afero"
)

// FS provides read and write access to the file system.
type FS interface {
	fs.ReadDirFS
	// We should use [ReadOnlyFS] directly, but mock generator doesn't support embedding alias to an interface.

	WriteOnlyFS
}

type (
	// ReadOnlyFS provides read-only access to the file system.
	// Currently, it contains only two methods (Open and ReadDir).
	ReadOnlyFS = fs.ReadDirFS

	// WritableFile is similar to [fs.File] which also includes almost all
	// methods from [os.File]. This interface allows to perform write operations
	// in streaming manner.
	WritableFile = afero.File
)

// WriteOnlyFS provides write access to the file system.
type WriteOnlyFS interface {
	OpenFile(name string, flag int, perm fs.FileMode) (WritableFile, error)
	MkdirAll(path string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
	Rename(src, dst string) error
	Remove(name string) error
	RemoveAll(path string) error
}

// ListFS is an optional interface that can be implemented by [FS] implementations
// Currently only [FS] returned by [NewMapFS] implements this interface.
// You should not use this interface in production code.
type ListFS interface {
	List() ([]fs.FileInfo, error)
}

// Option is a function that configures an [FS].
type Option func(fs FS) FS

// NewFS returns FS with user-defined options
//
// Note: options order matters.
// Options are applied in natural order, which means that
// the first option will be the outermost, while the last option
// will be the innermost wrapper around the FS.
func NewFS(fs FS, options ...Option) FS {
	for i := len(options) - 1; i >= 0; i-- {
		fs = options[i](fs)
	}

	return fs
}

// NewRealFS returns [FS] built on top of OS file system.
func NewRealFS() FS {
	return &aferoFS{a: afero.NewOsFs()}
}

// NewMapFS returns [FS] built on top of in-memory map file system.
func NewMapFS() FS {
	return &mapFS{&aferoFS{a: afero.NewMemMapFs()}}
}

// NewRecommended returns [FS] with recommended and user-defined options.
func NewRecommended(f FS, options ...Option) FS {
	return NewFS(f, append(options,
		WithDirCreate(fs.ModePerm),
		WithAtomicWrite(),
		WithUnique(),
	)...)
}

// NewRecommendedReal returns [FS] built on top of OS file system,
// with recommended and user-defined options.
func NewRecommendedReal(options ...Option) FS {
	return NewRecommended(NewRealFS(), options...)
}

type aferoFS struct{ a afero.Fs }

var _ FS = (*aferoFS)(nil)

type fsOnly struct{ fs.FS }

func (a *aferoFS) Open(name string) (fs.File, error) {
	return a.a.Open(filepath.Clean(name))
}

func (a *aferoFS) OpenFile(name string, flag int, perm fs.FileMode) (WritableFile, error) {
	return a.a.OpenFile(filepath.Clean(name), flag, perm)
}

func (a *aferoFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(fsOnly{a}, filepath.Clean(name))
}

func (a *aferoFS) MkdirAll(path string, perm fs.FileMode) error {
	return a.a.MkdirAll(filepath.Clean(path), perm)
}

func (a *aferoFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return afero.WriteFile(a.a, filepath.Clean(name), data, perm)
}

func (a *aferoFS) Rename(src, dst string) error {
	return a.a.Rename(filepath.Clean(src), filepath.Clean(dst))
}

func (a *aferoFS) Remove(name string) error {
	return a.a.Remove(filepath.Clean(name))
}

func (a *aferoFS) RemoveAll(path string) error {
	return a.a.RemoveAll(filepath.Clean(path))
}

type mapFS struct{ *aferoFS }

var _ ListFS = (*mapFS)(nil)

type pathFile struct {
	fs.FileInfo
	path string
}

func (f pathFile) Name() string {
	return f.path
}

// List returns all the files in the FS. This function should
// be used only in tests, do not use this in production.
func (m *mapFS) List() ([]fs.FileInfo, error) {
	var files []fs.FileInfo

	aferoFs := m.a
	// error ignore is intended, we only care about getting all valid files
	walkFn := func(path string, info fs.FileInfo, _ error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		files = append(files, pathFile{FileInfo: info, path: path})
		return nil
	}
	// we run walk from two roots to get all the files,
	// because MapFS doesn't provide a single root
	_ = afero.Walk(aferoFs, ".", walkFn)
	_ = afero.Walk(aferoFs, "/", walkFn)

	slices.SortFunc(files, func(l, r fs.FileInfo) int {
		return strings.Compare(l.Name(), r.Name())
	})

	return files, nil
}

// ReadFile is an alias for [fs.ReadFile].
func ReadFile(f ReadOnlyFS, name string) ([]byte, error) {
	return fs.ReadFile(f, name)
}

// CopyFS copies all files and directories from src to dst.
// It returns an error if any operation fails.
// If dstFS already contains some files, they will not be removed or overwritten.
// Only regular files and directories are copied. Symlinks and other file types
// will cause an error.
// In other aspects it behaves similarly to [os.CopyFS].
func CopyFS(dst WriteOnlyFS, src fs.FS) error {
	return fs.WalkDir(src, ".", func(path string, d fs.DirEntry, e error) (rErr error) {
		if e != nil {
			return e
		}

		switch d.Type() { // Type returns only the type bits, without the permission bits
		case fs.ModeDir:
			return dst.MkdirAll(path, fs.ModePerm)
		case 0: // no type bits set, means regular file
			r, err := src.Open(path)
			if err != nil {
				return err
			}

			defer closeWithErr(r, &rErr)

			info, err := r.Stat()
			if err != nil {
				return err
			}

			w, err := dst.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o666|info.Mode().Perm())
			if err != nil {
				return err
			}

			if _, err = io.Copy(w, r); err != nil {
				return &fs.PathError{Op: "Copy", Path: path, Err: err}
			}
			return err
		default:
			return &fs.PathError{Op: "CopyFS", Path: path, Err: fs.ErrInvalid}
		}
	})
}

func closeWithErr(r io.Closer, e *error) {
	cerr := r.Close()
	switch {
	case cerr == nil:
		return
	case *e == nil:
		*e = fmt.Errorf("close error: %w", cerr)
	default:
		*e = fmt.Errorf("close error: %w: while returning error %w", cerr, *e)
	}
}
