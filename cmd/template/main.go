package main

import (
	ck2parser "ck2-parser/internal/app/ck2parser"
	"log"
	"os"
)

func main() {

	file, err := os.Open("data/text.txt")
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

	err = parser.Parse()
	if err != nil {
		panic(err)
		// log.Fatal(err)
	}
}
