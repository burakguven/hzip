package hzip

import "container/heap"

type node struct {
	val         byte
	freq        int
	left, right *node
}

// buildTree builds a Huffman coding tree based on the given symbol frequencies
// and returns a pointer to the root node.
func buildTree(freqs map[byte]int) *node {
	if len(freqs) == 0 {
		return nil
	}
	nodes := new(nodeHeap)
	// Start with a forest of nodes, each node being one symbol in the alphabet
	for val, freq := range freqs {
		heap.Push(nodes, node{val: val, freq: freq})
	}
	// Combine the two trees with the lowest frequency until only one tree left
	for len(*nodes) > 1 {
		a := heap.Pop(nodes).(node)
		b := heap.Pop(nodes).(node)
		parent := node{
			freq:  a.freq + b.freq,
			left:  &a,
			right: &b,
		}
		heap.Push(nodes, parent)
	}
	root := heap.Pop(nodes).(node)
	return &root
}

// buildCodeMap builds a Huffman coding map, mapping from the input symbol to
// the symbol encoding, which is represented as a string of ASCII '1' and '0'
// characters.
func buildCodeMap(freqs map[byte]int) map[byte]string {
	root := buildTree(freqs)
	codes := make(map[byte]string)
	buildCodeMapRec(root, "", &codes)
	return codes
}

func buildCodeMapRec(n *node, code string, codes *map[byte]string) {
	if n == nil {
		return
	}
	if n.left == nil && n.right == nil {
		(*codes)[n.val] = code
	} else {
		buildCodeMapRec(n.left, code+"0", codes)
		buildCodeMapRec(n.right, code+"1", codes)
	}
}

type nodeHeap []node

// Implement heap.Interface

func (h nodeHeap) Len() int           { return len(h) }
func (h nodeHeap) Less(i, j int) bool { return h[i].freq < h[j].freq }
func (h nodeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *nodeHeap) Push(x interface{}) {
	*h = append(*h, x.(node))
}
func (h *nodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
