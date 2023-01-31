package main

import (
	"fmt"
	"io"
)

type bitstream struct {
	source io.ByteReader
	buf    byte
	mask   byte // current bit position within buf; 8 is MSB
}

// nextBit is little endian (LSB to MSB)
func nextBit(stream *bitstream) byte {
	if stream.mask == 0 { // overflow, means need to get the next byte
		stream.mask = 0x01 // reset mask
		var err error
		stream.buf, err = stream.source.ReadByte()
		if err != nil {
			panic(err)
		}
	}

	var bit byte = 0
	if stream.buf&stream.mask > 0 {
		bit = 1
	}
	stream.mask <<= 1
	return bit
}

// readBitsInv read in little-endian form but interpreted in big-endian form
func readBitsInv(stream *bitstream, count int) (value int) {
	if count > 31 {
		panic("the buffer used (`value`) is of `int` type")
	}
	for i := 0; i < count; i++ {
		bit := nextBit(stream)
		value |= int(bit) << i // set as MSB
	}
	if explanationMode {
		fmt.Printf("-- reading %d bits, value is %x\n", count, value)
	}
	return value
}

func helperBitStringToBytes(bits string) []byte {
	bytes := make([]byte, 0)
	for i, c := range bits {
		if i%8 == 0 {
			bytes = append(bytes, 0x00)
		}
		if string(c) == "1" {
			bytes[len(bytes)-1] |= 1 << (i % 8)
		}
	}
	return bytes
}
