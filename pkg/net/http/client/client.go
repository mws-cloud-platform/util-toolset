// Package client provides HTTP client utilities.
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"unicode"

	"github.com/mws-cloud-platform/util-toolset/pkg/utils/consterr"
)

// ErrInvalidBodyTrail is returned when body trail is not blank.
const ErrInvalidBodyTrail = consterr.Error("invalid body trail")

// ReadJSON reads JSON from the body, decodes it into v and discards body trail
// if any. If body trail should be checked, use [ReadJSONSafe].
func ReadJSON(body io.Reader, v any) error {
	err := json.NewDecoder(body).Decode(v)
	if _, discardErr := io.Copy(io.Discard, body); discardErr != nil {
		err = errors.Join(err, fmt.Errorf("discard body trail: %w", discardErr))
	}
	return err
}

// ReadJSONSafe reads JSON from the body, decodes it into v and checks body
// trail if any. It returns [ErrInvalidBodyTrail] if body trail is not empty or
// contains non-whitespace characters.
func ReadJSONSafe(body io.Reader, v any) error {
	err := json.NewDecoder(body).Decode(v)
	if err != nil {
		if _, discardErr := io.Copy(io.Discard, body); discardErr != nil {
			err = errors.Join(err, fmt.Errorf("discard body trail: %w", discardErr))
		}
	} else {
		trail, readErr := io.ReadAll(body)
		switch {
		case readErr != nil:
			err = fmt.Errorf("read body trail: %w", readErr)
		case bytes.ContainsFunc(trail, func(r rune) bool {
			return !unicode.IsSpace(r)
		}):
			err = ErrInvalidBodyTrail
		}
	}
	return err
}
