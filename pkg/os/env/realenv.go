package env

import "os"

// Wrapper for default os env operations
type RealEnv struct {
}

var _ Env = RealEnv{}

func (RealEnv) Getenv(key string) string {
	return os.Getenv(key)
}

func (RealEnv) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (RealEnv) Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func (RealEnv) Environ() []string {
	return os.Environ()
}
