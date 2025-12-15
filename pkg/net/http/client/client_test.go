package client_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/util-toolset/pkg/net/http/client"
)

const jsonDecoderMinRead = 512

func TestReadJSON(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		v         any
		expected  any
		expectErr bool
	}{
		{
			name:     "ok",
			body:     `{"foo": "bar"}`,
			v:        map[string]any{},
			expected: map[string]any{"foo": "bar"},
		},
		{
			name:     "ok with valid trail",
			body:     `{"foo": "bar"}` + strings.Repeat(" ", jsonDecoderMinRead),
			v:        map[string]any{},
			expected: map[string]any{"foo": "bar"},
		},
		{
			name:     "ok with invalid trail",
			body:     `{"foo": "bar"}` + strings.Repeat("x", jsonDecoderMinRead),
			v:        map[string]any{},
			expected: map[string]any{"foo": "bar"},
		},
		{
			name:      "invalid json",
			body:      `}{`,
			v:         map[string]any{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			err := client.ReadJSON(body, &tt.v)
			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, tt.v)
		})
	}
}

func TestReadJSONSafe(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		v         any
		expected  any
		expectErr bool
	}{
		{
			name:     "ok",
			body:     `{"foo": "bar"}`,
			v:        map[string]any{},
			expected: map[string]any{"foo": "bar"},
		},
		{
			name:     "ok with valid trail",
			body:     `{"foo": "bar"}` + strings.Repeat(" ", jsonDecoderMinRead),
			v:        map[string]any{},
			expected: map[string]any{"foo": "bar"},
		},
		{
			name:      "ok with invalid trail",
			body:      `{"foo": "bar"}` + strings.Repeat("x", jsonDecoderMinRead),
			v:         map[string]any{},
			expectErr: true,
		},
		{
			name:      "invalid json",
			body:      `}{`,
			v:         map[string]any{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			err := client.ReadJSONSafe(body, &tt.v)
			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, tt.v)
		})
	}
}
