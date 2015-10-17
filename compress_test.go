package hzip

import (
	"bytes"
	"io"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyCompress(t *testing.T) {
	buf := new(bytes.Buffer)
	w := NewWriter(buf)
	w.Close()
	bytes := buf.Bytes()
	// Length of original file
	assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, bytes[:4])
	// Size of alphabet
	assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, bytes[4:8])
	assert.Empty(t, bytes[8:])
}

func TestSingleSymbolCompress(t *testing.T) {
	buf := new(bytes.Buffer)
	w := NewWriter(buf)
	io.WriteString(w, "XXXXXXXXXX")
	w.Close()
	bytes := buf.Bytes()
	// Length of original file
	assert.Equal(t, []byte{0x0a, 0x00, 0x00, 0x00}, bytes[:4])
	// Size of alphabet
	assert.Equal(t, []byte{0x01, 0x00, 0x00, 0x00}, bytes[4:8])
	// The symbol, number of encoded bits, bits
	// Note that a single-symbol alphabet doesn't require any data to encode
	assert.Equal(t, []byte{'X', 0x00}, bytes[8:])
}

func TestSimpleCompress(t *testing.T) {
	buf := new(bytes.Buffer)
	w := NewWriter(buf)
	io.WriteString(w, "Hello World")
	w.Close()
	bytes := buf.Bytes()
	// Length of original file
	assert.Equal(t, []byte{0x0b, 0x00, 0x00, 0x00}, bytes[:4])
	freqs := map[byte]int{
		'H': 1,
		'e': 1,
		'l': 3,
		'o': 2,
		' ': 1,
		'W': 1,
		'r': 1,
		'd': 1,
	}
	// Size of alphabet
	assert.Equal(t, []byte{0x08, 0x00, 0x00, 0x00}, bytes[4:8])

	// Since there are symbols with the same frequency, the code is going
	// to be ambiguous. Just test that the symbols are there
	var sortedAlphabet []byte
	for i := 0x00; i <= 0xff; i++ {
		if _, ok := freqs[byte(i)]; !ok {
			continue
		}
		sortedAlphabet = append(sortedAlphabet, byte(i))
	}
	for i := 0; i < len(sortedAlphabet); i++ {
		assert.Equal(t, []byte{sortedAlphabet[i]}, bytes[8+3*i:8+3*i+1])
	}
}

func genRandBytes(length int) []byte {
	arr := make([]byte, length)
	for i := 0; i < length; i++ {
		arr[i] = byte(rand.Intn(256))
	}
	return arr
}

func TestCompressRandom(t *testing.T) {
	for i := 0; i < 1000; i++ {
		randBytes := genRandBytes(i)
		compressBuf := new(bytes.Buffer)
		writer := NewWriter(compressBuf)
		if _, err := writer.Write(randBytes); err != nil {
			t.Fatal(err)
		}
		writer.Close()

		reader, err := NewReader(bytes.NewReader(compressBuf.Bytes()))
		if err != nil {
			t.Fatal(err)
		}
		decompressBuf := new(bytes.Buffer)
		if _, err := io.Copy(decompressBuf, reader); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, randBytes, decompressBuf.Bytes())
	}
}
