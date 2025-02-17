package main

import (
	"encoding/json"
	"fmt"
	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
	"github.com/codecrafters-io/bittorrent-starter-go/internal/torrent"
	"os"
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

	output := "Tracker URL: " + info.Announce + "\n" + "Length: " + fmt.Sprint(info.Length)
	return output, nil
}

func decode(args []string) (string, error) {
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
