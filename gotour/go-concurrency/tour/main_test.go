package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestUpdateMessage(t *testing.T) {
	wg.Add(1)
	updateMessage("Hello, universe!")
	wg.Wait()
	if msg != "Hello, universe!" {
		t.Errorf("Expected %s. Get %v", "Hello, universe!", msg)
	}
}

func TestPrintMessage(t *testing.T) {
	stdOut := os.Stdout
	r,w,_ := os.Pipe()
	os.Stdout = w

	wg.Add(1)
	printMessage()
	wg.Wait()

	_ = w.Close()
	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if (!strings.Contains(output, "universe!")) {
		t.Error("Expected to find universe, but it's not there!")
	}
}

func TestMain(t *testing.T) {
	stdOut := os.Stdout
	r,w,_ := os.Pipe()
	os.Stdout = w

	main()

	_ = w.Close()
	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if (!strings.Contains(output, "universe!")) {
		t.Error("Expected to find universe, but it's not there!")
	}

	if (!strings.Contains(output, "cosmos!")) {
		t.Error("Expected to find universe, but it's not there!")
	}

	if (!strings.Contains(output, "world!")) {
		t.Error("Expected to find universe, but it's not there!")
	}
}
/*
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
*/