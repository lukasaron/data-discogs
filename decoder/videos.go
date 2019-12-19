package decoder

import (
	"encoding/xml"
	"github.com/Twyer/discogs-parser/model"
)

func (x XMLDecoder) parseVideos() (videos []model.Video) {
	if x.err != nil {
		return videos
	}

	video := model.Video{}

	for {
		t, _ := x.d.Token()
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
				video.Title = x.parseValue()
			case "description":
				video.Description = x.parseValue()
			}
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "video" {
			videos = append(videos, video)
			video = model.Video{}
		}
	}

	return videos
}
