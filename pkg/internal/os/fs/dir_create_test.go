package fs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/mws-cloud-platform/util-toolset/pkg/internal/os/fs"
	fsmock "github.com/mws-cloud-platform/util-toolset/pkg/internal/os/fs/mock"
)

func TestDirCreate(t *testing.T) {
	gm := gomock.NewController(t)
	mock := fsmock.NewMockFS(gm)
	dirCreate := fs.NewFS(mock, fs.WithDirCreate(os.ModePerm))

	mock.EXPECT().MkdirAll("foo/bar", gomock.Any())
	mock.EXPECT().WriteFile("foo/bar/test.txt", gomock.Any(), gomock.Any())
	require.NoError(t, dirCreate.WriteFile("foo/bar/test.txt", nil, os.ModePerm))
}

func TestDirCreateRootDirRemoveAll(t *testing.T) {
	gm := gomock.NewController(t)
	mock := fsmock.NewMockFS(gm)
	fs := fs.NewFS(mock, fs.WithDirCreate(os.ModePerm))

	touch := func(name string) {
		t.Helper()
		assert.NoError(t, fs.WriteFile(name, []byte(""), 0644))
	}

	gomock.InOrder(
		mock.EXPECT().MkdirAll("foo/bar", gomock.Any()),
		mock.EXPECT().WriteFile("foo/bar/1.txt", gomock.Any(), gomock.Any()),
		mock.EXPECT().WriteFile("foo/bar/2.txt", gomock.Any(), gomock.Any()),
		mock.EXPECT().RemoveAll("foo"),
		mock.EXPECT().MkdirAll("foo/bar", gomock.Any()),
		mock.EXPECT().WriteFile("foo/bar/3.txt", gomock.Any(), gomock.Any()),
	)

	touch("foo/bar/1.txt")
	touch("foo/bar/2.txt")
	assert.NoError(t, fs.RemoveAll("foo"))
	touch("foo/bar/3.txt")
}

func TestDirCreateRootDirRemove(t *testing.T) {
	gm := gomock.NewController(t)
	mock := fsmock.NewMockFS(gm)
	fs := fs.NewFS(mock, fs.WithDirCreate(os.ModePerm))

	touch := func(name string) {
		t.Helper()
		assert.NoError(t, fs.WriteFile(name, []byte(""), 0644))
	}

	gomock.InOrder(
		mock.EXPECT().MkdirAll("foo/bar", gomock.Any()),
		mock.EXPECT().WriteFile("foo/bar/1.txt", gomock.Any(), gomock.Any()),
		mock.EXPECT().Remove("foo/bar/1.txt"),
		mock.EXPECT().Remove("foo/bar"),
		mock.EXPECT().MkdirAll("foo/bar", gomock.Any()),
		mock.EXPECT().WriteFile("foo/bar/2.txt", gomock.Any(), gomock.Any()),
	)

	touch("foo/bar/1.txt")
	assert.NoError(t, fs.Remove("foo/bar/1.txt"))
	assert.NoError(t, fs.Remove("foo/bar"))
	touch("foo/bar/2.txt")
}
