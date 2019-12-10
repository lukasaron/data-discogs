package main

import (
	"github.com/Twyer/discogs/cmd"
	"log"
)

func main() {
	err := cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
