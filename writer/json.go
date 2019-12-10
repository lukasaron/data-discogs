package writer

import (
	"encoding/json"
	"github.com/Twyer/discogs/model"
	"os"
)

type Json struct {
	f   *os.File
	err error
}

func NewJson(fileName string) Writer {
	j := Json{}
	j.f, j.err = os.Create(fileName)
	return j
}

func (j Json) Close() error {
	if j.err != nil {
		return j.err
	}

	return j.f.Close()
}

func (j Json) WriteArtist(a model.Artist) error {
	if j.err != nil {
		return j.err
	}

	return j.marshalAndWrite(a)
}

func (j Json) WriteArtists(as []model.Artist) (err error) {
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

func (j Json) WriteLabel(l model.Label) error {
	if j.err != nil {
		return j.err
	}

	return j.marshalAndWrite(l)
}

func (j Json) WriteLabels(ls []model.Label) (err error) {
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

func (j Json) WriteMaster(m model.Master) error {
	if j.err != nil {
		return j.err
	}

	return j.marshalAndWrite(m)
}

func (j Json) WriteMasters(ms []model.Master) (err error) {
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
func (j Json) WriteRelease(r model.Release) error {
	if j.err != nil {
		return j.err
	}

	return j.marshalAndWrite(r)
}
func (j Json) WriteReleases(rs []model.Release) (err error) {
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

func (j Json) marshalAndWrite(d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	_, err = j.f.Write(b)
	return err
}
