package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseLabel(se xml.StartElement, tr xml.TokenReader) model.Label {
	return model.Label{}
}
