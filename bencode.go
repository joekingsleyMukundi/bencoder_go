package bencoder

import (
	"github.com/joekingsleyMukundi/bencoder_go/decode"
	"github.com/joekingsleyMukundi/bencoder_go/encode"
)

type BencodeTypes interface {
	int | int64 | string | []any | map[string]any
}

func Encode[T BencodeTypes](data T) (string, error) {
	return encode.BencodeEncode(data)
}

func Decode[T any](data []byte) (T, error) {
	return decode.BencodeDecode[T](data)
}
