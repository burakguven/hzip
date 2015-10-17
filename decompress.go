package hzip

import (
	"encoding/binary"
	"io"
)

type Reader struct {
	r        *bitReader
	nRead    uint32          // number of symbols read
	fileSize uint32          // size of the decompressed file
	symbols  map[string]byte // map from code to symbol
	mem      []byte          // small slice of memory to avoid memory allocation in calls to Read
}

// NewReader returns an io.Reader that reads from the given io.Reader and
// decompresses it using the Huffman coding algorithm.
func NewReader(r io.Reader) (*Reader, error) {
	hr := &Reader{
		r:   newBitReader(r),
		mem: make([]byte, 1),
	}
	err := hr.readHeader()
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return hr, err
}

func (r *Reader) Read(p []byte) (int, error) {
	n := 0
	for i := 0; i < len(p); i++ {
		if r.nRead == r.fileSize {
			// Slurp up any remaining padding bytes
			var err error
			for err != io.EOF {
				_, err = r.r.ReadBit()
			}
			return n, io.EOF
		}
		// Read one bit at a time until there's a code match
		code := ""
		for {
			if _, ok := r.symbols[code]; ok {
				break
			}
			bit, err := r.r.ReadBit()
			if err == io.EOF {
				return n, io.ErrUnexpectedEOF
			}
			code += string(bit)
		}
		p[i] = r.symbols[code]
		r.nRead++
		n++
	}
	return n, nil
}

// See compress.go for documentation about the header format.
func (r *Reader) readHeader() error {
	// File size
	if err := binary.Read(r.r, binary.LittleEndian, &r.fileSize); err != nil {
		return err
	}

	// Alphabet size
	var alphabetSize uint32
	if err := binary.Read(r.r, binary.LittleEndian, &alphabetSize); err != nil {
		return err
	}
	r.symbols = make(map[string]byte)
	for i := uint32(0); i < alphabetSize; i++ {
		// Symbol
		if _, err := r.r.Read(r.mem[:1]); err != nil {
			return err
		}
		symbol := r.mem[0]

		// Number of bits in code
		if _, err := r.r.Read(r.mem[:1]); err != nil {
			return err
		}
		codeLen := r.mem[0]

		// The code itself
		var codeBits []byte
		for k := byte(0); k < codeLen; k++ {
			bit, err := r.r.ReadBit()
			if err != nil {
				return err
			}
			codeBits = append(codeBits, bit)
		}
		r.symbols[string(codeBits)] = symbol
		// Get rid of extra padding at the end, if any
		r.r.Reset()
	}
	return nil
}
