package env

import "sort"

// MapEnv is an [Env] implementation backed by a map.
type MapEnv map[string]string

var _ Env = MapEnv(nil)

// NewMapEnv creates a new [MapEnv].
func NewMapEnv() MapEnv {
	return make(map[string]string)
}

// Getenv retrieves the value of the environment variable named by the key.
// If the variable is not present in the map, Getenv returns the empty string.
func (e MapEnv) Getenv(key string) string {
	return e[key]
}

// LookupEnv retrieves the value of the environment variable named by the key.
// If the variable is present in the map, the value (which may be empty)
// is returned and the boolean is true. Otherwise, the returned value will be
// empty and the boolean will be false.
func (e MapEnv) LookupEnv(key string) (string, bool) {
	value, ok := e[key]
	return value, ok
}

// Environ returns a copy of strings from map representing the environment,
// in the form "key=value".
func (e MapEnv) Environ() []string {
	result := make([]string, 0, len(e))
	for k, v := range e {
		result = append(result, k+"="+v)
	}
	sort.Strings(result)
	return result
}
