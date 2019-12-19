package decoder

import (
	"encoding/xml"
	"github.com/Twyer/discogs-parser/model"
)

func (x *XMLDecoder) parseTrackList() (trackList []model.Track) {
	if x.err != nil {
		return trackList
	}

	track := model.Track{}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "tracklist" {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "position":
				track.Position = x.parseValue()
			case "title":
				track.Title = x.parseValue()
			case "duration":
				track.Duration = x.parseValue()
			}
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "track" {
			trackList = append(trackList, track)
			track = model.Track{}
		}
	}

	return trackList
}
