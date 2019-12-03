package main

import (
	"github.com/Twyer/discogs/decoder"
	"log"
	"os"
	"strconv"
)

func main() {
	// fmt.Println(util.DecompressAll("/Users/lukas/Downloads/Discogs", "discogs_20190101", "xml.gz"))
	// parser.Parse()
	d := decoder.NewDecoder("/Users/aronlukas/Downloads/discogs/discogs_20191101_artists.xml")
	n := 1

	for i := 0; n != 0; i++ {
		f, err := os.Create("/Users/aronlukas/Downloads/discogs/discogs_20191101_artists" + strconv.Itoa(i) + ".json")
		if err != nil {
			log.Fatal(err)
		}

		n, err = d.DecodeArtistJson(f, 100_000)
		_ = f.Close()
	}

	_ = d.Close()
}
