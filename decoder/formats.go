package decoder

import (
	"encoding/xml"
	"github.com/lukasaron/discogs-parser/model"
)

func (x *XMLDecoder) parseFormats() (formats []model.Format) {
	if x.err != nil {
		return formats
	}
	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {

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

			format.Descriptions = x.parseChildValues("descriptions", "description")
			formats = append(formats, format)
		}
	}
	return formats
}
