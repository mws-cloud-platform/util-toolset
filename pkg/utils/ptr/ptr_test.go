package ptr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	require.Equal(t, *Get("test"), "test")
}

func TestValue(t *testing.T) {
	var strPointer *string
	require.Equal(t, "", Value(strPointer))

	strPointer = Get("test")
	require.Equal(t, "test", Value(strPointer))
}

func TestClone(t *testing.T) {
	sPtr := new(string)
	newPtr := Clone(sPtr)
	require.NotSame(t, newPtr, sPtr)
}

func TestEqual(t *testing.T) {
	var nilPtr *string
	require.True(t, Equal(nilPtr, nilPtr))
	require.True(t, Equal(new(string), new(string)))
	require.False(t, Equal(nil, Get("test")))
	require.False(t, Equal(Get("foo"), Get("test")))
}
