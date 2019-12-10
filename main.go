package main

import (
	"github.com/Twyer/discogs/cmd"
	"log"
)

func main() {
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
