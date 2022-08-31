package ck2

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

// var printwidth = 80
// var tabsize = 4

// type CK2Formatter struct {
// 	file *os.File `json:"file"`
// }

func Parse(f *os.File) error {
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

	parser := &CK2Parser{
		Entities: []*Entity{},
		Depth:    0,
	}

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
	fmt.Printf("JSON Print - \n%s\n", string(aJSON))

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

func (parser *CK2Parser) ParseLine(line []byte) []byte {
	if bytes.Contains(line, []byte("=")) {
		kv := bytes.Split(line, []byte(" = "))
		key := kv[0]
		value := kv[1]
		fmt.Println("key:", strconv.Quote(string(key)), "value:", strconv.Quote(string(value)))

		if value[0] == byte('{') {
			// enter into entity scope
			if parser.Scope == nil {
				parser.Entities = append(parser.Entities, &Entity{
					Name:     string(key),
					Elements: []*EntityElement{},
				})
				parser.Scope = parser.Entities[len(parser.Entities)-1]
				parser.PrevScope = parser.Scope
			} else {
				// enter another scope
				fmt.Println("enter into scope of:", strconv.Quote(string(key)))
				parser.Scope.Elements = append(parser.Scope.Elements, &EntityElement{
					Type:  Block,
					Data:  []*EntityElement{},
					Key:   string(key),
					Value: "",
				})
				parser.Scope = parser.Entities[len(parser.Entities)-1]
				fmt.Println("scope:", parser.Scope)
				// last.Elements = append(last.Elements, &EntityElement{
				// 	_type: Property,
				// 	Data:  nil,
				// 	Key:   string(key),
				// 	Value: string(value),
				// })
			}
		} else {
			parser.Scope.Elements = append(parser.Scope.Elements, &EntityElement{
				Type:  Property,
				Data:  nil,
				Key:   string(key),
				Value: string(value),
			})
			fmt.Println("property:", strconv.Quote(string(key)), "value:", strconv.Quote(string(value)))
		}
	}
	if len(line) > 0 && line[0] == byte('}') {
		fmt.Println("END of scope")
		parser.Scope = parser.PrevScope
	}

	fmt.Print("\n")

	return line
}
