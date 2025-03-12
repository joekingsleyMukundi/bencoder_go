package decode

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
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
		return result, errors.New("unexpected end of data")
	}
	buffer.UnreadByte()

	switch {
	case b == 'i':
		decoded, err := decodeInteger(buffer)
		if err != nil {
			return result, err
		}
		return typeAssert[T](decoded)
	case b == 'l':
		decoded, err := decodeList[T](buffer)
		if err != nil {
			return result, err
		}
		return typeAssert[T](decoded)
	case b == 'd':
		decoded, err := decodeDictionary[T](buffer)
		if err != nil {
			return result, err
		}
		return typeAssert[T](decoded)
	case b >= '0' && b <= '9':
		decoded, err := decodeString(buffer)
		if err != nil {
			return result, err
		}
		return typeAssert[T](decoded)
	default:
		return result, fmt.Errorf("invalid bencode format, unexpected byte: %c", b)
	}
}

func typeAssert[T any](value any) (T, error) {
	if result, ok := value.(T); ok {
		return result, nil
	}
	var zeroValue T
	return zeroValue, fmt.Errorf("type mismatch: expected %T, got %T", zeroValue, value)
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
	return strconv.ParseInt(numStr, 10, 64)
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
	result := []T{}
	if _, err := buffer.ReadByte(); err != nil {
		return result, err
	}
	for {
		b, err := buffer.ReadByte()
		if err != nil {
			return result, errors.New("unterminated list")
		}
		if b == 'e' {
			return result, nil
		}
		buffer.UnreadByte()
		item, err := decodeElement[T](buffer)
		if err != nil {
			return result, err
		}
		result = append(result, item)
	}
}
func decodeDictionary[T any](buffer *bytes.Buffer) (any, error) {
	if _, err := buffer.ReadByte(); err != nil {
		return nil, err
	}
	bencodeMap := make(map[string]any)

	for {
		b, err := buffer.ReadByte()
		if err != nil {
			return nil, errors.New("unterminated dictionary")
		}
		if b == 'e' {
			break
		}
		buffer.UnreadByte()
		key, err := decodeString(buffer)
		if err != nil {
			return nil, err
		}
		value, err := decodeElement[any](buffer)
		if err != nil {
			return nil, err
		}
		bencodeMap[key] = value
	}
	var result T
	v := reflect.ValueOf(&result).Elem()
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			bencodeTag := field.Tag.Get("bencode")
			if bencodeTag == "" {
				bencodeTag = field.Name
			}
			if value, exists := bencodeMap[bencodeTag]; exists {
				fieldValue := v.Field(i)
				if fieldValue.CanSet() {
					fieldValue.Set(reflect.ValueOf(value))
				}
			}
		}
		return result, nil
	}
	return bencodeMap, nil
}
