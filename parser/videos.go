package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseVideos(tr xml.TokenReader) []model.Video {
	videos := make([]model.Video, 0, 0)
	video := model.Video{}
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "videos" {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "video":
				for _, attr := range se.Attr {
					switch attr.Name.Local {
					case "duration":
						video.Duration = attr.Value
					case "embed":
						video.Embed = attr.Value
					case "src":
						video.Src = attr.Value
					}
				}
			case "title":
				video.Title = parseValue(tr)
			case "description":
				video.Description = parseValue(tr)
			}
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "video" {
			videos = append(videos, video)
			video = model.Video{}
		}
	}
	return videos
}
