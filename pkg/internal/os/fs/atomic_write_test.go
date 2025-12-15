package fs_test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"go.mws.cloud/util-toolset/pkg/internal/os/fs"
	fsmock "go.mws.cloud/util-toolset/pkg/internal/os/fs/mock"
)

type atomicWriteTestSuite struct {
	suite.Suite
	mock   *fsmock.MockFS
	atomic fs.FS
}

func TestAtomicWrite(t *testing.T) {
	suite.Run(t, &atomicWriteTestSuite{})
}

func (s *atomicWriteTestSuite) SetupTest() {
	gm := gomock.NewController(s.T())
	s.mock = fsmock.NewMockFS(gm)
	s.atomic = fs.NewFS(s.mock, fs.WithAtomicWrite())
}

func (s *atomicWriteTestSuite) TestWrite() {
	gomock.InOrder(
		s.mock.EXPECT().WriteFile("test.txt.tmp", gomock.Any(), gomock.Any()),
		s.mock.EXPECT().Rename("test.txt.tmp", "test.txt"),
	)
	s.NoError(s.atomic.WriteFile("test.txt", nil, os.ModePerm))
}

func (s *atomicWriteTestSuite) TestFailedWrite() {
	gomock.InOrder(
		s.mock.EXPECT().WriteFile("test.txt.tmp", gomock.Any(), gomock.Any()).Return(os.ErrInvalid),
		s.mock.EXPECT().Remove("test.txt.tmp"),
	)
	s.ErrorIs(s.atomic.WriteFile("test.txt", nil, os.ModePerm), os.ErrInvalid)
}

func (s *atomicWriteTestSuite) TestFailedRename() {
	gomock.InOrder(
		s.mock.EXPECT().WriteFile("test.txt.tmp", gomock.Any(), gomock.Any()),
		s.mock.EXPECT().Rename("test.txt.tmp", "test.txt").Return(os.ErrInvalid),
		s.mock.EXPECT().Remove("test.txt.tmp"),
	)
	s.ErrorIs(s.atomic.WriteFile("test.txt", nil, os.ModePerm), os.ErrInvalid)
}

func (s *atomicWriteTestSuite) TestCustomDir() {
	s.atomic = fs.NewFS(s.mock, fs.WithAtomicWriteCustomDir("/base"))
	gomock.InOrder(
		s.mock.EXPECT().WriteFile("/base/test.txt.tmp", gomock.Any(), gomock.Any()),
		s.mock.EXPECT().Rename("/base/test.txt.tmp", "test.txt"),
	)
	s.NoError(s.atomic.WriteFile("test.txt", nil, os.ModePerm))
}

func (s *atomicWriteTestSuite) TestMapFS() {
	data := []byte("hello")
	f := fs.NewFS(fs.NewMapFS(), fs.WithAtomicWrite())
	s.Require().NoError(f.WriteFile("test.txt", []byte("hello"), os.ModePerm))

	file, err := f.Open("test.txt")
	s.Require().NoError(err)

	bytes, err := io.ReadAll(file)
	s.Require().NoError(err)
	s.Equal(data, bytes)
	s.Require().NoError(file.Close())

	_, err = f.Open("test.txt.tmp")
	s.ErrorIs(err, os.ErrNotExist)
}
