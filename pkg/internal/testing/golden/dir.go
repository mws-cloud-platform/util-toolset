// Package golden provides helpers for working with golden files in tests.
package golden

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"sync"

	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/mws-cloud-platform/util-toolset/pkg/internal/os/fs"
)

// Dir is a helper for managing golden files in a directory
// under [fs.FS].
type Dir struct {
	mu               sync.Mutex
	registry         map[string]bool
	path             string
	recreateOnUpdate bool
	update           bool
	fs               fs.FS
}

// NewDir creates a directory where golden files will be stored.
// update field is initialized with current updateFlag value.
//
// Possible options:
// golden.WithRecreateOnUpdate() will recreate the golden directory when updating;
// golden.WithPath(customPath) set the path to the golden directory, the default is ./testdata/
//
//	dir := golden.NewDir(t, golden.WithPath("testdata/example/golden"), golden.WithRecreateOnUpdate())
func NewDir(t testing.T, opts ...DirOption) *Dir {
	x := &Dir{
		path:             "./testdata/",
		recreateOnUpdate: false,
		registry:         map[string]bool{},
		update:           *updateFlag,
		fs:               fs.NewRealFS(),
	}
	for _, v := range opts {
		v(x)
	}
	if x.recreateOnUpdate && x.update {
		require.NoError(t, x.fs.RemoveAll(x.path))
		require.NoError(t, x.fs.MkdirAll(x.path, os.ModePerm))
	}
	return x
}

func (d *Dir) fileName(x string) string {
	return path.Join(d.path, x)
}

func (d *Dir) writeFile(t testing.T, fileName string, actual []byte) {
	t.Helper()
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.registry[fileName] {
		require.Fail(t, "trying to write golden with name "+fileName+" multiple times")
	}
	d.registry[fileName] = true
	require.NoError(t, d.fs.WriteFile(fileName, actual, 0o644))
}

// MkDir creates a directory, fails the test in case of an error.
//
//	dir.MkDir(t, "sub_dir_name")
func (d *Dir) MkDir(t testing.T, name string) {
	t.Helper()
	require.NoError(t, d.fs.MkdirAll(d.fileName(name), os.ModePerm))
}

// Bytes requires that actual is equal to file content or updates file content with actual when flag -update is used.
//
//	dir.Bytes(t, "expected.txt", myBytes)
func (d *Dir) Bytes(t testing.T, fn string, actual []byte) {
	t.Helper()
	fileName := d.fileName(fn)
	if d.update {
		d.writeFile(t, fileName, actual)
		return
	}
	expected := readFileFromFS(t, d.fs, fileName)
	require.True(t, bytes.Equal(expected, actual), updateMessage)
}

// String requires that actual is equal to file content or updates file content with actual when flag -update is used.
//
//	dir.String(t, "expected.txt", myString)
func (d *Dir) String(t testing.T, fn, actual string) {
	t.Helper()
	fileName := d.fileName(fn)
	if d.update {
		d.writeFile(t, fileName, []byte(actual))
		return
	}
	expected := readFileFromFS(t, d.fs, fileName)
	require.Equal(t, string(expected), actual, updateMessage)
}

// JSONBytes formats actual json and requires that result is equal to the contents of the file.
//
//	dir.JSONBytes(t, "expected.json", myJsonBytes)
func (d *Dir) JSONBytes(t testing.T, fn string, actual []byte) {
	t.Helper()
	fileName := d.fileName(fn)
	if d.update {
		indented, err := json.MarshalIndent(json.RawMessage(actual), "", " ")
		require.NoError(t, err)
		d.writeFile(t, fileName, indented)
		return
	}
	expected := readFileFromFS(t, d.fs, fileName)
	require.JSONEq(t, string(expected), string(actual), updateMessage)
}

// JSON Marshals x and uses JSONBytes to compare result and contents of the file.
//
//	dir.JSON(t, "expected.json", myMarshalableStruct)
func (d *Dir) JSON(t testing.T, fn string, x any) {
	t.Helper()
	actual, err := json.Marshal(x)
	require.NoError(t, err)
	d.JSONBytes(t, fn, actual)
}
