package write

import (
	"bytes"
	"github.com/lukasaron/data-discogs/model"
	"io"
)

type RDFWriter struct {
	o   Options
	w   io.Writer
	b   *bytes.Buffer
	err error
}

func NewRDFWriter(output io.Writer, options *Options) Writer {
	if options == nil {
		options = &Options{}
	}

	return &RDFWriter{
		b: &bytes.Buffer{},
		o: *options,
		w: output,
	}
}

func (r *RDFWriter) WriteArtist(artist model.Artist) error {
	panic("implement me")
}

func (r *RDFWriter) WriteArtists(artists []model.Artist) error {
	panic("implement me")
}

func (r *RDFWriter) WriteLabel(label model.Label) error {
	panic("implement me")
}

func (r *RDFWriter) WriteLabels(labels []model.Label) error {
	panic("implement me")
}

func (r *RDFWriter) WriteMaster(master model.Master) error {
	panic("implement me")
}

func (r *RDFWriter) WriteMasters(masters []model.Master) error {
	panic("implement me")
}

func (r *RDFWriter) WriteRelease(release model.Release) error {
	panic("implement me")
}

func (r *RDFWriter) WriteReleases(releases []model.Release) error {
	panic("implement me")
}

func (r *RDFWriter) Options() Options {
	panic("implement me")
}
