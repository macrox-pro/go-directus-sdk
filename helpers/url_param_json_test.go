package helpers_test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

func TestURLParamJSON_EncodeValues(t *testing.T) {
	tests := []struct {
		name        string
		data        any
		key         string
		expectError bool
		expectValue string
	}{
		{
			name:        "nil data",
			data:        nil,
			key:         "filter",
			expectError: false,
			expectValue: "",
		},
		{
			name:        "simple map",
			data:        map[string]string{"name": "john"},
			key:         "filter",
			expectError: false,
			expectValue: `{"name":"john"}`,
		},
		{
			name:        "struct",
			data:        struct{ ID int }{ID: 42},
			key:         "fields",
			expectError: false,
			expectValue: `{"ID":42}`,
		},
		{
			name:        "slice",
			data:        []string{"a", "b", "c"},
			key:         "ids",
			expectError: false,
			expectValue: `["a","b","c"]`,
		},
		{
			name:        "empty map",
			data:        map[string]interface{}{},
			key:         "filter",
			expectError: false,
			expectValue: `{}`,
		},
		{
			name:        "string",
			data:        "plain string",
			key:         "q",
			expectError: false,
			expectValue: `"plain string"`,
		},
		{
			name:        "number",
			data:        123.45,
			key:         "value",
			expectError: false,
			expectValue: `123.45`,
		},
		{
			name:        "boolean",
			data:        true,
			key:         "active",
			expectError: false,
			expectValue: `true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param := helpers.URLParamJSON{Data: tt.data}
			values := url.Values{}

			err := param.EncodeValues(tt.key, &values)
			if tt.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expectValue == "" {
				if len(values) > 0 {
					t.Errorf("expected no values, got %v", values)
				}
				return
			}

			got := values.Get(tt.key)
			if got != tt.expectValue {
				t.Errorf("EncodeValues() = %q, want %q", got, tt.expectValue)
			}
		})
	}
}

func TestURLParamJSON_MarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		data        any
		expectJSON  string
		expectError bool
	}{
		{
			name:       "nil data",
			data:       nil,
			expectJSON: `null`,
		},
		{
			name:       "simple struct",
			data:       struct{ Name string }{Name: "test"},
			expectJSON: `{"Name":"test"}`,
		},
		{
			name:       "map",
			data:       map[string]int{"x": 1},
			expectJSON: `{"x":1}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param := helpers.URLParamJSON{Data: tt.data}
			bytes, err := param.MarshalJSON()

			if tt.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError {
				var got, want interface{}
				if err := json.Unmarshal(bytes, &got); err != nil {
					t.Errorf("failed to unmarshal result: %v", err)
				}
				if err := json.Unmarshal([]byte(tt.expectJSON), &want); err != nil {
					t.Errorf("failed to unmarshal expected: %v", err)
				}

				// Compare as JSON strings for simplicity
				gotStr := string(bytes)
				if gotStr != tt.expectJSON {
					t.Errorf("MarshalJSON() = %s, want %s", gotStr, tt.expectJSON)
				}
			}
		})
	}
}

func TestURLParamJSON_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		expected bool
	}{
		{
			name:     "nil data",
			data:     nil,
			expected: true,
		},
		{
			name:     "non-nil data",
			data:     "something",
			expected: false,
		},
		{
			name:     "empty map",
			data:     map[string]interface{}{},
			expected: false, // Data is not nil, so IsEmpty returns false
		},
		{
			name:     "zero struct",
			data:     struct{}{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param := helpers.URLParamJSON{Data: tt.data}
			empty := param.IsEmpty()
			if empty != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", empty, tt.expected)
			}
		})
	}
}
