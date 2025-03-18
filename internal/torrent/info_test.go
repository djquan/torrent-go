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

	if len(info.PieceHashes) != 3 {
		t.Errorf("expected 3 piece hashes, got %d", len(info.PieceHashes))
	}

	if info.PieceHashes[0] != "e876f67a2a8886e8f36b136726c30fa29703022d" {
		t.Errorf("expected e876f67a2a8886e8f36b136726c30fa29703022d, got %s", info.PieceHashes[0])
	}

	if info.PieceHashes[1] != "6e2275e604a0766656736e81ff10b55204ad8d35" {
		t.Errorf("expected 6e2275e604a0766656736e81ff10b55204ad8d35, got %s", info.PieceHashes[1])
	}

	if info.PieceHashes[2] != "f00d937a0213df1982bc8d097227ad9e909acc17" {
		t.Errorf("expected f00d937a0213df1982bc8d097227ad9e909acc17, got %s", info.PieceHashes[2])
	}
}
