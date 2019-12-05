package main

import (
	"fmt"
	"github.com/Twyer/discogs/decoder"
	"github.com/Twyer/discogs/writer"
	"log"
)

func main() {
	// fmt.Println(util.DecompressAll("/Users/lukas/Downloads/Discogs", "discogs_20190101", "xml.gz"))
	// parser.Parse()
	d := decoder.NewDecoder("/Users/aronlukas/Downloads/discogs/discogs_20191101_releases.xml")
	defer d.Close()

	pg := writer.NewPostgres("localhost", "discogs", "user", "password", "disable", 5432)
	if pg.Error != nil {
		log.Fatal(pg.Error)
	}
	defer pg.Close()
	fmt.Print(pg.Ping())
	num, data, err := d.Releases(1000)
	if num != 1000 {

		log.Fatalf("not exact number: %d", num)
	}
	if err != nil {
		log.Fatal(err)
	}
	for _, release := range data {
		err = pg.WriteRelease(release)
		if err != nil {
			fmt.Println(err)
		}
	}

}
