package decoder

import (
	"encoding/json"
	"github.com/Twyer/discogs/parser"
	"io"
)

func (d *Decoder) DecodeArtistJson(w io.Writer, limit int) (int, error) {
	if d.Error == nil {
		artists := parser.ParseArtists(d.decoder, limit)
		n := len(artists)
		if n == 0 {
			return 0, nil
		}

		b, err := json.Marshal(artists)
		if err != nil {
			return 0, err
		}
		_, err = w.Write(b)
		return n, err
	}

	return 0, d.Error
}
