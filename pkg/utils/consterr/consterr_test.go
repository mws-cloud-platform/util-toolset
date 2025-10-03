package consterr_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mws-cloud-platform/util-toolset/pkg/utils/consterr"
)

const (
	ErrOne = consterr.Error("one")
	ErrTwo = consterr.Error("two")
)

var _ error = ErrOne

func TestError(t *testing.T) {
	err := fmt.Errorf("%w: %w", ErrOne, fmt.Errorf("test"))
	require.ErrorIs(t, err, ErrOne)
	require.NotErrorIs(t, err, ErrTwo)
	require.Equal(t, "one: test", err.Error())
}
