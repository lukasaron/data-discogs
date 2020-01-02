package decoder

import (
	"encoding/xml"
	"github.com/lukasaron/discogs-parser/model"
)

func (x *XMLDecoder) parseImages(se xml.StartElement) (images []model.Image) {
	if x.err != nil {
		return images
	}

	if se.Name.Local != "images" {
		x.err = notCorrectStarElement
		return images
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "image" {
			img := x.parseImage(se)
			if x.err != nil {
				return images
			}

			images = append(images, img)
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "images" {
			break
		}
	}

	return images
}

func (x *XMLDecoder) parseImage(se xml.StartElement) (img model.Image) {
	if x.err != nil {
		return img
	}

	if se.Name.Local != "image" {
		x.err = notCorrectStarElement
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
