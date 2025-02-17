package main

import (
	"testing"
)

func TestDecodeBencode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
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
			name:    "empty string",
			input:   "",
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
