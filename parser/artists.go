package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseArtists(d *xml.Decoder, limit int) (artists []model.Artist, err error) {
	var t xml.Token
	cnt := 0
	artists = make([]model.Artist, 0, 0)
	for t, err = d.Token(); t != nil && err == nil && cnt != limit; t, err = d.Token() {
		if IsStartElementName(t, "artist") {
			artist, err := ParseArtist(t.(xml.StartElement), d)
			if err != nil {
				return artists, err
			}
			artists = append(artists, artist)
			cnt++
		}
	}

	return artists, err
}

func ParseArtist(se xml.StartElement, tr xml.TokenReader) (model.Artist, error) {
	artist := model.Artist{}
	if se.Name.Local != "artist" {
		return artist, notCorrectStarElement
	}

	for {
		t, _ := tr.Token()
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				imgs, err := ParseImages(se, tr)
				if err != nil {
					return artist, err
				}

				artist.Images = imgs
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

	return artist, nil
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
