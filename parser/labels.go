package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseLabel(se xml.StartElement, tr xml.TokenReader) model.Label {
	label := model.Label{}
	if se.Name.Local != "label" {
		return label
	}

	for {
		t, _ := tr.Token()
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "image":
				label.Images = append(label.Images, ParseImage(se))
			case "id":
				label.Id = parseValue(tr)
			case "name":
				label.Name = parseValue(tr)
			case "contactinfo":
				label.ContactInfo = parseValue(tr)
			case "profile":
				label.Profile = parseValue(tr)
			case "sublabels":
				label.SubLabels = parseChildValue("sublabels", "label", tr)
			case "data_quality":
				label.DataQuality = parseValue(tr)
			case "parentLabel":
				label.ParentLabel = parseValue(tr)
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "label" {
			break
		}
	}

	return label
}
