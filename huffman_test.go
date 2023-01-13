package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildHuffmanTree(t *testing.T) {
	hRanges := []huffmanRange{
		{1, 4},
		{4, 6},
		{6, 4},
		{14, 5},
		{18, 6},
		{21, 4},
		{26, 6},
	}
	expected := []string{
		"0000",
		"0001",
		"101100",
		"101101",
		"101110",
		"0010",
		"0011",
		"01110",
		"01111",
		"10000",
		"10001",
		"10010",
		"10011",
		"10100",
		"10101",
		"101111",
		"110000",
		"110001",
		"110010",
		"0100",
		"0101",
		"0110",
		"110011",
		"110100",
		"110101",
		"110110",
		"110111",
	}
	root := buildHuffmanTree(hRanges)
	codeTable := make([]string, hRanges[len(hRanges)-1].end+1) // symbol 0 is valid
	traverseHuffmanTree(root, "", codeTable)
	for i, v := range codeTable {
		assert.Equal(t, expected[i], v)
	}
}

func traverseHuffmanTree(node *huffmanNode, prefix string, codeTable []string) {
	if node.code != -1 {
		codeTable[node.code] = prefix
	}
	if node.one != nil {
		traverseHuffmanTree(node.one, prefix + "1", codeTable)
	}
	if node.zero != nil {
		traverseHuffmanTree(node.zero, prefix + "0", codeTable)
	}
}
