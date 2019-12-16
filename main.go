package main

import (
	"github.com/Twyer/discogs-parser/cmd"
	"log"
)

func main() {
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
