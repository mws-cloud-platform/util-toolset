// Package env provides an interface for accessing environment variables.
package env

// Env is an interface for accessing environment variables.
type Env interface {
	// Getenv retrieves the value of the environment variable named by the key.
	// If the variable is not present, Getenv returns the empty string.
	Getenv(key string) string

	// LookupEnv retrieves the value of the environment variable named by the key.
	// If the variable is present in the environment the value (which may be empty)
	// is returned and the boolean is true. Otherwise, the returned value will be
	// empty and the boolean will be false.
	LookupEnv(key string) (string, bool)

	// Environ returns a copy of strings representing the environment,
	// in the form "key=value".
	Environ() []string
}
