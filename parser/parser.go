package parser

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
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
	f, _ := os.Open("/Users/aron_lukas/Workspace/Go/src/github.com/Twyer/discogs/data/discogs_20190101_labels.xml")
	defer f.Close()
	d := xml.NewDecoder(f)
	for t, err := d.Token(); t != nil && err == nil; t, err = d.Token() {
		if IsStartElement(t) {
			element := t.(xml.StartElement)
			if element.Name.Local == "label" {
				release := ParseLabel(element, d)
				j, _ := json.Marshal(release)
				fmt.Println(string(j))
			}
		}
	}
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

func parseChildValues(parentName, childName string, tr xml.TokenReader) []string {
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
