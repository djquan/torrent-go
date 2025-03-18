package torrent

import (
	"encoding/hex"
	"testing"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/testutil"
)

func TestHandshake(t *testing.T) {
	// Create a mock TCP connection
	mockConn := testutil.NewMockTCPConn()

	// Create test metadata
	metadata := &Metadata{
		InfoHash: [20]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
			0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14},
	}

	// Set up the peer's response
	// The response should be: <pstrlen><pstr><reserved><info_hash><peer_id>
	// where pstr is "BitTorrent protocol" (19 bytes)
	peerResponse := []byte{
		19,                                                                                            // pstrlen
		'B', 'i', 't', 'T', 'o', 'r', 'r', 'e', 'n', 't', ' ', 'p', 'r', 'o', 't', 'o', 'c', 'o', 'l', // pstr
		0, 0, 0, 0, 0, 0, 0, 0, // reserved bytes
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, // info_hash
		0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14,
		'P', 'E', 'E', 'R', 'I', 'D', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', // peer_id
		0, 0, 0, 0, // padding to make it 20 bytes
	}
	mockConn.SetReadData(peerResponse)

	// Perform the handshake
	peerID, err := Handshake(mockConn, metadata)
	if err != nil {
		t.Fatalf("Handshake failed: %v", err)
	}

	// Verify the peer ID was returned correctly in hex format
	expectedPeerID := hex.EncodeToString(append([]byte("PEERID1234567890"), 0, 0, 0, 0))
	if peerID != expectedPeerID {
		t.Errorf("expected peer ID %q, got %q", expectedPeerID, peerID)
	}

	// Verify the handshake message was written correctly
	writtenData := mockConn.GetWrittenData()
	if len(writtenData) != 68 { // 1 + 19 + 8 + 20 + 20 bytes
		t.Errorf("expected 68 bytes written, got %d", len(writtenData))
	}

	// Verify the protocol string length
	if writtenData[0] != 19 {
		t.Errorf("expected protocol string length 19, got %d", writtenData[0])
	}

	// Verify the protocol string
	protocol := string(writtenData[1:20])
	if protocol != "BitTorrent protocol" {
		t.Errorf("expected protocol string 'BitTorrent protocol', got %q", protocol)
	}

	// Verify the info hash
	for i := 0; i < 20; i++ {
		if writtenData[28+i] != metadata.InfoHash[i] {
			t.Errorf("info hash mismatch at position %d: expected %x, got %x", i, metadata.InfoHash[i], writtenData[28+i])
		}
	}

	// Verify the connection was closed
	if !mockConn.CloseCalled {
		t.Error("expected connection to be closed")
	}
}
