package torrent

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/testutil"
)

func TestDownloadPiece(t *testing.T) {
	// Create a mock TCP connection
	mockConn := testutil.NewMockTCPConn()
	defer mockConn.Close()

	// Create a buffer to capture written data
	var outputBuffer bytes.Buffer

	// Create test metadata
	metadata := &Metadata{
		InfoHash: [20]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
			0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14},
		PieceLength: 16384, // Match the block size used in the request message
		Length:      16384, // One piece worth of data
		PieceHashes: []string{
			"0123456789abcdef0123456789abcdef01234567", // Hash for piece 0
		},
	}

	// Set up the complete peer message sequence
	// 1. Bitfield message
	// 2. Interested message
	// 3. Unchoke message
	// 4. Request message
	// 5. Piece message
	peerMessages := []byte{
		// Bitfield message
		0x00, 0x00, 0x00, 0x03, // length prefix (3 bytes)
		0x05,       // message type (bitfield)
		0x80, 0x00, // bitfield (piece 0 available - 1000 0000)

		// Unchoke message
		0x00, 0x00, 0x00, 0x01, // length prefix (1 byte)
		0x01, // message type (unchoke)

		// Piece message
		0x00, 0x00, 0x00, 0x0A, // length prefix (10 bytes)
		0x07,                   // message type (piece)
		0x00, 0x00, 0x00, 0x00, // piece index
		0x00, 0x00, 0x00, 0x00, // block offset
		0x01, // data (1 byte)
	}
	mockConn.SetReadData(peerMessages)

	expectedInput := []byte{
		// Interested message
		0x00, 0x00, 0x00, 0x01, // length prefix (1 byte)
		0x02, // message type (interested)

		// Request message
		0x00, 0x00, 0x00, 0x0D, // length prefix (13 bytes)
		0x06,                   // message type (request)
		0x00, 0x00, 0x00, 0x00, // piece index
		0x00, 0x00, 0x00, 0x00, // block offset
		0x00, 0x00, 0x40, 0x00, // block length (16384 bytes)
	}

	// Perform the peer message handling
	err := DownloadPiece(mockConn, &outputBuffer, metadata, 0)
	if err != nil {
		t.Fatalf("Peer message handling failed: %v", err)
	}

	// Verify the messages sent to the peer
	writtenData := mockConn.GetWrittenData()

	if !bytes.Equal(writtenData, expectedInput) {
		t.Errorf("Expected written data %v, got %v", expectedInput, writtenData)
	}

	// Verify the piece data was written to the output buffer
	expectedOutput := []byte{0x01} // The data from the piece message
	if !bytes.Equal(outputBuffer.Bytes(), expectedOutput) {
		t.Errorf("Expected output buffer to contain %v, got %v", expectedOutput, outputBuffer.Bytes())
	}
}

func TestPieceBlockRequesting(t *testing.T) {
	// Create a mock TCP connection
	mockConn := testutil.NewMockTCPConn()
	defer mockConn.Close()

	// Create test metadata for a 32KB piece (2 blocks)
	metadata := &Metadata{
		InfoHash: [20]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
			0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14},
		PieceLength: 32768, // 32KB piece
		Length:      32768,
		PieceHashes: []string{
			"0123456789abcdef0123456789abcdef01234567", // Hash for piece 0
		},
	}

	// Create the piece message data (16384 bytes for each block)
	blockSize := 16384
	firstBlockData := make([]byte, blockSize)
	secondBlockData := make([]byte, blockSize)
	for i := 0; i < blockSize; i++ {
		firstBlockData[i] = 0x01
		secondBlockData[i] = 0x02
	}

	// Set up peer messages:
	// 1. Bitfield message (has piece 0)
	// 2. Unchoke message
	// 3. Two piece messages (one for each block)
	peerMessages := make([]byte, 0)

	// Bitfield message
	peerMessages = append(peerMessages,
		0x00, 0x00, 0x00, 0x03, // length prefix (3 bytes)
		0x05,       // message type (bitfield)
		0x80, 0x00, // bitfield (piece 0 available - 1000 0000)
	)

	// Unchoke message
	peerMessages = append(peerMessages,
		0x00, 0x00, 0x00, 0x01, // length prefix (1 byte)
		0x01, // message type (unchoke)
	)

	// First piece message
	firstPieceHeader := []byte{
		0x00, 0x00, 0x40, 0x09, // length prefix (16393 bytes)
		0x07,                   // message type (piece)
		0x00, 0x00, 0x00, 0x00, // piece index
		0x00, 0x00, 0x00, 0x00, // block offset
	}
	peerMessages = append(peerMessages, firstPieceHeader...)
	peerMessages = append(peerMessages, firstBlockData...)

	// Second piece message
	secondPieceHeader := []byte{
		0x00, 0x00, 0x40, 0x09, // length prefix (16393 bytes)
		0x07,                   // message type (piece)
		0x00, 0x00, 0x00, 0x00, // piece index
		0x00, 0x00, 0x40, 0x00, // block offset (16384)
	}
	peerMessages = append(peerMessages, secondPieceHeader...)
	peerMessages = append(peerMessages, secondBlockData...)

	mockConn.SetReadData(peerMessages)

	// Expected messages to be sent:
	// 1. Interested message
	// 2. Two request messages (one for each block)
	expectedInput := []byte{
		// Interested message
		0x00, 0x00, 0x00, 0x01, // length prefix (1 byte)
		0x02, // message type (interested)

		// First request message
		0x00, 0x00, 0x00, 0x0D, // length prefix (13 bytes)
		0x06,                   // message type (request)
		0x00, 0x00, 0x00, 0x00, // piece index
		0x00, 0x00, 0x00, 0x00, // block offset
		0x00, 0x00, 0x40, 0x00, // block length (16384 bytes)

		// Second request message
		0x00, 0x00, 0x00, 0x0D, // length prefix (13 bytes)
		0x06,                   // message type (request)
		0x00, 0x00, 0x00, 0x00, // piece index
		0x00, 0x00, 0x40, 0x00, // block offset (16384)
		0x00, 0x00, 0x40, 0x00, // block length (16384 bytes)
	}

	// Create a buffer to capture written data
	var outputBuffer bytes.Buffer

	// Perform the peer message handling
	err := DownloadPiece(mockConn, &outputBuffer, metadata, 0)
	if err != nil {
		t.Fatalf("Peer message handling failed: %v", err)
	}

	// Verify the messages sent to the peer
	writtenData := mockConn.GetWrittenData()
	if !bytes.Equal(writtenData, expectedInput) {
		t.Errorf("Expected written data %v, got %v", expectedInput, writtenData)
	}

	// Verify the piece data was written to the output buffer
	// We expect both blocks to be concatenated
	expectedOutput := make([]byte, 32768)
	for i := 0; i < blockSize; i++ {
		expectedOutput[i] = 0x01
		expectedOutput[i+blockSize] = 0x02
	}
	if !bytes.Equal(outputBuffer.Bytes(), expectedOutput) {
		t.Errorf("Expected output buffer to contain %v, got %v", expectedOutput, outputBuffer.Bytes())
	}
}

