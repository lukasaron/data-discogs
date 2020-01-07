package decode

import "testing"

func TestToQualityLevel_InValid(t *testing.T) {
	ql := ToQualityLevel("")
	if ql != All {
		t.Error("when the value is invalid, should be All")
	}
}

func TestToQualityLevel_All(t *testing.T) {
	value := "All"
	ql := ToQualityLevel(value)
	if ql != All {
		t.Errorf("should be All for value %s", value)
	}
}

func TestToQualityLevel_NeedsVote(t *testing.T) {
	value := "Needs Vote"
	ql := ToQualityLevel(value)
	if ql != NeedsVote {
		t.Errorf("should be NeedsVote for value %s", value)
	}
}

func TestToQualityLevel_EntirelyIncorrect(t *testing.T) {
	value := "Entirely Incorrect"
	ql := ToQualityLevel(value)
	if ql != EntirelyIncorrect {
		t.Errorf("should be Entirely Incorrect for value %s", value)
	}
}

func TestToQualityLevel_NeedsMajorChanges(t *testing.T) {
	value := "Needs Major Changes"
	ql := ToQualityLevel(value)
	if ql != NeedsMajorChanges {
		t.Errorf("should be NeedsMajorChanges for value %s", value)
	}
}

func TestToQualityLevel_NeedsMinorChanges(t *testing.T) {
	value := "Needs Minor Changes"
	ql := ToQualityLevel(value)
	if ql != NeedsMinorChanges {
		t.Errorf("should be NeedsMinorChanges for value %s", value)
	}
}

func TestToQualityLevel_Correct(t *testing.T) {
	value := "Correct"
	ql := ToQualityLevel(value)
	if ql != Correct {
		t.Errorf("should be Correct for value %s", value)
	}
}

func TestToQualityLevel_CompleteAndCorrect(t *testing.T) {
	value := "Complete and Correct"
	ql := ToQualityLevel(value)
	if ql != CompleteAndCorrect {
		t.Errorf("should be CompleteAndCorrect for value %s", value)
	}
}

func TestQualityLevel_Includes_All(t *testing.T) {
	ql := All
	if !ql.Includes(CompleteAndCorrect) {
		t.Error("all should contain Complete and Correct")
	}

	if !ql.Includes(Correct) {
		t.Error("all should contain Correct")
	}

	if !ql.Includes(NeedsMinorChanges) {
		t.Error("all should contain Needs Minor Changes")
	}

	if !ql.Includes(NeedsMajorChanges) {
		t.Error("all should contain Needs Major Changes")
	}

	if !ql.Includes(NeedsVote) {
		t.Error("all should contain Needs Vode")
	}

	if !ql.Includes(EntirelyIncorrect) {
		t.Error("all should contain Entirely Incorrect")
	}

	if !ql.Includes(All) {
		t.Error("all should contain itself")
	}
}

func TestQualityLevel_Includes_EntirelyIncorrect(t *testing.T) {
	ql := EntirelyIncorrect
	if !ql.Includes(CompleteAndCorrect) {
		t.Error("entirely incorrect should contain Complete and Correct")
	}

	if !ql.Includes(Correct) {
		t.Error("entirely incorrect should contain Correct")
	}

	if !ql.Includes(NeedsMinorChanges) {
		t.Error("entirely incorrect should contain Needs Minor Changes")
	}

	if !ql.Includes(NeedsMajorChanges) {
		t.Error("entirely incorrect should contain Needs Major Changes")
	}

	if !ql.Includes(NeedsVote) {
		t.Error("entirely incorrect should contain Needs Vote")
	}

	if !ql.Includes(EntirelyIncorrect) {
		t.Error("entirely incorrect should contain itself")
	}

	if ql.Includes(All) {
		t.Error("entirely incorrect shouldn't contain All")
	}
}

