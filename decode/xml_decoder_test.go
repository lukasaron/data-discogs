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

func TestXMLDecoder_Artists_Block_ItemSize(t *testing.T) {
	d := NewXmlDecoder(strings.NewReader(artists), &Options{
		FileType: Artists,
		Block: Block{
			ItemSize: 1,
		},
	})

	num, a, err := d.Artists()
	if num != 1 || len(a) != 1 {
		t.Error("there should be 1 artist parsed")
	}

	if err != nil {
		t.Errorf("no error expected when there are still some data to process, got %v", err)
	}

	num, a, err = d.Artists()
	if num != 1 || len(a) != 1 {
		t.Error("there should be 1 artist parsed")
	}

	if err != nil {
		t.Errorf("no error expected when there are still some data to process, got %v", err)
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
	d := NewXmlDecoder(strings.NewReader(labels), &Options{
		FileType: Labels,
	})

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

func TestXMLDecoder_Labels_Block_ItemSize(t *testing.T) {
	d := NewXmlDecoder(strings.NewReader(labels), &Options{
		FileType: Labels,
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
