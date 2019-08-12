package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseArtist(se xml.StartElement, tr xml.TokenReader) *model.Artist {
	if se.Name.Local != "artist" {
		return nil
	}

	artist := &model.Artist{}
	for {
		t, _ := tr.Token()
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "image":
				img := ParseImage(se)
				artist.Images = append(artist.Images, img)
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
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "artist" {
			break
		}
	}

	return artist
}

func parseAliases(tr xml.TokenReader) []*model.ArtistAlias {
	aliases := make([]*model.ArtistAlias, 0, 0)
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "aliases" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "name" {
			alias := &model.ArtistAlias{
				Id:   se.Attr[0].Value,
				Name: parseValue(tr),
			}
			aliases = append(aliases, alias)
		}
	}
	return aliases
}
