package main

import "os"

var (
	debugMode = false
)

// TODO create cli interface
func main() {
	if os.Getenv("DEBUG") == "true" {
		debugMode = true
	}
	file, err := os.Open("attachment/genesis.txt.gz")
	if err != nil {
		panic(err)
	}
	readGzipFile(file)
}
