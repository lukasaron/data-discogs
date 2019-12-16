package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs-parser/model"
)

func ParseMasters(d *xml.Decoder, limit int) (masters []model.Master, err error) {
	var t xml.Token
	cnt := 0
	for t, err = d.Token(); t != nil && err == nil && cnt != limit; t, err = d.Token() {
		if IsStartElementName(t, "master") {
			m, err := ParseMaster(t.(xml.StartElement), d)
			if err != nil {
				return masters, err
			}

			masters = append(masters, m)
			cnt++
		}
	}

	return masters, err
}

func ParseMaster(se xml.StartElement, tr xml.TokenReader) (master model.Master, err error) {

	if se.Name.Local != "master" {
		return master, notCorrectStarElement
	}

	master.Id = se.Attr[0].Value
	for {
		t, _ := tr.Token()
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				imgs, err := ParseImages(se, tr)
				if err != nil {
					return master, err
				}
				master.Images = imgs
			case "main_release":
				master.MainRelease = parseValue(tr)
			case "artists":
				master.Artists = parseArtists("artists", tr)
			case "genres":
				master.Genres = parseChildValues("genres", "genre", tr)
			case "styles":
				master.Styles = parseChildValues("styles", "style", tr)
			case "year":
				master.Year = parseValue(tr)
			case "title":
				master.Title = parseValue(tr)
			case "data_quality":
				master.DataQuality = parseValue(tr)
			case "videos":
				master.Videos, _ = ParseVideos(tr)
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "master" {
			break
		}
	}

	return master, nil
}
