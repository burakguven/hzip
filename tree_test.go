package hzip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTree(t *testing.T) {
	var root *node
	root = buildTree(nil)
	assert.Nil(t, root)

	root = buildTree(map[byte]int{0x12: 5})
	assert.Equal(t, byte(0x12), root.val)
	assert.Equal(t, 5, root.freq)
	assert.Nil(t, root.left)
	assert.Nil(t, root.right)
}

func TestBuildCodeMapBase(t *testing.T) {
	var codes map[byte]string
	codes = buildCodeMap(nil)
	assert.Empty(t, codes)

	codes = buildCodeMap(map[byte]int{0x12: 5})
	assert.Equal(t, 1, len(codes))
	assert.Equal(t, "", codes[0x12])
}

func TestBuildCodeMap(t *testing.T) {
	freqs := map[byte]int{
		'A': 6, 'B': 4, 'C': 5, 'G': 1, 'H': 2,
	}
	wantCodes := map[byte]string{
		'A': "11", 'B': "01", 'C': "10", 'G': "000", 'H': "001",
	}
	codes := buildCodeMap(freqs)
	if len(codes) != len(wantCodes) {
		t.Errorf("len(codes) == %d; want %d", len(codes), len(wantCodes))
	}
	for sym := range codes {
		if codes[sym] != wantCodes[sym] {
			t.Errorf("codes[%q] == %q; want %q", sym, codes[sym], wantCodes[sym])
		}
	}
}
