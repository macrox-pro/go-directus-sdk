package helpers_test

import (
	"testing"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

func TestExtractFieldsJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected []string
	}{
		{
			name: "simple struct with json tags",
			input: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Age  int    `json:"age,omitempty"`
			}{},
			expected: []string{"id", "name", "age"},
		},
		{
			name: "nested struct",
			input: struct {
				ID      string `json:"id"`
				Address struct {
					Street string `json:"street"`
					City   string `json:"city"`
				} `json:"address"`
			}{},
			expected: []string{"id", "address.street", "address.city"},
		},
		{
			name: "named nested field",
			input: struct {
				Base struct {
					ID string `json:"id"`
				} `json:"base"`
				Name string `json:"name"`
			}{},
			expected: []string{"base.id", "name"},
		},
		{
			name: "ignore dash tag",
			input: struct {
				Public  string `json:"public"`
				Private string `json:"-"`
			}{},
			expected: []string{"public"},
		},
		{
			name: "empty tag uses field name",
			input: struct {
				Field1 string `json:""`
				Field2 string `json:",omitempty"`
			}{},
			expected: []string{"Field1", "Field2"},
		},
		{
			name:     "empty struct",
			input:    struct{}{},
			expected: nil,
		},
		{
			name: "slice element",
			input: []struct {
				Value string `json:"value"`
			}{},
			expected: []string{"value"},
		},
		{
			name: "multiple levels",
			input: struct {
				A struct {
					B struct {
						C string `json:"c"`
					} `json:"b"`
				} `json:"a"`
			}{},
			expected: []string{"a.b.c"},
		},
		{
			name:     "non-struct type (string)",
			input:    "hello",
			expected: nil,
		},
		{
			name:     "non-struct type (int)",
			input:    42,
			expected: nil,
		},
		{
			name: "field without json tag",
			input: struct {
				Tagged   string `json:"tagged"`
				Untagged string
			}{},
			expected: []string{"tagged", "Untagged"},
		},
		{
			name: "array of structs",
			input: [2]struct {
				Value string `json:"value"`
			}{},
			expected: []string{"value"},
		},
		{
			name: "anonymous field with non-struct type",
			input: struct {
				int
				Name string `json:"name"`
			}{},
			expected: []string{"name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := helpers.ExtractFieldsJSON(tt.input)

			if tt.expected == nil && fields != nil {
				t.Errorf("expected nil, got %v", fields)
				return
			}

			if len(fields) != len(tt.expected) {
				t.Errorf("field count mismatch: expected %d, got %d", len(tt.expected), len(fields))
				return
			}

			for i, f := range fields {
				if f != tt.expected[i] {
					t.Errorf("field %d mismatch: expected %q, got %q", i, tt.expected[i], f)
				}
			}
		})
	}
}
