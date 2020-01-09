package discogs

import (
	"io"
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
		t.Errorf("block size should be set to default value: %v, when the set value is invalid", defaultBlockSize)
	}

	if opt.Block.Limit != defaultBlockLimit {
		t.Errorf("block limit should be set to default value: %v, when the set value is invalid", defaultBlockLimit)
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
	d := NewXMLDecoder(strings.NewReader(artists), nil)

	num, a, err := d.Artists()
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if num != 2 {
		t.Error("wrong number of parsed artists")
	}

	artist := a[0]

	if artist.ID != "1" {
		t.Errorf("wrong artist id, expected: %s, got: %s", "1", artist.ID)
	}

	if artist.Name != "The Persuader" {
		t.Errorf("wrong name, expected: %s, got: %s", "The Persuader", artist.Name)
	}

	if artist.RealName != "Jesper Dahlbäck" {
		t.Errorf("wrong real name, expected: %s, got: %s", "Jesper Dahlbäck", artist.RealName)
	}

	if artist.DataQuality != "Needs Vote" {
		t.Errorf("wrong data quality, expected: %s, got: %s", "Needs Vote", artist.DataQuality)
	}

	if len(artist.Urls) != 1 {
		t.Errorf("wrong number of urls, expected: %d, got: %d", 1, len(artist.Urls))
	}

	if artist.Urls[0] != "https://en.wikipedia.org/wiki/Jesper_Dahlbäck" {
		t.Errorf("wrong url, expected: %s, got: %s", "https://en.wikipedia.org/wiki/Jesper_Dahlbäck", artist.Urls[0])
	}

	if len(artist.NameVariations) != 2 {
		t.Errorf("wrong number of name variations, expected: %d, got: %d", 2, len(artist.NameVariations))
	}

	if artist.NameVariations[0] != "Persuader" {
		t.Errorf("wrong name variation, expected: %s, got: %s", "Persuader", artist.NameVariations[0])
	}

	if artist.NameVariations[1] != "The Presuader" {
		t.Errorf("wrong name variation, expected: %s, got: %s", "The Presuader", artist.NameVariations[1])
	}

	if len(artist.Images) != 2 {
		t.Errorf("wrong number of images, expected: %d, got: %d", 2, len(artist.Images))
	}

	if len(artist.Aliases) != 7 {
		t.Errorf("wrong number of aliases, expected: %d, got: %d", 7, len(artist.Aliases))
	}

	if artist.Aliases[0].ID != "239" || artist.Aliases[0].Name != "Jesper Dahlbäck" {
		t.Errorf("wrong alias, expected id: %s, got: %s, expected name: %s, got: %s",
			"239", artist.Aliases[0].ID,
			"Jesper Dahlbäck", artist.Aliases[0].Name)
	}

	if artist.Aliases[1].ID != "16055" || artist.Aliases[1].Name != "Groove Machine" {
		t.Errorf("wrong alias, expected id: %s, got: %s, expected name: %s, got: %s",
			"16055", artist.Aliases[1].ID,
			"Groove Machine", artist.Aliases[1].Name)
	}

	if artist.Aliases[2].ID != "19541" || artist.Aliases[2].Name != "Dick Track" {
		t.Errorf("wrong alias, expected id: %s, got: %s, expected name: %s, got: %s",
			"19541", artist.Aliases[2].ID,
			"Dick Track", artist.Aliases[2].Name)
	}

	if artist.Aliases[3].ID != "25227" || artist.Aliases[3].Name != "Lenk" {
		t.Errorf("wrong alias, expected id: %s, got: %s, expected name: %s, got: %s",
			"25227", artist.Aliases[3].ID,
			"Lenk", artist.Aliases[3].Name)
	}

	if artist.Aliases[4].ID != "196957" || artist.Aliases[4].Name != "Janne Me' Amazonen" {
		t.Errorf("wrong alias, expected id: %s, got: %s, expected name: %s, got: %s",
			"196957", artist.Aliases[4].ID,
			"Janne Me' Amazonen", artist.Aliases[4].Name)
	}

	if artist.Aliases[5].ID != "278760" || artist.Aliases[5].Name != "Faxid" {
		t.Errorf("wrong alias, expected id: %s, got: %s, expected name: %s, got: %s",
			"278760", artist.Aliases[5].ID,
			"Faxid", artist.Aliases[5].Name)
	}

	if artist.Aliases[6].ID != "439150" || artist.Aliases[6].Name != "The Pinguin Man" {
		t.Errorf("wrong alias, expected id: %s, got: %s, expected name: %s, got: %s",
			"439150", artist.Aliases[6].ID,
			"The Pinguin Man", artist.Aliases[6].Name)
	}

}

