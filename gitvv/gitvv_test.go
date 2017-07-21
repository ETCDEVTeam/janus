package gitvv

import (
	"path/filepath"
	"strings"
	"testing"
	"os"
	"log"
)

const testDir = "testdata"

var (
	baseProjectDir string
	noTagsDir string
	onTagDir string
	aboveTagDir string
)

// Initialize multiple repo testdirs.
// Test repos have different git states
func TestMain(m *testing.M) {
	projectDir, e := os.Getwd()
	if e != nil {
		log.Fatalln(e)
	}
	baseProjectDir = projectDir
	noTagsDir = filepath.Join(baseProjectDir, testDir, "no-tags")
	onTagDir = filepath.Join(baseProjectDir, testDir, "on-tag")
	aboveTagDir = filepath.Join(baseProjectDir, testDir, "above-tag")
	os.Exit(m.Run())
}

func Test_getHEADHash(t *testing.T) {
	table := []struct{
		dir string
		want string
		hashLength int
	}{
		{noTagsDir, "8673a80f120d8e11d607f1580da41c717e13863f", 7},
		{noTagsDir, "8673a80f120d8e11d607f1580da41c717e13863f", 9},
		{noTagsDir, "8673a80f120d8e11d607f1580da41c717e13863f", 2},
		{aboveTagDir, "fe53b1e838d2fa761b3ce11d9fec683209f093a4", 7},
		{onTagDir, "e35b683e9d2b32c444976484472980582b4c68a9", 7},
	}

	for _, repo := range table {
		if e := os.Chdir(repo.dir); e != nil {
			t.Fatal(e)
		}
		// Print that we're in the right directory.
		cwd, e := os.Getwd()
		if e != nil {
			t.Error(e)
		}
		t.Log(cwd)

		// Clear cache
		cacheHEADHash = ""

		h := getHEADHash(repo.hashLength)
		if len(h) != repo.hashLength {
			t.Errorf("want: %d, got: %d, h: %s", repo.hashLength, len(h), h)
		}
		if !isHash(h) {
			t.Error("unexpected: not a hash")
		}

		if h != repo.want {
			t.Errorf("unexpected hash: want: %s, got: %s", repo.want, h)
		}

		if e := os.Chdir(baseProjectDir); e != nil {
			t.Fatal(e)
		}
	}
}

// Test helpers
//
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
