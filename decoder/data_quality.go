package decoder

type DataQualityGetter interface {
	Quality() string
}

type QualityLevel int

const (
	All QualityLevel = iota
	EntirelyIncorrect
	NeedsVote
	NeedsMajorChanges
	NeedsMinorChanges
	Correct
	CompleteAndCorrect
)

func StrToQualityLevel(str string) (ql QualityLevel) {
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
	default:
		ql = All
	}

	return ql
}

func (ql QualityLevel) Includes(q QualityLevel) bool {
	return ql <= q
}
