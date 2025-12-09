package golden

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

const dirOptionsTestBaseDir = "testdata/dir_options"

func (suite *dirOptionsTestSuite) TestGoldenDirWithPath() {
	dir := NewDir(suite.T(), WithPath(dirOptionsTestBaseDir))
	suite.Require().Equal(dirOptionsTestBaseDir, dir.path)
}

func (suite *dirOptionsTestSuite) TestGoldenDirWithRecreateOnUpdate() {
	dir := NewDir(suite.T(), WithPath(dirOptionsTestBaseDir), WithRecreateOnUpdate())
	suite.Require().True(dir.recreateOnUpdate)
}

func TestDirOptions(t *testing.T) {
	suite.Run(t, &dirOptionsTestSuite{})
}

type dirOptionsTestSuite struct {
	suite.Suite
}

func (suite *dirOptionsTestSuite) SetupTest() {
	createDir(suite.T(), dirOptionsTestBaseDir)
}

func (suite *dirOptionsTestSuite) TearDownTest() {
	removeDir(suite.T(), dirOptionsTestBaseDir)
}
