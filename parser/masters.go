package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseMaster(se xml.StartElement, tr xml.TokenReader) *model.Master {
	if se.Name.Local != "master" {
		return nil
	}

	master := &model.Master{}
	master.Id = se.Attr[0].Value

	for {
		t, _ := tr.Token()
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "image":
				img := ParseImage(se)
				master.Images = append(master.Images, img)
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
				master.Videos = parseVideos(tr)
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "master" {
			break
		}
	}

	return master
}
