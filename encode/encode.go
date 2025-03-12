package encode

import (
	"bytes"
	"errors"
	"fmt"
)

type BencodeTypes interface {
	int | int64 | string | []any | map[string]any
}

func BencodeEncode[T BencodeTypes](data T) (string, error) {
	return bencodeHelper(any(data))
}
func bencodeHelper(data any) (string, error) {
	var buf bytes.Buffer

	switch v := data.(type) {
	case int, int64:
		buf.WriteString(fmt.Sprintf("i%de", v))

	case string:
		buf.WriteString(fmt.Sprintf("%d:%s", len(v), v))

	case []any:
		buf.WriteString("l")
		for _, item := range v {
			encoded, err := bencodeHelper(item)
			if err != nil {
				return "", err
			}
			buf.WriteString(encoded)
		}
		buf.WriteString("e")

	case map[string]any:
		buf.WriteString("d")
		for k, val := range v {
			keyEncoded, err := bencodeHelper(k)
			if err != nil {
				return "", err
			}
			valEncoded, err := bencodeHelper(val)
			if err != nil {
				return "", err
			}
			buf.WriteString(keyEncoded + valEncoded)
		}
		buf.WriteString("e")
	default:
		return "", errors.New("unsupported data type for Bencode")
	}
	return buf.String(), nil
}
