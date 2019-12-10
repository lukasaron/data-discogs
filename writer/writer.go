package writer

import (
	"github.com/Twyer/discogs/model"
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
	Close() error
}
