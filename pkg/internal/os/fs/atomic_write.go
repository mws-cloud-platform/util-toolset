package fs

import (
	"errors"
	"io/fs"
	"path"
)

type atomicWrite struct {
	FS
	dir string
}

// WithAtomicWrite is an option for [NewFS] that wraps the [FS] so that WriteFile
// is atomic. It writes every file to NAME.tmp first, and then moves it to
// dst location.
func WithAtomicWrite() Option {
	return func(fs FS) FS {
		return &atomicWrite{
			FS: fs,
		}
	}
}

// WithAtomicWriteCustomDir is an option for [NewFS] that wraps the [FS] so that WriteFile
// is atomic. It writes every file to NAME.tmp first, and then moves it to
// dst location. It also allows to define a directory which will be used for
// temporary files. The selected directory must belong to the same mounted file
// system as the files that will be created using this [FS] instance.
//
// See https://linux.die.net/man/2/rename for more details.
func WithAtomicWriteCustomDir(dir string) Option {
	return func(fs FS) FS {
		return &atomicWrite{
			FS:  fs,
			dir: dir,
		}
	}
}

func (a *atomicWrite) WriteFile(name string, data []byte, perm fs.FileMode) error {
	tmpName := name + ".tmp"
	if a.dir != "" {
		tmpName = path.Join(a.dir, tmpName)
	}
	if err := a.FS.WriteFile(tmpName, data, perm); err != nil {
		return errors.Join(err, a.Remove(tmpName))
	}

	// Even within the same directory, on non-Unix platforms Rename is not an atomic operation.
	if err := a.Rename(tmpName, name); err != nil {
		return errors.Join(err, a.Remove(tmpName))
	}

	return nil
}
