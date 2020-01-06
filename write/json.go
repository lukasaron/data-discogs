package write

import (
	"bytes"
	"encoding/json"
	"github.com/lukasaron/data-discogs/model"
	"os"
)

type JsonWriter struct {
	o   Options
	f   *os.File
	b   bytes.Buffer
	err error
}

func NewJsonWriter(fileName string, options *Options) Writer {
	j := &JsonWriter{
		b: bytes.Buffer{},
	}

	j.f, j.err = os.Create(fileName)

	if options != nil {
		j.o = *options
	}

	return j
}

func (j *JsonWriter) Reset() error {
	j.b.Reset()
	return nil
}

func (j *JsonWriter) Close() error {
	return j.f.Close()
}

func (j JsonWriter) Options() Options {
	return j.o
}

func (j *JsonWriter) WriteArtist(a model.Artist) error {
	if j.o.ExcludeImages {
		a.Images = nil
	}

	j.marshalAndWrite(a)
	j.flush()
	j.clean()

	return j.err
}

func (j *JsonWriter) WriteArtists(artists []model.Artist) error {
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

func (j *JsonWriter) WriteLabel(label model.Label) error {
	if j.o.ExcludeImages {
		label.Images = nil
	}

	j.marshalAndWrite(label)
	j.flush()
	j.clean()

	return j.err
}

func (j *JsonWriter) WriteLabels(labels []model.Label) error {
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

func (j *JsonWriter) WriteMaster(master model.Master) error {
	if j.o.ExcludeImages {
		master.Images = nil
	}

	j.marshalAndWrite(master)
	j.flush()
	j.clean()

	return j.err
}

func (j *JsonWriter) WriteMasters(masters []model.Master) error {
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
func (j *JsonWriter) WriteRelease(release model.Release) error {
	if j.o.ExcludeImages {
		release.Images = nil
	}

	j.marshalAndWrite(release)
	j.flush()
	j.clean()
	return j.err
}

func (j *JsonWriter) WriteReleases(releases []model.Release) error {
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

func (j *JsonWriter) marshalAndWrite(d interface{}) {
	if j.err != nil {
		return
	}

	b, err := json.Marshal(d)
	if err != nil {
		j.err = err
		return
	}

	_, j.err = j.f.Write(b)
}

func (j *JsonWriter) writeDelimiter() {
	if j.err == nil && j.b.Len() > 0 {
		_, j.err = j.b.WriteString(",")
	}
}

func (j *JsonWriter) writeInitial() {
	if j.err != nil {
		return
	}
	_, j.err = j.b.WriteString("[")
}

func (j *JsonWriter) writeClosing() {
	if j.err != nil {
		return
	}

	_, j.err = j.b.WriteString("]")
}

func (j *JsonWriter) flush() {
	if j.err != nil {
		return
	}

	_, j.err = j.f.Write(j.b.Bytes())
}

func (j *JsonWriter) clean() {
	j.b.Reset()
}
