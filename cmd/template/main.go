package main

import (
	ck2parser "ck2-parser/internal/app/ck2parser"
	"encoding/json"
	"log"
	"os"
)

func main() {
	// start := time.Now()

	file, err := os.Open("data/simple.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	parser, err := ck2parser.New(file)
	if err != nil {
		log.Fatal(err)
	}

	// parsing
	p, err := parser.Parse()
	if err != nil {
		panic(err)
	}
	Save(p, "parsed.json")

	// linting
	linter := ck2parser.NewLinter(p.Filepath, p.Data)
	linter.Lint()

	w, err := os.Create("tmp/linted.txt")
	if err != nil {
		panic(err)
	}
	w.Write(linter.LintedData())

	// elapsed := time.Since(start)
	// log.Println("Saving took", elapsed.Seconds())
}

func Save(p any, filename string) error {
	w, err := os.Create("tmp/" + filename)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	enc.Encode(p)
	return nil
}
