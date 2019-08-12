package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseRelease(se xml.StartElement, tr xml.TokenReader) *model.Release {
	if se.Name.Local != "release" {
		return nil
	}

	release := &model.Release{}
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
			case "image":
				img := ParseImage(se)
				release.Images = append(release.Images, img)
			case "artists":
				release.Artists = parseArtists("artists", tr)
			case "extraartists":
				release.ExtraArtists = parseArtists("extraartists", tr)
			case "title":
				release.Title = parseValue(tr)
			case "labels":
				release.Labels = parseLabels(tr)
			case "formats":
				release.Formats = parseFormats(tr)
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
				release.TrackList = parseTrackList(tr)
			case "identifiers":
				release.Identifiers = parseIdentifiers(tr)
			case "videos":
				release.Videos = parseVideos(tr)
			case "companies":
				release.Companies = parseCompanies(tr)
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "release" {
			break
		}
	}

	return release
}

func parseArtists(wrapperName string, tr xml.TokenReader) []*model.ReleaseArtist {
	artists := make([]*model.ReleaseArtist, 0, 0)
	artist := &model.ReleaseArtist{}
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
			artist = &model.ReleaseArtist{}
		}
	}
	return artists
}

func parseLabels(tr xml.TokenReader) []*model.ReleaseLabel {
	labels := make([]*model.ReleaseLabel, 0, 0)
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "labels" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "label" {
			label := &model.ReleaseLabel{}

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

func parseFormats(tr xml.TokenReader) []*model.ReleaseFormat {
	formats := make([]*model.ReleaseFormat, 0, 0)
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "formats" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "format" {
			format := &model.ReleaseFormat{}
			for _, attr := range se.Attr {
				switch attr.Name.Local {
				case "qty":
					format.Quantity = attr.Value
				case "name":
					format.Name = attr.Value
				case "text":
					format.Text = attr.Value
				}
			}

			format.Descriptions = parseChildValues("descriptions", "description", tr)
			formats = append(formats, format)
		}
	}
	return formats
}

func parseCompanies(tr xml.TokenReader) []*model.ReleaseCompany {
	companies := make([]*model.ReleaseCompany, 0, 0)
	company := &model.ReleaseCompany{}
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "companies" {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "id":
				company.Id = parseValue(tr)
			case "name":
				company.Name = parseValue(tr)
			case "catno":
				company.Category = parseValue(tr)
			case "entity_type":
				company.EntityType = parseValue(tr)
			case "entity_type_name":
				company.EntityTypeName = parseValue(tr)
			case "resource_url":
				company.ResourceUrl = parseValue(tr)
			}
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "company" {
			companies = append(companies, company)
			company = &model.ReleaseCompany{}
		}
	}
	return companies
}

func parseVideos(tr xml.TokenReader) []*model.ReleaseVideo {
	videos := make([]*model.ReleaseVideo, 0, 0)
	video := &model.ReleaseVideo{}
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "videos" {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "video":
				for _, attr := range se.Attr {
					switch attr.Name.Local {
					case "duration":
						video.Duration = attr.Value
					case "embed":
						video.Embed = attr.Value
					case "src":
						video.Src = attr.Value
					}
				}
			case "title":
				video.Title = parseValue(tr)
			case "description":
				video.Description = parseValue(tr)
			}
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "video" {
			videos = append(videos, video)
			video = &model.ReleaseVideo{}
		}
	}
	return videos
}

func parseIdentifiers(tr xml.TokenReader) []*model.ReleaseIdentifier {
	identifiers := make([]*model.ReleaseIdentifier, 0, 0)
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "identifiers" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "identifier" {
			identifier := &model.ReleaseIdentifier{}
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

func parseTrackList(tr xml.TokenReader) []*model.ReleaseTrack {
	trackList := make([]*model.ReleaseTrack, 0, 0)
	track := &model.ReleaseTrack{}
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "tracklist" {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "position":
				track.Position = parseValue(tr)
			case "title":
				track.Title = parseValue(tr)
			case "duration":
				track.Duration = parseValue(tr)
			}
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "track" {
			trackList = append(trackList, track)
			track = &model.ReleaseTrack{}
		}
	}
	return trackList
}
