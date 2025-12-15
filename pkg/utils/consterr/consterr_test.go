package consterr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/util-toolset/pkg/utils/consterr"
)

const (
	ErrOne = consterr.Error("one")
	ErrTwo = consterr.Error("two")
)

var _ error = ErrOne

var errTest = errors.New("test")

func TestError(t *testing.T) {
	err := fmt.Errorf("%w: %w", ErrOne, errTest)
	require.ErrorIs(t, err, ErrOne)
	require.NotErrorIs(t, err, ErrTwo)
	require.Equal(t, "one: test", err.Error())
}
