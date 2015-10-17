package hzip

import (
	"fmt"
	"io"
)

const byteSize = 8

type bitWriter struct {
	w      io.Writer
	buf    byte   // bit buffer
	mem    []byte // slice of memory to avoid allocation in calls to w.Write
	shift  int    // number of bit positions to shift for the next write
	closed bool
}

// newBitWriter returns an io.Writer that proxies Write calls to the underlying
// io.Writer, except it has methods that are designed to write individual bits.
// The bits are encoded from left to right, and the final byte is padded on the
// right with zeroes.
func newBitWriter(w io.Writer) *bitWriter {
	return &bitWriter{
		w:     w,
		shift: byteSize - 1,
	}
}

// Write writes the given slice to the underlying io.Writer. If the Write
// method is called while the bit buffer has bits in it, then it will panic.
// Call Flush to write and clear the buffered bits.
func (w *bitWriter) Write(p []byte) (int, error) {
	if w.shift != byteSize-1 {
		panic("hzip: invalid write call - bit buffer not empty")
	}
	return w.w.Write(p)
}

// WriteBit encodes an ASCII '1' or '0' character as a binary bit and writes it to
// the underlying io.Writer. Note that this is buffered so nothing will be
// written to the underlying io.Writer until the total number of bits reaches
// the size of a byte or Flush is called.
//
// Will panic if the bit isn't an ASCII '1' or '0' character.
func (w *bitWriter) WriteBit(bit byte) (int, error) {
	if bit != '0' && bit != '1' {
		panic(fmt.Errorf("hzip: invalid bit %q", bit))
	}
	if w.shift == -1 {
		err := w.Flush()
		if err != nil {
			return 0, err
		}
	}
	w.buf |= (bit - '0') << uint(w.shift)
	w.shift--
	return 1, nil
}

// Calls WriteBit on each byte in the slice.
func (w *bitWriter) WriteBits(p []byte) (int, error) {
	ntotal := 0
	for i := 0; i < len(p); i++ {
		n, err := w.WriteBit(p[i])
		if err != nil {
			return ntotal, err
		}
		ntotal += n
	}
	return ntotal, nil
}

// Calls WriteBit on each byte in the string.
func (w *bitWriter) WriteBitString(s string) (int, error) {
	ntotal := 0
	for i := 0; i < len(s); i++ {
		n, err := w.WriteBit(s[i])
		if err != nil {
			return ntotal, err
		}
		ntotal += n
	}
	return ntotal, nil
}

func (w *bitWriter) Close() error {
	if w.closed {
		return nil
	}
	w.closed = true
	return w.Flush()
}

// Flush ends the current byte, padding it on the right with zeroes, and
// writes it to the underlying io.Writer.
func (w *bitWriter) Flush() error {
	if w.shift != byteSize-1 {
		if w.mem != nil {
			w.mem = w.mem[:0]
		}
		w.mem = append(w.mem, w.buf)
		_, err := w.w.Write(w.mem)
		if err != nil {
			return err
		}
		w.shift = byteSize - 1
		w.buf = 0x00
	}
	return nil
}

type bitReader struct {
	r    io.Reader
	buf  byte   // bit buffer
	mem  []byte // slice of memory to avoid allocation
	mask byte   // current bit mask
}

// newBitReader returns an io.Reader that proxies Read calls to the underlying
// io.Reader, except that it has methods that are designed to read individual
// bits. Data is buffered a byte at a time.
func newBitReader(r io.Reader) *bitReader {
	return &bitReader{
		r:   r,
		mem: make([]byte, 1),
	}
}

// Read reads len(p) bytes from the underlying io.Reader. If the Read method is
// called while the bit buffer has bits in it, then it will panic. Call Reset
// to clear the bit buffer.
func (r *bitReader) Read(p []byte) (int, error) {
	if r.mask != 0x00 {
		panic("hzip: invalid read call - bit buffer not empty")
	}
	return r.r.Read(p)
}

// ReadBit reads a single bit and returns it as an ASCII '0' or '1'.
func (r *bitReader) ReadBit() (byte, error) {
	if r.mask == 0x00 {
		r.mask = 0x80
		_, err := r.r.Read(r.mem)
		if err != nil {
			return 0x00, err
		}
		r.buf = r.mem[0]
	}
	ret := byte('0')
	if r.buf&r.mask != 0 {
		ret = '1'
	}
	r.mask >>= 1
	return ret, nil
}

// Reset clears the bit buffer.
func (r *bitReader) Reset() {
	r.buf = 0x00
	r.mask = 0x00
}
