package main

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadGzipHeader(t *testing.T) {
	expectedGzipHeader := GzipHeader{
		ID:                [2]byte{0x01, 0x02},
		CompressionMethod: 0x03,
		Flags:             0x04,
		Mtime:             [4]byte{0x05, 0x06, 0x07, 0x08},
		ExtraFlags:        0x09,
		OS:                0x0A,
	}
	file := bytes.NewReader([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A})
	gzipHeader := GzipHeader{}
	assert.NoError(t, binary.Read(file, binary.LittleEndian, &gzipHeader))
	assert.Equal(t, expectedGzipHeader, gzipHeader)
}

func TestReadGzipMetaData(t *testing.T) {
	expectedGzipMetaData := GzipMetaData{
		Header: GzipHeader{
			ID:                [2]byte{0x1F, 0x8B},
			CompressionMethod: 0x08,
			Flags:             0x1F, // enable everything
			Mtime:             [4]byte{0x05, 0x06, 0x07, 0x08},
			ExtraFlags:        0x09,
			OS:                0x0A,
		},
		Xlen:     0x0008, // 8 bytes of extra data
		Extra:    []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
		Fname:    []byte{0x61, 0x62, 0x63},
		Fcomment: []byte{0x65, 0x66, 0x67},
		Crc16:    uint16(0xABCD),
	}
	metadata := []byte{
		0x1F, 0x8B, 0x08, 0x1F, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, // header
		0x08, 0x00, // xlen uint16
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, // extra data (xlen bytes)
		0x61, 0x62, 0x63, 0x00, // Fname (c-string) 'abc\0'
		0x65, 0x66, 0x67, 0x00, // Fcomment (c-string) 'efg\0'
		0xCD, 0xAB, // CRC16 (uint16)
	}
	gzipFile := readGzipMetaData(bytes.NewReader(metadata))
	assert.Equal(t, expectedGzipMetaData.Header, gzipFile.Header)
	assert.Equal(t, expectedGzipMetaData.Xlen, gzipFile.Xlen)
	assert.Equal(t, expectedGzipMetaData.Extra, gzipFile.Extra)
	assert.Equal(t, expectedGzipMetaData.Fname, gzipFile.Fname)
	assert.Equal(t, expectedGzipMetaData.Fcomment, gzipFile.Fcomment)
	assert.Equal(t, expectedGzipMetaData.Crc16, gzipFile.Crc16)
}
