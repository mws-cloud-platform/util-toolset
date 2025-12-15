package fs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"go.mws.cloud/util-toolset/pkg/internal/os/fs"
	fsmock "go.mws.cloud/util-toolset/pkg/internal/os/fs/mock"
)

func TestBaseDir(t *testing.T) {
	gm := gomock.NewController(t)
	mock := fsmock.NewMockFS(gm)
	baseDir := fs.NewFS(mock, fs.WithBaseDir("/foo"))

	mock.EXPECT().Open("/foo/baz/qux")
	_, _ = baseDir.Open("baz/qux")

	mock.EXPECT().ReadDir("/foo/bar")
	_, _ = baseDir.ReadDir("bar")

	mock.EXPECT().MkdirAll("/foo/test", gomock.Any())
	_ = baseDir.MkdirAll("test/hello/..", os.ModePerm)

	require.Error(t, baseDir.WriteFile("../../x.txt", nil, os.ModePerm))

	mock.EXPECT().Rename("/foo/x.txt", "/foo/y.txt")
	_ = baseDir.Rename("not/../what/../x.txt", "./y.txt")

	mock.EXPECT().Remove("/foo/bar/x.txt")
	_ = baseDir.Remove("not/../bar/x.txt")

	mock.EXPECT().RemoveAll("/foo/y.txt")
	_ = baseDir.RemoveAll("../bar/../foo/y.txt")
}
