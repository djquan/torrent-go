package torrent

import (
	"encoding/binary"
	"fmt"
	"io"
)

// MessageType represents the type of a peer message
type MessageType byte

const (
	MessageTypeChoke MessageType = iota
	MessageTypeUnchoke
	MessageTypeInterested
	MessageTypeNotInterested
	MessageTypeHave
	MessageTypeBitfield
	MessageTypeRequest
	MessageTypePiece
	MessageTypeCancel
)

// Message represents a parsed peer message
type Message struct {
	Type    MessageType
	Payload []byte
}

// readMessage reads a complete peer message from the connection
func readMessage(conn io.Reader) (*Message, error) {
	// Read message length (4 bytes)
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lengthBuf); err != nil {
		return nil, fmt.Errorf("failed to read message length: %w", err)
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// Read message type (1 byte)
	typeBuf := make([]byte, 1)
	if _, err := io.ReadFull(conn, typeBuf); err != nil {
		return nil, fmt.Errorf("failed to read message type: %w", err)
	}
	msgType := MessageType(typeBuf[0])

	// Read payload if any
	var payload []byte
	if length > 1 {
		payload = make([]byte, length-1)
		if _, err := io.ReadFull(conn, payload); err != nil {
			return nil, fmt.Errorf("failed to read message payload: %w", err)
		}
	}

	return &Message{
		Type:    msgType,
		Payload: payload,
	}, nil
}

// writeMessage writes a complete peer message to the connection
func writeMessage(conn io.Writer, msgType MessageType, payload []byte) error {
	// Calculate total message length (1 byte for type + payload length)
	length := uint32(1 + len(payload))

	// Write length prefix (4 bytes)
	lengthBuf := []byte{
		byte(length >> 24),
		byte(length >> 16),
		byte(length >> 8),
		byte(length),
	}
	if _, err := conn.Write(lengthBuf); err != nil {
		return fmt.Errorf("failed to write message length: %w", err)
	}

	// Write message type (1 byte)
	if _, err := conn.Write([]byte{byte(msgType)}); err != nil {
		return fmt.Errorf("failed to write message type: %w", err)
	}

	// Write payload if any
	if len(payload) > 0 {
		if _, err := conn.Write(payload); err != nil {
			return fmt.Errorf("failed to write message payload: %w", err)
		}
	}

	return nil
}

// handlePieceMessage handles a piece message by writing its data to the writer
func handlePieceMessage(msg *Message, writer io.Writer) error {
	if len(msg.Payload) < 8 {
		return fmt.Errorf("piece message payload too short: %d bytes", len(msg.Payload))
	}

	// Extract piece index and block offset
	// pieceIndex := binary.BigEndian.Uint32(msg.Payload[0:4])
	// blockOffset := binary.BigEndian.Uint32(msg.Payload[4:8])

	// Get the actual data (everything after the first 8 bytes)
	data := msg.Payload[8:]
	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("failed to write piece data: %w", err)
	}

	return nil
}

// hasPiece checks if a peer has a specific piece based on their bitfield message
func hasPiece(bitfield []byte, pieceIndex int) bool {
	if pieceIndex >= len(bitfield)*8 {
		return false
	}
	byteIndex := pieceIndex / 8
	bitIndex := 7 - (pieceIndex % 8)
	return (bitfield[byteIndex] & (1 << bitIndex)) != 0
}

// handleBitfieldMessage processes a bitfield message
func handleBitfieldMessage(msg *Message, pieceIndex int) error {
	if !hasPiece(msg.Payload, pieceIndex) {
		return fmt.Errorf("peer does not have piece %d", pieceIndex)
	}
	return nil
}

// sendRequest sends a request message for a specific block
func sendRequest(conn io.Writer, pieceIndex, blockOffset, blockLength uint32) error {
	requestPayload := make([]byte, 12)
	binary.BigEndian.PutUint32(requestPayload[0:4], pieceIndex)
	binary.BigEndian.PutUint32(requestPayload[4:8], blockOffset)
	binary.BigEndian.PutUint32(requestPayload[8:12], blockLength)
	return writeMessage(conn, MessageTypeRequest, requestPayload)
}

// handleUnchokeMessage processes an unchoke message and requests blocks
func handleUnchokeMessage(conn io.ReadWriter, writer io.Writer, metadata *Metadata, pieceIndex int) error {
	blockSize := uint32(16 * 1024) // 16KB blocks
	pieceLength := uint32(metadata.PieceLength)

	// Calculate the actual length of this piece
	startOffset := uint32(pieceIndex) * pieceLength
	remainingFileLength := uint32(metadata.Length) - startOffset
	actualPieceLength := pieceLength
	if remainingFileLength < pieceLength {
		actualPieceLength = remainingFileLength
	}

	numBlocks := (actualPieceLength + blockSize - 1) / blockSize // Ceiling division

	for blockIndex := uint32(0); blockIndex < numBlocks; blockIndex++ {
		blockOffset := blockIndex * blockSize
		remainingBytes := actualPieceLength - blockOffset
		remainingBytes = uint32(min(int(remainingBytes), int(blockSize)))

		// Send request message
		if err := sendRequest(conn, uint32(pieceIndex), blockOffset, remainingBytes); err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}

		// Read piece message
		pieceMsg, err := readMessage(conn)
		if err != nil {
			return fmt.Errorf("failed to read piece message: %w", err)
		}

		if pieceMsg.Type != MessageTypePiece {
			return fmt.Errorf("expected piece message (type 7), got type %d", pieceMsg.Type)
		}

		// Handle piece message
		if err := handlePieceMessage(pieceMsg, writer); err != nil {
			return fmt.Errorf("failed to handle piece message: %w", err)
		}
	}

	return nil
}

// DownloadPiece handles downloading a specific piece from a peer using the peer protocol
func DownloadPiece(conn io.ReadWriter, writer io.Writer, metadata *Metadata, pieceIndex int) error {
	// Send interested message
	if err := writeMessage(conn, MessageTypeInterested, nil); err != nil {
		return fmt.Errorf("failed to send interested message: %w", err)
	}

	// Read messages until we get the piece we want
	for {
		msg, err := readMessage(conn)
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		switch msg.Type {
		case MessageTypeBitfield:
			if err := handleBitfieldMessage(msg, pieceIndex); err != nil {
				return err
			}
		case MessageTypeHave:
			// Just continue the loop, waiting for unchoke
			continue
		case MessageTypeUnchoke:
			return handleUnchokeMessage(conn, writer, metadata, pieceIndex)
		case MessageTypeChoke:
			return fmt.Errorf("peer choked us")
		case MessageTypePiece:
			return fmt.Errorf("unexpected piece message received")
		default:
			return fmt.Errorf("invalid message type: %d", msg.Type)
		}
	}
}
