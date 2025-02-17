package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

func decodeString(bencodedString string) (string, int, error) {
	var firstColonIndex int
	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[:firstColonIndex]
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid string length: %v", err)
	}

	if firstColonIndex+1+length > len(bencodedString) {
		return "", 0, fmt.Errorf("string length exceeds input length")
	}

	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], firstColonIndex + 1 + length, nil
}

func decodeInteger(bencodedString string) (int, int, error) {
	lastIndex := strings.Index(bencodedString, "e")
	if lastIndex == -1 {
		return 0, 0, fmt.Errorf("invalid integer format: missing 'e'")
	}

	// Get the full number string without skipping the first character
	numStr := bencodedString[:lastIndex]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid integer: %v", err)
	}

	return num, lastIndex + 1, nil
}

func decodeBencode(bencodedString string) (interface{}, int, error) {
	if len(bencodedString) == 0 {
		return "", 0, fmt.Errorf("empty string is not valid bencode")
	}

	switch {
	case unicode.IsDigit(rune(bencodedString[0])):
		value, i, err := decodeString(bencodedString)
		return value, i, err

	case bencodedString[0] == 'i':
		value, i, err := decodeInteger(bencodedString[1:])
		return value, i, err

	case bencodedString[0] == 'l':
		var list []interface{}
		consumed := 1 // Start with 1 for the 'l'
		remaining := bencodedString[1:]

		for len(remaining) > 0 && remaining[0] != 'e' {
			decoded, i, err := decodeBencode(remaining)
			if err != nil {
				return nil, 0, err
			}
			list = append(list, decoded)
			consumed += i
			remaining = remaining[i:]
		}

		if remaining[0] != 'e' {
			return nil, 0, fmt.Errorf("unterminated list")
		}

		return list, consumed + 1, nil // +1 for the 'e'

	default:
		return "", 0, fmt.Errorf("unsupported bencode type")
	}
}

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		decoded, _, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
