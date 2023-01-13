package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestNextBit(t *testing.T) {
	source := bytes.NewReader([]byte{0xA3, 0xF2})
	stream := &bitstream{
		source: io.ByteReader(source),
	}
	expected := []byte{
		// 3
		0x01,
		0x01,
		0x00,
		0x00,

		// A
		0x00,
		0x01,
		0x00,
		0x01,

		// 2
		0x00,
		0x01,
		0x00,
		0x00,

		// F (only the first bit)
		0x01,
	}
	for i, expBit := range expected {
		assert.Equal(t, expBit, nextBit(stream), fmt.Sprintf("{%d}-th bit", i+1))
	}
}

func TestReadBitsInv(t *testing.T) {
	//   10100011 11110010
	// 1    _____
	// 2 ___            __
	// 3          ______
	source := bytes.NewReader([]byte{0xA3, 0xF2})
	stream := &bitstream{
		source: io.ByteReader(source),
	}
	expected := []struct{
		length int
		bitsValue int
	}{
		{5, 0b00011},
		{5, 0b10101},
		{6, 0b111100},
	}
	for _, tc := range expected {
		assert.Equal(t, tc.bitsValue, readBitsInv(stream, tc.length), tc)
	}
}