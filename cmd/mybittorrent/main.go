package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
)

func run(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("Usage: mybittorrent decode <bencoded-value>")
	}

	command := args[1]
	if command != "decode" {
		return "", fmt.Errorf("Unknown command: %s", command)
	}

	if len(args) < 3 {
		return "", fmt.Errorf("Missing bencoded value")
	}

	decoded, _, err := bencode.Decode(args[2])
	if err != nil {
		return "", err
	}

	jsonOutput, err := json.Marshal(decoded)
	if err != nil {
		return "", fmt.Errorf("Error converting to JSON: %v", err)
	}

	return string(jsonOutput), nil
}

func main() {
	output, err := run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(output)
}
