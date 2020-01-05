package writer

import (
	"github.com/lukasaron/discogs-parser/decode"
)

type Writer interface {
	WriteArtist(artist decode.Artist) error
	WriteArtists(artists []decode.Artist) error
	WriteLabel(label decode.Label) error
	WriteLabels(labels []decode.Label) error
	WriteMaster(master decode.Master) error
	WriteMasters(masters []decode.Master) error
	WriteRelease(release decode.Release) error
	WriteReleases(releases []decode.Release) error
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
