package fs_test

import (
	"io"
	iofs "io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/mws-cloud-platform/util-toolset/pkg/internal/os/fs"
	fsmock "github.com/mws-cloud-platform/util-toolset/pkg/internal/os/fs/mock"
	"github.com/mws-cloud-platform/util-toolset/pkg/internal/testing/fstest"
)

func TestUniqueDirCreate(t *testing.T) {
	gm := gomock.NewController(t)
	mock := fsmock.NewMockFS(gm)

	unique, dirCreate := fs.WithUnique(), fs.WithDirCreate(os.ModePerm)
	suites := [][]fs.Option{
		{unique, dirCreate},
		{dirCreate, unique},
	}

	for _, options := range suites {
		f := fs.NewFS(mock, options...)

		mock.EXPECT().WriteFile("/foo/bar/x.txt", gomock.Any(), gomock.Any())
		mock.EXPECT().WriteFile("/foo/bar/y.txt", gomock.Any(), gomock.Any())
		mock.EXPECT().MkdirAll("/foo/bar", gomock.Any())

		require.NoError(t, f.WriteFile("/foo/bar/x.txt", nil, os.ModePerm))
		require.NoError(t, f.WriteFile("/foo/bar/y.txt", nil, os.ModePerm))
	}
}

func TestMapFS(t *testing.T) {
	f := fs.NewMapFS()

	require.NoError(t, f.WriteFile("./x.txt", []byte("x"), os.ModePerm))
	require.NoError(t, f.WriteFile("/y.txt", []byte("y"), os.ModePerm))

	suites := []struct {
		name     string
		expected []byte
	}{
		{"x.txt", []byte("x")},
		{"./x.txt", []byte("x")},
		{"/y.txt", []byte("y")},
	}
	for _, s := range suites {
		file, err := f.Open(s.name)
		require.NoError(t, err)

		actual, err := io.ReadAll(file)
		require.NoError(t, err)
		require.Equal(t, s.expected, actual)
		require.NoError(t, file.Close())
	}

	for _, name := range []string{"../x.txt", "./../y.txt"} {
		_, err := f.Open(name)
		require.ErrorIs(t, err, os.ErrNotExist)
	}

	require.Equal(t, []string{"/y.txt", "x.txt"}, fstest.Names(t, f))
}

func TestCopyMapFSReal(t *testing.T) {
	testFS(
		t,
		fs.NewFS(fs.NewRealFS(), fs.WithBaseDir(t.TempDir()), fs.WithDirCreate(0700)),
		fs.NewFS(fs.NewRealFS(), fs.WithBaseDir(t.TempDir()), fs.WithDirCreate(os.ModePerm)),
	)
}

func testFS(t *testing.T, srcFS, dstFS fs.FS) {
	// Test file writing with different paths (including relative and absolute)
	require.NoError(t, srcFS.WriteFile("./a.txt", []byte("a"), os.ModePerm))
	require.NoError(t, srcFS.WriteFile("/a.txt", []byte("a"), os.ModePerm))
	require.NoError(t, srcFS.WriteFile("a.txt", []byte("aa"), os.ModePerm))
	require.NoError(t, srcFS.WriteFile("source/b.txt", []byte("b"), os.ModePerm))
	require.NoError(t, srcFS.WriteFile("/source/b.txt", []byte("b"), os.ModePerm))
	require.NoError(t, srcFS.WriteFile("./source/b.txt", []byte("b"), os.ModePerm))

	// Wrap the source FS with a base directory
	sourceFs := fs.NewFS(srcFS, fs.WithBaseDir("/source"))

	// This writes should be allowed and correctly resolved
	require.NoError(t, sourceFs.WriteFile("c.txt", []byte("c"), os.ModePerm))
	require.NoError(t, sourceFs.WriteFile("/c.txt", []byte("c"), os.ModePerm))
	require.NoError(t, sourceFs.WriteFile("./c.txt", []byte("c"), os.ModePerm))

	// This writes should fail due to escaping the base directory
	require.Error(t, sourceFs.WriteFile("../d.txt", []byte("d"), os.ModePerm))
	require.Error(t, sourceFs.WriteFile("../../d.txt", []byte("d"), os.ModePerm))
	require.Error(t, sourceFs.WriteFile("source/../../d.txt", []byte("d"), os.ModePerm))

	// This write should be allowed as it resolves within the base directory
	require.NoError(t, sourceFs.WriteFile("source/../d.txt", []byte("d"), os.ModePerm))

	require.Equal(
		t,
		[]string{"a.txt", "source/b.txt", "source/c.txt", "source/d.txt"},
		collectElements(t, srcFS),
	)

	mountedFS := fs.NewFS(dstFS, fs.WithBaseDir("/mounted"))

	err := fs.CopyFS(mountedFS, sourceFs)
	require.NoError(t, err)

	suites := []struct {
		name     string
		expected []byte
	}{
		{"b.txt", []byte("b")},
		{"/c.txt", []byte("c")},
		{"./random/../d.txt", []byte("d")},
	}
	for _, s := range suites {
		file := mustValue(mountedFS.OpenFile(s.name, os.O_RDONLY, 0))(t)

		actual := mustValue(io.ReadAll(file))(t)
		require.Equal(t, s.expected, actual)
		require.NoError(t, file.Close())
	}

	require.NoError(t, mountedFS.Rename("d.txt", "newfolder/e.txt"))
	require.Equal(
		t,
		[]string{"mounted/b.txt", "mounted/c.txt", "mounted/newfolder/e.txt"},
		collectElements(t, dstFS),
	)

	require.NoError(t, mountedFS.Remove("b.txt"))
	require.Equal(
		t,
		[]string{"mounted/c.txt", "mounted/newfolder/e.txt"},
		collectElements(t, dstFS),
	)

	require.NoError(t, mountedFS.RemoveAll("."))
	require.Zero(t, collectElements(t, dstFS))
}

func collectElements(t *testing.T, dfs fs.FS) []string {
	var elements []string

	require.NoError(t, iofs.WalkDir(dfs, ".", func(path string, d iofs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			elements = append(elements, path)
		}

		return nil
	}))

	return elements
}

func mustValue[V any](v V, err error) func(t *testing.T) V {
	return func(t *testing.T) V {
		t.Helper()

		require.NoError(t, err)

		return v
	}
}
