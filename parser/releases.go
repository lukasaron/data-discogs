package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs-parser/model"
)

func ParseReleases(d *xml.Decoder, limit int) (releases []model.Release, err error) {
	var t xml.Token
	cnt := 0
	for t, err = d.Token(); t != nil && err == nil && cnt != limit; t, err = d.Token() {
		if IsStartElementName(t, "release") {
			rls, err := ParseRelease(t.(xml.StartElement), d)
			if err != nil {
				return releases, err
			}

			releases = append(releases, rls)
			cnt++
		}
	}

	return releases, err
}

func ParseRelease(se xml.StartElement, tr xml.TokenReader) (release model.Release, err error) {
	if se.Name.Local != "release" {
		return release, notCorrectStarElement
	}

	for _, attr := range se.Attr {
		switch attr.Name.Local {
		case "id":
			release.Id = attr.Value
		case "status":
			release.Status = attr.Value
		}
	}

	for {
		t, _ := tr.Token()
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				imgs, err := ParseImages(se, tr)
				if err != nil {
					return release, err
				}
				release.Images = imgs
			case "artists":
				release.Artists = parseArtists("artists", tr)
			case "extraartists":
				release.ExtraArtists = parseArtists("extraartists", tr)
			case "title":
				release.Title = parseValue(tr)
			case "labels":
				release.Labels = parseLabels(tr)
			case "formats":
				release.Formats = ParseFormats(tr)
			case "genres":
				release.Genres = parseChildValues("genres", "genre", tr)
			case "styles":
				release.Styles = parseChildValues("styles", "style", tr)
			case "country":
				release.Country = parseValue(tr)
			case "released":
				release.Released = parseValue(tr)
			case "notes":
				release.Notes = parseValue(tr)
			case "data_quality":
				release.DataQuality = parseValue(tr)
			case "master_id":
				release.MainRelease = se.Attr[0].Value
				release.MasterId = parseValue(tr)
			case "tracklist":
				release.TrackList, _ = ParseTrackList(tr)
			case "identifiers":
				release.Identifiers = parseIdentifiers(tr)
			case "videos":
				release.Videos, _ = ParseVideos(tr)
			case "companies":
				release.Companies = ParseCompanies(tr)
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "release" {
			break
		}
	}

	return release, nil
}

func parseArtists(wrapperName string, tr xml.TokenReader) (artists []model.ReleaseArtist) {
	artist := model.ReleaseArtist{}
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == wrapperName {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "id":
				artist.Id = parseValue(tr)
			case "name":
				artist.Name = parseValue(tr)
			case "anv":
				artist.Anv = parseValue(tr)
			case "join":
				artist.Join = parseValue(tr)
			case "role":
				artist.Role = parseValue(tr)
			case "tracks":
				artist.Tracks = parseValue(tr)
			}
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "artist" {
			artists = append(artists, artist)
			artist = model.ReleaseArtist{}
		}
	}
	return artists
}

func parseLabels(tr xml.TokenReader) (labels []model.ReleaseLabel) {
	for {
		t, _ := tr.Token()
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

func parseIdentifiers(tr xml.TokenReader) (identifiers []model.Identifier) {
	for {
		t, _ := tr.Token()
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
