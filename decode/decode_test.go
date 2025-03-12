package decode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBencodeDecodeInteger(t *testing.T) {
	data := []byte("i42e")
	result, err := BencodeDecode[any](data)
	assert.Nil(t, err)
	assert.Equal(t, int64(42), result)

	data = []byte("i-15e")
	result, err = BencodeDecode[any](data)
	assert.Nil(t, err)
	assert.Equal(t, int64(-15), result)
}
func TestBencodeDecodeString(t *testing.T) {
	data := []byte("5:hello")
	result, err := BencodeDecode[string](data)
	assert.Nil(t, err)
	assert.Equal(t, "hello", result)

	data = []byte("0:")
	result, err = BencodeDecode[string](data)
	assert.Nil(t, err)
	assert.Equal(t, "", result)
}
func TestBencodeDecodeList(t *testing.T) {
	data := []byte("li42e4:teste")
	result, err := BencodeDecode[any](data)
	assert.Nil(t, err)
	assert.Equal(t, []any{int64(42), "test"}, result)

	data = []byte("le")
	result, err = BencodeDecode[any](data)
	assert.Nil(t, err)
	assert.Equal(t, []any{}, result)
}
func TestBencodeDecodeDictionary(t *testing.T) {
	data := []byte("d3:agei25e4:name5:Alicee")
	result, err := BencodeDecode[any](data)
	expected := map[string]any{
		"age":  int64(25),
		"name": "Alice",
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestBencodeDecodeInvalidCases(t *testing.T) {
	data := []byte("i42")
	_, err := BencodeDecode[any](data)
	assert.NotNil(t, err)

	data = []byte("5hello")
	_, err = BencodeDecode[any](data)
	assert.NotNil(t, err)

	data = []byte("d3:age25e")
	_, err = BencodeDecode[any](data)
	assert.NotNil(t, err)
}
