package decoder

import (
	"encoding/xml"
	"github.com/lukasaron/discogs-parser/model"
)

func (x *XMLDecoder) parseReleases(limit int) (releases []model.Release) {
	if x.err != nil {
		return releases
	}

	var t xml.Token
	cnt := 0
	for t, x.err = x.d.Token(); t != nil && x.err == nil && cnt != limit; t, x.err = x.d.Token() {
		if x.isStartElementName(t, "release") {
			rls := x.parseRelease(t.(xml.StartElement))
			if x.err != nil {
				return releases
			}

			releases = append(releases, rls)
			cnt++
		}
	}

	return releases
}

func (x *XMLDecoder) parseRelease(se xml.StartElement) (release model.Release) {
	if x.err != nil {
		return release
	}

	if se.Name.Local != "release" {
		x.err = notCorrectStarElement
		return release
	}

	for _, attr := range se.Attr {
		switch attr.Name.Local {
		case "id":
			release.Id = attr.Value
		case "status":
			release.Status = attr.Value
		}
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				imgs := x.parseImages(se)
				if x.err != nil {
					return release
				}
				release.Images = imgs
			case "artists":
				release.Artists = x.parseReleaseArtists("artists")
			case "extraartists":
				release.ExtraArtists = x.parseReleaseArtists("extraartists")
			case "title":
				release.Title = x.parseValue()
			case "labels":
				release.Labels = x.parseReleaseLabels()
			case "formats":
				release.Formats = x.parseFormats()
			case "genres":
				release.Genres = x.parseChildValues("genres", "genre")
			case "styles":
				release.Styles = x.parseChildValues("styles", "style")
			case "country":
				release.Country = x.parseValue()
			case "released":
				release.Released = x.parseValue()
			case "notes":
				release.Notes = x.parseValue()
			case "data_quality":
				release.DataQuality = x.parseValue()
			case "master_id":
				release.MainRelease = se.Attr[0].Value
				release.MasterId = x.parseValue()
			case "tracklist":
				release.TrackList = x.parseTrackList()
			case "identifiers":
				release.Identifiers = x.parseIdentifiers()
			case "videos":
				release.Videos = x.parseVideos()
			case "companies":
				release.Companies = x.parseCompanies()
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "release" {
			break
		}
	}

	return release
}

func (x *XMLDecoder) parseReleaseArtists(wrapperName string) (artists []model.ReleaseArtist) {
	if x.err != nil {
		return artists
	}

	artist := model.ReleaseArtist{}
	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == wrapperName {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "id":
				artist.Id = x.parseValue()
			case "name":
				artist.Name = x.parseValue()
			case "anv":
				artist.Anv = x.parseValue()
			case "join":
				artist.Join = x.parseValue()
			case "role":
				artist.Role = x.parseValue()
			case "tracks":
				artist.Tracks = x.parseValue()
			}
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "artist" {
			artists = append(artists, artist)
			artist = model.ReleaseArtist{}
		}
	}
	return artists
}

func (x *XMLDecoder) parseReleaseLabels() (labels []model.ReleaseLabel) {
	if x.err != nil {
		return labels
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "labels" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "label" {
			label := model.ReleaseLabel{}

			for _, attr := range se.Attr {
				switch attr.Name.Local {
				case "id":
					label.Id = attr.Value
				case "name":
					label.Name = attr.Value
				case "catno":
					label.Category = attr.Value
				}
			}

			labels = append(labels, label)
		}
	}
	return labels
}

func (x *XMLDecoder) parseIdentifiers() (identifiers []model.Identifier) {
	if x.err != nil {
		return identifiers
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "identifiers" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "identifier" {
			identifier := model.Identifier{}
			for _, attr := range se.Attr {
				switch attr.Name.Local {
				case "description":
					identifier.Description = attr.Value
				case "type":
					identifier.Type = attr.Value
				case "value":
					identifier.Value = attr.Value
				}
			}

			identifiers = append(identifiers, identifier)
		}
	}
	return identifiers
}
