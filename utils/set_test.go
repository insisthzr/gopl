package utils

import "testing"

const (
	format = "want: %d got: %s"
)

func TestSet(t *testing.T) {
	set := NewSet()
	set.Add("hello")
	set.Add("hello")
	set.Add("world")
	want := `["hello","world"]`
	got := set.String()
	if want != got {
		t.Fatalf(format, want, got)
	}

	set.Delete("world")
	want = `["hello"]`
	got = set.String()
	if want != got {
		t.Fatalf(format, want, got)
	}
}
