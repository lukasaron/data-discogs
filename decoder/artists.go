package decoder

import (
	"encoding/xml"
	"github.com/lukasaron/discogs-parser/model"
)

func (x *XMLDecoder) parseArtists(limit int) (artists []model.Artist) {
	if x.err != nil {
		return artists
	}

	var t xml.Token
	cnt := 0
	for t, x.err = x.d.Token(); t != nil && x.err == nil && cnt != limit; t, x.err = x.d.Token() {
		if x.isStartElementName(t, "artist") {
			artist := x.parseArtist(t.(xml.StartElement))
			if x.err != nil {
				return artists
			}
			artists = append(artists, artist)
			cnt++
		}
	}

	return artists
}

func (x *XMLDecoder) parseArtist(se xml.StartElement) (artist model.Artist) {
	if x.err != nil {
		return artist
	}

	if se.Name.Local != "artist" {
		x.err = notCorrectStarElement
		return artist
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil && !x.isEndElementName(t, "artist"); t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				imgs := x.parseImages(se)
				if x.err != nil {
					return artist
				}

				artist.Images = imgs
			case "id":
				artist.Id = x.parseValue()
			case "name":
				artist.Name = x.parseValue()
			case "realname":
				artist.RealName = x.parseValue()
			case "namevariations":
				artist.NameVariations = x.parseChildValues("namevariations", "name")
			case "members":
				artist.Members = x.parseMembers()
			case "aliases":
				artist.Aliases = x.parseAliases()
			case "profile":
				artist.Profile = x.parseValue()
			case "data_quality":
				artist.DataQuality = x.parseValue()
			case "urls":
				artist.Urls = x.parseChildValues("urls", "url")
			}
		}
	}

	return artist
}

func (x *XMLDecoder) parseAliases() (aliases []model.Alias) {
	if x.err != nil {
		return
	}
	var t xml.Token

	for t, x.err = x.d.Token(); x.err == nil && !x.isEndElementName(t, "aliases"); t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "name" {
			alias := model.Alias{
				Id:   se.Attr[0].Value,
				Name: x.parseValue(),
			}
			aliases = append(aliases, alias)
		}
	}

	return aliases
}

func (x *XMLDecoder) parseMembers() (members []model.Member) {
	if x.err != nil {
		return
	}
	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil && !x.isEndElementName(t, "members"); t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "name" {
			member := model.Member{
				Id:   se.Attr[0].Value,
				Name: x.parseValue(),
			}
			members = append(members, member)
		}
	}

	return members

}
