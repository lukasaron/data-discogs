package decoder

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
	"github.com/Twyer/discogs/parser"
	"os"
)

type Decoder struct {
	file    *os.File
	decoder *xml.Decoder
	Error   error
}

func NewDecoder(fileName string) *Decoder {
	d := &Decoder{}
	d.file, d.Error = os.Open(fileName)
	if d.Error == nil {
		d.decoder = xml.NewDecoder(d.file)
	}
	return d
}

func (d *Decoder) Close() error {
	return d.file.Close()
}

func (d *Decoder) Artists(limit int) (int, []model.Artist, error) {
	if d.Error == nil {
		artists, err := parser.ParseArtists(d.decoder, limit)
		return len(artists), artists, err
	}

	return 0, nil, d.Error
}

func (d *Decoder) Labels(limit int) (int, []model.Label, error) {
	if d.Error == nil {
		labels, err := parser.ParseLabels(d.decoder, limit)
		return len(labels), labels, err
	}

	return 0, nil, d.Error
}

func (d *Decoder) Masters(limit int) (int, []model.Master, error) {
	if d.Error == nil {
		masters, err := parser.ParseMasters(d.decoder, limit)
		return len(masters), masters, err
	}

	return 0, nil, d.Error
}

func (d *Decoder) Releases(limit int) (int, []model.Release, error) {
	if d.Error == nil {
		releases, err := parser.ParseReleases(d.decoder, limit)
		return len(releases), releases, err
	}

	return 0, nil, d.Error
}
