package fs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"go.mws.cloud/util-toolset/pkg/internal/os/fs"
	fsmock "go.mws.cloud/util-toolset/pkg/internal/os/fs/mock"
)

type readFile struct {
	fs.FS
}

func (*readFile) ReadFile(string) ([]byte, error) {
	return []byte("hello"), nil
}

func TestChangedOnly(t *testing.T) {
	gm := gomock.NewController(t)
	mock := fsmock.NewMockFS(gm)
	changedOnly := fs.NewFS(&readFile{FS: mock}, fs.WithChangedOnly())

	require.NoError(t, changedOnly.WriteFile("test.txt", []byte("hello"), os.ModePerm))

	mock.EXPECT().WriteFile("test.txt", []byte("world"), gomock.Any())
	require.NoError(t, changedOnly.WriteFile("test.txt", []byte("world"), os.ModePerm))
}
