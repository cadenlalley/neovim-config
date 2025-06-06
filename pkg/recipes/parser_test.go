package recipes

import (
	"testing"

	"gopkg.in/guregu/null.v4"
)

func TestParseNullString(t *testing.T) {
	testCases := []struct {
		name     string
		input    null.String
		expected null.String
	}{
		{
			name:     "null string",
			input:    null.StringFrom("null"),
			expected: null.String{},
		},
		{
			name:     "dash string",
			input:    null.StringFrom("-"),
			expected: null.String{},
		},
		{
			name:     "empty string",
			input:    null.StringFrom(""),
			expected: null.String{},
		},
		{
			name:     "n/a string",
			input:    null.StringFrom("n/a"),
			expected: null.String{},
		},
		{
			name:     "valid string",
			input:    null.StringFrom("test"),
			expected: null.StringFrom("test"),
		},
	}

	for _, tc := range testCases {
		output := ParseNullString(tc.input)
		if output != tc.expected {
			t.Errorf("%s: expected %v, got %v", tc.name, tc.expected, output)
		}
	}
}

func TestParseNullFloat(t *testing.T) {
	testCases := []struct {
		name     string
		input    null.Float
		expected null.Float
	}{
		{
			name:     "zero float",
			input:    null.FloatFrom(0),
			expected: null.Float{},
		},
		{
			name:     "non zero float",
			input:    null.FloatFrom(1),
			expected: null.FloatFrom(1),
		},
	}

	for _, tc := range testCases {
		output := ParseNullFloat(tc.input)
		if output != tc.expected {
			t.Errorf("%s: expected %v, got %v", tc.name, tc.expected, output)
		}
	}
}
