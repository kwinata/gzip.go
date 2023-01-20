package main

import (
	"encoding/binary"
	"fmt"
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

type GzipFile struct {
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

func readGzipMetaData(file io.Reader) GzipFile {
	gzip := GzipFile{}
	if err := binary.Read(file, binary.LittleEndian, &gzip.Header); err != nil {
		panic(err)
	}
	fmt.Printf("%v+\n", gzip)
	if !(gzip.Header.ID[0] == 0x1f && gzip.Header.ID[1] == 0x8b) {
		panic("Not a gzip format")
	}
	if gzip.Header.CompressionMethod != 8 {
		panic("Unrecognized compression method")
	}
	if (gzip.Header.Flags & FEXTRA) != 0 {
		if err := binary.Read(file, binary.LittleEndian, &gzip.Xlen); err != nil {
			panic(err)
		}
		gzip.Extra = make([]byte, gzip.Xlen)
		if err := binary.Read(file, binary.LittleEndian, &gzip.Extra); err != nil {
			panic(err)
		}
		// for now we just ignore the extra data
	}
	if (gzip.Header.Flags & FNAME) != 0 {
		gzip.Fname = readCString(file)
	}
	if (gzip.Header.Flags & FCOMMENT) != 0 {
		gzip.Fcomment = readCString(file)
	}
	if (gzip.Header.Flags & FHCRC) != 0 {
		if err := binary.Read(file, binary.LittleEndian, &gzip.Crc16); err != nil {
			panic(err)
		}
	}
	return gzip
}