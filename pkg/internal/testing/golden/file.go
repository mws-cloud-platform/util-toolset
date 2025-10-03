package golden

import (
	"bytes"
	"os"

	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/mws-cloud-platform/util-toolset/pkg/internal/os/fs"
)

func readFileFromFS(t testing.T, f fs.FS, fileName string) []byte {
	t.Helper()
	data, err := fs.ReadFile(f, fileName)
	if os.IsNotExist(err) {
		require.Fail(t, "File "+fileName+" does not exist on disk. You should probably run your test with -update flag to create it")
	}
	require.NoError(t, err)
	return data
}

func readFile(t testing.T, fileName string) []byte {
	t.Helper()
	return readFileFromFS(t, fs.NewRealFS(), fileName)
}

// Bytes requires that actual is equal to file content or updates file content with actual when flag -update is used.
//
//	golden.Bytes(t, "expected.txt", myBytes)
func Bytes(t testing.T, fileName string, actual []byte) {
	t.Helper()
	if IsUpdate() {
		require.NoError(t, os.WriteFile(fileName, actual, 0644))
		return
	}
	expected := readFile(t, fileName)
	require.True(t, bytes.Equal(expected, actual), updateMessage)
}

// String requires that actual is equal to file content or updates file content with actual when flag -update is used.
//
//	golden.String(t, "expected.txt", myString)
func String(t testing.T, fileName string, actual string) {
	t.Helper()
	if IsUpdate() {
		require.NoError(t, os.WriteFile(fileName, []byte(actual), 0644))
		return
	}
	expected := readFile(t, fileName)
	require.Equal(t, string(expected), actual, updateMessage)
}
