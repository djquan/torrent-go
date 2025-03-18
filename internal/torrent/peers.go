package torrent

import (
	"encoding/binary"
	"fmt"
	"io"
	"net/http"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
)

// HTTPClient represents an HTTP client capable of making requests to the tracker.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Peers contacts the tracker and returns a list of peers in the format "IP:port".
func Peers(httpClient HTTPClient, metadata *Metadata) ([]string, error) {
	request, err := http.NewRequest("GET", metadata.Announce, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	q := request.URL.Query()
	q.Add("info_hash", string(metadata.InfoHash[:]))
	q.Add("peer_id", "99999999999999999999")
	q.Add("port", "6881")
	q.Add("uploaded", "0")
	q.Add("downloaded", "0")
	q.Add("left", fmt.Sprintf("%d", metadata.Length))
	q.Add("compact", "1")
	request.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	decoded, _, err := bencode.Decode(body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	root, ok := decoded.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected decoded value to be a dictionary, but got %T", decoded)
	}

	peers, ok := root["peers"].([]byte)
	if !ok {
		return nil, fmt.Errorf("expected peers to be a byte string, but got %T", root["peers"])
	}

	peersList := make([]string, 0, len(peers)/6)
	for i := 0; i < len(peers); i += 6 {
		ip := fmt.Sprintf("%d.%d.%d.%d", peers[i], peers[i+1], peers[i+2], peers[i+3])
		port := binary.BigEndian.Uint16(peers[i+4 : i+6])
		peersList = append(peersList, fmt.Sprintf("%s:%d", ip, port))
	}

	return peersList, nil
}
