package decoder

import (
	"encoding/json"
	"github.com/Twyer/discogs/parser"
)

func (d *Decoder) DecodeArtistJson(limit int) (int, []byte, error) {
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
