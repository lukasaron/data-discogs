package writer

import (
	"github.com/Twyer/discogs-parser/model"
)

type Writer interface {
	WriteArtist(artist model.Artist) error
	WriteArtists(artists []model.Artist) error
	WriteLabel(label model.Label) error
	WriteLabels(labels []model.Label) error
	WriteMaster(master model.Master) error
	WriteMasters(masters []model.Master) error
	WriteRelease(release model.Release) error
	WriteReleases(releases []model.Release) error
	Reset() error
	Close() error
}

type Options struct {
	ExcludeImages bool
}

type Type int

const (
	PostgresType Type = iota
	JsonType
	SqlType
)

func StrToWriterType(str string) (t Type) {
	switch str {
	case "json":
		t = JsonType
	case "postgres":
		t = PostgresType
	case "sql":
		t = SqlType
	default:
		t = SqlType
	}
	return t
}
