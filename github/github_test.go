package github

import (
	"testing"
)

func listIs(got []string, want ...string) bool {
	if len(got) != len(want) {
		return false
	}

	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}

	return true
}

func TestLatest(t *testing.T) {
	var v []string

	v = Latest([]string{})
	if !listIs(v) {
		t.Errorf("expected Latest([]) to be [], not %v", v)
	}

	v = Latest([]string{ "1.0.0" })
	if !listIs(v, "1.0.0") {
		t.Errorf("expected Latest([1.0.0]) to be [1.0.0], not %v", v)
	}

	v = Latest([]string{ "3.0.0", "2.0.0", "1.0.0" })
	if !listIs(v, "3.0.0", "2.0.0", "1.0.0") {
		t.Errorf("expected Latest([3.0.0 2.0.0 1.0.0]) to be [3.0.0 2.0.0 1.0.0], not %v", v)
	}

	v = Latest([]string{ "1.0.4", "1.0.3", "1.0.2" })
	if !listIs(v, "1.0.4", "1.0.3", "1.0.2") {
		t.Errorf("expected Latest([1.0.4 1.0.3 1.0.2]) to be [1.0.4 1.0.3 1.0.2], not %v", v)
	}

	v = Latest([]string{ "1.0.4", "1.0.3", "1.0.2", "0.9.7", "0.9.6", "0.8.1", "0.8.0" })
	if !listIs(v, "1.0.4", "1.0.3", "1.0.2", "0.9.7", "0.8.1") {
		t.Errorf("Latest() should keep all patch versions of the latest major.minor, but only keeping latest of older revs.\ngot %v", v)
	}

	v = Latest([]string{ "1.1.0", "1.0.3", "1.0.2", "0.9.7", "0.9.6", "0.8.1", "0.8.0" })
	if !listIs(v, "1.1.0", "1.0.3", "0.9.7", "0.8.1") {
		t.Errorf("Latest() should keep all patch versions of the latest major.minor, but only keeping latest of older revs.\ngot %v", v)
	}
}
