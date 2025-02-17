package bencode

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

func decodeString(bencodedString []byte) ([]byte, int, error) {
	var firstColonIndex int
	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[:firstColonIndex]
	length, err := strconv.Atoi(string(lengthStr))
	if err != nil {
		return []byte{}, 0, fmt.Errorf("invalid string length: %v", err)
	}

	if firstColonIndex+1+length > len(bencodedString) {
		return []byte{}, 0, fmt.Errorf("string length exceeds input length")
	}

	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], firstColonIndex + 1 + length, nil
}

func decodeInteger(bencodedString []byte) (int, int, error) {
	// Expect 'i' prefix
	if len(bencodedString) == 0 || bencodedString[0] != 'i' {
		return 0, 0, fmt.Errorf("invalid integer format: missing 'i' prefix")
	}

	lastIndex := bytes.IndexByte(bencodedString, 'e')
	if lastIndex == -1 {
		return 0, 0, fmt.Errorf("invalid integer format: missing 'e'")
	}

	// Get the number string between 'i' and 'e'
	numStr := bencodedString[1:lastIndex]
	num, err := strconv.Atoi(string(numStr))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid integer: %v", err)
	}

	return num, lastIndex + 1, nil
}

func decodeList(bencodedString []byte) ([]interface{}, int, error) {
	if len(bencodedString) == 0 || bencodedString[0] != 'l' {
		return nil, 0, fmt.Errorf("invalid list format: missing 'l' prefix")
	}

	list := make([]interface{}, 0)
	consumed := 1 // Start with 1 for the 'l'
	remaining := bencodedString[1:]

	for len(remaining) > 0 && remaining[0] != 'e' {
		decoded, i, err := Decode(remaining)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, decoded)
		consumed += i
		remaining = remaining[i:]
	}

	if len(remaining) == 0 || remaining[0] != 'e' {
		return nil, 0, fmt.Errorf("unterminated list")
	}

	return list, consumed + 1, nil // +1 for the 'e'
}

func decodeDictionary(bencodedString []byte) (map[string]interface{}, int, error) {
	if len(bencodedString) == 0 || bencodedString[0] != 'd' {
		return nil, 0, fmt.Errorf("invalid dictionary format: missing 'd' prefix")
	}

	dict := make(map[string]interface{})
	consumed := 1 // Start with 1 for the 'd'
	remaining := bencodedString[1:]
	for len(remaining) > 0 && remaining[0] != 'e' {
		key, i, err := decodeString(remaining)
		if err != nil {
			return nil, 0, err
		}
		consumed += i
		remaining = remaining[i:]

		value, i, err := Decode(remaining)
		if err != nil {
			return nil, 0, err
		}
		consumed += i
		remaining = remaining[i:]

		dict[string(key)] = value
	}

	return dict, consumed + 1, nil
}

// Decode takes a bencoded string and decodes it into the go value.
func Decode(bencodedString []byte) (interface{}, int, error) {
	if len(bencodedString) == 0 {
		return "", 0, fmt.Errorf("empty string is not valid bencode")
	}

	switch {
	case unicode.IsDigit(rune(bencodedString[0])):
		value, i, err := decodeString(bencodedString)
		return value, i, err

	case bencodedString[0] == 'i':
		value, i, err := decodeInteger(bencodedString[0:])
		return value, i, err

	case bencodedString[0] == 'l':
		value, i, err := decodeList(bencodedString)
		return value, i, err

	case bencodedString[0] == 'd':
		value, i, err := decodeDictionary(bencodedString)
		return value, i, err

	default:
		return "", 0, fmt.Errorf("unsupported bencode type")
	}
}
