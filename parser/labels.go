package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseLabels(d *xml.Decoder, limit int) ([]model.Label, error) {
	cnt := 0
	labels := make([]model.Label, 0, 0)
	for t, err := d.Token(); t != nil && err == nil && cnt != limit; t, err = d.Token() {
		if IsStartElementName(t, "label") {
			l, err := ParseLabel(t.(xml.StartElement), d)
			if err != nil {
				return labels, err
			}

			labels = append(labels, l)
			cnt++
		}
	}

	return labels, nil
}

func ParseLabel(se xml.StartElement, tr xml.TokenReader) (model.Label, error) {
	label := model.Label{}
	if se.Name.Local != "label" {
		return label, notCorrectStarElement
	}

	for {
		t, _ := tr.Token()
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				imgs, err := ParseImages(se, tr)
				if err != nil {
					return label, err
				}

				label.Images = imgs
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

	return label, nil
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
