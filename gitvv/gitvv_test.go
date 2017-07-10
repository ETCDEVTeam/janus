package gitvv

import "testing"
import "strings"

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

func Test_getCommitCountFromDescription(t *testing.T) {
	table := []struct {
		s    string
		want string
	}{
		{"v0.1.7", "0"},
		{"v3.5.0-64-g7658d85", "64"},
	}

	for _, tt := range table {
		if got := getCommitCountFromDescription(tt.s); got != tt.want {
			t.Errorf("hash: %s, got: %v, want: %v", tt.s, got, tt.want)
		}
	}
}

func Test_getSemverFromDescription(t *testing.T) {
	table := []struct {
		s    string
		want []string
	}{
		{"v0.1.7", []string{"0", "1", "7"}},
		{"v3.5.0-64-g7658d85", []string{"3", "5", "0"}},
		{"dc1dfcc", []string{"0", "0", "0"}},
	}

	for _, tt := range table {
		if got := getSemverFromDescription(tt.s); strings.Join(got, "") != strings.Join(tt.want, "") {
			t.Errorf("hash: %s, got: %v, want: %v", tt.s, got, tt.want)
		}
	}
}
