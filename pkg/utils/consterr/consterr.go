package consterr

// Error is used for sentinel errors
// Example: const ErrMy = consterr.Error("something happened")
type Error string

func (e Error) Error() string {
	return string(e)
}
