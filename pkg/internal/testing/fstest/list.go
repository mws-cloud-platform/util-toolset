package fstest

import (
	iofs "io/fs"

	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"go.mws.cloud/util-toolset/pkg/internal/os/fs"
)

// List returns all the files in the FS.
// Usable only on FS that implements ListFS.
func List(t testing.T, f fs.FS) []iofs.FileInfo {
	t.Helper()

	listFS, ok := f.(fs.ListFS)
	require.True(t, ok)

	list, err := listFS.List()
	require.NoError(t, err)

	return list
}

// Names returns full names of all the files in the FS.
// Usable only on FS that implements ListFS.
func Names(t testing.T, f fs.FS) []string {
	t.Helper()

	files := List(t, f)
	names := make([]string, len(files))
	for i, file := range files {
		names[i] = file.Name()
	}

	return names
}
