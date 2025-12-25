// Package must provide helper functions that simplify error handling by panicking on errors.
package must

// NoError panics if the provided error is not nil.
func NoError(err error) {
	if err != nil {
		panic(err)
	}
}

// Value returns the value if the error is nil; otherwise, it panics with the error.
func Value[T any](v T, err error) T {
	NoError(err)

	return v
}

// Values returns the two values if the error is nil; otherwise, it panics with the error.
func Values[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	NoError(err)

	return v1, v2
}
