package hzip

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func stripExtension(filename string) string {
	sep := strings.LastIndexByte(filename, '.')
	if sep == -1 {
		return filename
	}
	return filename[:sep]
}

func TestDecompressDataFiles(t *testing.T) {
	matches, err := filepath.Glob("testdata/*.hz")
	if err != nil {
		t.Fatal(err)
	}
	for _, hzName := range matches {
		t.Log(hzName)
		hzFile, err := os.Open(hzName)
		if err != nil {
			t.Fatal(err)
		}
		hzReader, err := NewReader(hzFile)
		if err != nil {
			t.Fatal(err)
		}
		rawData, err := ioutil.ReadFile(stripExtension(hzName))
		if err != nil {
			t.Fatal(err)
		}
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, hzReader); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, rawData, buf.Bytes())
		hzFile.Close()
	}
}

func tryDecompress(t *testing.T, data []byte) error {
	reader, err := NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, reader); err != nil {
		return err
	}
	return nil
}

func TestMalformedHeader(t *testing.T) {
	var err error
	data, err := ioutil.ReadFile("testdata/hello.hz")
	if err != nil {
		t.Fatal(err)
	}
	var mangled []byte
	// Say the original file size is larger than it actually is
	mangled = append([]byte(nil), data...)
	mangled[0] = 0xff
	assert.Equal(t, io.ErrUnexpectedEOF, tryDecompress(t, mangled))

	// Say the alphabet size is larger than it actually is
	mangled = append([]byte(nil), data...)
	mangled[4] = 0xff
	assert.Equal(t, io.ErrUnexpectedEOF, tryDecompress(t, mangled))

	// Say the alphabet size is smaller than it actually is
	mangled = append([]byte(nil), data...)
	mangled[4] = 0x00
	assert.Equal(t, io.ErrUnexpectedEOF, tryDecompress(t, mangled))

	// Empty file
	mangled = nil
	assert.Equal(t, io.ErrUnexpectedEOF, tryDecompress(t, mangled))

	// Has length but no alphabet size
	mangled = append([]byte(nil), 0xa, 0x00, 0x00, 0x00)
	assert.Equal(t, io.ErrUnexpectedEOF, tryDecompress(t, mangled))
}
