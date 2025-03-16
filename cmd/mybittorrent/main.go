package main

import (
	"encoding/json"
	"fmt"
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
	default:
		return "", fmt.Errorf("Unknown command: %s", command)
	}
}

func info(args []string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("Missing torrent file")
	}

	filenameArg := args[2]
	content, err := os.ReadFile(filenameArg)
	if err != nil {
		return "", fmt.Errorf("Error reading file: %v", err)
	}
	info, err := torrent.Info(content)
	if err != nil {
		return "", fmt.Errorf("Error parsing torrent file: %v", err)
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
func convertBytesToStrings(v interface{}) interface{} {
	switch v := v.(type) {
	case []byte:
		return string(v)
	case []interface{}:
		for i, val := range v {
			v[i] = convertBytesToStrings(val)
		}
		return v
	case map[string]interface{}:
		for k, val := range v {
			v[k] = convertBytesToStrings(val)
		}
		return v
	default:
		return v
	}
}

func customMarshal(v interface{}) ([]byte, error) {
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
