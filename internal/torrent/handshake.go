package torrent

import (
	"encoding/hex"
	"fmt"
)

// TCPConn represents a TCP connection interface for dependency injection.
type TCPConn interface {
	// Read reads data from the connection.
	Read(b []byte) (n int, err error)
	// Write writes data to the connection.
	Write(b []byte) (n int, err error)
	// Close closes the connection.
	Close() error
}

const (
	// ProtocolString is the BitTorrent protocol identifier
	ProtocolString = "BitTorrent protocol"
	// HandshakeLength is the total length of a handshake message
	HandshakeLength = 68 // 1 + 19 + 8 + 20 + 20 bytes
)

// Handshake performs the BitTorrent handshake protocol with a peer.
// It writes the handshake message and reads the peer's response.
// Returns the peer ID in hexadecimal format and any error encountered.
func Handshake(tcpConn TCPConn, metadata *Metadata) (string, error) {
	defer tcpConn.Close()

	// Create the handshake message
	handshake := make([]byte, HandshakeLength)
	handshake[0] = byte(len(ProtocolString))      // Protocol string length
	copy(handshake[1:20], []byte(ProtocolString)) // Protocol string
	// Reserved bytes (8 zeros) are already zeroed
	copy(handshake[28:48], metadata.InfoHash[:]) // Info hash
	// Generate a random peer ID (20 bytes)
	peerID := make([]byte, 20)
	for i := range peerID {
		peerID[i] = byte(i + 1) // Simple deterministic peer ID for testing
	}
	copy(handshake[48:], peerID)

	// Write the handshake message
	if _, err := tcpConn.Write(handshake); err != nil {
		return "", err
	}

	// Read the peer's response
	response := make([]byte, HandshakeLength)
	if _, err := tcpConn.Read(response); err != nil {
		return "", err
	}

	// Validate the response
	if response[0] != byte(len(ProtocolString)) {
		return "", fmt.Errorf("invalid protocol string length: expected %d, got %d", len(ProtocolString), response[0])
	}

	if string(response[1:20]) != ProtocolString {
		return "", fmt.Errorf("invalid protocol string: expected %q, got %q", ProtocolString, string(response[1:20]))
	}

	// Verify the info hash matches
	for i := 0; i < 20; i++ {
		if response[28+i] != metadata.InfoHash[i] {
			return "", fmt.Errorf("info hash mismatch")
		}
	}

	// Return the peer ID in hexadecimal format
	return hex.EncodeToString(response[48:]), nil
}
