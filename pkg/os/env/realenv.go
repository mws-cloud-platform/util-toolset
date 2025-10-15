package env

import "os"

// RealEnv is a wrapper for default os env operations.
type RealEnv struct{}

var _ Env = RealEnv{}

// Getenv retrieves the value of the environment variable named by the key.
// If the variable is not present, Getenv returns the empty string.
func (RealEnv) Getenv(key string) string {
	return os.Getenv(key)
}

// LookupEnv retrieves the value of the environment variable named by the key.
// If the variable is present in the environment the value (which may be empty)
// is returned and the boolean is true. Otherwise, the returned value will be
// empty and the boolean will be false.
func (RealEnv) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

// Environ returns a copy of strings representing the environment,
// in the form "key=value".
func (RealEnv) Environ() []string {
	return os.Environ()
}

// Setenv sets the value of the environment variable named by the key.
func (RealEnv) Setenv(key, value string) error {
	return os.Setenv(key, value)
}
