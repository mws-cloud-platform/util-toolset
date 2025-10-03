package ptr

// Get returns a pointer for the passed value
func Get[T any](v T) *T {
	return &v
}

func Value[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}

	return *p
}

// Clone returns a new pointer object
func Clone[T any](p *T) *T {
	if p == nil {
		return nil
	}

	clone := *p
	return &clone
}

// Equal compares the values in the pointers
func Equal[T comparable](p1, p2 *T) bool {
	if p1 == nil || p2 == nil {
		return p1 == p2
	}
	return *p1 == *p2
}
