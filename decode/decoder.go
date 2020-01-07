// Package decode contains all necessary functionality to cover decoding of the Discogs data.
package decode

import (
	"github.com/lukasaron/data-discogs/model"
	"github.com/lukasaron/data-discogs/write"
)

// Decoder is the interface that wraps the basic decoding method.
//
// Artists, Labels, Masters and Releases are method that parse and decode data from discogs and return the appropriate
// slice of structured data. When error occurs it is returned as a second parameter.
//
// Close method clean all data related to decoding.
type Decoder interface {
	Decode(write.Writer) error
	Artists() (int, []model.Artist, error)
	Labels() (int, []model.Label, error)
	Masters() (int, []model.Master, error)
	Releases() (int, []model.Release, error)
	Options() Options
	SetOptions(Options)
	Error() error
}

// Quality Level specifies the required data to be parsed based on the Discogs marking.
type QualityLevel int

// Quality Level constants defined from a Discogs data.
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
	ItemSize int // Item size of the block - how many items will be parsed at one block
	Limit    int // Limit of blocks - how many blocks will be parsed - the upper limit
	Skip     int // Skip blocks - how many blocks will be skipped from the beginning
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
