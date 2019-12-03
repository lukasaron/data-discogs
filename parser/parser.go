package parser

import (
	"encoding/xml"
	"errors"
	"github.com/Twyer/discogs/model"
	"strings"
)

var (
	notStartElement error = errors.New("token is not start element")
)

func IsStartElement(token xml.Token) bool {
	_, ok := token.(xml.StartElement)
	return ok
}

func IsEndElement(token xml.Token) bool {
	_, ok := token.(xml.EndElement)
	return ok
}

func IsStartElementName(token xml.Token, name string) bool {
	se, ok := token.(xml.StartElement)
	return ok && se.Name.Local == name
}

func IsEndElementName(token xml.Token, name string) bool {
	ee, ok := token.(xml.EndElement)
	return ok && ee.Name.Local == name
}

func ParseArtists(d *xml.Decoder, limit int) []model.Artist {
	cnt := 0
	artists := make([]model.Artist, 0, 0)
	for t, err := d.Token(); t != nil && err == nil; t, err = d.Token() {
		if IsStartElementName(t, "artist") {
			artists = append(artists, ParseArtist(t.(xml.StartElement), d))
			cnt++
		}

		if cnt+1 == limit { // when we have limit 0/negative value it will return everything
			break
		}
	}

	return artists
}

func parseValue(tr xml.TokenReader) string {
	sb := strings.Builder{}
	for {
		t, _ := tr.Token()
		if IsEndElement(t) {
			break
		}

		if cr, ok := t.(xml.CharData); ok {
			sb.Write(cr)
		}
	}
	return sb.String()
}

func parseChildValues(parentName, childName string, tr xml.TokenReader) []string {
	nv := make([]string, 0, 0)
	for {
		t, _ := tr.Token()
		if IsStartElementName(t, childName) {
			nv = append(nv, parseValue(tr))
		}
		if IsEndElementName(t, parentName) {
			break
		}
	}
	return nv
}
