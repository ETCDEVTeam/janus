package gitvv

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testDir = "testdata"

var (
	baseProjectDir string
	noTagsDir      string
	onTagDir       string
	aboveTagDir    string
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

func Test_getLastTag(t *testing.T) {
	table := []struct {
		dir   string
		wants string
		wantb bool
	}{
		{noTagsDir, "", false},
		{aboveTagDir, "v0.0.1", true},
		{onTagDir, "v0.0.1", true},
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
		cacheLastTagName = ""

		tag, ok := getLastTag(repo.dir)
		if ok != repo.wantb {
			t.Errorf("got: %v, want: %v", ok, repo.wantb)
		}
		if tag != repo.wants {
			t.Errorf("got: %v, want: %v", tag, repo.wants)
		}

		if e := os.Chdir(baseProjectDir); e != nil {
			t.Fatal(e)
		}
	}
}

func Test_getTagIfTagOnHEADCommit(t *testing.T) {
	table := []struct {
		dir   string
		wants string
		wantb bool
	}{
		{noTagsDir, "", false},
		{aboveTagDir, "", false},
		{onTagDir, "v0.0.1", true},
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
		cacheLastTagName = ""

		tag, ok := getTagIfTagOnHEADCommit(repo.dir)
		if ok != repo.wantb {
			t.Errorf("got: %v, want: %v", ok, repo.wantb)
		}
		if tag != repo.wants {
			t.Errorf("got: %v, want: %v", tag, repo.wants)
		}

		if e := os.Chdir(baseProjectDir); e != nil {
			t.Fatal(e)
		}
	}
}

func Test_getHEADHash(t *testing.T) {
	table := []struct {
		dir        string
		want       string
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

		h := getHEADHash(repo.hashLength, repo.dir)
		if len(h) != repo.hashLength {
			t.Errorf("want: %d, got: %d, h: %s", repo.hashLength, len(h), h)
		}
		if !isHash(h) {
			t.Error("unexpected: not a hash")
		}

		if h != repo.want[:repo.hashLength] {
			t.Errorf("unexpected hash: want: %s, got: %s", repo.want, h)
		}

		if e := os.Chdir(baseProjectDir); e != nil {
			t.Fatal(e)
		}
	}
}

func Test_getCommitCountFrom(t *testing.T) {
	table := []struct {
		dir     string
		fromTag string
		want    string
	}{
		{noTagsDir, "", "1"},
		{noTagsDir, "v9.9.9", "0"},
		{aboveTagDir, "", "3"},
		{aboveTagDir, "v0.0.1", "1"},
		{onTagDir, "", "1"},
		{onTagDir, "v0.0.1", "0"},
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
		cacheCommitCount = ""
		cacheCommitCountFromTagName = ""

		count := getCommitCountFrom(repo.fromTag, repo.dir)
		if count != repo.want {
			t.Errorf("got: %v, want: %v", count, repo.want)
		}

		if e := os.Chdir(baseProjectDir); e != nil {
			t.Fatal(e)
		}
	}
}

func Test_getB(t *testing.T) {
	table := []struct {
		dir   string
		wants string
		wantb bool
	}{
		{noTagsDir, "", false},
		{aboveTagDir, "101", true},
		{onTagDir, "100", true},
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
		cacheCommitCount = ""
		cacheCommitCountFromTagName = ""
		cacheLastTagName = ""

		b, ok := getB(repo.dir)
		if ok != repo.wantb {
			t.Errorf("got: %v, want: %v", ok, repo.wantb)
		}
		if b != repo.wants {
			t.Errorf("got: %v, want: %v", b, repo.wants)
		}

		if e := os.Chdir(baseProjectDir); e != nil {
			t.Fatal(e)
		}
	}
}

func TestGetVersion(t *testing.T) {
	table := []struct {
		dir    string
		format string
		want   string
	}{
		{noTagsDir, "v", "v"},
		{onTagDir, "v", "v"},
		{aboveTagDir, "v", "v"},

		{noTagsDir, "v%M.%m.%P", "v?.?.?"},
		{onTagDir, "v%M.%m.%P", "v0.0.1"},
		{aboveTagDir, "v%M.%m.%P", "v0.0.1"},

		{noTagsDir, "%T", ""},
		{onTagDir, "%T-stable", "v0.0.1-stable"},
		{aboveTagDir, "%T-unstable", "v0.0.1-unstable"},

		{noTagsDir, "%C", "1"},
		{onTagDir, "%C", "0"},
		{aboveTagDir, "%C", "1"},

		{noTagsDir, "%S", "8673a80f120d8e11d607f1580da41c717e13863f"[:defaultHashLength]},
		{onTagDir, "%S", "e35b683e9d2b32c444976484472980582b4c68a9"[:defaultHashLength]},
		{aboveTagDir, "%S", "fe53b1e838d2fa761b3ce11d9fec683209f093a4"[:defaultHashLength]},

		{noTagsDir, "%S4", "8673a80f120d8e11d607f1580da41c717e13863f"[:4]},
		{onTagDir, "%S5", "e35b683e9d2b32c444976484472980582b4c68a9"[:5]},
		{aboveTagDir, "%S6", "fe53b1e838d2fa761b3ce11d9fec683209f093a4"[:6]},

		{noTagsDir, "v%M.%m.%P+%C-%S5", "v?.?.?+1-8673a"},
		{onTagDir, "v%M.%m.%P+%C-%S5", "v0.0.1+0-e35b6"},
		{aboveTagDir, "v%M.%m.%P+%C-%S5", "v0.0.1+1-fe53b"},

		{noTagsDir, "TAG_OR_NIGHTLY", "v?.?.?+1-8673a80"},
		{onTagDir, "TAG_OR_NIGHTLY", "v0.0.1-e35b683"},
		{aboveTagDir, "TAG_OR_NIGHTLY", "v0.0.1+1-fe53b1e"},

		{noTagsDir, "v%M.%m.%B-%S5", "v?.?.?-8673a"},
		{onTagDir, "v%M.%m.%B-%S5", "v0.0.100-e35b6"},
		{aboveTagDir, "v%M.%m.%B-%S5", "v0.0.101-fe53b"},
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
		cacheCommitCount = ""
		cacheLastTagName = ""
		cacheCommitCountFromTagName = ""

		got := GetVersion(repo.format, repo.dir)
		if got != repo.want {
			t.Errorf("got: %v, want: %v", got, repo.want)
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

func Test_parseSemverFromTag(t *testing.T) {
	table := []struct {
		s    string
		want []string
	}{
		{"v0.1.7", []string{"0", "1", "7"}},
		{"0.1.7", []string{"0", "1", "7"}},
		{"v0.1.7-stable", []string{"0", "1", "7"}},
		{"vv0.1.7.very-unstable", []string{"0", "1", "7"}},
		{"V0.1.7.very...unstable", []string{"0", "1", "7"}},
		{"1.0.is-this-a-feature-or-bug?", []string{"1","0"}}, // PTAL IDK.
	}

	for _, tt := range table {
		if got := parseSemverFromTag(tt.s); strings.Join(got, "") != strings.Join(tt.want, "") {
			t.Errorf("tag: %s, got: %v, want: %v", tt.s, got, tt.want)
		}
	}
}

func Test_parseHashLength(t *testing.T) {
	table := []struct {
		s     string
		wantl int
		wante error
	}{
		{"v%M.%m.%P+%C-%S", defaultHashLength, nil}, // test in big format string w/o specifier
		{"v%M.%m.%P+%C-%S8", 8, nil},                // test in big format string with specifier
		{"%S2", 2, nil},                             // test alone
		{"v%S2", 2, nil},                            // test irrelevant prefix
		{"%S109", 109, nil},                         // test not dependent on actual length of hash
		{"%s109", defaultHashLength, nil},           // test not %S (no matching format)
		{"%S17%S11", 17, nil},                       // test only works for single occurrence
	}

	for _, tt := range table {
		i, e := parseHashLength(tt.s)
		if e != tt.wante {
			t.Fatalf("unexpected error: %v", e)
		}
		if i != tt.wantl {
			t.Errorf("want: %v, got: %v", tt.wantl, i)
		}
	}
}
