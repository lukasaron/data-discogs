package parser

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Twyer/discogs/model"
	"os"
	"strings"
)

var (
	notStartElement error = errors.New("token is not start element")
)

func IsStartElement(token xml.Token) bool {
	_, ok := token.(xml.StartElement)
	return ok
}

func Parse() {
	f, _ := os.Open("/Users/lukas/Downloads/Discogs/artists.xml")
	defer f.Close()
	d := xml.NewDecoder(f)
	artists := make([]model.Artist, 0, 0)
	for t, err := d.Token(); t != nil && err == nil; t, err = d.Token() {
		if IsStartElement(t) {
			element := t.(xml.StartElement)
			if element.Name.Local == "artist" {
				artists = append(artists, ParseArtist(element, d))
			}
		}
	}
	j, _ := json.Marshal(artists)
	fmt.Println(string(j))
}

func parseValue(tr xml.TokenReader) string {
	sb := strings.Builder{}
	for {
		t, _ := tr.Token()
		if _, ok := t.(xml.EndElement); ok {
			break
		}
		if cr, ok := t.(xml.CharData); ok {
			sb.Write(cr)
		}
	}
	return sb.String()
}

func parseChildValue(parentName, childName string, tr xml.TokenReader) []string {
	nv := make([]string, 0, 0)
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == parentName {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == childName {
			nv = append(nv, parseValue(tr))
		}
	}
	return nv
}
