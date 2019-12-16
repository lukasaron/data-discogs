package decoder

import "github.com/Twyer/discogs-parser/model"

type Decoder interface {
	Close() error
	Artists(limit int) (int, []model.Artist, error)
	Labels(limit int) (int, []model.Label, error)
	Masters(limit int) (int, []model.Master, error)
	Releases(limit int) (int, []model.Release, error)
}

type Options struct {
	QualityLevel QualityLevel
}
