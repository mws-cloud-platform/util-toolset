package fs

import (
	"bytes"
	"io/fs"
	"os"
)

type changedOnly struct{ FS }

// WithChangedOnly is an option for [NewFS] that wraps the [FS] so that WriteFile
// only writes if the content has changed.
func WithChangedOnly() Option {
	return func(fs FS) FS {
		return &changedOnly{FS: fs}
	}
}

func (c *changedOnly) WriteFile(name string, data []byte, m fs.FileMode) error {
	actual, err := fs.ReadFile(c.FS, name)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if bytes.Equal(actual, data) {
		return nil
	}
	return c.FS.WriteFile(name, data, m)
}
