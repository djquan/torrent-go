package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
	"github.com/codecrafters-io/bittorrent-starter-go/internal/torrent"
)

func run(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("Usage: mybittorrent decode <bencoded-value>")
	}

	command := args[1]

	switch command {
	case "decode":
		return decode(args)
	case "info":
		return info(args)
	case "peers":
		return peers(args)
	case "handshake":
		return handshake(args)
	case "download_piece":
		return downloadPiece(args)
	case "download":
		return download(args)
	default:
		return "", fmt.Errorf("Unknown command: %s", command)
	}
}

// Create channels for coordinating downloads
type pieceResult struct {
	index int
	err   error
}

func download(args []string) (string, error) {
	if len(args) < 4 {
		return "", fmt.Errorf("Usage: mybittorrent download -o <output-file> <torrent-file>")
	}
	outputFile := args[3]
	torrentFile := args[4]
	info, err := torrent.ReadFromFile(torrentFile)
	if err != nil {
		return "", err
	}

	peers, err := torrent.Peers(http.DefaultClient, info)
	if err != nil {
		return "", fmt.Errorf("Error getting peers: %v", err)
	}

	// Create a buffer for each piece
	pieceBuffers := make([]*bytes.Buffer, len(info.PieceHashes))
	for i := range pieceBuffers {
		pieceBuffers[i] = &bytes.Buffer{}
	}

	resultChan := make(chan pieceResult, len(info.PieceHashes))
	peerChan := make(chan string, len(peers))

	// Fill peer channel with available peers
	for _, peer := range peers {
		peerChan <- peer
	}

	// Start parallel downloads
	for i := range info.PieceHashes {
		go func(pieceIndex int) {
			// Get a peer from the channel
			peerAddr := <-peerChan
			// Return the peer to the channel when done
			defer func() { peerChan <- peerAddr }()

			// Create a new connection for this piece
			conn, err := net.Dial("tcp", peerAddr)
			if err != nil {
				resultChan <- pieceResult{pieceIndex, fmt.Errorf("Failed to connect to peer %s: %v", peerAddr, err)}
				return
			}
			defer conn.Close()

			_, err = torrent.Handshake(conn, info)
			if err != nil {
				resultChan <- pieceResult{pieceIndex, fmt.Errorf("Handshake failed: %v", err)}
				return
			}

			err = torrent.DownloadPiece(conn, pieceBuffers[pieceIndex], info, pieceIndex)
			if err != nil {
				resultChan <- pieceResult{pieceIndex, fmt.Errorf("Failed to download piece %d: %v", pieceIndex, err)}
				return
			}

			resultChan <- pieceResult{pieceIndex, nil}
		}(i)
	}

	// Collect results
	for range info.PieceHashes {
		result := <-resultChan
		if result.err != nil {
			return "", result.err
		}
	}

	// Write all pieces to the output file in order
	outputFileHandle, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFileHandle.Close()

	for i, buffer := range pieceBuffers {
		_, err := io.Copy(outputFileHandle, buffer)
		if err != nil {
			return "", fmt.Errorf("failed to write piece %d to file: %v", i, err)
		}
	}

	return "download complete", nil
}

func downloadPiece(args []string) (string, error) {
	if len(args) < 5 {
		return "", fmt.Errorf("Usage: mybittorrent download_piece -o <output-file> <torrent-file> <piece-index>")
	}

	outputFile := args[3]
	torrentFile := args[4]
	pieceIndex, err := strconv.Atoi(args[5])
	if err != nil {
		return "", fmt.Errorf("Invalid piece index: %v", err)
	}

	info, err := torrent.ReadFromFile(torrentFile)
	if err != nil {
		return "", err
	}

	peers, err := torrent.Peers(http.DefaultClient, info)
	if err != nil {
		return "", fmt.Errorf("Error getting peers: %v", err)
	}

	peerAddr := peers[0]

	// Create TCP connection to peer
	conn, err := net.Dial("tcp", peerAddr)
	if err != nil {
		return "", fmt.Errorf("Failed to connect to peer %s: %v", peerAddr, err)
	}
	defer conn.Close()

	_, err = torrent.Handshake(conn, info)
	if err != nil {
		return "", fmt.Errorf("Handshake failed: %v", err)
	}

	outputFileHandle, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFileHandle.Close()

	// Use the file writer for the piece data
	err = torrent.DownloadPiece(conn, outputFileHandle, info, pieceIndex)
	if err != nil {
		return "", fmt.Errorf("Failed to download piece: %v", err)
	}

	return "Piece downloaded successfully", nil
}

func handshake(args []string) (string, error) {
	info, err := torrent.ReadFromFile(args[2])
	if err != nil {
		return "", err
	}

	// Split the peer address into host and port
	peerAddr := args[3]
	if !strings.Contains(peerAddr, ":") {
		return "", fmt.Errorf("Invalid peer address format. Expected <ip>:<port>, got %s", peerAddr)
	}

	// Create TCP connection to peer
	conn, err := net.Dial("tcp", peerAddr)
	if err != nil {
		return "", fmt.Errorf("Failed to connect to peer %s: %v", peerAddr, err)
	}

	defer conn.Close()

	// Perform handshake
	peerID, err := torrent.Handshake(conn, info)
	if err != nil {
		return "", fmt.Errorf("Handshake failed: %v", err)
	}

	return "Peer ID: " + peerID, nil
}

func peers(args []string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("Missing torrent file")
	}

	info, err := torrent.ReadFromFile(args[2])
	if err != nil {
		return "", err
	}

	peers, err := torrent.Peers(http.DefaultClient, info)
	if err != nil {
		return "", fmt.Errorf("Error getting peers: %v", err)
	}

	return strings.Join(peers, "\n"), nil
}

func info(args []string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("Missing torrent file")
	}

	info, err := torrent.ReadFromFile(args[2])
	if err != nil {
		return "", err
	}

	output := "Tracker URL: " + info.Announce + "\n" +
		"Length: " + fmt.Sprint(info.Length) + "\n" +
		"Info Hash: " + fmt.Sprintf("%x", info.InfoHash) + "\n" +
		"Piece Length: " + fmt.Sprint(info.PieceLength) + "\n" +
		"Piece Hashes:\n" + strings.Join(info.PieceHashes, "\n")
	return output, nil
}

func decode(args []string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("Missing bencoded value")
	}
	decoded, _, err := bencode.Decode([]byte(args[2]))
	if err != nil {
		return "", err
	}
	jsonOutput, err := customMarshal(decoded)
	if err != nil {
		return "", fmt.Errorf("Error converting to JSON: %v", err)
	}

	return string(jsonOutput), nil
}

// convertBytesToStrings recursively converts all []byte to strings in the interface{}
func convertBytesToStrings(v any) any {
	switch v := v.(type) {
	case []byte:
		return string(v)
	case []any:
		for i, val := range v {
			v[i] = convertBytesToStrings(val)
		}
		return v
	case map[string]any:
		for k, val := range v {
			v[k] = convertBytesToStrings(val)
		}
		return v
	default:
		return v
	}
}

func customMarshal(v any) ([]byte, error) {
	converted := convertBytesToStrings(v)
	return json.Marshal(converted)
}

func main() {
	output, err := run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(output)
}
