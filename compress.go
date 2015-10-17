// Package hzip implements reading and writing of files compressed with the
// Huffman coding algorithm.
package hzip

import (
	"encoding/binary"
	"io"
)

type Writer struct {
	w      *bitWriter
	buf    []byte
	freqs  map[byte]int
	codes  map[byte]string
	closed bool
}

// NewWriter returns an io.Writer that compresses the data written to it using
// the Huffman coding algorithm and writes it to the given io.Writer.
//
// No data is written to the underlying io.Writer until Close is called.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:     newBitWriter(w),
		freqs: make(map[byte]int),
	}
}

func (w *Writer) Write(p []byte) (int, error) {
	for _, b := range p {
		w.freqs[b]++
		w.buf = append(w.buf, b)
	}
	return len(p), nil
}

// Close writes the compressed data to the underlying io.Writer and closes the
// Writer. It does not close the underlying io.Writer.  Nothing is written to
// the underlying io.Writer until Close is called because of the nature of the
// algorithm.
func (w *Writer) Close() error {
	if w.closed {
		return nil
	}
	w.codes = buildCodeMap(w.freqs)
	if err := w.writeHeader(); err != nil {
		return err
	}
	if _, err := w.writeData(); err != nil {
		return err
	}
	w.w.Flush()
	w.closed = true
	return nil
}

// File Format: Header followed by compressed data
// Header:
//	- 4 bytes (uint32): the number of bytes in the original file
//	- 4 bytes (uint32): the size of the alphabet
//	- For each symbol in the alphabet (sorted by symbol value):
//		- 1 byte: the symbol itself
//		- 1 byte: the number of bits in its codeword
//		- 0 or more bytes: the codeword, padding to the right with 0 bits
// Compressed Data:
//	- 1 or more bytes: raw bytes padded to the right with 0 bits
// All multi-byte values are in little endian.

func (w *Writer) writeHeader() error {
	// The number of bytes in the original file
	if err := binary.Write(w.w, binary.LittleEndian, uint32(len(w.buf))); err != nil {
		return err
	}
	// The size of the alphabet
	if err := binary.Write(w.w, binary.LittleEndian, uint32(len(w.freqs))); err != nil {
		return err
	}
	for i := 0; i <= 0xff; i++ {
		code, ok := w.codes[byte(i)]
		if !ok {
			continue
		}
		// The symbol itself
		if err := binary.Write(w.w, binary.LittleEndian, byte(i)); err != nil {
			return err
		}
		// The number of bits in its codeword
		if err := binary.Write(w.w, binary.LittleEndian, byte(len(code))); err != nil {
			return err
		}
		// The codeword, padded on the right with 0 bits. Note that
		// it's legal to have an empty codeword, but it only happens
		// when the alphabet has a single symbol.
		if _, err := w.w.WriteBitString(code); err != nil {
			return err
		}
		if err := w.w.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeData() (int, error) {
	ntotal := 0
	for i := range w.buf {
		_, err := w.w.WriteBitString(w.codes[w.buf[i]])
		if err != nil {
			return ntotal, err
		}
		ntotal++
	}
	return ntotal, nil
}
