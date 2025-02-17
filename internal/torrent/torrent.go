package torrent

import (
	"crypto/sha1"
	"fmt"
	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
)

// TorrentMetadata represents the metadata extracted from a torrent file.
type TorrentMetadata struct {
	Name        string   // The name of the file or folder.
	Length      int      // The length of the file in bytes (for single files).
	PieceLength int      // The length of each piece in bytes.
	Pieces      []string // The hash of each piece, typically in 20-byte SHA-1 hash strings.
	Announce    string   // The URL of the tracker for the torrent.
	InfoHash    [20]byte // hash of the info
}

// Info parses the given list of bytes and returns a TorrentMetadata object.
func Info(data []byte) (*TorrentMetadata, error) {
	// Implementation of parsing logic should be added here.
	// This is a placeholder example, and it's assumed that the `data` contains
	// properly formatted torrent metadata information.

	decoded, _, err := bencode.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse torrent file: %v", err)
	}

	root, ok := decoded.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected decoded value to be a dictionary, but got %T", decoded)
	}

	info, ok := root["info"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected info field to be a dictionary, but got %T", root["info"])
	}

	encodedInfo, err := bencode.Encode(info)
	if err != nil {
		return nil, fmt.Errorf("failed to encode info: %v", err)
	}
	hash := sha1.Sum(encodedInfo)

	// For now, returning an example TorrentMetadata for demonstration.
	return &TorrentMetadata{
		Name:        string(info["name"].([]byte)),
		Length:      info["length"].(int),
		PieceLength: 16384,
		Pieces:      []string{"piece1_hash", "piece2_hash"},
		Announce:    string(root["announce"].([]byte)),
		InfoHash:    hash,
	}, nil
}
