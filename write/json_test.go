package write

import "testing"

func TestJSONWriter_Options(t *testing.T) {
	j := NewJSONWriter(nil, nil)
	opt := j.Options()

	if opt.ExcludeImages {
		t.Error("exclude images should be false as a default value")
	}
}
