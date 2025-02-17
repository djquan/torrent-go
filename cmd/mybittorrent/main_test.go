package main

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{
			name:    "decode simple string",
			args:    []string{"program", "decode", "5:hello"},
			want:    `"hello"`,
			wantErr: false,
		},
		{
			name:    "decode integer",
			args:    []string{"program", "decode", "i42e"},
			want:    "42",
			wantErr: false,
		},
		{
			name:    "decode list",
			args:    []string{"program", "decode", "l5:helloi42ee"},
			want:    `["hello",42]`,
			wantErr: false,
		},
		{
			name:    "no command provided",
			args:    []string{"program"},
			wantErr: true,
		},
		{
			name:    "invalid command",
			args:    []string{"program", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := run(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && strings.TrimSpace(got) != tt.want {
				t.Errorf("run() = %v, want %v", got, tt.want)
			}
		})
	}
}
