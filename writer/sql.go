package writer

import (
	"github.com/Twyer/discogs-parser/model"
	"os"
)

type SqlWriter struct {
	option Options
	f      *os.File
	err    error
}

func NewSqlWriter(fileName string, options ...Options) Writer {
	j := SqlWriter{}

	j.f, j.err = os.Create(fileName)

	if options != nil && len(options) > 0 {
		j.option = options[0]
	}

	return j
}

func (sql SqlWriter) WriteArtist(artist model.Artist) error {
	return nil
}

func (sql SqlWriter) WriteArtists(artists []model.Artist) error {
	return nil
}
func (sql SqlWriter) WriteLabel(label model.Label) error {
	return nil
}
func (sql SqlWriter) WriteLabels(labels []model.Label) error {
	return nil
}
func (sql SqlWriter) WriteMaster(master model.Master) error {
	return nil
}
func (sql SqlWriter) WriteMasters(masters []model.Master) error {
	return nil
}
func (sql SqlWriter) WriteRelease(release model.Release) error {
	return nil
}
func (sql SqlWriter) WriteReleases(releases []model.Release) error {
	return nil
}

func (sql SqlWriter) Close() error {
	return sql.f.Close()
}
