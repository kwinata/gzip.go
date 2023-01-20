package main

import "os"

// TODO create cli interface
func main() {
	file, err := os.Open("attachment/gunzip.c.gz")
	if err != nil {
		panic(err)
	}
	readGzipFile(file)
}
