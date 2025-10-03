package golden

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	mockfs "github.com/mws-cloud-platform/util-toolset/pkg/internal/os/fs/mock"
)

func TestDirBytesUpdate(t *testing.T) {
	suite.Run(
		t,
		&dirUpdateTestSuite{
			GoldenDirpath: "cool_dir",
			Filename:      "cool_file",
		},
	)
}

func (suite *dirUpdateTestSuite) TestDirBytesUpdate() {
	data := []byte("no loreum ipsum :c")

	suite.MFS.
		EXPECT().
		WriteFile(suite.goldenFilepath(), data, os.FileMode(0644)).
		Return(nil)

	dir := NewDir(suite.T(), WithPath(suite.GoldenDirpath), WithFS(suite.MFS))
	dir.update = true

	dir.Bytes(suite.T(), suite.Filename, data)
}

func (suite *dirUpdateTestSuite) TestDirStringUpdate() {
	data := "bad movie"

	suite.MFS.
		EXPECT().
		WriteFile(suite.goldenFilepath(), []byte(data), os.FileMode(0644)).
		Return(nil)

	dir := NewDir(suite.T(), WithPath(suite.GoldenDirpath), WithFS(suite.MFS))
	dir.update = true

	dir.String(suite.T(), suite.Filename, data)
}

func (suite *dirUpdateTestSuite) TestDirJSONBytesUpdate() {
	obj := dummyPerson{
		Name: "Joe",
		Age:  27,
		Hobbies: []dummyHobby{
			{
				Name:  "dancing",
				Skill: 2,
			},
		},
	}
	data, err := json.MarshalIndent(obj, "", " ")
	suite.Require().NoError(err)

	suite.MFS.
		EXPECT().
		WriteFile(suite.goldenFilepath(), data, os.FileMode(0644)).
		Return(nil)

	dir := NewDir(suite.T(), WithPath(suite.GoldenDirpath), WithFS(suite.MFS))
	dir.update = true

	dir.JSONBytes(suite.T(), suite.Filename, data)
}

func (suite *dirUpdateTestSuite) TestDirJSONUpdate() {
	data := dummyPerson{
		Name: "Marie",
		Age:  41,
		Hobbies: []dummyHobby{
			{
				Name:  "gardening",
				Skill: 3,
			},
			{
				Name:  "skiing",
				Skill: 6,
			},
		},
	}
	rawData, err := json.MarshalIndent(data, "", " ")
	require.NoError(suite.T(), err)

	suite.MFS.
		EXPECT().
		WriteFile(suite.goldenFilepath(), rawData, os.FileMode(0644)).
		Return(nil)

	dir := NewDir(suite.T(), WithPath(suite.GoldenDirpath), WithFS(suite.MFS))
	dir.update = true

	dir.JSON(suite.T(), suite.Filename, data)
}

func (suite *dirUpdateTestSuite) SetupTest() {
	suite.MockCtrl = gomock.NewController(suite.T())
	suite.MFS = mockfs.NewMockFS(suite.MockCtrl)
}

type dirUpdateTestSuite struct {
	suite.Suite
	GoldenDirpath string             // target golden dir (mocked)
	Filename      string             // name of the file to be checked
	MockCtrl      *gomock.Controller // mocked controller
	MFS           *mockfs.MockFS     // mocked FS
}

func (suite *dirUpdateTestSuite) goldenFilepath() string {
	return path.Join(suite.GoldenDirpath, suite.Filename)
}
