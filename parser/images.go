package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseImages(se xml.StartElement, tr xml.TokenReader) (images []model.Image, err error) {
	if se.Name.Local != "images" {
		return nil, notCorrectStarElement
	}

	for {
		t, _ := tr.Token()
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "image" {
			img, err := ParseImage(se)
			if err != nil {
				return images, err
			}

			images = append(images, img)
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "images" {
			break
		}
	}

	return images, nil
}

func ParseImage(se xml.StartElement) (img model.Image, err error) {
	if se.Name.Local != "image" {
		return img, notCorrectStarElement
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

	return img, nil
}
