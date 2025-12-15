package fs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"go.mws.cloud/util-toolset/pkg/internal/os/fs"
	fsmock "go.mws.cloud/util-toolset/pkg/internal/os/fs/mock"
)

type uniqueTestSuite struct {
	suite.Suite
	mock   *fsmock.MockFS
	unique fs.FS
}

func TestUnique(t *testing.T) {
	suite.Run(t, &uniqueTestSuite{})
}

func (s *uniqueTestSuite) SetupTest() {
	gm := gomock.NewController(s.T())
	s.mock = fsmock.NewMockFS(gm)
	s.unique = fs.NewFS(s.mock, fs.WithUnique())
}

func (s *uniqueTestSuite) TestWrite() {
	s.mock.EXPECT().WriteFile("x.txt", gomock.Any(), gomock.Any())
	s.NoError(s.unique.WriteFile("x.txt", nil, os.ModePerm))
}

func (s *uniqueTestSuite) TestWriteDuplicateName() {
	s.mock.EXPECT().WriteFile("x.txt", gomock.Any(), gomock.Any())
	s.Require().NoError(s.unique.WriteFile("x.txt", nil, os.ModePerm))
	s.nonUniqueError(s.unique.WriteFile("x.txt", nil, os.ModePerm))
}

func (s *uniqueTestSuite) TestWriteAfterRename() {
	gomock.InOrder(
		s.mock.EXPECT().WriteFile("y.txt", gomock.Any(), gomock.Any()),
		s.mock.EXPECT().Rename("y.txt", "z.txt"),
		s.mock.EXPECT().WriteFile("y.txt", gomock.Any(), gomock.Any()),
	)

	s.NoError(s.unique.WriteFile("y.txt", nil, os.ModePerm))
	s.NoError(s.unique.Rename("y.txt", "z.txt"))
	s.NoError(s.unique.WriteFile("y.txt", nil, os.ModePerm))
	s.nonUniqueError(s.unique.WriteFile("z.txt", nil, os.ModePerm))
}

func (s *uniqueTestSuite) TestWriteAfterRemove() {
	gomock.InOrder(
		s.mock.EXPECT().WriteFile("z.txt", gomock.Any(), gomock.Any()),
		s.mock.EXPECT().Remove("z.txt"),
		s.mock.EXPECT().WriteFile("z.txt", gomock.Any(), gomock.Any()),
	)
	s.NoError(s.unique.WriteFile("z.txt", nil, os.ModePerm))
	s.NoError(s.unique.Remove("z.txt"))
	s.NoError(s.unique.WriteFile("z.txt", nil, os.ModePerm))
}

func (s *uniqueTestSuite) TestRenameDifferentName() {
	gomock.InOrder(
		s.mock.EXPECT().WriteFile("y.txt", gomock.Any(), gomock.Any()),
		s.mock.EXPECT().Rename("y.txt", "z.txt"),
	)
	s.NoError(s.unique.WriteFile("y.txt", nil, os.ModePerm))
	s.NoError(s.unique.Rename("y.txt", "z.txt"))
}

func (s *uniqueTestSuite) TestRenameSameName() {
	s.mock.EXPECT().WriteFile("x.txt", gomock.Any(), gomock.Any())
	s.Require().NoError(s.unique.WriteFile("x.txt", nil, os.ModePerm))
	s.nonUniqueError(s.unique.Rename("x.txt", "x.txt"))
}

func (s *uniqueTestSuite) TestRenameDuplicateName() {
	gomock.InOrder(
		s.mock.EXPECT().WriteFile("y.txt", gomock.Any(), gomock.Any()),
		s.mock.EXPECT().WriteFile("x.txt", gomock.Any(), gomock.Any()),
	)
	s.NoError(s.unique.WriteFile("y.txt", nil, os.ModePerm))
	s.NoError(s.unique.WriteFile("x.txt", nil, os.ModePerm))
	s.nonUniqueError(s.unique.Rename("y.txt", "x.txt"))
}

func (s *uniqueTestSuite) nonUniqueError(err error) {
	var nue *fs.NonUniqueError
	s.ErrorAs(err, &nue)
}