func TestXMLDecoder_Artists_Second(t *testing.T) {
	d := NewXMLDecoder(strings.NewReader(artists), nil)

	num, a, err := d.Artists()
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if num != 2 {
		t.Error("wrong number of parsed artists")
	}

	artist := a[1]

	if artist.ID != "2" {
		t.Errorf("wrong artist id, expected: %s, got: %s", "2", artist.ID)
	}

	if artist.Name != "Mr. James Barth & A.D." {
		t.Errorf("wrong name, expected: %s, got: %s", "Mr. James Barth & A.D.", artist.Name)
	}

	if artist.RealName != "Cari Lekebusch & Alexi Delano" {
		t.Errorf("wrong real name, expected: %s, got: %s", "Cari Lekebusch & Alexi Delano", artist.RealName)
	}

	if artist.DataQuality != "Correct" {
		t.Errorf("wrong data quality, expected: %s, got: %s", "Correct", artist.DataQuality)
	}

	if len(artist.NameVariations) != 4 {
		t.Errorf("wrong number of name variations, expected: %d, got: %d", 4, len(artist.NameVariations))

	}

	if artist.NameVariations[0] != "Mr Barth & A.D." {
		t.Errorf("wrong name variation, expected: %s, got: %s", "Mr Barth & A.D.", artist.NameVariations[0])
	}

	if artist.NameVariations[1] != "MR JAMES BARTH & A. D." {
		t.Errorf("wrong name variation, expected: %s, got: %s", "MR JAMES BARTH & A. D.", artist.NameVariations[1])
	}

	if artist.NameVariations[2] != "Mr. Barth & A.D." {
		t.Errorf("wrong name variation, expected: %s, got: %s", "Mr. Barth & A.D.", artist.NameVariations[2])
	}

	if artist.NameVariations[3] != "Mr. James Barth & A. D." {
		t.Errorf("wrong name variation, expected: %s, got: %s", "Mr. James Barth & A. D.", artist.NameVariations[3])
	}

	if len(artist.Aliases) != 5 {
		t.Errorf("wrong number of aliases, expected: %d, got: %d", 5, len(artist.Aliases))

	}

	if artist.Aliases[0].ID != "2470" || artist.Aliases[0].Name != "Puente Latino" {
		t.Errorf("wrong alias, expected id: %s, got: %s, expected name: %s, got: %s",
			"2470", artist.Aliases[0].ID,
			"Puente Latino", artist.Aliases[0].Name)
	}

	if artist.Aliases[1].ID != "19536" || artist.Aliases[1].Name != "Yakari & Delano" {
		t.Errorf("wrong alias, expected id: %s, got: %s, expected name: %s, got: %s",
			"19536", artist.Aliases[1].ID,
			"Yakari & Delano", artist.Aliases[1].Name)
	}

	if artist.Aliases[2].ID != "103709" || artist.Aliases[2].Name != "Crushed Insect & The Sick Puppy" {
		t.Errorf("wrong alias, expected id: %s, got: %s, expected name: %s, got: %s",
			"103709", artist.Aliases[2].ID,
			"Crushed Insect & The Sick Puppy", artist.Aliases[2].Name)
	}

	if artist.Aliases[3].ID != "384581" || artist.Aliases[3].Name != "ADCL" {
		t.Errorf("wrong alias, expected id: %s, got: %s, expected name: %s, got: %s",
			"384581", artist.Aliases[3].ID,
			"ADCL", artist.Aliases[3].Name)
	}

	if artist.Aliases[4].ID != "1779857" || artist.Aliases[4].Name != "Alexi Delano & Cari Lekebusch" {
		t.Errorf("wrong alias, expected id: %s, got: %s, expected name: %s, got: %s",
			"1779857", artist.Aliases[4].ID,
			"Alexi Delano & Cari Lekebusch", artist.Aliases[4].Name)
	}

	if len(artist.Members) != 2 {
		t.Errorf("wrong number of members, expected: %d, got: %d", 5, len(artist.Members))

	}

	if artist.Members[0].ID != "26" || artist.Members[0].Name != "Alexi Delano" {
		t.Errorf("wrong member, expected id: %s, got: %s, expected name: %s, got: %s",
			"26", artist.Members[0].ID,
			"Alexi Delano", artist.Members[0].Name)
	}

	if artist.Members[1].ID != "27" || artist.Members[1].Name != "Cari Lekebusch" {
		t.Errorf("wrong member, expected id: %s, got: %s, expected name: %s, got: %s",
			"27", artist.Members[1].ID,
			"Cari Lekebusch", artist.Members[1].Name)
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
	d := NewXMLDecoder(strings.NewReader(labels), nil)
	num, l, err := d.Labels()
	if err != nil && err != io.EOF {
		t.Error(err)

	}

	if num != 2 {
		t.Error("wrong number of parsed labels")
	}

	label := l[0]

	if label.ID != "1" {
		t.Errorf("wrong label id, expected: %s, got: %s", "1", label.ID)
	}

	if label.Name != "Planet E" {
		t.Errorf("wrong label name, expected: %s, got: %s", "Planet E", label.Name)
	}

	if label.ContactInfo != "Planet E Communications\r\nP.O. Box 27218\r\nDetroit, Michigan, MI 48227\r\nUSA\r\n\r\nPhone: +1 313 874 8729\r\nFax: +1 313 874 8732\r\nEmail: info@Planet-e.net" {
		t.Error("wrong contact info")
	}

	if label.Profile != "[a=Carl Craig]'s classic techno label founded in 1991.\r\n\r\nOn at least 1 release, Planet E is listed as publisher." {
		t.Error("wrong profile")
	}

	if label.DataQuality != "Correct" {
		t.Errorf("wrong data quality, expected: %s, got: %s", "Correct", label.DataQuality)
	}

	if len(label.Images) != 7 {
		t.Error("wrong number of images")
	}

	if len(label.Urls) != 13 {
		t.Error("wrong number of urls")
	}

	if label.Urls[0] != "http://planet-e.net" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://planet-e.net", label.Urls[0])
	}

	if label.Urls[1] != "http://planetecommunications.bandcamp.com" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://planetecommunications.bandcamp.com", label.Urls[1])
	}

	if label.Urls[2] != "http://www.facebook.com/planetedetroit" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://www.facebook.com/planetedetroit", label.Urls[2])
	}

	if label.Urls[3] != "http://www.flickr.com/photos/planetedetroit" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://www.flickr.com/photos/planetedetroit", label.Urls[3])
	}

	if label.Urls[4] != "http://plus.google.com/100841702106447505236" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://plus.google.com/100841702106447505236", label.Urls[4])
	}

	if label.Urls[5] != "http://www.instagram.com/carlcraignet" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://www.instagram.com/carlcraignet", label.Urls[5])
	}

	if label.Urls[6] != "http://myspace.com/planetecom" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://myspace.com/planetecom", label.Urls[6])
	}

	if label.Urls[7] != "http://myspace.com/planetedetroit" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://myspace.com/planetedetroit", label.Urls[7])
	}

	if label.Urls[8] != "http://soundcloud.com/planetedetroit" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://soundcloud.com/planetedetroit", label.Urls[8])
	}

	if label.Urls[9] != "http://twitter.com/planetedetroit" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://twitter.com/planetedetroit", label.Urls[9])
	}

	if label.Urls[10] != "http://vimeo.com/user1265384" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://vimeo.com/user1265384", label.Urls[10])
	}

	if label.Urls[11] != "http://en.wikipedia.org/wiki/Planet_E_Communications" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://en.wikipedia.org/wiki/Planet_E_Communications", label.Urls[11])
	}

	if label.Urls[12] != "http://www.youtube.com/user/planetedetroit" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://www.youtube.com/user/planetedetroit", label.Urls[12])
	}

	if len(label.SubLabels) != 8 {
		t.Error("wrong number of sub labels")
	}

	if label.SubLabels[0].ID != "86537" || label.SubLabels[0].Name != "Antidote (4)" {
		t.Error("wrong sublabel 0")
	}

	if label.SubLabels[1].ID != "41841" || label.SubLabels[1].Name != "Community Projects" {
		t.Error("wrong sublabel 1")
	}

	if label.SubLabels[2].ID != "153760" || label.SubLabels[2].Name != "Guilty Pleasures" {
		t.Error("wrong sublabel 2")
	}

	if label.SubLabels[3].ID != "31405" || label.SubLabels[3].Name != "I Ner Zon Sounds" {
		t.Error("wrong sublabel 3")
	}

	if label.SubLabels[4].ID != "277579" || label.SubLabels[4].Name != "Planet E Communications" {
		t.Error("wrong sublabel 4")
	}

	if label.SubLabels[5].ID != "294738" || label.SubLabels[5].Name != "Planet E Communications, Inc." {
		t.Error("wrong sublabel 5")
	}

	if label.SubLabels[6].ID != "1560615" || label.SubLabels[6].Name != "Planet E Productions" {
		t.Error("wrong sublabel 6")
	}

	if label.SubLabels[7].ID != "488315" || label.SubLabels[7].Name != "TWPENTY" {
		t.Error("wrong sublabel 7")
	}
}

