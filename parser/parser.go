package parser

import (
	"encoding/xml"
	"errors"
	"strings"
)

var (
	notStartElement       error = errors.New("token is not start element")
	notCorrectStarElement error = errors.New("token is not a correct start element")
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
