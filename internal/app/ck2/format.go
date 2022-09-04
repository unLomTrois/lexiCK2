package ck2

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// var printwidth = 80
// var tabsize = 4

// type CK2Formatter struct {
// 	file *os.File `json:"file"`
// }

func Parse(f *os.File) error {
	file_path, err := filepath.Abs(f.Name())
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	tee := io.TeeReader(f, buf)

	raw, _ := io.ReadAll(tee)
	fmt.Println(string(raw))

	// * init scanner
	scanner := bufio.NewScanner(buf)
	scanner.Split(ScanLinesLF)

	new_file, err := os.Create("./tmp/newfile.txt")
	if err != nil {
		return err
	}
	defer new_file.Close()

	w := bufio.NewWriter(new_file)

	// * start scanning
	linenumber := 1
	// depth := 0

	parser := NewParser(file_path)

	for scanner.Scan() { // internally, it advances token based on sperator
		text_bytes := scanner.Bytes()
		text_bytes = FormatLine(text_bytes)
		parser.ParseLine(text_bytes)

		// text := string(text_bytes)
		// text := string(text_bytes[:len(text_bytes)-1]) // trim whitespace

		// fmt.Println(linenumber, strconv.Quote(text)) // token in unicode-char
		// fmt.Println(linenumber, text_bytes)          // token in text_bytes
		linenumber += 1

		_, err := w.Write(append(text_bytes, '\n'))
		if err != nil {
			return err
		}

		// fmt.Printf("wrote %d bytes\n", n4)
	}
	aJSON, _ := json.MarshalIndent(parser, "", "  ")
	// fmt.Printf("JSON Print - \n%s\n", string(aJSON))
	_ = os.WriteFile("tmp/simple.json", aJSON, 0644)
	w.Flush()

	return nil
}

func FormatLine(line []byte) []byte {
	line = bytes.TrimSpace(line)
	line = bytes.ReplaceAll(line, []byte(" = "), []byte("="))
	line = bytes.ReplaceAll(line, []byte("= {"), []byte("={"))
	line = bytes.ReplaceAll(line, []byte("="), []byte(" = "))
	return line
}
