/*
Package decoder implements a simple library for parsing Discogs dump files.
*/
package decoder

import (
	"errors"
	"github.com/lukasaron/discogs-parser/model"
	"github.com/lukasaron/discogs-parser/writer"
	"io"
	"log"
)

// Decoder is the interface that wraps the basic decoding method.
//
// Artists, Labels, Masers and Releases are method that parse and decode data from discogs and return the appropriate
// slice of structured data. When error occurs it is returned as a second parameter.
//
// Close method clean all data related to decoding.
type Decoder interface {
	Close() error
	Options() Options
	SetOptions(Options)
	Artists(limit int) (int, []model.Artist, error)
	Labels(limit int) (int, []model.Label, error)
	Masters(limit int) (int, []model.Master, error)
	Releases(limit int) (int, []model.Release, error)
}

// Quality Level specifies the required data to be parsed based on the Discogs marking.
type QualityLevel int

// Quality Level constants defined by a Discogs data.
const (
	All QualityLevel = iota
	EntirelyIncorrect
	NeedsVote
	NeedsMajorChanges
	NeedsMinorChanges
	Correct
	CompleteAndCorrect
)

// ToQualityLevel transforms string representation of Quality Level into the appropriate data type.
func ToQualityLevel(str string) (ql QualityLevel) {
	switch str {
	case "Entirely Incorrect":
		ql = EntirelyIncorrect
	case "Needs Vote":
		ql = NeedsVote
	case "Needs Major Changes":
		ql = NeedsMajorChanges
	case "Needs Minor Changes":
		ql = NeedsMinorChanges
	case "Correct":
		ql = Correct
	case "Complete and Correct":
		ql = CompleteAndCorrect
	case "All":
		fallthrough
	default:
		ql = All
	}

	return
}

// Includes decide if the Quality Level in the parameter has lower priority and thus the current Quality Level already
// contains the parameter value or not.
func (ql QualityLevel) Includes(q QualityLevel) bool {
	return ql <= q
}

// Block option structure implements the pagination principle for decoding the stream of data.
type Block struct {
	Size  int
	Limit int
	Skip  int
}

// File Type determines the input type of data to be decoded.
type FileType int

// Various types of File types are considered including Unknown type when the type is not specified.
const (
	Unknown FileType = iota
	Artists
	Labels
	Masters
	Releases
)

// Options consist of QualityLevel and Block settings
type Options struct {
	QualityLevel QualityLevel
	Block        Block
	FileType     FileType
}

// Error code returned by failure to parse a type specification.
var wrongTypeSpecified = errors.New("wrong file type specified")

// Decode data process input data via Decoder
func DecodeData(decoder Decoder, writer writer.Writer) error {
	opts := decoder.Options()
	if opts.Block.Limit < 1 {
		opts.Block.Limit = int(^uint(0) >> 1)
	}

	// get decode function based on the file type
	fn, err := decodeFunction(opts.FileType)
	if err != nil {
		return err
	}

	for blockCount := 1; blockCount <= opts.Block.Limit; blockCount++ {
		// call appropriate decoder function
		num, err := fn(decoder, writer, opts.Block.Size, blockCount > opts.Block.Skip)
		if err != nil && err != io.EOF {
			log.Printf("Block %d failed [%d]", blockCount, num)
			return err
		}

		if num == 0 && err == io.EOF {
			break
		}

		if blockCount > opts.Block.Skip {
			log.Printf("Block %d written [%d]", blockCount, num)
		} else {
			log.Printf("Block %d skipped [%d]", blockCount, num)
		}
	}

	return nil
}

func decodeFunction(ft FileType) (func(Decoder, writer.Writer, int, bool) (int, error), error) {
	switch ft {
	case Artists:
		return decodeArtists, nil
	case Labels:
		return decodeLabels, nil
	case Masters:
		return decodeMasters, nil
	case Releases:
		return decodeReleases, nil
	case Unknown:
		fallthrough
	default:
		return nil, wrongTypeSpecified
	}
}

func decodeArtists(d Decoder, w writer.Writer, blockSize int, write bool) (int, error) {
	num, a, err := d.Artists(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteArtists(a)
	}

	return num, err
}

func decodeLabels(d Decoder, w writer.Writer, blockSize int, write bool) (int, error) {
	num, l, err := d.Labels(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteLabels(l)
	}

	return num, err
}

func decodeMasters(d Decoder, w writer.Writer, blockSize int, write bool) (int, error) {
	num, m, err := d.Masters(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteMasters(m)
	}

	return num, err
}

func decodeReleases(d Decoder, w writer.Writer, blockSize int, write bool) (int, error) {
	num, r, err := d.Releases(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteReleases(r)
	}

	return num, err
}
