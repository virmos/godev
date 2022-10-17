package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type logWritter struct{}

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	lw := logWritter{}

	io.Copy(lw, file)
}

func (logWritter) Write(bs []byte) (int, error) {
	fmt.Println(string(bs))

	return len(bs), nil
}
