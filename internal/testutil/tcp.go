package testutil

import (
	"bytes"
)

// MockTCPConn is a mock implementation of a TCP connection for testing.
type MockTCPConn struct {
	// ReadBuffer contains the data that will be returned by Read calls
	ReadBuffer *bytes.Buffer
	// WriteBuffer contains the data that was written via Write calls
	WriteBuffer *bytes.Buffer
	// CloseCalled tracks if Close was called
	CloseCalled bool
}

// NewMockTCPConn creates a new MockTCPConn with empty buffers.
func NewMockTCPConn() *MockTCPConn {
	return &MockTCPConn{
		ReadBuffer:  &bytes.Buffer{},
		WriteBuffer: &bytes.Buffer{},
	}
}

// Read implements the TCPConn interface by reading from the ReadBuffer.
func (m *MockTCPConn) Read(b []byte) (n int, err error) {
	return m.ReadBuffer.Read(b)
}

// Write implements the TCPConn interface by writing to the WriteBuffer.
func (m *MockTCPConn) Write(b []byte) (n int, err error) {
	return m.WriteBuffer.Write(b)
}

// Close implements the TCPConn interface by marking Close as called.
func (m *MockTCPConn) Close() error {
	m.CloseCalled = true
	return nil
}

// SetReadData sets the data that will be returned by subsequent Read calls.
func (m *MockTCPConn) SetReadData(data []byte) {
	m.ReadBuffer.Reset()
	m.ReadBuffer.Write(data)
}

// GetWrittenData returns the data that was written via Write calls.
func (m *MockTCPConn) GetWrittenData() []byte {
	return m.WriteBuffer.Bytes()
}
