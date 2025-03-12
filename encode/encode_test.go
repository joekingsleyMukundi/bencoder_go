package encode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBencodeEncodeInteger(t *testing.T) {
	result, err := BencodeEncode(42)
	assert.Nil(t, err)
	assert.Equal(t, "i42e", result)

	result, err = BencodeEncode(-5)
	assert.Nil(t, err)
	assert.Equal(t, "i-5e", result)
}
func TestBencodeEncodeString(t *testing.T) {
	result, err := BencodeEncode("hello")
	assert.Nil(t, err)
	assert.Equal(t, "5:hello", result)

	result, err = BencodeEncode("")
	assert.Nil(t, err)
	assert.Equal(t, "0:", result)
}
func TestBencodeEncodeList(t *testing.T) {
	result, err := BencodeEncode([]any{42, "test"})
	assert.Nil(t, err)
	assert.Equal(t, "li42e4:teste", result)

	result, err = BencodeEncode([]any{})
	assert.Nil(t, err)
	assert.Equal(t, "le", result)
}
func TestBencodeEncodeDictionary(t *testing.T) {
	result, err := BencodeEncode(map[string]any{
		"name": "Alice",
		"age":  25,
	})
	assert.Nil(t, err)
	assert.Equal(t, "d3:agei25e4:name5:Alicee", result)
}
