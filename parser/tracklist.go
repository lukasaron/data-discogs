package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseTrackList(tr xml.TokenReader) ([]model.Track, error) {
	trackList := make([]model.Track, 0, 0)
	track := model.Track{}

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
			track = model.Track{}
		}
	}
	return trackList, nil
}
