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
				label.SubLabels = parseLabelLabels("sublabels", tr)
			case "data_quality":
				label.DataQuality = parseValue(tr)
			case "parentLabel":
				pl := parseLabelLabels("parentLabel", tr)
				if len(pl) > 0 {
					label.ParentLabel = pl[0]
				}
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "label" {
			break
		}
	}

	return label
}

func parseLabelLabels(wrapperName string, tr xml.TokenReader) []*model.LabelLabels {
	labels := make([]*model.LabelLabels, 0, 0)
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == wrapperName {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "label" {
			label := &model.LabelLabels{}
			label.Id = se.Attr[0].Value
			label.Name = parseValue(tr)
			labels = append(labels, label)
		}
	}
	return labels
}
