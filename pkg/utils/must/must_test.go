package must_test

import (
	"errors"
	"testing"

	"go.mws.cloud/util-toolset/pkg/utils/consterr"
	"go.mws.cloud/util-toolset/pkg/utils/must"
)

func TestNoError(t *testing.T) {
	const expectedErr = consterr.Error("test error")

	defer func() {
		r := recover()
		if err, ok := r.(error); !ok || !errors.Is(err, expectedErr) {
			t.Fatal("must.NoError did not panic with the expected error")
		}
	}()

	must.NoError(expectedErr)
}

func TestValue(t *testing.T) {
	fn := func() (string, error) {
		return "Hello World", nil
	}

	if must.Value(fn()) != "Hello World" {
		t.Fatal("must.Value did not return the expected value")
	}
}

func TestValues(t *testing.T) {
	fn := func() (string, int, error) {
		return "Hello", 42, nil
	}

	v1, v2 := must.Values(fn())
	if v1 != "Hello" || v2 != 42 {
		t.Fatal("must.Values did not return the expected values")
	}
}
