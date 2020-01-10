package write

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONWriter_Options(t *testing.T) {
	j := NewJSONWriter(nil, nil)
	opt := j.Options()

	if opt.ExcludeImages {
		t.Error("exclude images should be false as a default value")
	}
}

func TestJSONWriter_WriteArtist(t *testing.T) {
	b := &strings.Builder{}
	j := NewJSONWriter(b, nil)
	err := j.WriteArtist(artists[0])
	if err != nil {
		t.Error(err)
	}

	ma, _ := json.Marshal(artists[0])
	expected := string(ma)
	get := b.String()
	if expected != get {
		t.Error("json artist differs from json marshal expected solution")
	}
}

func TestJSONWriter_WriteArtists(t *testing.T) {
	b := &strings.Builder{}
	j := NewJSONWriter(b, nil)
	err := j.WriteArtists(artists)
	if err != nil {
		t.Error(err)
	}

	ma, _ := json.Marshal(artists)
	expected := string(ma)
	get := b.String()
	if expected != get {
		t.Error("json artists differ from json marshal expected solution")
	}
}

func TestJSONWriter_WriteLabel(t *testing.T) {
	b := &strings.Builder{}
	j := NewJSONWriter(b, nil)
	err := j.WriteLabel(labels[0])
	if err != nil {
		t.Error(err)
	}

	ma, _ := json.Marshal(labels[0])
	expected := string(ma)
	get := b.String()
	if expected != get {
		t.Error("json label differs from json marshal expected solution")
	}
}

func TestJSONWriter_WriteLabels(t *testing.T) {
	b := &strings.Builder{}
	j := NewJSONWriter(b, nil)
	err := j.WriteLabels(labels)
	if err != nil {
		t.Error(err)
	}

	ma, _ := json.Marshal(labels)
	expected := string(ma)
	get := b.String()
	if expected != get {
		t.Error("json labels differ from json marshal expected solution")
	}
}

func TestJSONWriter_WriteMaster(t *testing.T) {
	b := &strings.Builder{}
	j := NewJSONWriter(b, nil)
	err := j.WriteMaster(masters[0])
	if err != nil {
		t.Error(err)
	}

	ma, _ := json.Marshal(masters[0])
	expected := string(ma)
	get := b.String()
	if expected != get {
		t.Error("json master differs from json marshal expected solution")
	}
}

func TestJSONWriter_WriteMasters(t *testing.T) {
	b := &strings.Builder{}
	j := NewJSONWriter(b, nil)
	err := j.WriteMasters(masters)
	if err != nil {
		t.Error(err)
	}

	ma, _ := json.Marshal(masters)
	expected := string(ma)
	get := b.String()
	if expected != get {
		t.Error("json masters differ from json marshal expected solution")
	}
}

func TestJSONWriter_WriteRelease(t *testing.T) {
	b := &strings.Builder{}
	j := NewJSONWriter(b, nil)
	err := j.WriteRelease(releases[0])
	if err != nil {
		t.Error(err)
	}

	ma, _ := json.Marshal(releases[0])
	expected := string(ma)
	get := b.String()
	if expected != get {
		t.Error("json release differs from json marshal expected solution")
	}
}

func TestJSONWriter_WriteReleases(t *testing.T) {
	b := &strings.Builder{}
	j := NewJSONWriter(b, nil)
	err := j.WriteReleases(releases)
	if err != nil {
		t.Error(err)
	}

	ma, _ := json.Marshal(releases)
	expected := string(ma)
	get := b.String()
	if expected != get {
		t.Error("json releases differ from json marshal expected solution")
	}
}
