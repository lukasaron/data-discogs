// Copyright (c) 2020 Lukas Aron. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package write

import (
	"bytes"
	"encoding/json"
	"github.com/lukasaron/data-discogs/model"
	"io"
)

// JSONWriter is one of few provided writers that implements the Writer interface and provides the ability to save
// decoded data directly into JSON output.
type JSONWriter struct {
	o   Options
	w   io.Writer
	b   bytes.Buffer
	err error
}

// NewJSONWriter creates a new Writer instance based on the provided output writer (for instance a file).
// Options with ExcludeImages can be set when we don't want images as part of the final solution.
// When this is not the case and we want images in the result JSON the Option can be omitted.
func NewJSONWriter(output io.Writer, options *Options) Writer {

	if options == nil {
		options = &Options{}
	}

	return &JSONWriter{
		b: bytes.Buffer{},
		o: *options,
		w: output,
	}
}

// Options function returns the current options. It could be useful to get the default options.
func (j JSONWriter) Options() Options {
	return j.o
}

// WriteArtist function writes an artist to the JSON output.
func (j *JSONWriter) WriteArtist(a model.Artist) error {
	if j.o.ExcludeImages {
		a.Images = nil
	}

	j.marshalAndWrite(a)
	j.flush()
	j.clean()

	return j.err
}

// WriteArtists function writes a slice of artists to the JSON output.
func (j *JSONWriter) WriteArtists(artists []model.Artist) error {
	j.writeInitial()

	for _, a := range artists {
		j.writeDelimiter()

		if j.o.ExcludeImages {
			a.Images = nil
		}

		j.marshalAndWrite(a)
		if j.err != nil {
			return j.err
		}
	}

	j.writeClosing()
	j.flush()
	j.clean()

	return j.err
}

// WriteLabel function writes a label to the JSON output.
func (j *JSONWriter) WriteLabel(label model.Label) error {
	if j.o.ExcludeImages {
		label.Images = nil
	}

	j.marshalAndWrite(label)
	j.flush()
	j.clean()

	return j.err
}

// WriteLabels function writes a slice of labels to the JSON output.
func (j *JSONWriter) WriteLabels(labels []model.Label) error {
	j.writeInitial()

	for _, l := range labels {
		j.writeDelimiter()

		if j.o.ExcludeImages {
			l.Images = nil
		}

		j.marshalAndWrite(l)
		if j.err != nil {
			return j.err
		}
	}

	j.writeClosing()
	j.flush()
	j.clean()

	return j.err
}

// WriteMaster function writes a master to the JSON output.
func (j *JSONWriter) WriteMaster(master model.Master) error {
	if j.o.ExcludeImages {
		master.Images = nil
	}

	j.marshalAndWrite(master)
	j.flush()
	j.clean()

	return j.err
}

// WriteMasters function writes a slice of masters to the JSON output.
func (j *JSONWriter) WriteMasters(masters []model.Master) error {
	j.writeInitial()
	for _, m := range masters {
		j.writeDelimiter()

		if j.o.ExcludeImages {
			m.Images = nil
		}

		j.marshalAndWrite(m)
		if j.err != nil {
			return j.err
		}
	}

	j.writeClosing()
	j.flush()
	j.clean()

	return j.err
}

// WriteRelease function writes a release to the JSON output.
func (j *JSONWriter) WriteRelease(release model.Release) error {
	if j.o.ExcludeImages {
		release.Images = nil
	}

	j.marshalAndWrite(release)
	j.flush()
	j.clean()
	return j.err
}

// WriteReleases function writes a slice of releases to the JSON output.
func (j *JSONWriter) WriteReleases(releases []model.Release) error {
	j.writeInitial()

	for _, r := range releases {
		j.writeDelimiter()

		if j.o.ExcludeImages {
			r.Images = nil
		}

		j.marshalAndWrite(r)
		if j.err != nil {
			return j.err
		}
	}

	j.writeClosing()
	j.flush()
	j.clean()

	return j.err
}

// ----------------------------------------------- UNPUBLISHED FUNCTIONS -----------------------------------------------

func (j *JSONWriter) marshalAndWrite(d interface{}) {
	if j.err != nil {
		return
	}

	b, err := json.Marshal(d)
	if err != nil {
		j.err = err
		return
	}

	_, j.err = j.b.Write(b)
}

func (j *JSONWriter) writeDelimiter() {
	if j.err == nil && j.b.Len() > 1 {
		_, j.err = j.b.WriteString(",")
	}
}

func (j *JSONWriter) writeInitial() {
	if j.err != nil {
		return
	}
	_, j.err = j.b.WriteString("[")
}

func (j *JSONWriter) writeClosing() {
	if j.err != nil {
		return
	}

	_, j.err = j.b.WriteString("]")
}

func (j *JSONWriter) flush() {
	if j.err != nil {
		return
	}

	_, j.err = j.w.Write(j.b.Bytes())
}

func (j *JSONWriter) clean() {
	j.b.Reset()
}
