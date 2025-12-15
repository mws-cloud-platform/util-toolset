package fstest_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/util-toolset/pkg/internal/os/fs"
	"go.mws.cloud/util-toolset/pkg/internal/testing/fstest"
)

func TestList(t *testing.T) {
	f := fs.NewMapFS()
	data := []byte("hello")
	require.NoError(t, f.WriteFile("x.txt", data, os.ModePerm))
	require.NoError(t, f.WriteFile("/y.txt", data, os.ModePerm))
	require.NoError(t, f.WriteFile("../z.txt", data, os.ModePerm))
	require.NoError(t, f.WriteFile("foo/bar/z.txt", nil, os.ModePerm))
	require.NoError(t, f.WriteFile("/baz/qux/w.txt", data, os.ModePerm))
	require.NoError(t, f.WriteFile("empty", nil, os.ModePerm))
	require.NoError(t, f.Rename("/baz/qux/w.txt", "/baz/qux/v.txt"))
	require.NoError(t, f.WriteFile("empty2", nil, os.ModePerm))
	require.NoError(t, f.Rename("empty2", "empty"))
	require.NoError(t, f.Remove("../z.txt"))

	names := fstest.Names(t, f)
	require.Equal(t, []string{"/baz/qux/v.txt", "/y.txt", "empty", "foo/bar/z.txt", "x.txt"}, names)
}
