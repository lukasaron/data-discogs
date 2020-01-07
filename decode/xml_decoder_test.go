package decode

import (
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

	_, _, err = d.Artists(10)
	if err != readerError {
		t.Error("there should be the same state error occurring")
	}

	_, _, err = d.Labels(10)
	if err != readerError {
		t.Error("there should be the same state error occurring")
	}

	_, _, err = d.Masters(10)
	if err != readerError {
		t.Error("there should be the same state error occurring")
	}

	_, _, err = d.Releases(10)
	if err != readerError {
		t.Error("there should be the same state error occurring")
	}

	_, _, err = d.Artists(10)
	if err != readerError {
		t.Error("there should be the same state error occurring")
	}
}

//func TestNewXmlDecoderNoOptions(t *testing.T) {
//	d := NewXmlDecoder(nil, nil)
//	if d.Error() != nil {
//		t.Errorf("there shouldn't be an error: %v", d.Error())
//	}
//
//	opt := d.Options()
//	if opt.QualityLevel != All {
//		t.Error("there should be All quality level")
//	}
//
//	if opt.FileType != Unknown {
//		t.Error("there should be Unknown file type")
//	}
//
//	if opt.Block.Size != defaultBlockSize {
//		t.Error("block size should be set to default value")
//	}
//
//	if opt.Block.Limit != defaultBlockLimit {
//		t.Error("block limit should be set to default value")
//	}
//}
