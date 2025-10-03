package env

type Env interface {
	Getenv(key string) string
	LookupEnv(key string) (string, bool)

	// Environ returns a copy of strings representing the environment,
	// in the form "key=value".
	Environ() []string
}
