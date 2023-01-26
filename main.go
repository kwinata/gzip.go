package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	slowPrintMode = false
	fileName      string
)

// TODO create cli interface
func main() {
	flag.StringVar(&fileName, "f", "", "-f [path to file name]")
	flag.BoolVar(&slowPrintMode, "s", false, "-s to enable slow print mode")
	flag.Parse()

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	out := readGzipFile(file)
	if !slowPrintMode {
		fmt.Print(string(out))
	}
}
