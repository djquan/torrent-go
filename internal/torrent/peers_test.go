package torrent

import (
	"os"
	"testing"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
	"github.com/codecrafters-io/bittorrent-starter-go/internal/testutil"
)

func TestPeers(t *testing.T) {
	content, err := os.ReadFile("../../sample.torrent")
	if err != nil {
		t.Fatalf("failed to read sample.torrent: %v", err)
	}

	info, err := Info(content)
	if err != nil {
		t.Fatalf("failed to parse torrent file: %v", err)
	}

	responseData := map[string]any{
		"interval": 0,
		"peers":    []byte{165, 232, 41, 73, 201, 84},
	}

	encodedResponse, err := bencode.Encode(responseData)
	if err != nil {
		t.Fatalf("failed to encode response: %v", err)
	}

	mockHTTPClient := &testutil.MockHTTPClient{
		Response: encodedResponse,
	}

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

	if len(peers) != 1 || peers[0] != "165.232.41.73:51540" {
		t.Errorf("expected peers to be ['165.232.41.73:51540'], but got %v", peers)
	}
}
