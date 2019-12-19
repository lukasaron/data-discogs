package decoder

import (
	"encoding/xml"
	"github.com/Twyer/discogs-parser/model"
)

func (x *XMLDecoder) parseMasters(limit int) (masters []model.Master) {
	if x.err != nil {
		return masters
	}

	var t xml.Token
	cnt := 0
	for t, x.err = x.d.Token(); t != nil && x.err == nil && cnt != limit; t, x.err = x.d.Token() {
		if x.isStartElementName(t, "master") {
			m := x.parseMaster(t.(xml.StartElement))
			if x.err != nil {
				return masters
			}

			masters = append(masters, m)
			cnt++
		}
	}

	return masters
}

func (x *XMLDecoder) parseMaster(se xml.StartElement) (master model.Master) {
	if x.err != nil {
		return master
	}

	if se.Name.Local != "master" {
		x.err = notCorrectStarElement
		return master
	}

	master.Id = se.Attr[0].Value
	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				imgs := x.parseImages(se)
				if x.err != nil {
					return master
				}
				master.Images = imgs
			case "main_release":
				master.MainRelease = x.parseValue()
			case "artists":
				master.Artists = x.parseReleaseArtists("artists")
			case "genres":
				master.Genres = x.parseChildValues("genres", "genre")
			case "styles":
				master.Styles = x.parseChildValues("styles", "style")
			case "year":
				master.Year = x.parseValue()
			case "title":
				master.Title = x.parseValue()
			case "data_quality":
				master.DataQuality = x.parseValue()
			case "videos":
				master.Videos = x.parseVideos()
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "master" {
			break
		}
	}

	return master
}
