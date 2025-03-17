package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	default:
		return "", fmt.Errorf("Unknown command: %s", command)
	}
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
