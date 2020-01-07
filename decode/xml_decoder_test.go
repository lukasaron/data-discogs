package decode

import (
	"io"
	"strings"
	"testing"
)

func TestNewXmlDecoder(t *testing.T) {
	d := NewXmlDecoder(nil, nil)
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
	d := NewXmlDecoder(strings.NewReader(artists), nil)
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
	d := NewXmlDecoder(strings.NewReader(artists), &Options{
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
	d := NewXmlDecoder(nil, &Options{
		QualityLevel: NeedsVote,
		Block: Block{
			ItemSize: 0,
			Limit:    -1,
			Skip:     -1,
		},
		FileType: Artists,
	})

	opt := d.Options()
	if opt.QualityLevel != NeedsVote {
		t.Error("there should be Needs Vote quality level")
	}

	if opt.FileType != Artists {
		t.Error("there should be Artists file type")
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
	d := NewXmlDecoder(strings.NewReader(artists), &Options{
		FileType: Artists,
	})

	num, artists, err := d.Artists()
	if num != 2 || len(artists) != 2 {
		t.Error("there should be 2 artists parsed")
	}

	if err == nil {
		t.Error("expecting EOF error")
	}

	if err != io.EOF {
		t.Errorf("there should be EOF error instead of %v", err)
	}
}

func TestXMLDecoder_Artists_Block_Size(t *testing.T) {
	d := NewXmlDecoder(strings.NewReader(artists), &Options{
		FileType: Artists,
		Block: Block{
			ItemSize: 1,
		},
	})

	num, artists, err := d.Artists()
	if num != 1 || len(artists) != 1 {
		t.Error("there should be 1 artist parsed")
	}

	if err != nil {
		t.Errorf("no error expected when there are still some data to process, got %v", err)
	}

	num, artists, err = d.Artists()
	if num != 1 || len(artists) != 1 {
		t.Error("there should be 1 artist parsed")
	}

	if err != nil {
		t.Errorf("no error expected when there are still some data to process, got %v", err)
	}

	num, artists, err = d.Artists()
	if num != 0 || len(artists) != 0 {
		t.Error("there shouldn't be any artist parsed")
	}

	if err == nil {
		t.Error("EOF  error expected when there is nothing else to parse")
	}
}

var artists = `
<artists>
<artist><images><image height="450" type="primary" uri="" uri150="" width="600"/><image height="771" type="secondary" uri="" uri150="" width="600"/></images><id>1</id><name>The Persuader</name><realname>Jesper Dahlbäck</realname><profile></profile><data_quality>Needs Vote</data_quality><urls><url>https://en.wikipedia.org/wiki/Jesper_Dahlbäck</url></urls><namevariations><name>Persuader</name><name>The Presuader</name></namevariations><aliases><name id="239">Jesper Dahlbäck</name><name id="16055">Groove Machine</name><name id="19541">Dick Track</name><name id="25227">Lenk</name><name id="196957">Janne Me' Amazonen</name><name id="278760">Faxid</name><name id="439150">The Pinguin Man</name></aliases></artist>
<artist><id>2</id><name>Mr. James Barth &amp; A.D.</name><realname>Cari Lekebusch &amp; Alexi Delano</realname><profile></profile><data_quality>Correct</data_quality><namevariations><name>Mr Barth &amp; A.D.</name><name>MR JAMES BARTH &amp; A. D.</name><name>Mr. Barth &amp; A.D.</name><name>Mr. James Barth &amp; A. D.</name></namevariations><aliases><name id="2470">Puente Latino</name><name id="19536">Yakari &amp; Delano</name><name id="103709">Crushed Insect &amp; The Sick Puppy</name><name id="384581">ADCL</name><name id="1779857">Alexi Delano &amp; Cari Lekebusch</name></aliases><members><id>26</id><name id="26">Alexi Delano</name><id>27</id><name id="27">Cari Lekebusch</name></members></artist>
</artists>`
