package fstest

import (
	iofs "io/fs"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/util-toolset/pkg/internal/os/fs"
	"go.mws.cloud/util-toolset/pkg/testing/golden"
)

const (
	testDir     = "testdata/dir_ops"
	copyDir     = "copy_dir"
	compareDirs = "compare_dirs"
	inputDir    = "input"
	goldenDir   = "golden"
	expectedDir = "expected"
)

type stubT struct {
	failed bool
}

func (*stubT) Errorf(string, ...any) {}

func (s *stubT) FailNow() {
	s.failed = true
}

func (*stubT) Helper() {}

func TestCopyDir(t *testing.T) {
	for _, tc := range []struct {
		name string
		dir  string
	}{
		{
			name: "RealToMapSingleFile",
			dir:  "single_file",
		},
		{
			name: "RealToMapSeveralFiles",
			dir:  "several_files",
		},
		{
			name: "RealToMapComplex",
			dir:  "complex",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			inputPath := path.Join(testDir, copyDir, inputDir, tc.dir)
			fromFS := fs.NewRecommendedReal(fs.WithBaseDir(inputPath))
			toFS := fs.NewMapFS()
			CopyDir(t, fromFS, toFS)

			goldenPath := path.Join(testDir, copyDir, goldenDir, tc.dir)
			goldenDir := golden.NewDir(t, golden.WithPath(goldenPath))
			validateDir(t, goldenDir, toFS, fromFS)
		})
	}
}

func TestCopyDirWithPath(t *testing.T) {
	for _, tc := range []struct {
		name     string
		dir      string
		fromFS   func(baseDir string) fs.FS
		toFS     func() fs.FS
		fromPath string
		toPath   string
	}{
		{
			name:     "SourceFSHasSpecificRoot",
			dir:      "source_fs_has_specific_root",
			fromPath: "root",
			toPath:   ".",
		},

		{
			name:     "DestinationFSHasSpecificRoot",
			dir:      "destination_fs_has_specific_root",
			fromPath: ".",
			toPath:   "root/another_dir",
		},

		{
			name:     "BothFSHaveSpecificRoot",
			dir:      "both_fs_have_specific_root",
			fromPath: "root/another_source_dir",
			toPath:   "root/another_dest_dir",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			inputPath := path.Join(testDir, copyDir, inputDir, tc.dir)
			fromFS := fs.NewRecommendedReal(fs.WithBaseDir(inputPath))
			toFS := fs.NewMapFS()
			CopyDirWithPath(t, fromFS, toFS, tc.fromPath, tc.toPath)

			goldenPath := path.Join(testDir, copyDir, goldenDir, tc.dir)
			goldenDirWithPath := golden.NewDir(t, golden.WithPath(goldenPath))
			validateDirWithPath(t, goldenDirWithPath, toFS, fromFS, tc.fromPath, tc.toPath)
		})
	}
}

func TestCompareDirs(t *testing.T) {
	for _, tc := range []struct {
		name    string
		baseDir string
		dir     string
	}{
		{
			name:    "CompareDirsOneFile",
			baseDir: "one_file",
			dir:     "dir",
		},
		{
			name:    "CompareDirsComplex",
			baseDir: "complex",
			dir:     "dir",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gPath := path.Join(testDir, compareDirs, goldenDir, tc.baseDir)
			gFS := fs.NewRecommendedReal(fs.WithBaseDir(gPath))

			actualBaseDir := path.Join(testDir, compareDirs, inputDir, tc.baseDir)
			actualFS := fs.NewRecommendedReal(fs.WithBaseDir(actualBaseDir))

			CompareDirs(t, gFS, actualFS, tc.dir, tc.dir)
		})
	}
}

func TestCompareDirsError(t *testing.T) {
	tcDirName := "dir"

	for _, tc := range []struct {
		name    string
		baseDir string
	}{
		{
			name:    "CompareDirsErrorNamesDiffer",
			baseDir: "error_files_names_differ",
		},
		{
			name:    "CompareDirsErrorFilesContentsDiffer",
			baseDir: "error_files_contents_differ",
		},
		{
			name:    "CompareDirsErrorExpectedHasUniqueFile",
			baseDir: "error_expected_has_unique",
		},
		{
			name:    "CompareDirsErrorActualHasUniqueFile",
			baseDir: "error_actual_has_unique",
		},
		{
			name:    "CompareDirsErrorBothHaveUniqueFiles",
			baseDir: "error_both_have_unique",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			expectedPath := path.Join(testDir, compareDirs, expectedDir, tc.baseDir)
			expFS := fs.NewRecommendedReal(fs.WithBaseDir(expectedPath))

			sutDir := path.Join(testDir, compareDirs, inputDir, tc.baseDir)
			sutFS := fs.NewRecommendedReal(fs.WithBaseDir(sutDir))

			stub := stubT{}
			CompareDirs(&stub, expFS, sutFS, tcDirName, tcDirName)

			require.True(t, stub.failed, "error was expected")
		})
	}
}

func validateDir(t *testing.T, goldenDir *golden.Dir, toFS, fromFS fs.FS) {
	t.Helper()

	validateDirWithPath(t, goldenDir, toFS, fromFS, ".", ".")
}

func validateDirWithPath(t *testing.T, goldenDir *golden.Dir, toFS, fromFS fs.FS, fromPath, toPath string) {
	t.Helper()

	err := iofs.WalkDir(toFS, toPath, func(ePath string, e iofs.DirEntry, err error) error {
		require.NoError(t, err, "error walking directory")

		if e.IsDir() {
			return nil
		}

		actual, err := fs.ReadFile(toFS, ePath)
		require.NoError(t, err, "reading file error")
		goldenDir.Bytes(t, ePath, actual)
		relPath, err := filepath.Rel(toPath, ePath)
		require.NoError(t, err, "error getting relative path")
		fromEPath := path.Join(fromPath, relPath)
		inputFileStat, err := iofs.Stat(fromFS, fromEPath)
		require.NoError(t, err, "error getting file stat")
		outputFileStat, err := iofs.Stat(toFS, ePath)
		require.NoError(t, err, "error getting file stat")
		require.Equal(t, inputFileStat.Mode(), outputFileStat.Mode(), "files modes are not equal")

		return nil
	})

	require.NoError(t, err, "error in walk dir function")
}
