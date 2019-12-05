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
	d := decoder.NewDecoder("/Users/aronlukas/Downloads/discogs/discogs_20191101_releases.xml")
	defer d.Close()

	//for i := 0; ; i++ {
	//	n, b, err := d.ArtistJson(100_000)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	if n == 0 {
	//		break
	//	}
	//
	//	f, err := os.Create("/Users/aronlukas/Downloads/discogs/discogs_20191101_artists" + strconv.Itoa(i) + ".json")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	_, err = f.Write(b)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	_ = f.Close()
	//}

	//for i := 0; ; i++ {
	//	n, b, err := d.LabelJson(100_000)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	if n == 0 {
	//		break
	//	}
	//
	//	f, err := os.Create("/Users/aronlukas/Downloads/discogs/discogs_20191101_labels_" + strconv.Itoa(i) + ".json")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	_, err = f.Write(b)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	_ = f.Close()
	//}

	//for i := 0; i < 5; i++ {
	//	n, b, err := d.MasterJson(10_000)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	if n == 0 {
	//		break
	//	}
	//
	//	f, err := os.Create("/Users/aronlukas/Downloads/discogs/discogs_20191101_masters_" + strconv.Itoa(i) + ".json")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	_, err = f.Write(b)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	_ = f.Close()
	//}

	for i := 0; i < 5; i++ {
		n, b, err := d.ReleaseJson(3_000)
		if err != nil {
			log.Fatal(err)
		}

		if n == 0 {
			break
		}

		f, err := os.Create("/Users/aronlukas/Downloads/discogs/discogs_20191101_releases_" + strconv.Itoa(i) + ".json")
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
