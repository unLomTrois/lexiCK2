package ck2

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

var printwidth = 80
var tabsize = 4

type CK2Formatter struct {
	file *os.File
}

func Parse(f *os.File) (*CK2Formatter, error) {
	buf := &bytes.Buffer{}
	tee := io.TeeReader(f, buf)

	b, _ := io.ReadAll(tee)
	b = bytes.Replace(b, []byte("={"), []byte("= {"), 1)
	fmt.Println(string(b))

	// * init scanner
	scanner := bufio.NewScanner(buf)
	scanner.Split(ScanLinesLF)

	// * start scanning
	for scanner.Scan() { // internally, it advances token based on sperator
		bytes := scanner.Bytes()
		text := string(bytes[:len(bytes)-1]) // trim whitespace

		fmt.Println(strconv.Quote(text)) // token in unicode-char
		fmt.Println(bytes)               // token in bytes
	}

	return &CK2Formatter{file: f}, nil
}
