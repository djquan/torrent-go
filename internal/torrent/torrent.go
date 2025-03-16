package torrent

import (
	"crypto/sha1"
	"fmt"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
)

// Metadata represents the metadata extracted from a torrent file.
type Metadata struct {
	Name        string   // The name of the file or folder.
	Length      int      // The length of the file in bytes (for single files).
	PieceLength int      // The length of each piece in bytes.
	PieceHashes []string // The hash of each piece, typically in 20-byte SHA-1 hash strings.
	Announce    string   // The URL of the tracker for the torrent.
	InfoHash    [20]byte // hash of the info
}

// Info parses the given list of bytes and returns a Metadata object.
func Info(data []byte) (*Metadata, error) {
	// Implementation of parsing logic should be added here.
	// This is a placeholder example, and it's assumed that the `data` contains
	// properly formatted torrent metadata information.

	decoded, _, err := bencode.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse torrent file: %v", err)
	}

	root, ok := decoded.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected decoded value to be a dictionary, but got %T", decoded)
	}

	info, ok := root["info"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected info field to be a dictionary, but got %T", root["info"])
	}

	// Get piece length and pieces
	pieceLength, ok := info["piece length"].(int)
	if !ok {
		return nil, fmt.Errorf("expected piece length to be an integer, but got %T", info["piece length"])
	}

	pieces, ok := info["pieces"].([]byte)
	if !ok {
		return nil, fmt.Errorf("expected pieces to be a byte string, but got %T", info["pieces"])
	}

	// Split pieces into 20-byte hashes
	pieceHashes := make([]string, 0, len(pieces)/20)
	for i := 0; i < len(pieces); i += 20 {
		pieceHashes = append(pieceHashes, fmt.Sprintf("%x", pieces[i:i+20]))
	}

	encodedInfo, err := bencode.Encode(info)
	if err != nil {
		return nil, fmt.Errorf("failed to encode info: %v", err)
	}
	hash := sha1.Sum(encodedInfo)

	return &Metadata{
		Name:        string(info["name"].([]byte)),
		Length:      info["length"].(int),
		PieceLength: pieceLength,
		PieceHashes: pieceHashes,
		Announce:    string(root["announce"].([]byte)),
		InfoHash:    hash,
	}, nil
}
