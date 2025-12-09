package golden

import (
	"encoding/json"
	"fmt"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const goldenValuesTestBaseDir = "testdata/values"

func TestDirBytesSame(t *testing.T) {
	goldenDir := path.Join(goldenValuesTestBaseDir, "bytes")
	inputData := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit")

	dir := NewDir(t, WithPath(goldenDir))
	dir.update = false

	dir.Bytes(t, "text.txt", inputData)
}

func TestDirBytesDifferent(t *testing.T) {
	goldenDir := path.Join(goldenValuesTestBaseDir, "bytes")
	inputData := []byte("no lorem ipsum :c")
	mT := &mockT{}

	dir := NewDir(t, WithPath(goldenDir))
	dir.update = false

	dir.Bytes(mT, "text.txt", inputData)

	require.True(t, mT.failed, mT.LogString())
}

func TestDirStringSame(t *testing.T) {
	goldenDir := path.Join(goldenValuesTestBaseDir, "string")
	inputData := "It's the story of a terrifying green ogre by the name of Shrek, who lives in a swamp."

	dir := NewDir(t, WithPath(goldenDir))
	dir.update = false

	dir.String(t, "text.txt", inputData)
}

func TestDirStringDifferent(t *testing.T) {
	goldenDir := path.Join(goldenValuesTestBaseDir, "string")
	inputData := "bad movie"
	mT := &mockT{}

	dir := NewDir(t, WithPath(goldenDir))
	dir.update = false

	dir.String(mT, "text.txt", inputData)

	require.True(t, mT.failed, mT.LogString())
}

func TestDirJSONBytesSame(t *testing.T) {
	goldenDir := path.Join(goldenValuesTestBaseDir, "jsonbytes")
	obj := dummyPerson{
		Name: "Alex",
		Age:  34,
		Hobbies: []dummyHobby{
			{
				Name:  "swimming",
				Skill: 3,
			},
			{
				Name:  "hiking",
				Skill: 5,
			},
		},
	}
	inputData, err := json.MarshalIndent(obj, "", " ")
	require.NoError(t, err)

	dir := NewDir(t, WithPath(goldenDir))
	dir.update = false

	dir.JSONBytes(t, "data.json", inputData)
}

func TestDirJSONBytesDifferent(t *testing.T) {
	goldenDir := path.Join(goldenValuesTestBaseDir, "jsonbytes")
	obj := dummyPerson{
		Name: "Dima",
		Age:  35,
		Hobbies: []dummyHobby{
			{
				Name:  "drawing",
				Skill: 10,
			},
		},
	}
	inputData, err := json.MarshalIndent(obj, "", " ")
	require.NoError(t, err)
	mT := &mockT{}

	dir := NewDir(t, WithPath(goldenDir))
	dir.update = false

	dir.JSONBytes(mT, "data.json", inputData)

	require.True(t, mT.failed, mT.LogString())
}

func TestDirJSONSame(t *testing.T) {
	goldenDir := path.Join(goldenValuesTestBaseDir, "json")
	inputData := dummyPerson{
		Name: "Vlad",
		Age:  31,
		Hobbies: []dummyHobby{
			{
				Name: "reading",
			},
		},
	}

	dir := NewDir(t, WithPath(goldenDir))
	dir.update = false

	dir.JSON(t, "data.json", inputData)
}

func TestDirJSONDifferent(t *testing.T) {
	goldenDir := path.Join(goldenValuesTestBaseDir, "json")
	inputData := dummyPerson{
		Name: "Steve",
		Age:  48,
		Hobbies: []dummyHobby{
			{
				Name:  "writing",
				Skill: 8,
			},
		},
	}
	mT := &mockT{}

	dir := NewDir(t, WithPath(goldenDir))
	dir.update = false

	dir.JSON(mT, "data.json", inputData)

	require.True(t, mT.failed, mT.LogString())
}

// Mock implementation of [testing.T].
type mockT struct {
	skipped bool
	failed  bool
	log     string
}

func (t *mockT) Error(args ...any) {
	t.log += fmt.Sprintln(args...)
	t.Fail()
}

func (t *mockT) Errorf(format string, args ...any) {
	t.log += fmt.Sprintf(format, args...)
	t.Fail()
}

func (t *mockT) Fail() {
	t.failed = true
}

func (t *mockT) FailNow() {
	t.log += "FailNow() call"
	t.failed = true
}

func (t *mockT) Failed() bool {
	return t.failed
}

func (t *mockT) Fatal(args ...any) {
	t.log += fmt.Sprint(args...)
	t.FailNow()
}

func (t *mockT) Fatalf(format string, args ...any) {
	t.log += fmt.Sprintf(format, args...)
	t.FailNow()
}

func (t *mockT) Log(args ...any) {
	t.log += fmt.Sprintln(args...)
}

func (t *mockT) Logf(format string, args ...any) {
	t.log += fmt.Sprintf(format, args...)
}

func (*mockT) Name() string { return "" }
func (*mockT) Parallel()    {}

func (t *mockT) Skip(args ...any) {
	t.log += fmt.Sprint(args...)
	t.SkipNow()
}

func (t *mockT) SkipNow() {
	t.skipped = true
}

func (t *mockT) Skipf(format string, args ...any) {
	t.log += fmt.Sprintf(format, args...)
	t.SkipNow()
}

func (t *mockT) Skipped() bool {
	return t.skipped
}

func (*mockT) Helper() {}

func (*mockT) Cleanup(func()) {}

// LogString returns accumulated logs.
func (t *mockT) LogString() string {
	return t.log
}

type dummyHobby struct {
	Name  string `json:"name"`
	Skill int    `json:"skill"`
}

type dummyPerson struct {
	Name    string       `json:"name"`
	Age     int          `json:"age"`
	Hobbies []dummyHobby `json:"hobbies"`
}
