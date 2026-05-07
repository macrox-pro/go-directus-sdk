package helpers_test

import (
	"testing"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

func TestJoinPartsURL(t *testing.T) {
	tests := []struct {
		name     string
		parts    []helpers.PartURL
		expected string
	}{
		{
			name:     "empty parts",
			parts:    []helpers.PartURL{},
			expected: "",
		},
		{
			name: "single part",
			parts: []helpers.PartURL{
				{Value: "api"},
			},
			expected: "api",
		},
		{
			name: "multiple parts",
			parts: []helpers.PartURL{
				{Value: "api"},
				{Value: "v1"},
				{Value: "users"},
			},
			expected: "api/v1/users",
		},
		{
			name: "parts with leading/trailing slashes",
			parts: []helpers.PartURL{
				{Value: "/api/"},
				{Value: "/v1/"},
				{Value: "/users/"},
			},
			expected: "api/v1/users",
		},
		{
			name: "skip some parts",
			parts: []helpers.PartURL{
				{Value: "api"},
				{Value: "v1", Skip: true},
				{Value: "users"},
			},
			expected: "api/users",
		},
		{
			name: "skip all parts",
			parts: []helpers.PartURL{
				{Value: "api", Skip: true},
				{Value: "v1", Skip: true},
			},
			expected: "",
		},
		{
			name: "empty string parts",
			parts: []helpers.PartURL{
				{Value: ""},
				{Value: "api"},
			},
			expected: "/api",
		},
		{
			name: "only slashes",
			parts: []helpers.PartURL{
				{Value: "///"},
				{Value: "/"},
			},
			expected: "/",
		},
		{
			name: "mixed skip and slashes",
			parts: []helpers.PartURL{
				{Value: "/api/", Skip: false},
				{Value: "//v1//", Skip: true},
				{Value: "users"},
			},
			expected: "api/users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.JoinPartsURL(tt.parts...)
			if result != tt.expected {
				t.Errorf("JoinPartsURL() = %q, want %q", result, tt.expected)
			}
		})
	}
}
