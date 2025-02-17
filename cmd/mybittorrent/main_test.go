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
		{
			name: "info of torrent file",
			args: []string{"program", "info", "../../sample.torrent"},
			want: "Tracker URL: http://bittorrent-test-tracker.codecrafters.io/announce\nLength: 92063\nInfo Hash: d69f91e6b2ae4c542468d1073a71d4ea13879a7f",
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
