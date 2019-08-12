package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseImage(se xml.StartElement) *model.Image {
	if se.Name.Local != "image" {
		return nil
	}

	img := &model.Image{}

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
