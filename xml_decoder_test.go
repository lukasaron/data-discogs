// Copyright (c) 2020 Lukas Aron. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package discogs

import (
	"github.com/lukasaron/data-discogs/model"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestNewXmlDecoder(t *testing.T) {
	d := NewXMLDecoder(nil, nil)
	readerError := d.Error()
	if readerError == nil {
		t.Error("there should be error when the reader is null")
	}

	err := d.Decode(nil)
	if err != readerError {
		t.Error("there should be the same state error occurring")
	}

	_, _, err = d.Artists()
	if err != readerError {
		t.Error("there should be the same state error occurring")
	}

	_, _, err = d.Labels()
	if err != readerError {
		t.Error("there should be the same state error occurring")
	}

	_, _, err = d.Masters()
	if err != readerError {
		t.Error("there should be the same state error occurring")
	}

	_, _, err = d.Releases()
	if err != readerError {
		t.Error("there should be the same state error occurring")
	}

	_, _, err = d.Artists()
	if err != readerError {
		t.Error("there should be the same state error occurring")
	}
}

func TestNewXmlDecoderNoOptions(t *testing.T) {
	d := NewXMLDecoder(strings.NewReader(artists), nil)
	if d.Error() != nil {
		t.Errorf("there shouldn't be the error: %v", d.Error())
	}

	opt := d.Options()
	if opt.QualityLevel != All {
		t.Error("there should be All quality level")
	}

	if opt.FileType != Unknown {
		t.Error("there should be Unknown file type")
	}

	if opt.Block.ItemSize != defaultBlockSize {
		t.Errorf("block size should be set to default value: %v", defaultBlockSize)
	}

	if opt.Block.Limit != defaultBlockLimit {
		t.Errorf("block limit should be set to default value: %v", defaultBlockLimit)
	}

	if opt.Block.Skip != 0 {
		t.Errorf("block skip should be set to default value: %v", 0)
	}
}

func TestNewXmlDecoderWithOptions(t *testing.T) {
	d := NewXMLDecoder(strings.NewReader(artists), &Options{
		QualityLevel: Correct,
		Block: Block{
			ItemSize: 100,
			Limit:    2,
			Skip:     1,
		},
		FileType: Artists,
	})

	if d.Error() != nil {
		t.Errorf("there should be the error: %v", d.Error())
	}

	opt := d.Options()
	if opt.QualityLevel != Correct {
		t.Error("there should be Correct quality level")
	}

	if opt.FileType != Artists {
		t.Error("there should be Artists file type")
	}

	if opt.Block.ItemSize != 100 {
		t.Errorf("block item size should be set to value: %v", 100)
	}

	if opt.Block.Limit != 2 {
		t.Errorf("block limit should be set to value: %v", 2)
	}

	if opt.Block.Skip != 1 {
		t.Errorf("block skip should be set to value: %v", 1)
	}
}

func TestXMLDecoder_SetOptions(t *testing.T) {
	d := NewXMLDecoder(nil, &Options{
		QualityLevel: NeedsVote,
		Block: Block{
			ItemSize: 0,
			Limit:    -1,
			Skip:     -1,
		},
		FileType: Releases,
	})

	opt := d.Options()
	if opt.QualityLevel != NeedsVote {
		t.Error("there should be Needs Vote quality level")
	}

	if opt.FileType != Releases {
		t.Error("there should be releases file type")
	}

	if opt.Block.ItemSize != defaultBlockSize {
		t.Errorf("block size should be set to default value: %v, when the set value is invalid",
			defaultBlockSize)
	}

	if opt.Block.Limit != defaultBlockLimit {
		t.Errorf("block limit should be set to default value: %v, when the set value is invalid",
			defaultBlockLimit)
	}

	if opt.Block.Skip != 0 {
		t.Errorf("block skip should be set to default value: %v, when the set value is invalid", 0)
	}

	d.SetOptions(Options{
		QualityLevel: Correct,
		Block: Block{
			ItemSize: 10,
			Limit:    1000,
			Skip:     3,
		},
		FileType: Labels,
	})

	opt = d.Options()
	if opt.QualityLevel != Correct {
		t.Error("there should be Correct quality level")
	}

	if opt.FileType != Labels {
		t.Error("there should be Labels file type")
	}

	if opt.Block.ItemSize != 10 {
		t.Errorf("block size should be set to value: %v", 10)
	}

	if opt.Block.Limit != 1000 {
		t.Errorf("block limit should be set to value: %v", 1000)
	}

	if opt.Block.Skip != 3 {
		t.Errorf("block skip should be set to value: %v", 3)
	}
}

func TestXMLDecoder_Artists(t *testing.T) {
	d := NewXMLDecoder(strings.NewReader(artists), nil)
	num, a, err := d.Artists()
	if num != 2 || len(a) != 2 {
		t.Error("there should be 2 artists parsed")
	}

	if err == nil {
		t.Error("expecting EOF error")
	}

	if err != io.EOF {
		t.Errorf("there should be EOF error instead of %v", err)
	}
}

func TestXMLDecoder_Artists_First(t *testing.T) {
	expectedArtist := model.Artist{
		ID:       "1",
		Name:     "The Persuader",
		RealName: "Jesper Dahlbäck",
		Images: []model.Image{
			{
				Height: "450",
				Width:  "600",
				Type:   "primary",
			},
			{
				Height: "771",
				Width:  "600",
				Type:   "secondary",
			},
		},
		DataQuality:    "Needs Vote",
		NameVariations: []string{"Persuader", "The Presuader"},
		Urls:           []string{"https://en.wikipedia.org/wiki/Jesper_Dahlbäck"},
		Aliases: []model.Alias{
			{
				ID:   "239",
				Name: "Jesper Dahlbäck",
			},
			{
				ID:   "16055",
				Name: "Groove Machine",
			},
			{
				ID:   "19541",
				Name: "Dick Track",
			},
			{
				ID:   "25227",
				Name: "Lenk",
			},
			{
				ID:   "196957",
				Name: "Janne Me' Amazonen",
			},
			{
				ID:   "278760",
				Name: "Faxid",
			},
			{
				ID:   "439150",
				Name: "The Pinguin Man",
			},
		},
	}

	d := NewXMLDecoder(strings.NewReader(artists), nil)
	num, a, err := d.Artists()
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if num != 2 {
		t.Error("wrong number of parsed artists")
	}

	if !reflect.DeepEqual(a[0], expectedArtist) {
		t.Error("expected artist differs from parsed artist")
	}
}

func TestXMLDecoder_Artists_Second(t *testing.T) {
	expectedArtist := model.Artist{
		ID:          "2",
		Name:        "Mr. James Barth & A.D.",
		RealName:    "Cari Lekebusch & Alexi Delano",
		DataQuality: "Correct",
		NameVariations: []string{"Mr Barth & A.D.", "MR JAMES BARTH & A. D.",
			"Mr. Barth & A.D.", "Mr. James Barth & A. D."},
		Aliases: []model.Alias{
			{
				ID:   "2470",
				Name: "Puente Latino",
			},
			{
				ID:   "19536",
				Name: "Yakari & Delano",
			},
			{
				ID:   "103709",
				Name: "Crushed Insect & The Sick Puppy",
			},
			{
				ID:   "384581",
				Name: "ADCL",
			},
			{
				ID:   "1779857",
				Name: "Alexi Delano & Cari Lekebusch",
			},
		},
		Members: []model.Member{
			{
				ID:   "26",
				Name: "Alexi Delano",
			},
			{
				ID:   "27",
				Name: "Cari Lekebusch",
			},
		},
	}

	d := NewXMLDecoder(strings.NewReader(artists), nil)
	num, a, err := d.Artists()
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if num != 2 {
		t.Error("wrong number of parsed artists")
	}

	if !reflect.DeepEqual(a[1], expectedArtist) {
		t.Error("expected artist differs from parsed artist")
	}
}

func TestXMLDecoder_Artists_Block_ItemSize(t *testing.T) {
	d := NewXMLDecoder(strings.NewReader(artists), &Options{
		Block: Block{
			ItemSize: 1,
		},
	})

	num, a, err := d.Artists()
	if err != nil {
		t.Errorf("no error expected when there are still some data to process, got %v", err)

	}

	if num != 1 || len(a) != 1 {
		t.Error("there should be 1 artist parsed")
	}

	num, a, err = d.Artists()
	if err != nil {
		t.Errorf("no error expected when there are still some data to process, got %v", err)

	}

	if num != 1 || len(a) != 1 {
		t.Error("there should be 1 artist parsed")
	}

	num, a, err = d.Artists()
	if num != 0 || len(a) != 0 {
		t.Error("there shouldn't be any artist parsed")
	}

	if err == nil {
		t.Error("EOF  error expected when there is nothing else to parse")
	}
}

func TestXMLDecoder_Labels(t *testing.T) {
	d := NewXMLDecoder(strings.NewReader(labels), nil)

	num, l, err := d.Labels()
	if num != 2 || len(l) != 2 {
		t.Error("there should be 2 labels parsed")
	}

	if err == nil {
		t.Error("expecting EOF error")
	}

	if err != io.EOF {
		t.Errorf("there should be EOF error instead of %v", err)
	}
}

func TestXMLDecoder_Labels_First(t *testing.T) {
	expectedLabel := model.Label{
		ID:   "1",
		Name: "Planet E",
		Images: []model.Image{
			{
				Height: "24",
				Width:  "132",
				Type:   "primary",
			},
			{
				Height: "126",
				Width:  "587",
				Type:   "secondary",
			},
			{
				Height: "196",
				Width:  "600",
				Type:   "secondary",
			},
			{
				Height: "121",
				Width:  "275",
				Type:   "secondary",
			},
			{
				Height: "720",
				Width:  "382",
				Type:   "secondary",
			},
			{
				Height: "398",
				Width:  "500",
				Type:   "secondary",
			},
			{
				Height: "189",
				Width:  "600",
				Type:   "secondary",
			},
		},
		ContactInfo: "Planet E Communications",
		Profile:     "[a=Carl Craig]'s classic techno label founded in 1991.",
		DataQuality: "Correct",
		Urls: []string{"http://planet-e.net", "http://planetecommunications.bandcamp.com",
			"http://www.facebook.com/planetedetroit", "http://www.flickr.com/photos/planetedetroit",
			"http://plus.google.com/100841702106447505236", "http://www.instagram.com/carlcraignet",
			"http://myspace.com/planetecom", "http://myspace.com/planetedetroit",
			"http://soundcloud.com/planetedetroit", "http://twitter.com/planetedetroit", "http://vimeo.com/user1265384",
			"http://en.wikipedia.org/wiki/Planet_E_Communications", "http://www.youtube.com/user/planetedetroit"},
		SubLabels: []model.LabelLabel{
			{
				ID:   "86537",
				Name: "Antidote (4)",
			},
			{
				ID:   "41841",
				Name: "Community Projects",
			},
			{
				ID:   "153760",
				Name: "Guilty Pleasures",
			},
			{
				ID:   "31405",
				Name: "I Ner Zon Sounds",
			},
			{
				ID:   "277579",
				Name: "Planet E Communications",
			},
			{
				ID:   "294738",
				Name: "Planet E Communications, Inc.",
			},
			{
				ID:   "1560615",
				Name: "Planet E Productions",
			},
			{
				ID:   "488315",
				Name: "TWPENTY",
			},
		},
	}

	d := NewXMLDecoder(strings.NewReader(labels), nil)
	num, l, err := d.Labels()
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if num != 2 {
		t.Error("wrong number of parsed labels")
	}

	if !reflect.DeepEqual(l[0], expectedLabel) {
		t.Error("expected label differs from parsed label")
	}

}

func TestXMLDecoder_Labels_Second(t *testing.T) {
	expectedLabel := model.Label{
		ID:          "2",
		Name:        "Earthtones Recordings",
		ContactInfo: "Seasons Recordings 2236 Pacific Avenue Suite D",
		Profile: "California deep house label founded by [a=Jamie Thinnes]. " +
			"Now defunct and continued as [l=Seasons Recordings].",
		DataQuality: "Correct",
		Urls:        []string{"http://www.seasonsrecordings.com/"},
	}

	d := NewXMLDecoder(strings.NewReader(labels), nil)
	num, l, err := d.Labels()
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if num != 2 {
		t.Error("wrong number of parsed labels")
	}

	if !reflect.DeepEqual(l[1], expectedLabel) {
		t.Error("expected label differs from parsed label")
	}
}

func TestXMLDecoder_Labels_Block_ItemSize(t *testing.T) {
	d := NewXMLDecoder(strings.NewReader(labels), &Options{
		Block: Block{
			ItemSize: 1,
		},
	})

	num, l, err := d.Labels()
	if num != 1 || len(l) != 1 {
		t.Error("there should be 1 label parsed")
	}

	if err != nil {
		t.Errorf("no error expected when there are still some data to process, got %v", err)
	}

	num, l, err = d.Labels()
	if num != 1 || len(l) != 1 {
		t.Error("there should be 1 label parsed")
	}

	if err != nil {
		t.Errorf("no error expected when there are still some data to process, got %v", err)
	}

	num, l, err = d.Labels()
	if num != 0 || len(l) != 0 {
		t.Error("there shouldn't be any label parsed")
	}

	if err == nil {
		t.Error("EOF  error expected when there is nothing else to parse")
	}
}

func TestXMLDecoder_Masters(t *testing.T) {
	d := NewXMLDecoder(strings.NewReader(masters), nil)

	num, m, err := d.Masters()
	if num != 2 || len(m) != 2 {
		t.Error("there should be 2 masters parsed")
	}

	if err == nil {
		t.Error("expecting EOF error")
	}

	if err != io.EOF {
		t.Errorf("there should be EOF error instead of %v", err)
	}
}

func TestXMLDecoder_Masters_Block_ItemSize(t *testing.T) {
	d := NewXMLDecoder(strings.NewReader(masters), &Options{
		Block: Block{
			ItemSize: 1,
		},
	})

	num, m, err := d.Masters()
	if num != 1 || len(m) != 1 {
		t.Error("there should be 1 master parsed")
	}

	if err != nil {
		t.Errorf("no error expected when there are still some data to process, got %v", err)
	}

	num, m, err = d.Masters()
	if num != 1 || len(m) != 1 {
		t.Error("there should be 1 master parsed")
	}

	if err != nil {
		t.Errorf("no error expected when there are still some data to process, got %v", err)
	}

	num, m, err = d.Masters()
	if num != 0 || len(m) != 0 {
		t.Error("there shouldn't be any masters parsed")
	}

	if err == nil {
		t.Error("EOF  error expected when there is nothing else to parse")
	}
}

func TestXMLDecoder_Masters_First(t *testing.T) {
	expectedMaster := model.Master{
		ID:          "18500",
		MainRelease: "155102",
		Images: []model.Image{
			{
				Height: "588",
				Width:  "600",
				Type:   "primary",
			},
		},
		Artists: []model.ReleaseArtist{
			{
				ID:   "212070",
				Name: "Samuel L Session",
				Anv:  "Samuel L",
			},
		},
		Genres:      []string{"Electronic"},
		Styles:      []string{"Techno"},
		Year:        "2001",
		Title:       "New Soil",
		DataQuality: "Correct",
		Videos: []model.Video{
			{
				Duration:    "489",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=f05Ai921itM",
				Title:       "Samuel L - Velvet",
				Description: "Samuel L - Velvet",
			},
			{
				Duration:    "348",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=v23rSPG_StA",
				Title:       "Samuel L - Danses D'Afrique",
				Description: "Samuel L - Danses D'Afrique",
			},
			{
				Duration:    "288",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=tHo82ha6p40",
				Title:       "Samuel L - Body N' Soul",
				Description: "Samuel L - Body N' Soul",
			},
			{
				Duration:    "331",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=KDcqzHca5dk",
				Title:       "Samuel L - Into The Groove",
				Description: "Samuel L - Into The Groove",
			},
			{
				Duration:    "334",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=3DIYjJFl8Dk",
				Title:       "Samuel L - Soul Syndrome",
				Description: "Samuel L - Soul Syndrome",
			},
			{
				Duration:    "325",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=_o8yZMPqvNg",
				Title:       "Samuel L - Lush",
				Description: "Samuel L - Lush",
			},
			{
				Duration:    "346",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=JPwwJSc_-30",
				Title:       "Samuel L - Velvet ( Direct Me )",
				Description: "Samuel L - Velvet ( Direct Me )",
			},
		},
	}

	d := NewXMLDecoder(strings.NewReader(masters), nil)
	num, m, err := d.Masters()
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if num != 2 {
		t.Error("there should be 2 masters decoded")
	}

	if !reflect.DeepEqual(m[0], expectedMaster) {
		t.Error("expected master differs from parsed master")
	}
}

func TestXMLDecoder_Masters_Second(t *testing.T) {
	expectedMaster := model.Master{
		ID:          "18512",
		MainRelease: "33699",
		Images: []model.Image{
			{
				Height: "150",
				Width:  "150",
				Type:   "primary",
			},
			{
				Height: "592",
				Width:  "600",
				Type:   "secondary",
			},
			{
				Height: "592",
				Width:  "600",
				Type:   "secondary",
			},
		},
		Artists: []model.ReleaseArtist{
			{
				ID:   "212070",
				Name: "Samuel L Session",
			},
		},
		Genres:      []string{"Electronic"},
		Styles:      []string{"Tribal", "Techno"},
		Year:        "2002",
		Title:       "Psyche EP",
		DataQuality: "Correct",
		Videos: []model.Video{
			{
				Duration:    "118",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=QYf4j0Pd2FU",
				Title:       "Samuel L. Session - Arrival",
				Description: "Samuel L. Session - Arrival",
			},
			{
				Duration:    "376",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=c_AfLqTdncI",
				Title:       "Samuel L. Session - Psyche Part 1",
				Description: "Samuel L. Session - Psyche Part 1",
			},
			{
				Duration:    "419",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=0nxvR8Zl9wY",
				Title:       "Samuel L. Session - Psyche Part 2",
				Description: "Samuel L. Session - Psyche Part 2",
			},
		},
	}

	d := NewXMLDecoder(strings.NewReader(masters), nil)
	num, m, err := d.Masters()
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if num != 2 {
		t.Error("there should be 2 masters decoded")
	}

	if !reflect.DeepEqual(m[1], expectedMaster) {
		t.Error("expected master differs from parsed master")
	}
}

func TestXMLDecoder_Releases(t *testing.T) {
	d := NewXMLDecoder(strings.NewReader(releases), nil)

	num, r, err := d.Releases()
	if num != 2 || len(r) != 2 {
		t.Error("there should be 2 releases parsed")
	}

	if err == nil {
		t.Error("expecting EOF error")
	}

	if err != io.EOF {
		t.Errorf("there should be EOF error instead of %v", err)
	}
}

func TestXMLDecoder_Releases_Block_ItemSize(t *testing.T) {
	d := NewXMLDecoder(strings.NewReader(releases), &Options{
		Block: Block{
			ItemSize: 1,
		},
	})

	num, r, err := d.Releases()
	if num != 1 || len(r) != 1 {
		t.Error("there should be 1 release parsed")
	}

	if err != nil {
		t.Errorf("no error expected when there are still some data to process, got %v", err)
	}

	num, r, err = d.Releases()
	if num != 1 || len(r) != 1 {
		t.Error("there should be 1 release parsed")
	}

	if err != nil {
		t.Errorf("no error expected when there are still some data to process, got %v", err)
	}

	num, r, err = d.Releases()
	if num != 0 || len(r) != 0 {
		t.Error("there shouldn't be any releases parsed")
	}

	if err == nil {
		t.Error("EOF  error expected when there is nothing else to parse")
	}
}

func TestXMLDecoder_Releases_First(t *testing.T) {
	expectedRelease := model.Release{
		ID:     "1",
		Status: "Accepted",
		Images: []model.Image{
			{
				Height: "600",
				Width:  "600",
				Type:   "primary",
			},
			{
				Height: "600",
				Width:  "600",
				Type:   "secondary",
			},
			{
				Height: "600",
				Width:  "600",
				Type:   "secondary",
			},
			{
				Height: "600",
				Width:  "600",
				Type:   "secondary",
			},
		},
		Artists: []model.ReleaseArtist{
			{
				ID:   "1",
				Name: "The Persuader",
			},
		},
		ExtraArtists: []model.ReleaseArtist{
			{
				ID:   "239",
				Name: "Jesper Dahlbäck",
				Role: "Music By [All Tracks By]",
			},
		},
		Title: "Stockholm",
		Formats: []model.Format{
			{
				Name:         "Vinyl",
				Quantity:     "2",
				Descriptions: []string{"12\"", "33 ⅓ RPM"},
			},
		},
		Genres:   []string{"Electronic"},
		Styles:   []string{"Deep House"},
		Country:  "Sweden",
		Released: "1999-03-00",
		Notes: "The song titles are the names of six of Stockholm's 82 districts.\n\n" +
			"Title on label: - Stockholm -\n\nRecorded at the Globe Studio, Stockholm\n\nFAX: +46 8 679 64 53",
		DataQuality: "Needs Vote",
		TrackList: []model.Track{
			{
				Position: "A",
				Title:    "Östermalm",
				Duration: "4:45",
			},
			{
				Position: "B1",
				Title:    "Vasastaden",
				Duration: "6:11",
			},
			{
				Position: "B2",
				Title:    "Kungsholmen",
				Duration: "2:49",
			},
			{
				Position: "C1",
				Title:    "Södermalm",
				Duration: "5:38",
			},
			{
				Position: "C2",
				Title:    "Norrmalm",
				Duration: "4:52",
			},
			{
				Position: "D",
				Title:    "Gamla Stan",
				Duration: "5:16",
			},
		},
		Identifiers: []model.Identifier{
			{
				Description: "A-Side Runout",
				Type:        "Matrix / Runout",
				Value:       "MPO SK 032 A1",
			},
			{
				Description: "B-Side Runout",
				Type:        "Matrix / Runout",
				Value:       "MPO SK 032 B1",
			},
			{
				Description: "C-Side Runout",
				Type:        "Matrix / Runout",
				Value:       "MPO SK 032 C1",
			},
			{
				Description: "D-Side Runout",
				Type:        "Matrix / Runout",
				Value:       "MPO SK 032 D1",
			},
			{
				Description: "Only On A-Side Runout",
				Type:        "Matrix / Runout",
				Value:       "G PHRUPMASTERGENERAL T27 LONDON",
			},
		},
		Videos: []model.Video{
			{
				Duration:    "296",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=MpmbntGDyNE",
				Title:       "The Persuader - Östermalm",
				Description: "The Persuader - Östermalm",
			},
			{
				Duration:    "376",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=Cawyll0pOI4",
				Title:       "The Persuader - Vasastaden",
				Description: "The Persuader - Vasastaden",
			},
			{
				Duration:    "176",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=XExCZfMCXdo",
				Title:       "The Persuader - Kungsholmen",
				Description: "The Persuader - Kungsholmen",
			},
			{
				Duration:    "341",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=WDZqiENap_U",
				Title:       "The Persuader - Södermalm",
				Description: "The Persuader - Södermalm",
			},
			{
				Duration:    "301",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=EBBHR3EMN50",
				Title:       "The Persuader - Norrmalm",
				Description: "The Persuader - Norrmalm",
			},
			{
				Duration:    "326",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=afMHNll9EVM",
				Title:       "The Persuader - Gamla Stan",
				Description: "The Persuader - Gamla Stan",
			},
		},
		Labels: []model.ReleaseLabel{
			{
				ID:       "5",
				Name:     "Svek",
				Category: "SK032",
			},
		},
		Companies: []model.Company{
			{
				ID:             "271046",
				Name:           "The Globe Studios",
				EntityType:     "23",
				EntityTypeName: "Recorded At",
				ResourceURL:    "https://api.discogs.com/labels/271046",
			},
			{
				ID:             "56025",
				Name:           "MPO",
				EntityType:     "17",
				EntityTypeName: "Pressed By",
				ResourceURL:    "https://api.discogs.com/labels/56025",
			},
		},
	}

	d := NewXMLDecoder(strings.NewReader(releases), nil)
	num, r, err := d.Releases()
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if num != 2 {
		t.Error("wrong number of decoded releases")
	}

	if !reflect.DeepEqual(r[0], expectedRelease) {
		t.Error("expected release differs from parsed release")
	}
}

func TestXMLDecoder_Releases_Second(t *testing.T) {
	expectedRelease := model.Release{
		ID:     "2",
		Status: "Accepted",
		Images: []model.Image{
			{
				Height: "394",
				Width:  "400",
				Type:   "primary",
			},
			{
				Height: "600",
				Width:  "600",
				Type:   "secondary",
			},
			{
				Height: "600",
				Width:  "600",
				Type:   "secondary",
			},
		},
		Artists: []model.ReleaseArtist{
			{
				ID:   "2",
				Name: "Mr. James Barth & A.D.",
			},
		},
		ExtraArtists: []model.ReleaseArtist{
			{
				ID:   "26",
				Name: "Alexi Delano",
				Role: "Producer, Recorded By",
			},
			{
				ID:   "27",
				Name: "Cari Lekebusch",
				Role: "Producer, Recorded By",
			},
			{
				ID:   "26",
				Name: "Alexi Delano",
				Anv:  "A. Delano",
				Role: "Written-By",
			},
			{
				ID:   "27",
				Name: "Cari Lekebusch",
				Anv:  "C. Lekebusch",
				Role: "Written-By",
			},
		},
		Title: "Knockin' Boots Vol 2 Of 2",
		Formats: []model.Format{
			{
				Name:         "Vinyl",
				Quantity:     "1",
				Descriptions: []string{"12\"", "33 ⅓ RPM"},
			},
		},
		Genres:      []string{"Electronic"},
		Styles:      []string{"Broken Beat", "Techno", "Tech House"},
		Country:     "Sweden",
		Released:    "1998-06-00",
		Notes:       "All joints recorded in NYC (Dec.97).",
		DataQuality: "Correct",
		MasterID:    "713738",
		MainRelease: "true",
		TrackList: []model.Track{
			{
				Position: "A1",
				Title:    "A Sea Apart",
				Duration: "5:08",
			},
			{
				Position: "A2",
				Title:    "Dutchmaster",
				Duration: "4:21",
			},
			{
				Position: "B1",
				Title:    "Inner City Lullaby",
				Duration: "4:22",
			},
			{
				Position: "B2",
				Title:    "Yeah Kid!",
				Duration: "4:46",
			},
		},
		Identifiers: []model.Identifier{
			{
				Description: "Side A Runout Etching",
				Type:        "Matrix / Runout",
				Value:       "MPO SK026-A -J.T.S.-",
			},
			{
				Description: "Side B Runout Etching",
				Type:        "Matrix / Runout",
				Value:       "MPO SK026-B -J.T.S.-",
			},
		},
		Videos: []model.Video{
			{
				Duration:    "310",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=MIgQNVhYILA",
				Title:       "Mr. James Barth & A.D. - A Sea Apart",
				Description: "Mr. James Barth & A.D. - A Sea Apart",
			},
			{
				Duration:    "265",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=LgLchSRehhc",
				Title:       "Mr. James Barth & A.D. - Dutchmaster",
				Description: "Mr. James Barth & A.D. - Dutchmaster",
			},
			{
				Duration:    "260",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=iaqHaULlqqg",
				Title:       "Mr. James Barth & A.D. - Inner City Lullaby",
				Description: "Mr. James Barth & A.D. - Inner City Lullaby",
			},
			{
				Duration:    "290",
				Embed:       "true",
				Src:         "https://www.youtube.com/watch?v=x_Os7b-iWKs",
				Title:       "Mr. James Barth & A.D. - Yeah Kid!",
				Description: "Mr. James Barth & A.D. - Yeah Kid!",
			},
		},
		Labels: []model.ReleaseLabel{
			{
				ID:       "5",
				Name:     "Svek",
				Category: "SK 026",
			},
			{
				ID:       "5",
				Name:     "Svek",
				Category: "SK026",
			},
		},
		Companies: []model.Company{
			{
				ID:             "266169",
				Name:           "JTS Studios",
				EntityType:     "29",
				EntityTypeName: "Mastered At",
				ResourceURL:    "https://api.discogs.com/labels/266169",
			},
			{
				ID:             "56025",
				Name:           "MPO",
				EntityType:     "17",
				EntityTypeName: "Pressed By",
				ResourceURL:    "https://api.discogs.com/labels/56025",
			},
		},
	}

	d := NewXMLDecoder(strings.NewReader(releases), nil)
	num, r, err := d.Releases()
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if num != 2 {
		t.Error("wrong number of decoded releases")
	}

	if !reflect.DeepEqual(r[1], expectedRelease) {
		t.Error("expected release differs from parsed release")
	}
}

// ------------------------------------------------------- DATA -------------------------------------------------------

var artists = `
<artists>
    <artist>
        <images>
            <image height="450" type="primary" uri="" uri150="" width="600" />
            <image height="771" type="secondary" uri="" uri150="" width="600" />
        </images>
        <id>1</id>
        <name>The Persuader</name>
        <realname>Jesper Dahlbäck</realname>
        <profile></profile>
        <data_quality>Needs Vote</data_quality>
        <urls>
            <url>https://en.wikipedia.org/wiki/Jesper_Dahlbäck</url>
        </urls>
        <namevariations>
            <name>Persuader</name>
            <name>The Presuader</name>
        </namevariations>
        <aliases>
            <name id="239">Jesper Dahlbäck</name>
            <name id="16055">Groove Machine</name>
            <name id="19541">Dick Track</name>
            <name id="25227">Lenk</name>
            <name id="196957">Janne Me' Amazonen</name>
            <name id="278760">Faxid</name>
            <name id="439150">The Pinguin Man</name>
        </aliases>
    </artist>
    <artist>
        <id>2</id>
        <name>Mr. James Barth &amp; A.D.</name>
        <realname>Cari Lekebusch &amp; Alexi Delano</realname>
        <profile></profile>
        <data_quality>Correct</data_quality>
        <namevariations>
            <name>Mr Barth &amp; A.D.</name>
            <name>MR JAMES BARTH &amp; A. D.</name>
            <name>Mr. Barth &amp; A.D.</name>
            <name>Mr. James Barth &amp; A. D.</name>
        </namevariations>
        <aliases>
            <name id="2470">Puente Latino</name>
            <name id="19536">Yakari &amp; Delano</name>
            <name id="103709">Crushed Insect &amp; The Sick Puppy</name>
            <name id="384581">ADCL</name>
            <name id="1779857">Alexi Delano &amp; Cari Lekebusch</name>
        </aliases>
        <members>
            <id>26</id>
            <name id="26">Alexi Delano</name>
            <id>27</id>
            <name id="27">Cari Lekebusch</name>
        </members>
    </artist>
</artists>`

var labels = `
<labels>
    <label>
        <images>
            <image height="24" type="primary" uri="" uri150="" width="132" />
            <image height="126" type="secondary" uri="" uri150="" width="587" />
            <image height="196" type="secondary" uri="" uri150="" width="600" />
            <image height="121" type="secondary" uri="" uri150="" width="275" />
            <image height="720" type="secondary" uri="" uri150="" width="382" />
            <image height="398" type="secondary" uri="" uri150="" width="500" />
            <image height="189" type="secondary" uri="" uri150="" width="600" />
        </images>
        <id>1</id>
        <name>Planet E</name>
        <contactinfo>Planet E Communications</contactinfo>
        <profile>[a=Carl Craig]'s classic techno label founded in 1991.</profile>
        <data_quality>Correct</data_quality>
        <urls>
            <url>http://planet-e.net</url>
            <url>http://planetecommunications.bandcamp.com</url>
            <url>http://www.facebook.com/planetedetroit</url>
            <url>http://www.flickr.com/photos/planetedetroit</url>
            <url>http://plus.google.com/100841702106447505236</url>
            <url>http://www.instagram.com/carlcraignet</url>
            <url>http://myspace.com/planetecom</url>
            <url>http://myspace.com/planetedetroit</url>
            <url>http://soundcloud.com/planetedetroit</url>
            <url>http://twitter.com/planetedetroit</url>
            <url>http://vimeo.com/user1265384</url>
            <url>http://en.wikipedia.org/wiki/Planet_E_Communications</url>
            <url>http://www.youtube.com/user/planetedetroit</url>
        </urls>
        <sublabels>
            <label id="86537">Antidote (4)</label>
            <label id="41841">Community Projects</label>
            <label id="153760">Guilty Pleasures</label>
            <label id="31405">I Ner Zon Sounds</label>
            <label id="277579">Planet E Communications</label>
            <label id="294738">Planet E Communications, Inc.</label>
            <label id="1560615">Planet E Productions</label>
            <label id="488315">TWPENTY</label>
        </sublabels>
    </label>
    <label>
        <id>2</id>
        <name>Earthtones Recordings</name>
        <contactinfo>Seasons Recordings 2236 Pacific Avenue Suite D</contactinfo>
        <profile>California deep house label founded by [a=Jamie Thinnes]. Now defunct and continued as [l=Seasons Recordings].</profile>
        <data_quality>Correct</data_quality>
        <urls>
            <url>http://www.seasonsrecordings.com/</url>
        </urls>
    </label>
</labels>`

var masters = `
<masters>
    <master id="18500">
        <main_release>155102</main_release>
        <images>
            <image height="588" type="primary" uri="" uri150="" width="600" />
        </images>
        <artists>
            <artist>
                <id>212070</id>
                <name>Samuel L Session</name>
                <anv>Samuel L</anv>
                <join></join>
                <role></role>
                <tracks></tracks>
            </artist>
        </artists>
        <genres>
            <genre>Electronic</genre>
        </genres>
        <styles>
            <style>Techno</style>
        </styles>
        <year>2001</year>
        <title>New Soil</title>
        <data_quality>Correct</data_quality>
        <videos>
            <video duration="489" embed="true" src="https://www.youtube.com/watch?v=f05Ai921itM">
                <title>Samuel L - Velvet</title>
                <description>Samuel L - Velvet</description>
            </video>
            <video duration="348" embed="true" src="https://www.youtube.com/watch?v=v23rSPG_StA">
                <title>Samuel L - Danses D'Afrique</title>
                <description>Samuel L - Danses D'Afrique</description>
            </video>
            <video duration="288" embed="true" src="https://www.youtube.com/watch?v=tHo82ha6p40">
                <title>Samuel L - Body N' Soul</title>
                <description>Samuel L - Body N' Soul</description>
            </video>
            <video duration="331" embed="true" src="https://www.youtube.com/watch?v=KDcqzHca5dk">
                <title>Samuel L - Into The Groove</title>
                <description>Samuel L - Into The Groove</description>
            </video>
            <video duration="334" embed="true" src="https://www.youtube.com/watch?v=3DIYjJFl8Dk">
                <title>Samuel L - Soul Syndrome</title>
                <description>Samuel L - Soul Syndrome</description>
            </video>
            <video duration="325" embed="true" src="https://www.youtube.com/watch?v=_o8yZMPqvNg">
                <title>Samuel L - Lush</title>
                <description>Samuel L - Lush</description>
            </video>
            <video duration="346" embed="true" src="https://www.youtube.com/watch?v=JPwwJSc_-30">
                <title>Samuel L - Velvet ( Direct Me )</title>
                <description>Samuel L - Velvet ( Direct Me )</description>
            </video>
        </videos>
    </master>
    <master id="18512">
        <main_release>33699</main_release>
        <images>
            <image height="150" type="primary" uri="" uri150="" width="150" />
            <image height="592" type="secondary" uri="" uri150="" width="600" />
            <image height="592" type="secondary" uri="" uri150="" width="600" />
        </images>
        <artists>
            <artist>
                <id>212070</id>
                <name>Samuel L Session</name>
                <anv></anv>
                <join></join>
                <role></role>
                <tracks></tracks>
            </artist>
        </artists>
        <genres>
            <genre>Electronic</genre>
        </genres>
        <styles>
            <style>Tribal</style>
            <style>Techno</style>
        </styles>
        <year>2002</year>
        <title>Psyche EP</title>
        <data_quality>Correct</data_quality>
        <videos>
            <video duration="118" embed="true" src="https://www.youtube.com/watch?v=QYf4j0Pd2FU">
                <title>Samuel L. Session - Arrival</title>
                <description>Samuel L. Session - Arrival</description>
            </video>
            <video duration="376" embed="true" src="https://www.youtube.com/watch?v=c_AfLqTdncI">
                <title>Samuel L. Session - Psyche Part 1</title>
                <description>Samuel L. Session - Psyche Part 1</description>
            </video>
            <video duration="419" embed="true" src="https://www.youtube.com/watch?v=0nxvR8Zl9wY">
                <title>Samuel L. Session - Psyche Part 2</title>
                <description>Samuel L. Session - Psyche Part 2</description>
            </video>
        </videos>
    </master>
</masters>`

var releases = `
<releases>
    <release id="1" status="Accepted">
        <images>
            <image height="600" type="primary" uri="" uri150="" width="600" />
            <image height="600" type="secondary" uri="" uri150="" width="600" />
            <image height="600" type="secondary" uri="" uri150="" width="600" />
            <image height="600" type="secondary" uri="" uri150="" width="600" />
        </images>
        <artists>
            <artist>
                <id>1</id>
                <name>The Persuader</name>
                <anv></anv>
                <join></join>
                <role></role>
                <tracks></tracks>
            </artist>
        </artists>
        <title>Stockholm</title>
        <labels>
            <label catno="SK032" id="5" name="Svek" />
        </labels>
        <extraartists>
            <artist>
                <id>239</id>
                <name>Jesper Dahlbäck</name>
                <anv></anv>
                <join></join>
                <role>Music By [All Tracks By]</role>
                <tracks></tracks>
            </artist>
        </extraartists>
        <formats>
            <format name="Vinyl" qty="2" text="">
                <descriptions>
                    <description>12"</description>
                    <description>33 ⅓ RPM</description>
                </descriptions>
            </format>
        </formats>
        <genres>
            <genre>Electronic</genre>
        </genres>
        <styles>
            <style>Deep House</style>
        </styles>
        <country>Sweden</country>
        <released>1999-03-00</released>
        <notes>The song titles are the names of six of Stockholm's 82 districts.

Title on label: - Stockholm -

Recorded at the Globe Studio, Stockholm

FAX: +46 8 679 64 53</notes>
        <data_quality>Needs Vote</data_quality>
        <tracklist>
            <track>
                <position>A</position>
                <title>Östermalm</title>
                <duration>4:45</duration>
            </track>
            <track>
                <position>B1</position>
                <title>Vasastaden</title>
                <duration>6:11</duration>
            </track>
            <track>
                <position>B2</position>
                <title>Kungsholmen</title>
                <duration>2:49</duration>
            </track>
            <track>
                <position>C1</position>
                <title>Södermalm</title>
                <duration>5:38</duration>
            </track>
            <track>
                <position>C2</position>
                <title>Norrmalm</title>
                <duration>4:52</duration>
            </track>
            <track>
                <position>D</position>
                <title>Gamla Stan</title>
                <duration>5:16</duration>
            </track>
        </tracklist>
        <identifiers>
            <identifier description="A-Side Runout" type="Matrix / Runout" value="MPO SK 032 A1" />
            <identifier description="B-Side Runout" type="Matrix / Runout" value="MPO SK 032 B1" />
            <identifier description="C-Side Runout" type="Matrix / Runout" value="MPO SK 032 C1" />
            <identifier description="D-Side Runout" type="Matrix / Runout" value="MPO SK 032 D1" />
            <identifier description="Only On A-Side Runout" type="Matrix / Runout" value="G PHRUPMASTERGENERAL T27 LONDON" />
        </identifiers>
        <videos>
            <video duration="296" embed="true" src="https://www.youtube.com/watch?v=MpmbntGDyNE">
                <title>The Persuader - Östermalm</title>
                <description>The Persuader - Östermalm</description>
            </video>
            <video duration="376" embed="true" src="https://www.youtube.com/watch?v=Cawyll0pOI4">
                <title>The Persuader - Vasastaden</title>
                <description>The Persuader - Vasastaden</description>
            </video>
            <video duration="176" embed="true" src="https://www.youtube.com/watch?v=XExCZfMCXdo">
                <title>The Persuader - Kungsholmen</title>
                <description>The Persuader - Kungsholmen</description>
            </video>
            <video duration="341" embed="true" src="https://www.youtube.com/watch?v=WDZqiENap_U">
                <title>The Persuader - Södermalm</title>
                <description>The Persuader - Södermalm</description>
            </video>
            <video duration="301" embed="true" src="https://www.youtube.com/watch?v=EBBHR3EMN50">
                <title>The Persuader - Norrmalm</title>
                <description>The Persuader - Norrmalm</description>
            </video>
            <video duration="326" embed="true" src="https://www.youtube.com/watch?v=afMHNll9EVM">
                <title>The Persuader - Gamla Stan</title>
                <description>The Persuader - Gamla Stan</description>
            </video>
        </videos>
        <companies>
            <company>
                <id>271046</id>
                <name>The Globe Studios</name>
                <catno></catno>
                <entity_type>23</entity_type>
                <entity_type_name>Recorded At</entity_type_name>
                <resource_url>https://api.discogs.com/labels/271046</resource_url>
            </company>
            <company>
                <id>56025</id>
                <name>MPO</name>
                <catno></catno>
                <entity_type>17</entity_type>
                <entity_type_name>Pressed By</entity_type_name>
                <resource_url>https://api.discogs.com/labels/56025</resource_url>
            </company>
        </companies>
    </release>
    <release id="2" status="Accepted">
        <images>
            <image height="394" type="primary" uri="" uri150="" width="400" />
            <image height="600" type="secondary" uri="" uri150="" width="600" />
            <image height="600" type="secondary" uri="" uri150="" width="600" />
        </images>
        <artists>
            <artist>
                <id>2</id>
                <name>Mr. James Barth &amp; A.D.</name>
                <anv></anv>
                <join></join>
                <role></role>
                <tracks></tracks>
            </artist>
        </artists>
        <title>Knockin' Boots Vol 2 Of 2</title>
        <labels>
            <label catno="SK 026" id="5" name="Svek" />
            <label catno="SK026" id="5" name="Svek" />
        </labels>
        <extraartists>
            <artist>
                <id>26</id>
                <name>Alexi Delano</name>
                <anv></anv>
                <join></join>
                <role>Producer, Recorded By</role>
                <tracks></tracks>
            </artist>
            <artist>
                <id>27</id>
                <name>Cari Lekebusch</name>
                <anv></anv>
                <join></join>
                <role>Producer, Recorded By</role>
                <tracks></tracks>
            </artist>
            <artist>
                <id>26</id>
                <name>Alexi Delano</name>
                <anv>A. Delano</anv>
                <join></join>
                <role>Written-By</role>
                <tracks></tracks>
            </artist>
            <artist>
                <id>27</id>
                <name>Cari Lekebusch</name>
                <anv>C. Lekebusch</anv>
                <join></join>
                <role>Written-By</role>
                <tracks></tracks>
            </artist>
        </extraartists>
        <formats>
            <format name="Vinyl" qty="1" text="">
                <descriptions>
                    <description>12"</description>
                    <description>33 ⅓ RPM</description>
                </descriptions>
            </format>
        </formats>
        <genres>
            <genre>Electronic</genre>
        </genres>
        <styles>
            <style>Broken Beat</style>
            <style>Techno</style>
            <style>Tech House</style>
        </styles>
        <country>Sweden</country>
        <released>1998-06-00</released>
        <notes>All joints recorded in NYC (Dec.97).</notes>
        <data_quality>Correct</data_quality>
        <master_id is_main_release="true">713738</master_id>
        <tracklist>
            <track>
                <position>A1</position>
                <title>A Sea Apart</title>
                <duration>5:08</duration>
            </track>
            <track>
                <position>A2</position>
                <title>Dutchmaster</title>
                <duration>4:21</duration>
            </track>
            <track>
                <position>B1</position>
                <title>Inner City Lullaby</title>
                <duration>4:22</duration>
            </track>
            <track>
                <position>B2</position>
                <title>Yeah Kid!</title>
                <duration>4:46</duration>
            </track>
        </tracklist>
        <identifiers>
            <identifier description="Side A Runout Etching" type="Matrix / Runout" value="MPO SK026-A -J.T.S.-" />
            <identifier description="Side B Runout Etching" type="Matrix / Runout" value="MPO SK026-B -J.T.S.-" />
        </identifiers>
        <videos>
            <video duration="310" embed="true" src="https://www.youtube.com/watch?v=MIgQNVhYILA">
                <title>Mr. James Barth &amp; A.D. - A Sea Apart</title>
                <description>Mr. James Barth &amp; A.D. - A Sea Apart</description>
            </video>
            <video duration="265" embed="true" src="https://www.youtube.com/watch?v=LgLchSRehhc">
                <title>Mr. James Barth &amp; A.D. - Dutchmaster</title>
                <description>Mr. James Barth &amp; A.D. - Dutchmaster</description>
            </video>
            <video duration="260" embed="true" src="https://www.youtube.com/watch?v=iaqHaULlqqg">
                <title>Mr. James Barth &amp; A.D. - Inner City Lullaby</title>
                <description>Mr. James Barth &amp; A.D. - Inner City Lullaby</description>
            </video>
            <video duration="290" embed="true" src="https://www.youtube.com/watch?v=x_Os7b-iWKs">
                <title>Mr. James Barth &amp; A.D. - Yeah Kid!</title>
                <description>Mr. James Barth &amp; A.D. - Yeah Kid!</description>
            </video>
        </videos>
        <companies>
            <company>
                <id>266169</id>
                <name>JTS Studios</name>
                <catno></catno>
                <entity_type>29</entity_type>
                <entity_type_name>Mastered At</entity_type_name>
                <resource_url>https://api.discogs.com/labels/266169</resource_url>
            </company>
            <company>
                <id>56025</id>
                <name>MPO</name>
                <catno></catno>
                <entity_type>17</entity_type>
                <entity_type_name>Pressed By</entity_type_name>
                <resource_url>https://api.discogs.com/labels/56025</resource_url>
            </company>
        </companies>
    </release>
</releases>`