func TestQualityLevel_Includes_NeedsVote(t *testing.T) {
	ql := NeedsVote
	if !ql.Includes(CompleteAndCorrect) {
		t.Error("needs vote should contain Complete and Correct")
	}

	if !ql.Includes(Correct) {
		t.Error("needs vote should contain Correct")
	}

	if !ql.Includes(NeedsMinorChanges) {
		t.Error("needs vote should contain Needs Minor Changes")
	}

	if !ql.Includes(NeedsMajorChanges) {
		t.Error("needs vote should contain Needs Major Changes")
	}

	if !ql.Includes(NeedsVote) {
		t.Error("needs vote should contain itself")
	}

	if ql.Includes(EntirelyIncorrect) {
		t.Error("needs vote shouldn't contain Entirely Incorrect")
	}

	if ql.Includes(All) {
		t.Error("needs vote shouldn't contain All")
	}
}

func TestQualityLevel_Includes_NeedsMajorChanges(t *testing.T) {
	ql := NeedsMajorChanges
	if !ql.Includes(CompleteAndCorrect) {
		t.Error("needs major changes should contain Complete and Correct")
	}

	if !ql.Includes(Correct) {
		t.Error("needs major changes should contain Correct")
	}

	if !ql.Includes(NeedsMinorChanges) {
		t.Error("needs major changes should contain Needs Minor Changes")
	}

	if !ql.Includes(NeedsMajorChanges) {
		t.Error("needs major changes should contain itself")
	}

	if ql.Includes(NeedsVote) {
		t.Error("needs major changes shouldn't contain Needs Vote")
	}

	if ql.Includes(EntirelyIncorrect) {
		t.Error("needs major changes shouldn't contain Entirely Incorrect")
	}

	if ql.Includes(All) {
		t.Error("needs major changes shouldn't contain All")
	}
}

func TestToQualityLevel_Includes_NeedsMinorChanges(t *testing.T) {
	ql := NeedsMinorChanges
	if !ql.Includes(CompleteAndCorrect) {
		t.Error("needs minor changes should contain Complete and Correct")
	}

	if !ql.Includes(Correct) {
		t.Error("needs minor changes should contain Correct")
	}

	if !ql.Includes(NeedsMinorChanges) {
		t.Error("needs minor changes should contain itself")
	}

	if ql.Includes(NeedsMajorChanges) {
		t.Error("needs minor changes shouldn't contain Needs Major Changes")
	}

	if ql.Includes(NeedsVote) {
		t.Error("needs minor changes shouldn't contain Needs Vote")
	}

	if ql.Includes(EntirelyIncorrect) {
		t.Error("needs minor changes shouldn't contain Entirely Incorrect")
	}

	if ql.Includes(All) {
		t.Error("needs minor changes shouldn't contain All")
	}
}

func TestQualityLevel_Includes_Correct(t *testing.T) {
	ql := Correct
	if !ql.Includes(CompleteAndCorrect) {
		t.Error("correct should contain Complete and Correct")
	}

	if !ql.Includes(Correct) {
		t.Error("correct should contain itself")
	}

	if ql.Includes(NeedsMinorChanges) {
		t.Error("correct shouldn't contain Needs Minor Changes")
	}

	if ql.Includes(NeedsMajorChanges) {
		t.Error("correct shouldn't contain Needs Major Changes")
	}

	if ql.Includes(NeedsVote) {
		t.Error("correct shouldn't contain Needs Vote")
	}

	if ql.Includes(EntirelyIncorrect) {
		t.Error("correct shouldn't contain Entirely Incorrect")
	}

	if ql.Includes(All) {
		t.Error("correct shouldn't contain All")
	}
}

func TestQualityLevel_Includes_CompleteAndCorrect(t *testing.T) {
	ql := CompleteAndCorrect
	if !ql.Includes(CompleteAndCorrect) {
		t.Error("complete and correct should contain itself")
	}

	if ql.Includes(Correct) {
		t.Error("complete and correct shouldn't contain Correct")
	}

	if ql.Includes(NeedsMinorChanges) {
		t.Error("complete and correct shouldn't contain Needs Minor Changes")
	}

	if ql.Includes(NeedsMajorChanges) {
		t.Error("complete and correct shouldn't contain Needs Major Changes")
	}

	if ql.Includes(NeedsVote) {
		t.Error("complete and correct shouldn't contain Needs Vote")
	}

	if ql.Includes(EntirelyIncorrect) {
		t.Error("complete and correct shouldn't contain Entirely Incorrect")
	}

	if ql.Includes(All) {
		t.Error("complete and correct shouldn't contain All")
	}
}
