package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	slowPrintMode   = false
	backPointerMode = false
	explanationMode = false
	fileName        string
)

var literalCount, backPointerCount, totalBytes int

func main() {
	literalCount, backPointerCount, totalBytes = 0, 0, 0

	flag.StringVar(&fileName, "f", "", "-f [path to file name]")
	flag.BoolVar(&slowPrintMode, "s", false, "-s to enable slow print mode")
	flag.BoolVar(&explanationMode, "e", false, "-e to enable explanation")
	flag.BoolVar(&backPointerMode, "bp", false, "-bp to enable back pointer (only effective in slow print mode")
	flag.Parse()

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	_ = readGzipFile(file)

	fmt.Printf("\n\nSummary Report: literalCount %d, backPointerCount %d, totalBytes %d\n", literalCount, backPointerCount, totalBytes)
}