func TestXMLDecoder_Labels_Second(t *testing.T) {
	d := NewXMLDecoder(strings.NewReader(labels), nil)
	num, l, err := d.Labels()
	if err != nil && err != io.EOF {
		t.Error(err)

	}

	if num != 2 {
		t.Error("wrong number of parsed labels")
	}

	label := l[1]

	if label.ID != "2" {
		t.Errorf("wrong label id, expected: %s, got: %s", "2", label.ID)
	}

	if label.Name != "Earthtones Recordings" {
		t.Errorf("wrong label name, expected: %s, got: %s", "Earthtones Recordings", label.Name)
	}

	if label.ContactInfo != "Seasons Recordings\r\n2236 Pacific Avenue\r\nSuite D\r\nCosta Mesa, CA  92627\r\n\r\ntel: +1.949.574.5255\r\nfax: +1.949.574.0255\r\n\r\nemail: jthinnes@seasonsrecordings.com\r\n" {
		t.Error("wrong contact info")
	}

	if label.Profile != "California deep house label founded by [a=Jamie Thinnes]. Now defunct and continued as [l=Seasons Recordings]." {
		t.Error("wrong profile")
	}

	if label.DataQuality != "Correct" {
		t.Errorf("wrong data quality, expected: %s, got: %s", "Correct", label.DataQuality)
	}

	if len(label.Urls) != 1 {
		t.Error("wrong number of urls")
	}

	if label.Urls[0] != "http://www.seasonsrecordings.com/" {
		t.Errorf("wrong url, expected: %s, got: %s", "http://www.seasonsrecordings.com/", label.Urls[0])
	}

	if len(label.Images) != 0 {
		t.Errorf("wrong number of images, expected: %d, got: %d", 0, len(label.Images))
	}

	if len(label.SubLabels) != 0 {
		t.Errorf("wrong number of sub labels expected: %d, got: %d", 0, len(label.SubLabels))

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
	d := NewXMLDecoder(strings.NewReader(masters), nil)
	num, m, err := d.Masters()
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if num != 2 {
		t.Error("there should be 2 masters decoded")
	}

	master := m[0]
	if master.ID != "18500" {
		t.Errorf("wrong master id, expected: %s, got: %s", "18500", master.ID)
	}

	if master.MainRelease != "155102" {
		t.Errorf("wrong main release, expected: %s, got: %s", "155102", master.MainRelease)
	}

	if len(master.Images) != 1 {
		t.Error("wrong number of images")
	}

	if master.Images[0].Height != "588" || master.Images[0].Width != "600" || master.Images[0].Type != "primary" {
		t.Error("wrong image parameters")
	}

	if len(master.Artists) != 1 {
		t.Error("wrong number of release artists")
	}

	if master.Artists[0].ID != "212070" ||
		master.Artists[0].Name != "Samuel L Session" ||
		master.Artists[0].Anv != "Samuel L" {
		t.Error("wrong release artist")
	}

	if len(master.Genres) != 1 {
		t.Error("wrong number of genres")
	}

	if master.Genres[0] != "Electronic" {
		t.Errorf("wrong genre, expected: %s, got: %s", "Electronic", master.Genres[0])
	}

	if len(master.Styles) != 1 {
		t.Error("wrong number of styles")
	}

	if master.Styles[0] != "Techno" {
		t.Errorf("wrong style, expected: %s, got: %s", "Techno", master.Styles[0])
	}

	if master.Year != "2001" {
		t.Errorf("wrong year, expected: %s, got: %s", "2001", master.Year)
	}

	if master.Title != "New Soil" {
		t.Errorf("wrong title, expected: %s, got: %s", "New Soil", master.Title)
	}

	if master.DataQuality != "Correct" {
		t.Errorf("wrong data quality, expected: %s, got: %s", "Correct", master.DataQuality)
	}

	if len(master.Videos) != 7 {
		t.Error("wrong number of videos")
	}

	if master.Videos[0].Title != "Samuel L - Velvet" ||
		master.Videos[0].Duration != "489" ||
		master.Videos[0].Embed != "true" ||
		master.Videos[0].Src != "https://www.youtube.com/watch?v=f05Ai921itM" ||
		master.Videos[0].Description != "Samuel L - Velvet" {
		t.Error("wrong video 0")
	}

	if master.Videos[1].Title != "Samuel L - Danses D'Afrique" ||
		master.Videos[1].Duration != "348" ||
		master.Videos[1].Embed != "true" ||
		master.Videos[1].Src != "https://www.youtube.com/watch?v=v23rSPG_StA" ||
		master.Videos[1].Description != "Samuel L - Danses D'Afrique" {
		t.Error("wrong video 1")
	}

	if master.Videos[2].Title != "Samuel L - Body N' Soul" ||
		master.Videos[2].Duration != "288" ||
		master.Videos[2].Embed != "true" ||
		master.Videos[2].Src != "https://www.youtube.com/watch?v=tHo82ha6p40" ||
		master.Videos[2].Description != "Samuel L - Body N' Soul" {
		t.Error("wrong video 2")
	}

	if master.Videos[3].Title != "Samuel L - Into The Groove" ||
		master.Videos[3].Duration != "331" ||
		master.Videos[3].Embed != "true" ||
		master.Videos[3].Src != "https://www.youtube.com/watch?v=KDcqzHca5dk" ||
		master.Videos[3].Description != "Samuel L - Into The Groove" {
		t.Error("wrong video 3")
	}

	if master.Videos[4].Title != "Samuel L - Soul Syndrome" ||
		master.Videos[4].Duration != "334" ||
		master.Videos[4].Embed != "true" ||
		master.Videos[4].Src != "https://www.youtube.com/watch?v=3DIYjJFl8Dk" ||
		master.Videos[4].Description != "Samuel L - Soul Syndrome" {
		t.Error("wrong video 4")
	}

	if master.Videos[5].Title != "Samuel L - Lush" ||
		master.Videos[5].Duration != "325" ||
		master.Videos[5].Embed != "true" ||
		master.Videos[5].Src != "https://www.youtube.com/watch?v=_o8yZMPqvNg" ||
		master.Videos[5].Description != "Samuel L - Lush" {
		t.Error("wrong video 5")
	}

	if master.Videos[6].Title != "Samuel L - Velvet ( Direct Me )" ||
		master.Videos[6].Duration != "346" ||
		master.Videos[6].Embed != "true" ||
		master.Videos[6].Src != "https://www.youtube.com/watch?v=JPwwJSc_-30" ||
		master.Videos[6].Description != "Samuel L - Velvet ( Direct Me )" {
		t.Error("wrong video 6")
	}
}

func TestXMLDecoder_Masters_Second(t *testing.T) {
	d := NewXMLDecoder(strings.NewReader(masters), nil)
	num, m, err := d.Masters()
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if num != 2 {
		t.Error("there should be 2 masters decoded")
	}

	master := m[1]
	if master.ID != "18512" {
		t.Errorf("wrong master id, expected: %s, got: %s", "18512", master.ID)
	}

	if master.MainRelease != "33699" {
		t.Errorf("wrong main release, expected: %s, got: %s", "33699", master.MainRelease)
	}

	if len(master.Images) != 3 {
		t.Error("wrong number of images")
	}

	if master.Images[0].Height != "150" || master.Images[0].Width != "150" || master.Images[0].Type != "primary" {
		t.Error("wrong image 0 parameters")
	}

	if master.Images[1].Height != "592" || master.Images[1].Width != "600" || master.Images[1].Type != "secondary" {
		t.Error("wrong image 1 parameters")
	}

	if master.Images[2].Height != "592" || master.Images[2].Width != "600" || master.Images[2].Type != "secondary" {
		t.Error("wrong image 2 parameters")
	}

	if len(master.Artists) != 1 {
		t.Error("wrong number of release artists")
	}

	if master.Artists[0].ID != "212070" || master.Artists[0].Name != "Samuel L Session" {
		t.Error("wrong release artist")
	}

	if len(master.Genres) != 1 {
		t.Error("wrong number of genres")
	}

	if master.Genres[0] != "Electronic" {
		t.Errorf("wrong genre, expected: %s, got: %s", "Electronic", master.Genres[0])
	}

	if len(master.Styles) != 2 {
		t.Error("wrong number of styles")
	}

	if master.Styles[0] != "Tribal" {
		t.Errorf("wrong style, expected: %s, got: %s", "Tribal", master.Styles[0])
	}

	if master.Styles[1] != "Techno" {
		t.Errorf("wrong style, expected: %s, got: %s", "Techno", master.Styles[1])
	}

	if master.Year != "2002" {
		t.Errorf("wrong year, expected: %s, got: %s", "2001", master.Year)
	}

	if master.Title != "Psyche EP" {
		t.Errorf("wrong title, expected: %s, got: %s", "Psyche EP", master.Title)
	}

	if master.DataQuality != "Correct" {
		t.Errorf("wrong data quality, expected: %s, got: %s", "Correct", master.DataQuality)
	}

	if len(master.Videos) != 3 {
		t.Error("wrong number of videos")
	}

	if master.Videos[0].Title != "Samuel L. Session - Arrival" ||
		master.Videos[0].Duration != "118" ||
		master.Videos[0].Embed != "true" ||
		master.Videos[0].Src != "https://www.youtube.com/watch?v=QYf4j0Pd2FU" ||
		master.Videos[0].Description != "Samuel L. Session - Arrival" {
		t.Error("wrong video 0")
	}

	if master.Videos[1].Title != "Samuel L. Session - Psyche Part 1" ||
		master.Videos[1].Duration != "376" ||
		master.Videos[1].Embed != "true" ||
		master.Videos[1].Src != "https://www.youtube.com/watch?v=c_AfLqTdncI" ||
		master.Videos[1].Description != "Samuel L. Session - Psyche Part 1" {
		t.Error("wrong video 1")
	}

	if master.Videos[2].Title != "Samuel L. Session - Psyche Part 2" ||
		master.Videos[2].Duration != "419" ||
		master.Videos[2].Embed != "true" ||
		master.Videos[2].Src != "https://www.youtube.com/watch?v=0nxvR8Zl9wY" ||
		master.Videos[2].Description != "Samuel L. Session - Psyche Part 2" {
		t.Error("wrong video 2")
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

// ---------------------------------------------------- DATA ----------------------------------------------------

var artists = `
<artists>
<artist><images><image height="450" type="primary" uri="" uri150="" width="600"/><image height="771" type="secondary" uri="" uri150="" width="600"/></images><id>1</id><name>The Persuader</name><realname>Jesper Dahlbäck</realname><profile></profile><data_quality>Needs Vote</data_quality><urls><url>https://en.wikipedia.org/wiki/Jesper_Dahlbäck</url></urls><namevariations><name>Persuader</name><name>The Presuader</name></namevariations><aliases><name id="239">Jesper Dahlbäck</name><name id="16055">Groove Machine</name><name id="19541">Dick Track</name><name id="25227">Lenk</name><name id="196957">Janne Me' Amazonen</name><name id="278760">Faxid</name><name id="439150">The Pinguin Man</name></aliases></artist>
<artist><id>2</id><name>Mr. James Barth &amp; A.D.</name><realname>Cari Lekebusch &amp; Alexi Delano</realname><profile></profile><data_quality>Correct</data_quality><namevariations><name>Mr Barth &amp; A.D.</name><name>MR JAMES BARTH &amp; A. D.</name><name>Mr. Barth &amp; A.D.</name><name>Mr. James Barth &amp; A. D.</name></namevariations><aliases><name id="2470">Puente Latino</name><name id="19536">Yakari &amp; Delano</name><name id="103709">Crushed Insect &amp; The Sick Puppy</name><name id="384581">ADCL</name><name id="1779857">Alexi Delano &amp; Cari Lekebusch</name></aliases><members><id>26</id><name id="26">Alexi Delano</name><id>27</id><name id="27">Cari Lekebusch</name></members></artist>
</artists>`

var labels = `
<labels>
<label><images><image height="24" type="primary" uri="" uri150="" width="132"/><image height="126" type="secondary" uri="" uri150="" width="587"/><image height="196" type="secondary" uri="" uri150="" width="600"/><image height="121" type="secondary" uri="" uri150="" width="275"/><image height="720" type="secondary" uri="" uri150="" width="382"/><image height="398" type="secondary" uri="" uri150="" width="500"/><image height="189" type="secondary" uri="" uri150="" width="600"/></images><id>1</id><name>Planet E</name><contactinfo>Planet E Communications&#13;
P.O. Box 27218&#13;
Detroit, Michigan, MI 48227&#13;
USA&#13;
&#13;
Phone: +1 313 874 8729&#13;
Fax: +1 313 874 8732&#13;
Email: info@Planet-e.net</contactinfo><profile>[a=Carl Craig]'s classic techno label founded in 1991.&#13;
&#13;
On at least 1 release, Planet E is listed as publisher.</profile><data_quality>Correct</data_quality><urls><url>http://planet-e.net</url><url>http://planetecommunications.bandcamp.com</url><url>http://www.facebook.com/planetedetroit</url><url>http://www.flickr.com/photos/planetedetroit</url><url>http://plus.google.com/100841702106447505236</url><url>http://www.instagram.com/carlcraignet</url><url>http://myspace.com/planetecom</url><url>http://myspace.com/planetedetroit</url><url>http://soundcloud.com/planetedetroit</url><url>http://twitter.com/planetedetroit</url><url>http://vimeo.com/user1265384</url><url>http://en.wikipedia.org/wiki/Planet_E_Communications</url><url>http://www.youtube.com/user/planetedetroit</url></urls><sublabels><label id="86537">Antidote (4)</label><label id="41841">Community Projects</label><label id="153760">Guilty Pleasures</label><label id="31405">I Ner Zon Sounds</label><label id="277579">Planet E Communications</label><label id="294738">Planet E Communications, Inc.</label><label id="1560615">Planet E Productions</label><label id="488315">TWPENTY</label></sublabels></label>
<label><id>2</id><name>Earthtones Recordings</name><contactinfo>Seasons Recordings&#13;
2236 Pacific Avenue&#13;
Suite D&#13;
Costa Mesa, CA  92627&#13;
&#13;
tel: +1.949.574.5255&#13;
fax: +1.949.574.0255&#13;
&#13;
email: jthinnes@seasonsrecordings.com&#13;
</contactinfo><profile>California deep house label founded by [a=Jamie Thinnes]. Now defunct and continued as [l=Seasons Recordings].</profile><data_quality>Correct</data_quality><urls><url>http://www.seasonsrecordings.com/</url></urls></label>
</labels>`

var masters = `
<masters>
<master id="18500"><main_release>155102</main_release><images><image height="588" type="primary" uri="" uri150="" width="600"/></images><artists><artist><id>212070</id><name>Samuel L Session</name><anv>Samuel L</anv><join></join><role></role><tracks></tracks></artist></artists><genres><genre>Electronic</genre></genres><styles><style>Techno</style></styles><year>2001</year><title>New Soil</title><data_quality>Correct</data_quality><videos><video duration="489" embed="true" src="https://www.youtube.com/watch?v=f05Ai921itM"><title>Samuel L - Velvet</title><description>Samuel L - Velvet</description></video><video duration="348" embed="true" src="https://www.youtube.com/watch?v=v23rSPG_StA"><title>Samuel L - Danses D'Afrique</title><description>Samuel L - Danses D'Afrique</description></video><video duration="288" embed="true" src="https://www.youtube.com/watch?v=tHo82ha6p40"><title>Samuel L - Body N' Soul</title><description>Samuel L - Body N' Soul</description></video><video duration="331" embed="true" src="https://www.youtube.com/watch?v=KDcqzHca5dk"><title>Samuel L - Into The Groove</title><description>Samuel L - Into The Groove</description></video><video duration="334" embed="true" src="https://www.youtube.com/watch?v=3DIYjJFl8Dk"><title>Samuel L - Soul Syndrome</title><description>Samuel L - Soul Syndrome</description></video><video duration="325" embed="true" src="https://www.youtube.com/watch?v=_o8yZMPqvNg"><title>Samuel L - Lush</title><description>Samuel L - Lush</description></video><video duration="346" embed="true" src="https://www.youtube.com/watch?v=JPwwJSc_-30"><title>Samuel L - Velvet ( Direct Me )</title><description>Samuel L - Velvet ( Direct Me )</description></video></videos></master>
<master id="18512"><main_release>33699</main_release><images><image height="150" type="primary" uri="" uri150="" width="150"/><image height="592" type="secondary" uri="" uri150="" width="600"/><image height="592" type="secondary" uri="" uri150="" width="600"/></images><artists><artist><id>212070</id><name>Samuel L Session</name><anv></anv><join></join><role></role><tracks></tracks></artist></artists><genres><genre>Electronic</genre></genres><styles><style>Tribal</style><style>Techno</style></styles><year>2002</year><title>Psyche EP</title><data_quality>Correct</data_quality><videos><video duration="118" embed="true" src="https://www.youtube.com/watch?v=QYf4j0Pd2FU"><title>Samuel L. Session - Arrival</title><description>Samuel L. Session - Arrival</description></video><video duration="376" embed="true" src="https://www.youtube.com/watch?v=c_AfLqTdncI"><title>Samuel L. Session - Psyche Part 1</title><description>Samuel L. Session - Psyche Part 1</description></video><video duration="419" embed="true" src="https://www.youtube.com/watch?v=0nxvR8Zl9wY"><title>Samuel L. Session - Psyche Part 2</title><description>Samuel L. Session - Psyche Part 2</description></video></videos></master>
</masters>`

var releases = `
<releases>
<release id="1" status="Accepted"><images><image height="600" type="primary" uri="" uri150="" width="600"/><image height="600" type="secondary" uri="" uri150="" width="600"/><image height="600" type="secondary" uri="" uri150="" width="600"/><image height="600" type="secondary" uri="" uri150="" width="600"/></images><artists><artist><id>1</id><name>The Persuader</name><anv></anv><join></join><role></role><tracks></tracks></artist></artists><title>Stockholm</title><labels><label catno="SK032" id="5" name="Svek"/></labels><extraartists><artist><id>239</id><name>Jesper Dahlbäck</name><anv></anv><join></join><role>Music By [All Tracks By]</role><tracks></tracks></artist></extraartists><formats><format name="Vinyl" qty="2" text=""><descriptions><description>12"</description><description>33 ⅓ RPM</description></descriptions></format></formats><genres><genre>Electronic</genre></genres><styles><style>Deep House</style></styles><country>Sweden</country><released>1999-03-00</released><notes>The song titles are the names of six of Stockholm's 82 districts.

Title on label: - Stockholm -

Recorded at the Globe Studio, Stockholm

FAX: +46 8 679 64 53</notes><data_quality>Needs Vote</data_quality><tracklist><track><position>A</position><title>Östermalm</title><duration>4:45</duration></track><track><position>B1</position><title>Vasastaden</title><duration>6:11</duration></track><track><position>B2</position><title>Kungsholmen</title><duration>2:49</duration></track><track><position>C1</position><title>Södermalm</title><duration>5:38</duration></track><track><position>C2</position><title>Norrmalm</title><duration>4:52</duration></track><track><position>D</position><title>Gamla Stan</title><duration>5:16</duration></track></tracklist><identifiers><identifier description="A-Side Runout" type="Matrix / Runout" value="MPO SK 032 A1"/><identifier description="B-Side Runout" type="Matrix / Runout" value="MPO SK 032 B1"/><identifier description="C-Side Runout" type="Matrix / Runout" value="MPO SK 032 C1"/><identifier description="D-Side Runout" type="Matrix / Runout" value="MPO SK 032 D1"/><identifier description="Only On A-Side Runout" type="Matrix / Runout" value="G PHRUPMASTERGENERAL T27 LONDON"/></identifiers><videos><video duration="296" embed="true" src="https://www.youtube.com/watch?v=MpmbntGDyNE"><title>The Persuader - Östermalm</title><description>The Persuader - Östermalm</description></video><video duration="376" embed="true" src="https://www.youtube.com/watch?v=Cawyll0pOI4"><title>The Persuader - Vasastaden</title><description>The Persuader - Vasastaden</description></video><video duration="176" embed="true" src="https://www.youtube.com/watch?v=XExCZfMCXdo"><title>The Persuader - Kungsholmen</title><description>The Persuader - Kungsholmen</description></video><video duration="341" embed="true" src="https://www.youtube.com/watch?v=WDZqiENap_U"><title>The Persuader - Södermalm</title><description>The Persuader - Södermalm</description></video><video duration="301" embed="true" src="https://www.youtube.com/watch?v=EBBHR3EMN50"><title>The Persuader - Norrmalm</title><description>The Persuader - Norrmalm</description></video><video duration="326" embed="true" src="https://www.youtube.com/watch?v=afMHNll9EVM"><title>The Persuader - Gamla Stan</title><description>The Persuader - Gamla Stan</description></video></videos><companies><company><id>271046</id><name>The Globe Studios</name><catno></catno><entity_type>23</entity_type><entity_type_name>Recorded At</entity_type_name><resource_url>https://api.discogs.com/labels/271046</resource_url></company><company><id>56025</id><name>MPO</name><catno></catno><entity_type>17</entity_type><entity_type_name>Pressed By</entity_type_name><resource_url>https://api.discogs.com/labels/56025</resource_url></company></companies></release>
<release id="2" status="Accepted"><images><image height="394" type="primary" uri="" uri150="" width="400"/><image height="600" type="secondary" uri="" uri150="" width="600"/><image height="600" type="secondary" uri="" uri150="" width="600"/></images><artists><artist><id>2</id><name>Mr. James Barth &amp; A.D.</name><anv></anv><join></join><role></role><tracks></tracks></artist></artists><title>Knockin' Boots Vol 2 Of 2</title><labels><label catno="SK 026" id="5" name="Svek"/><label catno="SK026" id="5" name="Svek"/></labels><extraartists><artist><id>26</id><name>Alexi Delano</name><anv></anv><join></join><role>Producer, Recorded By</role><tracks></tracks></artist><artist><id>27</id><name>Cari Lekebusch</name><anv></anv><join></join><role>Producer, Recorded By</role><tracks></tracks></artist><artist><id>26</id><name>Alexi Delano</name><anv>A. Delano</anv><join></join><role>Written-By</role><tracks></tracks></artist><artist><id>27</id><name>Cari Lekebusch</name><anv>C. Lekebusch</anv><join></join><role>Written-By</role><tracks></tracks></artist></extraartists><formats><format name="Vinyl" qty="1" text=""><descriptions><description>12"</description><description>33 ⅓ RPM</description></descriptions></format></formats><genres><genre>Electronic</genre></genres><styles><style>Broken Beat</style><style>Techno</style><style>Tech House</style></styles><country>Sweden</country><released>1998-06-00</released><notes>All joints recorded in NYC (Dec.97).</notes><data_quality>Correct</data_quality><master_id is_main_release="true">713738</master_id><tracklist><track><position>A1</position><title>A Sea Apart</title><duration>5:08</duration></track><track><position>A2</position><title>Dutchmaster</title><duration>4:21</duration></track><track><position>B1</position><title>Inner City Lullaby</title><duration>4:22</duration></track><track><position>B2</position><title>Yeah Kid!</title><duration>4:46</duration></track></tracklist><identifiers><identifier description="Side A Runout Etching" type="Matrix / Runout" value="MPO SK026-A -J.T.S.-"/><identifier description="Side B Runout Etching" type="Matrix / Runout" value="MPO SK026-B -J.T.S.-"/></identifiers><videos><video duration="310" embed="true" src="https://www.youtube.com/watch?v=MIgQNVhYILA"><title>Mr. James Barth &amp; A.D. - A Sea Apart</title><description>Mr. James Barth &amp; A.D. - A Sea Apart</description></video><video duration="265" embed="true" src="https://www.youtube.com/watch?v=LgLchSRehhc"><title>Mr. James Barth &amp; A.D. - Dutchmaster</title><description>Mr. James Barth &amp; A.D. - Dutchmaster</description></video><video duration="260" embed="true" src="https://www.youtube.com/watch?v=iaqHaULlqqg"><title>Mr. James Barth &amp; A.D. - Inner City Lullaby</title><description>Mr. James Barth &amp; A.D. - Inner City Lullaby</description></video><video duration="290" embed="true" src="https://www.youtube.com/watch?v=x_Os7b-iWKs"><title>Mr. James Barth &amp; A.D. - Yeah Kid!</title><description>Mr. James Barth &amp; A.D. - Yeah Kid!</description></video></videos><companies><company><id>266169</id><name>JTS Studios</name><catno></catno><entity_type>29</entity_type><entity_type_name>Mastered At</entity_type_name><resource_url>https://api.discogs.com/labels/266169</resource_url></company><company><id>56025</id><name>MPO</name><catno></catno><entity_type>17</entity_type><entity_type_name>Pressed By</entity_type_name><resource_url>https://api.discogs.com/labels/56025</resource_url></company></companies></release>
</releases>`
