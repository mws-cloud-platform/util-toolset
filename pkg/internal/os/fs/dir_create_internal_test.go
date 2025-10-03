package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isSubPathOf(t *testing.T) {
	assert.True(t, isSubPath("foo", "foo/bar"))
	assert.True(t, isSubPath("foo/", "foo/bar"))

	assert.False(t, isSubPath("fo", "foo/bar"))
	assert.False(t, isSubPath("foobar", "foo/bar"))
	assert.False(t, isSubPath("foo/zoo", "foo/bar"))
}
