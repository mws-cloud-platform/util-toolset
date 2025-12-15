package golden

import (
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/util-toolset/pkg/internal/os/fs"
)

const dirTestBaseDir = "testdata"

func TestDirMkdir(t *testing.T) {
	dirName := "awesome"
	fullPath := path.Join(dirTestBaseDir, dirName)

	dir := NewDir(t, WithPath(dirTestBaseDir))

	dir.MkDir(t, dirName)
	t.Cleanup(func() {
		removeDir(t, fullPath)
	})

	info, err := os.Stat(fullPath)
	require.NoError(t, err)
	require.True(t, info.IsDir())
}

func createDir(t *testing.T, path string) {
	t.Helper()
	err := os.MkdirAll(path, os.ModePerm)
	require.NoError(t, err)
}

func removeDir(t *testing.T, path string) {
	t.Helper()
	err := os.RemoveAll(path)
	require.NoError(t, err)
}

func Test_DirParallelUpdate(t *testing.T) {
	const goldenPath = "./in-memory"
	mfs := fs.NewMapFS()
	d := NewDir(t, WithFS(mfs), WithPath(goldenPath))
	d.update = true

	t.Run("parallel", func(t *testing.T) {
		for i := range 100 {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				t.Parallel()

				f := strconv.Itoa(i)
				err := mfs.WriteFile(path.Join(goldenPath, f), []byte(f+f), 0o644)
				require.NoError(t, err)
				d.String(t, f, f)
			})
		}
	})
}
