package writer

import (
	"encoding/json"
	"github.com/Twyer/discogs-parser/model"
	"os"
)

type JsonWriter struct {
	option  Options
	f       *os.File
	started bool
	err     error
}

func NewJson(fileName string, options ...Options) Writer {
	j := &JsonWriter{
		started: false,
	}

	j.f, j.err = os.Create(fileName)

	if options != nil && len(options) > 0 {
		j.option = options[0]
	}

	if j.err == nil {
		j.err = j.writeInitial()
	}

	return j
}

func (j *JsonWriter) Close() error {
	if j.err != nil {
		return j.err
	}

	_ = j.writeClosing()
	return j.f.Close()
}

func (j *JsonWriter) WriteArtist(a model.Artist) error {
	if j.err != nil {
		return j.err
	}

	if j.option.ExcludeImages {
		a.Images = nil
	}

	j.err = j.delimiterLogic()
	return j.marshalAndWrite(a)
}

func (j *JsonWriter) WriteArtists(as []model.Artist) (err error) {
	if j.err != nil {
		return j.err
	}

	for _, a := range as {
		err = j.WriteArtist(a)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *JsonWriter) WriteLabel(l model.Label) error {
	if j.err != nil {
		return j.err
	}

	if j.option.ExcludeImages {
		l.Images = nil
	}

	j.err = j.delimiterLogic()
	return j.marshalAndWrite(l)
}

func (j *JsonWriter) WriteLabels(ls []model.Label) (err error) {
	if j.err != nil {
		return j.err
	}

	for _, l := range ls {
		err = j.WriteLabel(l)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *JsonWriter) WriteMaster(m model.Master) error {
	if j.err != nil {
		return j.err
	}

	if j.option.ExcludeImages {
		m.Images = nil
	}

	j.err = j.delimiterLogic()
	return j.marshalAndWrite(m)
}

func (j *JsonWriter) WriteMasters(ms []model.Master) (err error) {
	if j.err != nil {
		return j.err
	}

	for _, m := range ms {
		err = j.WriteMaster(m)
		if err != nil {
			return err
		}
	}

	return nil
}
func (j *JsonWriter) WriteRelease(r model.Release) error {
	if j.err != nil {
		return j.err
	}

	if j.option.ExcludeImages {
		r.Images = nil
	}

	j.err = j.delimiterLogic()
	return j.marshalAndWrite(r)
}

func (j *JsonWriter) WriteReleases(rs []model.Release) (err error) {
	if j.err != nil {
		return j.err
	}

	for _, r := range rs {
		err = j.WriteRelease(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *JsonWriter) marshalAndWrite(d interface{}) error {
	if j.err != nil {
		return j.err
	}

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	_, err = j.f.Write(b)
	return err
}

func (j *JsonWriter) writeDelimiter() (err error) {
	_, err = j.f.WriteString(",")
	return err
}

func (j *JsonWriter) delimiterLogic() error {
	if j.started {
		return j.writeDelimiter()
	}

	j.started = true
	return nil
}

func (j *JsonWriter) writeInitial() (err error) {
	_, err = j.f.WriteString("[")
	return err
}

func (j *JsonWriter) writeClosing() (err error) {
	_, err = j.f.WriteString("]")
	return err
}
