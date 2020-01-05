package decoder

import (
	"encoding/xml"
	"errors"
	"github.com/lukasaron/discogs-parser/model"
	"github.com/lukasaron/discogs-parser/writer"
	"io"
	"log"
	"os"
	"strings"
)

var (
	// Errors returned when failure occurs
	wrongTypeSpecified    = errors.New("wrong file type specified")
	notCorrectStarElement = errors.New("token is not a correct start element")
)

// XML Decoder type is behaviour structure that implements Decoder interface and supports
// the Discogs XML dump data decoding.
type XMLDecoder struct {
	f   *os.File
	d   *xml.Decoder
	o   Options
	err error
}

func NewXmlDecoder(fileName string, options ...Options) Decoder {
	d := &XMLDecoder{}

	d.f, d.err = os.Open(fileName)
	d.d = xml.NewDecoder(d.f)

	if options != nil && len(options) > 0 {
		d.o = options[0]
	}

	return d
}

func (x *XMLDecoder) Close() error {
	return x.f.Close()
}

func (x *XMLDecoder) Options() Options {
	return x.o
}

func (x *XMLDecoder) SetOptions(opt Options) {
	x.o = opt
}

// Decode data
func (x *XMLDecoder) Decode(writer writer.Writer) error {
	if x.err != nil {
		return x.err
	}

	if x.o.Block.Limit < 1 {
		x.o.Block.Limit = int(^uint(0) >> 1)
	}

	// get decode function based on the file type
	fn, err := x.decodeFunction()
	if err != nil {
		return err
	}

	for blockCount := 1; blockCount <= x.o.Block.Limit; blockCount++ {
		// call appropriate decoder function
		num, err := fn(writer, x.o.Block.Size, blockCount > x.o.Block.Skip)
		if err != nil && err != io.EOF {
			log.Printf("Block %d failed [%d]\n", blockCount, num)
			return err
		}

		if num == 0 && err == io.EOF {
			break
		}

		if blockCount > x.o.Block.Skip {
			log.Printf("Block %d written [%d]\n", blockCount, num)
		} else {
			log.Printf("Block %d skipped [%d]\n", blockCount, num)
		}
	}

	return nil
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

//--------------------------------------------------- FILTERS ---------------------------------------------------

func (x *XMLDecoder) filterArtists(as []model.Artist) []model.Artist {
	fa := make([]model.Artist, 0, len(as))
	for _, a := range as {
		if x.o.QualityLevel.Includes(ToQualityLevel(a.DataQuality)) {
			fa = append(fa, a)
		}
	}

	return fa
}

func (x *XMLDecoder) filterLabels(ls []model.Label) []model.Label {
	fl := make([]model.Label, 0, len(ls))
	for _, l := range ls {
		if x.o.QualityLevel.Includes(ToQualityLevel(l.DataQuality)) {
			fl = append(fl, l)
		}
	}

	return fl
}

func (x *XMLDecoder) filterMasters(ms []model.Master) []model.Master {
	fm := make([]model.Master, 0, len(ms))
	for _, m := range ms {
		if x.o.QualityLevel.Includes(ToQualityLevel(m.DataQuality)) {
			fm = append(fm, m)
		}
	}

	return fm
}

func (x *XMLDecoder) filterReleases(rs []model.Release) []model.Release {
	fr := make([]model.Release, 0, len(rs))
	for _, r := range rs {
		if x.o.QualityLevel.Includes(ToQualityLevel(r.DataQuality)) {
			fr = append(fr, r)
		}
	}

	return fr
}

//--------------------------------------------------- Decoders ---------------------------------------------------

func (x *XMLDecoder) decodeFunction() (func(writer.Writer, int, bool) (int, error), error) {
	switch x.o.FileType {
	case Artists:
		return x.decodeArtists, nil
	case Labels:
		return x.decodeLabels, nil
	case Masters:
		return x.decodeMasters, nil
	case Releases:
		return x.decodeReleases, nil
	case Unknown:
		fallthrough
	default:
		return nil, wrongTypeSpecified
	}
}

func (x *XMLDecoder) decodeArtists(w writer.Writer, blockSize int, write bool) (int, error) {
	num, a, err := x.Artists(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteArtists(a)
	}

	return num, err
}

func (x *XMLDecoder) decodeLabels(w writer.Writer, blockSize int, write bool) (int, error) {
	num, l, err := x.Labels(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteLabels(l)
	}

	return num, err
}

func (x *XMLDecoder) decodeMasters(w writer.Writer, blockSize int, write bool) (int, error) {
	num, m, err := x.Masters(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteMasters(m)
	}

	return num, err
}

func (x *XMLDecoder) decodeReleases(w writer.Writer, blockSize int, write bool) (int, error) {
	num, r, err := x.Releases(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteReleases(r)
	}

	return num, err
}

//--------------------------------------------------- HELPERS ---------------------------------------------------

func (x *XMLDecoder) startElement(token xml.Token) bool {
	_, ok := token.(xml.StartElement)
	return ok
}

func (x *XMLDecoder) endElement(token xml.Token) bool {
	_, ok := token.(xml.EndElement)
	return ok
}

func (x *XMLDecoder) startElementName(token xml.Token, name string) bool {
	se, ok := token.(xml.StartElement)
	return ok && se.Name.Local == name
}

func (x *XMLDecoder) endElementName(token xml.Token, name string) bool {
	ee, ok := token.(xml.EndElement)
	return ok && ee.Name.Local == name
}

func (x *XMLDecoder) parseValue() string {
	sb := strings.Builder{}
	for {
		t, _ := x.d.Token()
		if x.endElement(t) {
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
		if x.startElementName(t, childName) {
			children = append(children, x.parseValue())
		}
		if x.endElementName(t, parentName) {
			break
		}
	}
	return children
}
