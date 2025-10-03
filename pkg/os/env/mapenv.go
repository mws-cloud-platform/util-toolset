package env

import "sort"

type MapEnv map[string]string

var _ Env = MapEnv{}

func NewMapEnv() MapEnv {
	return make(map[string]string)
}

func (e MapEnv) Getenv(key string) string {
	return e[key]
}

func (e MapEnv) LookupEnv(key string) (string, bool) {
	value, ok := e[key]
	return value, ok
}

func (e MapEnv) Environ() []string {
	result := make([]string, 0, len(e))
	for k, v := range e {
		result = append(result, k+"="+v)
	}
	sort.Strings(result)
	return result
}
