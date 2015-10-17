package hzip

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteBit(t *testing.T) {
	var (
		n   int
		err error
	)
	buf := new(bytes.Buffer)
	w := newBitWriter(buf)

	assert.Panics(t, func() {
		w.WriteBit('2')
	}, fmt.Sprintf("Calling WriteBit('2') should panic"))

	n, err = w.WriteBit('0')
	assert.Equal(t, 1, n)
	assert.Nil(t, err)
	assert.Empty(t, buf.Bytes())
	w.WriteBit('1')
	w.WriteBit('1')
	w.WriteBit('0')
	w.WriteBit('0')
	w.WriteBit('1')
	w.WriteBit('0')
	w.WriteBit('1')
	w.WriteBit('0')
	w.WriteBit('1')
	w.WriteBit('0')
	w.Flush()
	assert.Equal(t, []byte{0x65, 0x40}, buf.Bytes())
}

func TestWriteBits(t *testing.T) {
	var (
		n   int
		err error
	)
	buf := new(bytes.Buffer)
	w := newBitWriter(buf)

	n, err = w.WriteBits([]byte(""))
	assert.Equal(t, 0, n)
	assert.Nil(t, err)
	assert.Empty(t, buf.Bytes())

	w.Flush()
	buf.Reset()
	n, err = w.WriteBits([]byte("01100101010"))
	assert.Equal(t, 11, n)
	assert.Nil(t, err)
	w.Flush()
	assert.Equal(t, []byte{0x65, 0x40}, buf.Bytes())
}

func TestWriteBitString(t *testing.T) {
	var (
		n   int
		err error
	)
	buf := new(bytes.Buffer)
	w := newBitWriter(buf)

	n, err = w.WriteBitString("")
	assert.Equal(t, 0, n)
	assert.Nil(t, err)
	assert.Empty(t, buf.Bytes())

	n, err = w.WriteBitString("01100101010")
	assert.Equal(t, 11, n)
	assert.Nil(t, err)
	w.Flush()
	assert.Equal(t, []byte{0x65, 0x40}, buf.Bytes())
}

func TestBitFlush(t *testing.T) {
	var err error
	buf := new(bytes.Buffer)
	w := newBitWriter(buf)

	// Empty flush
	err = w.Flush()
	assert.Nil(t, err)
	assert.Empty(t, buf.Bytes())

	// Flush with one 0 bit
	buf.Reset()
	w.WriteBitString("0")
	assert.Empty(t, buf.Bytes())
	err = w.Flush()
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x00}, buf.Bytes())

	// Flush with 2 bytes
	buf.Reset()
	w.WriteBitString("01100101010")
	assert.Equal(t, []byte{0x65}, buf.Bytes())
	err = w.Flush()
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x65, 0x40}, buf.Bytes())
	// Test double flush
	err = w.Flush()
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x65, 0x40}, buf.Bytes())
}

func TestCloseBitWriter(t *testing.T) {
	var err error
	buf := new(bytes.Buffer)
	w := newBitWriter(buf)

	// Empty close
	w.Close()
	assert.Empty(t, buf.Bytes())

	w = newBitWriter(buf)
	w.WriteBitString("01100101010")
	assert.Equal(t, []byte{0x65}, buf.Bytes())
	err = w.Close()
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x65, 0x40}, buf.Bytes())
	// Test double close
	err = w.Close()
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x65, 0x40}, buf.Bytes())
}

func TestInterleavedBitWrites(t *testing.T) {
	buf := new(bytes.Buffer)
	w := newBitWriter(buf)

	w.Write([]byte{0x01, 0x02})
	assert.Equal(t, []byte{0x01, 0x02}, buf.Bytes())

	w.WriteBitString("01100101010")
	assert.Equal(t, []byte{0x01, 0x02, 0x65}, buf.Bytes())

	assert.Panics(t, func() {
		w.Write([]byte{0x03})
	}, "Calling Write with non-empty buffer should panic")

	w.Flush()
	w.Write([]byte{0x03})
	assert.Equal(t, []byte{0x01, 0x02, 0x65, 0x40, 0x03}, buf.Bytes())
}
