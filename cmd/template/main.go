package main

import (
	"ck2-parser/internal/app/ck2"
	"log"
	"os"
)

func main() {

	file, err := os.Open("data/elementary.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	err = ck2.Parse(file)
	if err != nil {
		log.Fatal(err)
	}
}
