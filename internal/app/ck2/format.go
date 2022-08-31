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

	raw, _ := io.ReadAll(tee)
	fmt.Println(string(raw))

	// * init scanner
	scanner := bufio.NewScanner(buf)
	scanner.Split(ScanLinesLF)

	new_file, err := os.Create("./tmp/newfile.txt")
	if err != nil {
		return nil, err
	}
	defer new_file.Close()

	w := bufio.NewWriter(new_file)

	// * start scanning
	linenumber := 1
	for scanner.Scan() { // internally, it advances token based on sperator
		text_bytes := scanner.Bytes()
		text_bytes = FormatLine(text_bytes)

		text := string(text_bytes)
		// text := string(text_bytes[:len(text_bytes)-1]) // trim whitespace

		// text = strings.TrimSpace(text)

		fmt.Println(linenumber, strconv.Quote(text)) // token in unicode-char
		fmt.Println(linenumber, text_bytes)          // token in text_bytes
		linenumber += 1

		n4, err := w.Write(append(text_bytes, '\n'))
		if err != nil {
			return nil, err
		}
		fmt.Printf("wrote %d bytes\n", n4)
	}
	w.Flush()

	return &CK2Formatter{file: f}, nil
}

func FormatLine(line []byte) []byte {
	line = bytes.TrimSpace(line)
	line = bytes.ReplaceAll(line, []byte(" = "), []byte("="))
	line = bytes.ReplaceAll(line, []byte("= {"), []byte("={"))
	line = bytes.ReplaceAll(line, []byte("="), []byte(" = "))

	// data
	if bytes.Contains(line, []byte("=")) {
		kv := bytes.Split(line, []byte(" = "))
		key := kv[0]
		value := kv[1]
		fmt.Println("key:", strconv.Quote(string(key)), "value:", strconv.Quote(string(value)))
	}
	// res := make(map[string]string)

	return line
}
