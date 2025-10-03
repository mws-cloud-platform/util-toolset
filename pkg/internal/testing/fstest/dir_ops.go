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
	ErrWalkDir        = consterr.Error("walking dir")
	ErrFileInfo       = consterr.Error("getting file info")
	ErrFileReading    = consterr.Error("reading file")
	ErrFileWriting    = consterr.Error("writing file")
	ErrRelPathGetting = consterr.Error("getting relative path")
	ErrPathCreating   = consterr.Error("creating path")
)

type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Helper()
}

func CopyDir(t TestingT, fromFS fs.ReadOnlyFS, toFS fs.WriteOnlyFS) {
	CopyDirWithPath(t, fromFS, toFS, ".", ".")
}

func CopyDirWithPath(t TestingT, fromFS fs.ReadOnlyFS, toFS fs.WriteOnlyFS, fromPath, toPath string) {
	t.Helper()
	CopyDirWithPathModify(t, fromFS, toFS, fromPath, func(p string) string {
		return path.Join(toPath, p)
	})
}

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

	require.NoError(t, err, fmt.Sprintf("%s: fromFS ('%s'): %s", ErrWalkDir, fromPath, err))
}

func CompareDirs(t TestingT, expectedFs, actualFs fs.FS, expectedDir, actualDir string) {
	t.Helper()

	expContent, err := expectedFs.ReadDir(expectedDir)
	require.NoError(t, err, fmt.Sprintf("compare dirs error, path '%s': %s", expectedDir, err))
	actualContent, err := actualFs.ReadDir(actualDir)
	require.NoError(t, err, fmt.Sprintf("compare dirs error, path '%s': %s", actualDir, err))

	uniqueExp, uniqueActual := findUniqueNames(expContent, actualContent)
	require.Equal(t, len(uniqueExp), 0, "expected no unique items among expected content")
	require.Equal(t, len(uniqueActual), 0, "expected no unique items among actual content")

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

func CompareFiles(t TestingT, expectedFs, actualFs fs.FS, expPath, actualPath string) {
	t.Helper()

	expContent, err := fs.ReadFile(expectedFs, expPath)
	require.NoError(t, err, fmt.Sprintf("compare files error, expectedFs ('%s'): %s", expPath, err))
	actualContent, err := fs.ReadFile(actualFs, actualPath)
	require.NoError(t, err, fmt.Sprintf("compare files error, actualFs ('%s'): %s", actualPath, err))

	require.Equal(
		t,
		string(expContent),
		string(actualContent),
		fmt.Sprintf("compare files error, files paths: expected: '%s', actual: '%s', contents: expected: %s; actual: %s",
			expPath,
			actualPath,
			expContent,
			actualContent))
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
