package decoder

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
	"github.com/Twyer/discogs/parser"
	"os"
)

type FileType int

const (
	Unknown FileType = iota
	Artists
	Labels
	Masters
	Releases
)

type XMLDecoder struct {
	file    *os.File
	decoder *xml.Decoder
	error   error
}

func NewDecoder(fileName string) Decoder {
	d := XMLDecoder{}

	d.file, d.error = os.Open(fileName)
	if d.error != nil {
		return d
	}

	d.decoder = xml.NewDecoder(d.file)
	return d
}

func (d XMLDecoder) Close() error {
	return d.file.Close()
}

func (d XMLDecoder) Artists(limit int) (int, []model.Artist, error) {
	if d.error != nil {
		return 0, nil, d.error
	}

	artists, err := parser.ParseArtists(d.decoder, limit)
	return len(artists), artists, err
}

func (d XMLDecoder) Labels(limit int) (int, []model.Label, error) {
	if d.error != nil {
		return 0, nil, d.error
	}

	labels, err := parser.ParseLabels(d.decoder, limit)
	return len(labels), labels, err
}

func (d XMLDecoder) Masters(limit int) (int, []model.Master, error) {
	if d.error != nil {
		return 0, nil, d.error
	}

	masters, err := parser.ParseMasters(d.decoder, limit)
	return len(masters), masters, err
}

func (d XMLDecoder) Releases(limit int) (int, []model.Release, error) {
	if d.error != nil {
		return 0, nil, d.error
	}

	releases, err := parser.ParseReleases(d.decoder, limit)
	return len(releases), releases, err
}
