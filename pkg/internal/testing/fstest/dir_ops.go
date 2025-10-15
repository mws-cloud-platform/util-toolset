// Package fstest provides utilities for tests working with file systems operations.
package fstest

import (
	"fmt"
	iofs "io/fs"
	"path"
	"path/filepath"

	"github.com/stretchr/testify/require"

	"github.com/mws-cloud-platform/util-toolset/pkg/internal/os/fs"
	"github.com/mws-cloud-platform/util-toolset/pkg/utils/consterr"
)

const (
	// ErrWalkDir indicates an error occurred while walking a directory.
	ErrWalkDir = consterr.Error("walking dir")

	// ErrFileInfo indicates an error occurred while getting file info.
	ErrFileInfo = consterr.Error("getting file info")

	// ErrFileReading indicates an error occurred while reading a file.
	ErrFileReading = consterr.Error("reading file")

	// ErrFileWriting indicates an error occurred while writing a file.
	ErrFileWriting = consterr.Error("writing file")

	// ErrRelPathGetting indicates an error occurred while getting a relative path.
	ErrRelPathGetting = consterr.Error("getting relative path")

	// ErrPathCreating indicates an error occurred while creating a path.
	ErrPathCreating = consterr.Error("creating path")
)

// TestingT is an interface that can be used in place of *testing.T
// in test helper functions.
type TestingT interface {
	Errorf(format string, args ...any)
	FailNow()
	Helper()
}

// CopyDir copies a directory from one file system to another,
// placing it under the root (".") in the target file system.
func CopyDir(t TestingT, fromFS fs.ReadOnlyFS, toFS fs.WriteOnlyFS) {
	CopyDirWithPath(t, fromFS, toFS, ".", ".")
}

// CopyDirWithPath copies a directory from one file system to another,
// placing it under the specified toPath in the target file system.
func CopyDirWithPath(t TestingT, fromFS fs.ReadOnlyFS, toFS fs.WriteOnlyFS, fromPath, toPath string) {
	t.Helper()
	CopyDirWithPathModify(t, fromFS, toFS, fromPath, func(p string) string {
		return path.Join(toPath, p)
	})
}

// CopyDirWithPathModify copies a directory from one file system to another,
// allowing modification of the target path using the modifyPath function.
func CopyDirWithPathModify(t TestingT, fromFS fs.ReadOnlyFS, toFS fs.WriteOnlyFS, fromPath string, modifyPath func(string) string) {
	t.Helper()

	err := iofs.WalkDir(fromFS, fromPath, func(ePath string, e iofs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("%w: fromFs (file '%s'): %w", ErrWalkDir, fromPath, err)
		}
		entryInfo, err := e.Info()
		if err != nil {
			return fmt.Errorf("%w: fromFs (file '%s'): %w", ErrFileInfo, fromPath, err)
		}

		relEPath, err := filepath.Rel(fromPath, ePath)
		if err != nil {
			return fmt.Errorf("%w: fromFs (file '%s'): %w", ErrRelPathGetting, relEPath, err)
		}
		targetPath := modifyPath(relEPath)
		if e.IsDir() {
			err = toFS.MkdirAll(targetPath, entryInfo.Mode())
			if err != nil {
				return fmt.Errorf("%w in toFS ('%s'): %w", ErrPathCreating, targetPath, err)
			}
		} else {
			bytes, err2 := fs.ReadFile(fromFS, ePath)
			if err2 != nil {
				return fmt.Errorf("%w in fromFS ('%s'): %w", ErrFileReading, ePath, err2)
			}
			err2 = toFS.WriteFile(targetPath, bytes, entryInfo.Mode())
			if err2 != nil {
				return fmt.Errorf("%w in toFS ('%s'): %w", ErrFileWriting, targetPath, err2)
			}
		}
		return nil
	})

	require.NoError(t, err, "%s: fromFS ('%s'): %s", ErrWalkDir, fromPath, err)
}

// CompareDirs compares directory contents in two file systems.
func CompareDirs(t TestingT, expectedFs, actualFs fs.FS, expectedDir, actualDir string) {
	t.Helper()

	expContent, err := expectedFs.ReadDir(expectedDir)
	require.NoError(t, err, "compare dirs error, path '%s': %s", expectedDir, err)
	actualContent, err := actualFs.ReadDir(actualDir)
	require.NoError(t, err, "compare dirs error, path '%s': %s", actualDir, err)

	uniqueExp, uniqueActual := findUniqueNames(expContent, actualContent)
	require.Empty(t, uniqueExp, "expected no unique items among expected content")
	require.Empty(t, uniqueActual, "expected no unique items among actual content")

	for _, e := range expContent {
		ePath := path.Join(expectedDir, e.Name())
		aPath := path.Join(actualDir, e.Name())

		if e.IsDir() {
			CompareDirs(t, expectedFs, actualFs, ePath, aPath)
		} else {
			CompareFiles(t, expectedFs, actualFs, ePath, aPath)
		}
	}
}

// CompareFiles compares file contents in two file systems.
func CompareFiles(t TestingT, expectedFs, actualFs fs.FS, expPath, actualPath string) {
	t.Helper()

	expContent, err := fs.ReadFile(expectedFs, expPath)
	require.NoError(t, err, "compare files error, expectedFs ('%s'): %s", expPath, err)
	actualContent, err := fs.ReadFile(actualFs, actualPath)
	require.NoError(t, err, "compare files error, actualFs ('%s'): %s", actualPath, err)

	require.Equal(
		t,
		string(expContent),
		string(actualContent),
		"compare files error, files paths: expected: '%s', actual: '%s', contents: expected: %s; actual: %s",
		expPath,
		actualPath,
		expContent,
		actualContent,
	)
}

func findUniqueNames(expEntries, actualEntries []iofs.DirEntry) (uniqueExp, uniqueActual []string) {
	expNamesMap := make(map[string]struct{})
	actualNamesMap := make(map[string]struct{})

	for _, entry := range expEntries {
		expNamesMap[entry.Name()] = struct{}{}
	}
	for _, entry := range actualEntries {
		actualNamesMap[entry.Name()] = struct{}{}
	}

	for name := range expNamesMap {
		if _, found := actualNamesMap[name]; !found {
			uniqueExp = append(uniqueExp, name)
		}
	}
	for name := range actualNamesMap {
		if _, found := expNamesMap[name]; !found {
			uniqueActual = append(uniqueActual, name)
		}
	}

	return uniqueExp, uniqueActual
}