func TestDownloadPiece1(t *testing.T) {
	// Create a mock TCP connection
	mockConn := testutil.NewMockTCPConn()
	defer mockConn.Close()

	// Create a buffer to capture written data
	var outputBuffer bytes.Buffer

	// Create test metadata
	metadata := &Metadata{
		InfoHash: [20]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
			0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14},
		PieceLength: 16384, // Match the block size used in the request message
		Length:      32768, // Two pieces worth of data
		PieceHashes: []string{
			"0123456789abcdef0123456789abcdef01234567", // Hash for piece 0
			"abcdef0123456789abcdef0123456789abcdef01", // Hash for piece 1
		},
	}

	// Set up the complete peer message sequence
	// 1. Bitfield message (has piece 1)
	// 2. Interested message
	// 3. Unchoke message
	// 4. Request message
	// 5. Piece message
	peerMessages := []byte{
		// Bitfield message
		0x00, 0x00, 0x00, 0x03, // length prefix (3 bytes)
		0x05,       // message type (bitfield)
		0x40, 0x00, // bitfield (piece 1 available - 0100 0000)

		// Unchoke message
		0x00, 0x00, 0x00, 0x01, // length prefix (1 byte)
		0x01, // message type (unchoke)

		// Piece message
		0x00, 0x00, 0x00, 0x0A, // length prefix (10 bytes)
		0x07,                   // message type (piece)
		0x00, 0x00, 0x00, 0x01, // piece index (1)
		0x00, 0x00, 0x00, 0x00, // block offset
		0x01, // data (1 byte)
	}
	mockConn.SetReadData(peerMessages)

	expectedInput := []byte{
		// Interested message
		0x00, 0x00, 0x00, 0x01, // length prefix (1 byte)
		0x02, // message type (interested)

		// Request message
		0x00, 0x00, 0x00, 0x0D, // length prefix (13 bytes)
		0x06,                   // message type (request)
		0x00, 0x00, 0x00, 0x01, // piece index (1)
		0x00, 0x00, 0x00, 0x00, // block offset
		0x00, 0x00, 0x40, 0x00, // block length (16384 bytes)
	}

	// Perform the peer message handling
	err := DownloadPiece(mockConn, &outputBuffer, metadata, 1)
	if err != nil {
		t.Fatalf("Peer message handling failed: %v", err)
	}

	// Verify the messages sent to the peer
	writtenData := mockConn.GetWrittenData()

	if !bytes.Equal(writtenData, expectedInput) {
		t.Errorf("Expected written data %v, got %v", expectedInput, writtenData)
	}

	// Verify the piece data was written to the output buffer
	expectedOutput := []byte{0x01} // The data from the piece message
	if !bytes.Equal(outputBuffer.Bytes(), expectedOutput) {
		t.Errorf("Expected output buffer to contain %v, got %v", expectedOutput, outputBuffer.Bytes())
	}
}
