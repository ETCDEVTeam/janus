package gitvv

import "testing"

func Test_isHash(t *testing.T) {
	table := []struct {
		hash string
		ok   bool
	}{
		{"asdf123", true},
		{"zxcv123", false},
		{"gasdf123", true},
		{"gqwer123", false},
		{"--nope", false},
		{"v0.1.5", false},
	}

	for _, tt := range table {
		if got := isHash(tt.hash); got != tt.ok {
			t.Errorf("hash: %s, got: %v, want: %v", tt.hash, got, tt.ok)
		}
	}
}

// TODO: more tests for helper funcs
