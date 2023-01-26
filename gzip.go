package main

import (
	"bufio"
	"encoding/binary"
	"io"
)

type GzipHeader struct {
	ID                [2]byte
	CompressionMethod byte
	Flags             byte
	Mtime             [4]byte
	ExtraFlags        byte
	OS                byte
}

type GzipMetaData struct {
	Header GzipHeader
	Xlen int16
	Extra []byte
	Fname []byte
	Fcomment []byte
	Crc16 uint16
}

const FTEXT byte = 0x01
const FHCRC byte = 0x02
const FEXTRA byte = 0x04
const FNAME byte = 0x08
const FCOMMENT byte = 0x10

func readCString(file io.Reader) (buf []byte) {
	var nextChar byte
	if err := binary.Read(file, binary.LittleEndian, &nextChar); err != nil {
		panic(err)
	}
	for nextChar != 0x00 {
		buf = append(buf, nextChar)
		if err := binary.Read(file, binary.LittleEndian, &nextChar); err != nil {
			panic(err)
		}
	}
	return
}

func readGzipMetaData(file io.Reader) GzipMetaData {
	gzipMetaData := GzipMetaData{}
	if err := binary.Read(file, binary.LittleEndian, &gzipMetaData.Header); err != nil {
		panic(err)
	}
	if !(gzipMetaData.Header.ID[0] == 0x1f && gzipMetaData.Header.ID[1] == 0x8b) {
		panic("Not a gzipMetaData format")
	}
	if gzipMetaData.Header.CompressionMethod != 8 {
		panic("Unrecognized compression method")
	}
	if (gzipMetaData.Header.Flags & FEXTRA) != 0 {
		if err := binary.Read(file, binary.LittleEndian, &gzipMetaData.Xlen); err != nil {
			panic(err)
		}
		gzipMetaData.Extra = make([]byte, gzipMetaData.Xlen)
		if err := binary.Read(file, binary.LittleEndian, &gzipMetaData.Extra); err != nil {
			panic(err)
		}
		// for now we just ignore the extra data
	}
	if (gzipMetaData.Header.Flags & FNAME) != 0 {
		gzipMetaData.Fname = readCString(file)
	}
	if (gzipMetaData.Header.Flags & FCOMMENT) != 0 {
		gzipMetaData.Fcomment = readCString(file)
	}
	if (gzipMetaData.Header.Flags & FHCRC) != 0 {
		if err := binary.Read(file, binary.LittleEndian, &gzipMetaData.Crc16); err != nil {
			panic(err)
		}
	}
	return gzipMetaData
}

func gzipInflate(file io.Reader) []byte {
	var lastBlock byte
	stream := &bitstream{source: file.(io.ByteReader)}
	var out []byte
	for lastBlock == 0 {
		lastBlock = nextBit(stream)
		blockFormat := readBitsInv(stream, 2)
		switch blockFormat {
		case 0b00:
			panic("uncompressed block type not supported")
		case 0b01:
			literalsRoot := readFixedHuffmanTree(stream)
			out = append(out, inflateHuffmanCodes(stream, literalsRoot, nil)...)
		case 0b10:
			literalsRoot, distancesRoot := readDynamicHuffmanTree(stream)
			out = append(out, inflateHuffmanCodes(stream, literalsRoot, distancesRoot)...)
		default:
			panic("unsupported block type")
		}
	}
	return out
}

func readGzipFile(file io.Reader) []byte {
	_ = readGzipMetaData(file)
	return gzipInflate(bufio.NewReader(file))
}
