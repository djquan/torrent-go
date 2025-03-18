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
			want: strings.Join([]string{
				"Tracker URL: http://bittorrent-test-tracker.codecrafters.io/announce",
				"Length: 92063",
				"Info Hash: d69f91e6b2ae4c542468d1073a71d4ea13879a7f",
				"Piece Length: 32768",
				"Piece Hashes:",
				"e876f67a2a8886e8f36b136726c30fa29703022d",
				"6e2275e604a0766656736e81ff10b55204ad8d35",
				"f00d937a0213df1982bc8d097227ad9e909acc17",
			}, "\n"),
		},
		{
			name: "peers of torrent file",
			args: []string{"program", "peers", "../../sample.torrent"},
			want: strings.Join([]string{
				"165.232.41.73:51556",
				"165.232.38.164:51493",
				"165.232.35.114:51476",
			}, "\n"),
		},
		{
			name: "handshake",
			args: []string{"program", "handshake", "../../sample.torrent", "165.232.41.73:51556"},
			want: "Peer ID: 2d524e302e302e302df4c05830bc2a1deb2944e6",
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
