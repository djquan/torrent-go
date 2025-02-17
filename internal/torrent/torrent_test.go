package torrent

import (
	"encoding/hex"
	"os"
	"testing"
)

func TestInfo(t *testing.T) {
	content, err := os.ReadFile("../../sample.torrent")
	if err != nil {
		t.Fatalf("failed to read sample.torrent: %v", err)
	}

	info, err := Info(content)
	if err != nil {
		t.Fatalf("failed to parse torrent file: %v", err)
	}
	t.Logf("info: %+v", info)
	if info.Name != "sample.txt" {
		t.Errorf("expected Name to be 'sample.txt', but got '%s'", info.Name)
	}

	if info.Announce != "http://bittorrent-test-tracker.codecrafters.io/announce" {
		t.Errorf("expected Announce to be 'http://bittorrent-test-tracker.codecrafters.io/announce', but got '%s'", info.Announce)
	}

	if info.Length != 92063 {
		t.Errorf("expected Length to be 92063, but got %d", info.Length)
	}

	if hex.EncodeToString(info.InfoHash[:]) != "d69f91e6b2ae4c542468d1073a71d4ea13879a7f" {
		t.Errorf("expected d69f91e6b2ae4c542468d1073a71d4ea13879a7f, got %v", hex.EncodeToString(info.InfoHash[:]))
	}
}
