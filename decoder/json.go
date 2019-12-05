package decoder

import (
	"encoding/json"
	"github.com/Twyer/discogs/parser"
)

func (d *Decoder) ArtistJson(limit int) (int, []byte, error) {
	if d.Error == nil {
		artists := parser.ParseArtists(d.decoder, limit)
		b, err := json.Marshal(artists)
		if err != nil {
			return 0, b, err
		}

		return len(artists), b, err
	}

	return 0, nil, d.Error
}

func (d *Decoder) LabelJson(limit int) (int, []byte, error) {
	if d.Error == nil {
		labels := parser.ParseLabels(d.decoder, limit)
		b, err := json.Marshal(labels)
		if err != nil {
			return 0, b, err
		}

		return len(labels), b, err
	}

	return 0, nil, d.Error
}

func (d *Decoder) MasterJson(limit int) (int, []byte, error) {
	if d.Error == nil {
		masters := parser.ParseMasters(d.decoder, limit)
		b, err := json.Marshal(masters)
		if err != nil {
			return 0, b, err
		}

		return len(masters), b, err
	}

	return 0, nil, d.Error
}

func (d *Decoder) ReleaseJson(limit int) (int, []byte, error) {
	if d.Error == nil {
		releases := parser.ParseReleases(d.decoder, limit)
		b, err := json.Marshal(releases)
		if err != nil {
			return 0, b, err
		}

		return len(releases), b, err
	}

	return 0, nil, d.Error
}
