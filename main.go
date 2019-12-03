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
	defer d.Close()

	for i := 0; ; i++ {
		n, b, err := d.DecodeArtistJson(100_000)
		if err != nil {
			log.Fatal(err)
		}

		if n == 0 {
			break
		}

		f, err := os.Create("/Users/aronlukas/Downloads/discogs/discogs_20191101_artists" + strconv.Itoa(i) + ".json")
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.Write(b)
		if err != nil {
			log.Fatal(err)
		}

		_ = f.Close()
	}
}
