package decode

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

func BencodeDecode[T any](data []byte) (T, error) {
	var result T
	reader := bytes.NewBuffer(data)
	decoded, err := decodeElement[T](reader)
	if err != nil {
		return result, err
	}
	return decoded, nil
}
func decodeElement[T any](buffer *bytes.Buffer) (T, error) {
	var result T
	b, err := buffer.ReadByte()
	if err != nil {
		return result, errors.New("unexpected endof data")
	}
	buffer.UnreadByte()
	switch {
	case b == 'i':
		decoded, err := decodeInteger(buffer)
		if err != nil {
			return result, err
		}
		if value, ok := any(decoded).(T); ok {
			return value, nil
		}
		return result, fmt.Errorf("type mismatch: expected %T, got int64", result)
	case b == 'l':
		decoded, err := decodeList[T](buffer)
		if err != nil {
			return result, err
		}
		if value, ok := any(decoded).(T); ok {
			return value, nil
		}
		return result, fmt.Errorf("type mismatch: expected %T, got list", result)
	case b == 'd':
		decoded, err := decodeDictionary[T](buffer)
		if err != nil {
			return result, err
		}
		if value, ok := any(decoded).(T); ok {
			return value, nil
		}
		return result, fmt.Errorf("type mismatch: expected %T, got dict", result)
	case b >= '0' && b <= '9':
		decoded, err := decodeString(buffer)
		if err != nil {
			return result, err
		}
		if value, ok := any(decoded).(T); ok {
			return value, nil
		}
		return result, fmt.Errorf("type Mismatch: expected %T, got string", result)
	default:
		return result, fmt.Errorf("invalid bencode format, unexpected byte: %c", b)
	}
}

func decodeInteger(buffer *bytes.Buffer) (int64, error) {
	if _, err := buffer.ReadByte(); err != nil {
		return 0, err
	}
	numBytes, err := buffer.ReadBytes('e')
	if err != nil {
		return 0, err
	}
	numStr := string(numBytes[:len(numBytes)-1])
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func decodeString(buffer *bytes.Buffer) (string, error) {
	lengthStr, err := buffer.ReadBytes(':')
	if err != nil {
		return "", err
	}
	length, err := strconv.Atoi(string(lengthStr[:len(lengthStr)-1]))
	if err != nil || length < 0 {
		return "", errors.New("invalid string length")
	}
	str := make([]byte, length)
	if _, err := buffer.Read(str); err != nil {
		return "", err
	}
	return string(str), nil
}

func decodeList[T any](buffer *bytes.Buffer) ([]T, error) {
	var result []T
	if _, err := buffer.ReadByte(); err != nil {
		return result, err
	}
	list := []T{}
	for {
		b, err := buffer.ReadByte()
		if err != nil {
			return result, errors.New("unterminated list")
		}
		if b == 'e' {
			return list, nil
		}
		buffer.UnreadByte()
		item, err := decodeElement[T](buffer)
		if err != nil {
			return result, err
		}
		list = append(list, item)
	}
}

func decodeDictionary[T any](buffer *bytes.Buffer) (map[string]T, error) {
	if _, err := buffer.ReadByte(); err != nil {
		return nil, err
	}
	dict := make(map[string]T)
	for {
		b, err := buffer.ReadByte()
		if err != nil {
			return nil, errors.New("unterminated dictionary")
		}
		if b == 'e' {
			return dict, nil
		}
		buffer.UnreadByte()
		key, err := decodeString(buffer)
		if err != nil {
			return nil, err
		}
		value, err := decodeElement[T](buffer)
		if err != nil {
			return nil, err
		}
		dict[key] = value
	}
}
