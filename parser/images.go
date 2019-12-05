package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseImages(se xml.StartElement, tr xml.TokenReader) []model.Image {
	images := make([]model.Image, 0, 0)
	if se.Name.Local != "images" {
		return nil
	}

	for {
		t, _ := tr.Token()
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "image" {
			images = append(images, ParseImage(se))
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "images" {
			break
		}
	}

	return images
}

func ParseImage(se xml.StartElement) model.Image {
	img := model.Image{}
	if se.Name.Local != "image" {
		return img
	}

	for _, attr := range se.Attr {
		switch attr.Name.Local {
		case "height":
			img.Height = attr.Value
		case "width":
			img.Width = attr.Value
		case "type":
			img.Type = attr.Value
		case "uri":
			img.Uri = attr.Value
		case "uri150":
			img.Uri150 = attr.Value
		}
	}

	return img
}
