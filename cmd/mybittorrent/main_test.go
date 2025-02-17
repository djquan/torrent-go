package main

import (
	"reflect"
	"testing"
)

func TestDecodeBencode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "basic string",
			input:    "5:hello",
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "longer string",
			input:    "10:helloworld",
			expected: "helloworld",
			wantErr:  false,
		},
		{
			name:     "single character",
			input:    "1:a",
			expected: "a",
			wantErr:  false,
		},
		{
			name:     "single number",
			input:    "i1e",
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "two digit number",
			input:    "i12e",
			expected: 12,
			wantErr:  false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "missing e",
			input:   "i12",
			wantErr: true,
		},
		{
			name:  "simple list",
			input: "l5:helloi1ee",
			expected: []interface{}{
				"hello",
				1,
			},
		},
		{
			name:     "empty list",
			input:    "le",
			expected: []interface{}{},
			wantErr:  false,
		},
		{
			name:  "nested list with integer and string",
			input: "lli956e5:appleee",
			expected: []interface{}{
				[]interface{}{
					956,
					"apple",
				},
			},
			wantErr: false,
		},
		{
			name:  "nested list",
			input: "ll5:helloi1eei2ee",
			expected: []interface{}{
				[]interface{}{
					"hello",
					1,
				},
				2,
			},
			wantErr: false,
		},
		{
			name:  "list with multiple types",
			input: "l5:helloi1ei2ee",
			expected: []interface{}{
				"hello",
				1,
				2,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := decodeBencode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeBencode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("decodeBencode() = %v, want %v", got, tt.expected)
			}
		})
	}
}
