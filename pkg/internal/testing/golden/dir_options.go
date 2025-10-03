package golden

import "github.com/mws-cloud-platform/util-toolset/pkg/internal/os/fs"

type DirOption func(*Dir)

// WithRecreateOnUpdate option makes that directories will be deleted and created when called with -update
func WithRecreateOnUpdate() DirOption {
	return func(d *Dir) {
		d.recreateOnUpdate = true
	}
}

// WithPath set the path, where golden files will be checked/created.
// Default ./testdata/
func WithPath(p string) DirOption {
	return func(d *Dir) {
		d.path = p
	}
}

// Set the FS
func WithFS(fs fs.FS) DirOption {
	return func(d *Dir) {
		d.fs = fs
	}
}
