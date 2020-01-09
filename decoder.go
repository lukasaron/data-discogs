// Copyright (C) 2020  Lukas Aron
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package discogs

import (
	"github.com/lukasaron/data-discogs/model"
	"github.com/lukasaron/data-discogs/write"
)

// Decoder is the interface that wraps the basic decoding method.
//
// Artists, Labels, Masters and Releases are method that parse and decode data from Discogs and return the appropriate
// slice of structured data. When error occurs it is returned as a second parameter.
//
// Close method cleans all data related to decoding.
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

// QualityLevel specifies the Quality Data field defined by Discogs in specific order to be able to define
// basic relationship and provide filtering ability.
type QualityLevel int

// QualityLevel constants mirror the Data Quality field defined by Discogs. Quality Level defines the order of these
// data quality values to be able to filter input items based on the quality level. Moreover, the All value has been
// added as a basic value (floor level), which covers all quality levels (this is a default value). When another value
// is set the chosen quality level will be part of decoded data set and all levels above as well. For instance when the
// value Correct is defined all decoded units with Data Quality level equals to Correct and Complete and Correct will
// be in the final set. On the other hand each item having Quality Field with value Needs Minor Changes,
// Needs Major Changes, Needs Vote and Entirely Incorrect will be ignored and will not be part of the result slice.
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

// Includes function decides if the Quality Level in the parameter has lower priority and thus the current Quality Level already
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

// FileType determines the input type of data to be decoded.
type FileType int

// Various types of input File types are considered based on the Discogs dump files, moreover this set
// includes Unknown type for cases when the type is not specified.
const (
	Unknown FileType = iota
	Artists
	Labels
	Masters
	Releases
)

// Options consist of QualityLevel, Block settings and FileType that will be decoded.
type Options struct {
	QualityLevel QualityLevel // Filters data based on the Data Quality field
	Block        Block        // Specifies the decoding Block values
	FileType     FileType     // Identifies XML file type
}
