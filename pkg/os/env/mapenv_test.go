package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapEnvGetenv(t *testing.T) {
	mapEnv := NewMapEnv()

	key, value := "keyring", "valuable"
	mapEnv[key] = value

	valueInMap := mapEnv.Getenv(key)
	require.Equal(t, value, valueInMap)
}

func TestMapEnvLookupEnv(t *testing.T) {
	mapEnv := NewMapEnv()

	key, value := "COOLKEY789", "value123"
	mapEnv[key] = value

	valueInMap, existsInMap := mapEnv.LookupEnv(key)
	require.True(t, existsInMap)
	require.Equal(t, value, valueInMap)
}

func TestMapEnvLookupEnvNonSet(t *testing.T) {
	mapEnv := NewMapEnv()

	key := "IDONTEXIST"

	valueInMap, existsInMap := mapEnv.LookupEnv(key)
	require.False(t, existsInMap)
	require.Equal(t, "", valueInMap)
}

func TestMapEnvEnviron(t *testing.T) {
	mapEnv := MapEnv{
		"ONE":         "foo",
		"HELLO_WORLD": "bar",
	}

	require.Equal(t, []string{
		"HELLO_WORLD=bar",
		"ONE=foo",
	}, mapEnv.Environ())
}
