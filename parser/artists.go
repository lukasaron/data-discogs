package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseArtist(se xml.StartElement, tr xml.TokenReader) model.Artist {
	artist := model.Artist{}
	if se.Name.Local != "artist" {
		return artist
	}

	for {
		t, _ := tr.Token()
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				artist.Images = ParseImages(se, tr)
			case "id":
				artist.Id = parseValue(tr)
			case "name":
				artist.Name = parseValue(tr)
			case "realname":
				artist.RealName = parseValue(tr)
			case "namevariations":
				artist.NameVariations = parseChildValues("namevariations", "name", tr)
			case "aliases":
				artist.Aliases = parseAliases(tr)
			case "profile":
				artist.Profile = parseValue(tr)
			case "data_quality":
				artist.DataQuality = parseValue(tr)
			case "urls":
				artist.Urls = parseChildValues("urls", "url", tr)
			}
		}
		if IsEndElementName(t, "artist") {
			break
		}
	}

	return artist
}

func parseAliases(tr xml.TokenReader) []model.Alias {
	aliases := make([]model.Alias, 0, 0)
	for {
		t, _ := tr.Token()
		if IsEndElementName(t, "aliases") {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "name" {
			alias := model.Alias{
				Id:   se.Attr[0].Value,
				Name: parseValue(tr),
			}
			aliases = append(aliases, alias)
		}
	}
	return aliases
}
