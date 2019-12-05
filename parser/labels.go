package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseLabels(d *xml.Decoder, limit int) []model.Label {
	cnt := 0
	labels := make([]model.Label, 0, 0)
	for t, err := d.Token(); t != nil && err == nil && cnt+1 != limit; t, err = d.Token() {
		if IsStartElementName(t, "label") {
			labels = append(labels, ParseLabel(t.(xml.StartElement), d))
			cnt++
		}
	}

	return labels
}

func ParseLabel(se xml.StartElement, tr xml.TokenReader) model.Label {
	label := model.Label{}
	if se.Name.Local != "label" {
		return label
	}

	for {
		t, _ := tr.Token()
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				label.Images = ParseImages(se, tr)
			case "id":
				label.Id = parseValue(tr)
			case "name":
				label.Name = parseValue(tr)
			case "contactinfo":
				label.ContactInfo = parseValue(tr)
			case "profile":
				label.Profile = parseValue(tr)
			case "sublabels":
				label.SubLabels = parseSubLabels(tr)
			case "data_quality":
				label.DataQuality = parseValue(tr)
			case "parentLabel":
				label.ParentLabel = &model.LabelLabel{
					Id:   se.Attr[0].Value,
					Name: parseValue(tr),
				}
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "label" {
			break
		}
	}

	return label
}

func parseSubLabels(tr xml.TokenReader) []model.LabelLabel {
	labels := make([]model.LabelLabel, 0, 0)
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "sublabels" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "label" {
			label := model.LabelLabel{}
			label.Id = se.Attr[0].Value
			label.Name = parseValue(tr)
			labels = append(labels, label)
		}
	}
	return labels
}
