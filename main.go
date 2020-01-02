package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/lukasaron/discogs-parser/decoder"
	"github.com/lukasaron/discogs-parser/writer"
	"log"
)

func main() {
	// Create XML decoder
	d := decoder.NewXmlDecoder("/Users/lukasaron/Downloads/discogs/releases.xml",
		decoder.Options{
			QualityLevel: decoder.All,
			Block: decoder.Block{
				Size:  100,
				Limit: 0,
				Skip:  0,
			},
			FileType: decoder.Releases},
	)

	// Setup a connection to DB
	connStr := "user=user password=password dbname=discogs host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Create a DB writer
	w := writer.NewDbWriter(db, writer.Options{ExcludeImages: true})

	// Decode Discogs data from XML into DB
	err = decoder.DecodeData(d, w)
	if err != nil {
		log.Fatal(err)
	}
}
