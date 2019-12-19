package decoder

import (
	"encoding/xml"
	"errors"
	"github.com/Twyer/discogs-parser/model"
	"io"
	"os"
	"strings"
)

type FileType int

const (
	Unknown FileType = iota
	Artists
	Labels
	Masters
	Releases
)

var (
	notCorrectStarElement error = errors.New("token is not a correct start element")
)

type XMLDecoder struct {
	f   *os.File
	d   *xml.Decoder
	o   Options
	err error
}

func NewDecoder(fileName string, options ...Options) Decoder {
	d := &XMLDecoder{}

	if options != nil && len(options) > 0 {
		d.o = options[0]
	}

	d.f, d.err = os.Open(fileName)
	if d.err != nil {
		return d
	}

	d.d = xml.NewDecoder(d.f)
	return d
}

func (x *XMLDecoder) Close() error {
	return x.f.Close()
}

func (x *XMLDecoder) Artists(limit int) (int, []model.Artist, error) {
	if x.err != nil {
		return 0, nil, x.err
	}

	artists := x.parseArtists(limit)
	if x.err == nil || x.err == io.EOF {
		artists = x.filterArtists(artists)
	}
	return len(artists), artists, x.err
}

func (x *XMLDecoder) filterArtists(as []model.Artist) []model.Artist {
	fa := make([]model.Artist, 0, len(as))
	for _, a := range as {
		if x.o.QualityLevel.Includes(StrToQualityLevel(a.DataQuality)) {
			fa = append(fa, a)
		}
	}

	return fa
}

func (x *XMLDecoder) Labels(limit int) (int, []model.Label, error) {
	if x.err != nil {
		return 0, nil, x.err
	}

	labels := x.parseLabels(limit)
	if x.err == nil || x.err == io.EOF {
		labels = x.filterLabels(labels)
	}
	return len(labels), labels, x.err
}

func (x *XMLDecoder) filterLabels(ls []model.Label) []model.Label {
	fl := make([]model.Label, 0, len(ls))
	for _, l := range ls {
		if x.o.QualityLevel.Includes(StrToQualityLevel(l.DataQuality)) {
			fl = append(fl, l)
		}
	}

	return fl
}

func (x *XMLDecoder) Masters(limit int) (int, []model.Master, error) {
	if x.err != nil {
		return 0, nil, x.err
	}

	masters := x.parseMasters(limit)
	if x.err == nil || x.err == io.EOF {
		masters = x.filterMasters(masters)
	}

	return len(masters), masters, x.err
}

func (x *XMLDecoder) filterMasters(ms []model.Master) []model.Master {
	fm := make([]model.Master, 0, len(ms))
	for _, m := range ms {
		if x.o.QualityLevel.Includes(StrToQualityLevel(m.DataQuality)) {
			fm = append(fm, m)
		}
	}

	return fm
}

func (x *XMLDecoder) Releases(limit int) (int, []model.Release, error) {
	if x.err != nil {
		return 0, nil, x.err
	}

	releases := x.parseReleases(limit)
	if x.err == nil || x.err == io.EOF {
		releases = x.filterReleases(releases)
	}
	return len(releases), releases, x.err
}

func (x *XMLDecoder) filterReleases(rs []model.Release) []model.Release {
	fr := make([]model.Release, 0, len(rs))
	for _, r := range rs {
		if x.o.QualityLevel.Includes(StrToQualityLevel(r.DataQuality)) {
			fr = append(fr, r)
		}
	}

	return fr
}

func (x *XMLDecoder) isStartElement(token xml.Token) bool {
	_, ok := token.(xml.StartElement)
	return ok
}

func (x *XMLDecoder) isEndElement(token xml.Token) bool {
	_, ok := token.(xml.EndElement)
	return ok
}

func (x *XMLDecoder) isStartElementName(token xml.Token, name string) bool {
	se, ok := token.(xml.StartElement)
	return ok && se.Name.Local == name
}

func (x *XMLDecoder) isEndElementName(token xml.Token, name string) bool {
	ee, ok := token.(xml.EndElement)
	return ok && ee.Name.Local == name
}

func (x *XMLDecoder) parseValue() string {
	sb := strings.Builder{}
	for {
		t, _ := x.d.Token()
		if x.isEndElement(t) {
			break
		}

		if cr, ok := t.(xml.CharData); ok {
			sb.Write(cr)
		}
	}
	return sb.String()
}

func (x *XMLDecoder) parseChildValues(parentName, childName string) (children []string) {
	for {
		t, _ := x.d.Token()
		if x.isStartElementName(t, childName) {
			children = append(children, x.parseValue())
		}
		if x.isEndElementName(t, parentName) {
			break
		}
	}
	return children
}
