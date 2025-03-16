package torrent

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
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

type MockHTTPClient struct {
	Requests *http.Request
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.Requests = req
	responseData := map[string]any{
		"interval": 0,
		"peers":    []byte{165, 232, 41, 73, 201, 84},
	}

	encodedResponse, err := bencode.Encode(responseData)
	if err != nil {
		return nil, err
	}

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBuffer(encodedResponse)),
	}
	return response, nil
}

func TestPeers(t *testing.T) {
	content, err := os.ReadFile("../../sample.torrent")
	if err != nil {
		t.Fatalf("failed to read sample.torrent: %v", err)
	}

	info, err := Info(content)
	if err != nil {
		t.Fatalf("failed to parse torrent file: %v", err)
	}

	mockHTTPClient := &MockHTTPClient{}

	peers, err := Peers(mockHTTPClient, info)
	if err != nil {
		t.Fatalf("failed to get peers: %v", err)
	}

	request := mockHTTPClient.Requests

	if request == nil {
		t.Errorf("expected request to be non-nil: did not call Do")
	}

	if request.URL.String() != "http://bittorrent-test-tracker.codecrafters.io/announce?compact=1&downloaded=0&info_hash=%D6%9F%91%E6%B2%AELT%24h%D1%07%3Aq%D4%EA%13%87%9A%7F&left=92063&peer_id=99999999999999999999&port=6881&uploaded=0" {
		t.Errorf("expected announce URL to be 'http://bittorrent-test-tracker.codecrafters.io/announce', but got '%s'", request.URL.String())
	}
	t.Logf("peers: %v", peers)
}
