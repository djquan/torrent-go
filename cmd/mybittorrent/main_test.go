package main

import (
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeBencode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeBencode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("decodeBencode() = %v, want %v", got, tt.expected)
			}
		})
	}
}
