package decoder

import (
	"encoding/xml"
	"os"
)

type Decoder struct {
	file    *os.File
	decoder *xml.Decoder
	Error   error
}

func NewDecoder(fileName string) *Decoder {
	d := &Decoder{}
	d.file, d.Error = os.Open(fileName)
	if d.Error == nil {
		d.decoder = xml.NewDecoder(d.file)
	}
	return d
}

func (d *Decoder) Close() error {
	return d.file.Close()
}
