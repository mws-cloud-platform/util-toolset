package golden

import (
	"path"
	"testing"
)

func TestBytes(t *testing.T) {
	dir := "testdata/file"

	a := path.Join(dir, "a.txt")
	b := path.Join(dir, "b.txt")

	aContents := readFile(t, a)
	bContents := readFile(t, b)

	Bytes(t, a, aContents)
	Bytes(t, b, bContents)
}

func TestString(t *testing.T) {
	dir := "testdata/file"

	a := path.Join(dir, "a.txt")
	b := path.Join(dir, "b.txt")

	aContents := string(readFile(t, a))
	String(t, a, aContents)

	bContents := string(readFile(t, b))
	String(t, b, bContents)
}
