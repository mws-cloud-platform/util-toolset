package golden

import (
	"os"
	"path"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mws-cloud-platform/util-toolset/pkg/internal/os/fs"
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

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f := strconv.Itoa(i)
			err := mfs.WriteFile(path.Join(goldenPath, f), []byte(f+f), 0644)
			require.NoError(t, err)
			d.String(t, f, f)
		}()
	}

	wg.Wait()
}
