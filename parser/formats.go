package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseFormats(tr xml.TokenReader) []model.Format {
	formats := make([]model.Format, 0, 0)
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "formats" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "format" {
			format := model.Format{}
			for _, attr := range se.Attr {
				switch attr.Name.Local {
				case "qty":
					format.Quantity = attr.Value
				case "name":
					format.Name = attr.Value
				case "text":
					format.Text = attr.Value
				}
			}

			format.Descriptions = parseChildValues("descriptions", "description", tr)
			formats = append(formats, format)
		}
	}
	return formats
}
