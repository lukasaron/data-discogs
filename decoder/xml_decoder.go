package decoder

import (
	"encoding/xml"
	"github.com/Twyer/discogs-parser/model"
	"github.com/Twyer/discogs-parser/parser"
	"io"
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
	options Options
	error   error
}

func NewDecoder(fileName string, options ...Options) Decoder {
	d := XMLDecoder{}

	if options != nil && len(options) > 0 {
		d.options = options[0]
	}

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
	if err == nil || err == io.EOF {
		artists = d.filterArtists(artists)
	}
	return len(artists), artists, err
}

func (d XMLDecoder) filterArtists(as []model.Artist) []model.Artist {
	fa := make([]model.Artist, 0, len(as))
	for _, a := range as {
		if d.options.QualityLevel.Includes(StrToQualityLevel(a.DataQuality)) {
			fa = append(fa, a)
		}
	}

	return fa
}

func (d XMLDecoder) Labels(limit int) (int, []model.Label, error) {
	if d.error != nil {
		return 0, nil, d.error
	}

	labels, err := parser.ParseLabels(d.decoder, limit)
	if err == nil || err == io.EOF {
		labels = d.filterLabels(labels)
	}
	return len(labels), labels, err
}

func (d XMLDecoder) filterLabels(ls []model.Label) []model.Label {
	fl := make([]model.Label, 0, len(ls))
	for _, l := range ls {
		if d.options.QualityLevel.Includes(StrToQualityLevel(l.DataQuality)) {
			fl = append(fl, l)
		}
	}

	return fl
}

func (d XMLDecoder) Masters(limit int) (int, []model.Master, error) {
	if d.error != nil {
		return 0, nil, d.error
	}

	masters, err := parser.ParseMasters(d.decoder, limit)
	if err == nil || err == io.EOF {
		masters = d.filterMasters(masters)
	}

	return len(masters), masters, err
}

func (d XMLDecoder) filterMasters(ms []model.Master) []model.Master {
	fm := make([]model.Master, 0, len(ms))
	for _, m := range ms {
		if d.options.QualityLevel.Includes(StrToQualityLevel(m.DataQuality)) {
			fm = append(fm, m)
		}
	}

	return fm
}

func (d XMLDecoder) Releases(limit int) (int, []model.Release, error) {
	if d.error != nil {
		return 0, nil, d.error
	}

	releases, err := parser.ParseReleases(d.decoder, limit)
	if err == nil || err == io.EOF {
		releases = d.filterReleases(releases)
	}
	return len(releases), releases, err
}

func (d XMLDecoder) filterReleases(rs []model.Release) []model.Release {
	fr := make([]model.Release, 0, len(rs))
	for _, r := range rs {
		if d.options.QualityLevel.Includes(StrToQualityLevel(r.DataQuality)) {
			fr = append(fr, r)
		}
	}

	return fr
}
