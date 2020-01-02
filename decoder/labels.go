package decoder

import (
	"encoding/xml"
	"github.com/lukasaron/discogs-parser/model"
)

func (x *XMLDecoder) parseLabels(limit int) (labels []model.Label) {
	if x.err != nil {
		return labels
	}

	var t xml.Token
	cnt := 0
	for t, x.err = x.d.Token(); t != nil && x.err == nil && cnt != limit; t, x.err = x.d.Token() {
		if x.isStartElementName(t, "label") {
			l := x.parseLabel(t.(xml.StartElement))
			if x.err != nil {
				return labels
			}

			labels = append(labels, l)
			cnt++
		}
	}

	return labels
}

func (x *XMLDecoder) parseLabel(se xml.StartElement) (label model.Label) {
	if x.err != nil {
		return label
	}

	if se.Name.Local != "label" {
		x.err = notCorrectStarElement
		return label
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				imgs := x.parseImages(se)
				if x.err != nil {
					return label
				}

				label.Images = imgs
			case "id":
				label.Id = x.parseValue()
			case "name":
				label.Name = x.parseValue()
			case "contactinfo":
				label.ContactInfo = x.parseValue()
			case "profile":
				label.Profile = x.parseValue()
			case "urls":
				label.Urls = x.parseChildValues("urls", "url")
			case "sublabels":
				label.SubLabels = x.parseSubLabels()
			case "data_quality":
				label.DataQuality = x.parseValue()
			case "parentLabel":
				label.ParentLabel = &model.LabelLabel{
					Id:   se.Attr[0].Value,
					Name: x.parseValue(),
				}
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "label" {
			break
		}
	}

	return label
}

func (x *XMLDecoder) parseSubLabels() (labels []model.LabelLabel) {
	if x.err != nil {
		return labels
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "sublabels" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "label" {
			label := model.LabelLabel{}
			label.Id = se.Attr[0].Value
			label.Name = x.parseValue()
			labels = append(labels, label)
		}
	}

	return labels
}
