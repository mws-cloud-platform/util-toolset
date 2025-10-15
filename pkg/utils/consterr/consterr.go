// Package consterr provides a simple way to create sentinel errors that
// cannot be changed during runtime.
package consterr

// Error is used for sentinel errors
//
//	const ErrMy = consterr.Error("something happened")
type Error string

func (e Error) Error() string {
	return string(e)
}
