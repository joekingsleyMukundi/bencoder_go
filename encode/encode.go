package encode

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"
)

type BencodeTypes interface {
	int | int64 | string | []any | map[string]any
}

func BencodeEncode[T BencodeTypes](data T) (string, error) {
	return bencodeHelper(any(data))
}
func bencodeHelper(data any) (string, error) {
	var buf bytes.Buffer
	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Int, reflect.Int64:
		buf.WriteString(fmt.Sprintf("i%de", v.Int()))

	case reflect.String:
		buf.WriteString(fmt.Sprintf("%d:%s", len(v.String()), v.String()))

	case reflect.Slice:
		buf.WriteString("l")
		for i := range v.Len() {
			encoded, err := bencodeHelper(v.Index(i).Interface())
			if err != nil {
				return "", err
			}
			buf.WriteString(encoded)
		}
		buf.WriteString("e")

	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return "", errors.New("bencode only supports maps with string keys")
		}
		buf.WriteString("d")
		keys := make([]string, 0, v.Len())
		for _, key := range v.MapKeys() {
			keys = append(keys, key.String())
		}
		sort.Strings(keys)
		for _, k := range keys {
			keyEncoded, err := bencodeHelper(k)
			if err != nil {
				return "", err
			}
			valEncoded, err := bencodeHelper(v.MapIndex(reflect.ValueOf(k)).Interface())
			if err != nil {
				return "", err
			}
			buf.WriteString(keyEncoded + valEncoded)
		}
		buf.WriteString("e")
	case reflect.Struct:
		buf.WriteString("d")
		t := v.Type()
		keys := make([]string, 0, v.NumField())
		fieldMap := make(map[string]string)

		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			bencodeTag := field.Tag.Get("bencode")
			if bencodeTag == "" {
				bencodeTag = field.Name // Use field name if no tag
			}
			keys = append(keys, bencodeTag)
			fieldEncoded, err := bencodeHelper(v.Field(i).Interface())
			if err != nil {
				return "", err
			}
			fieldMap[bencodeTag] = fieldEncoded
		}

		sort.Strings(keys) // Ensure dictionary keys are sorted
		for _, k := range keys {
			keyEncoded, err := bencodeHelper(k)
			if err != nil {
				return "", err
			}
			buf.WriteString(keyEncoded + fieldMap[k])
		}
		buf.WriteString("e")

	default:
		return "", errors.New("unsupported data type for Bencode")
	}
	return buf.String(), nil
}
