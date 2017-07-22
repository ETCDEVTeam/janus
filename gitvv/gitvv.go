package gitvv

import (
	"errors"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type gitDescription int

const defaultHashLength = 7

var (
	cacheLastTagName            string
	cacheCommitCountFromTagName string
	cacheCommitCount            string
	cacheHEADHash               string
)

func isHash(s string) bool {
	// Strip 'g' prefix for SHA1
	if strings.HasPrefix(s, "g") {
		s = s[1:]
	}
	// Must be 0-9 V a-f
	r := regexp.MustCompile(`\b([^g-z\W]\w+)\b`)
	if r.MatchString(s) {
		return true
	}
	return false
}

// getTagOnHEADCommit gets the tag on the current commit, else
// returns "" if no tag on current commit
func getTagIfTagOnHEADCommit() (string, bool) {
	//git describe --exact-match --abbrev=0
	c, e := exec.Command("git", "describe", "--exact-match", "--abbrev=0").CombinedOutput()
	if e != nil {
		//log.Println(e)
		return "", false
	}
	tag := strings.TrimSpace(string(c))
	if tag == "" {
		return tag, false
	}
	return tag, true
}

func getCommitCountFrom(fromTag string) string {
	if cacheCommitCount != "" && cacheCommitCountFromTagName == fromTag {
		return cacheCommitCount
	}

	reference := "HEAD"
	if fromTag != "" {
		reference = fromTag + "..HEAD"
	}

	c, e := exec.Command("git", "rev-list", reference, "--count").Output()
	if e != nil {
		// TODO: handle error better
		//log.Println(e)
		return "0"
	}

	// Save caches
	cacheCommitCount = strings.TrimSpace(string(c))
	cacheCommitCountFromTagName = fromTag

	return cacheCommitCount
}

func getHEADHash(length int) string {
	if cacheHEADHash != "" {
		return cacheHEADHash[:length]
	}

	c, e := exec.Command("git", "rev-parse", "HEAD").Output()
	// > b9d3d5da740b4ed748734565614b8fe7885d9714
	if e != nil {
		log.Fatalln(e)
		return "???????"
	}

	sha1 := strings.TrimSpace(string(c))
	cacheHEADHash = sha1 // cache

	if length > len(cacheHEADHash) {
		length = len(cacheHEADHash)
	}
	return cacheHEADHash[:length]
}

func getLastTag() (string, bool) {
	if cacheLastTagName != "" {
		return cacheLastTagName, true
	}

	vOut, verErr := exec.Command("git", "describe", "--tags", "--abbrev=0").CombinedOutput()
	if verErr != nil {
		//log.Println(verErr)
		return "", false
	}

	tag := strings.TrimSpace(string(vOut))

	// Has no tags
	if tag == "" {
		return tag, false
	}

	cacheLastTagName = tag
	return cacheLastTagName, true
}

// Assumes using semver format for tags, eg v3.5.0 or 3.4.0
func parseSemverFromTag(s string) []string {
	tag := strings.TrimPrefix(s, "v")
	vers := strings.Split(tag, ".")
	return vers
}

// parseHashLength parses desired hash length output with default for none set
// eg.
// %S8 -> 8
// %S123 -> 123
// %S -> defaultLen
// NOTE: only compatible with single use #TODO?
func parseHashLength(s string) (int, error) {
	re := regexp.MustCompile(`%S(\d+)`)
	m := re.MatchString(s)
	// no digits following %S, use default
	if !m {
		return defaultHashLength, nil
	}
	f := re.FindAllString(s, 1)
	if f == nil || len(f) == 0 {
		return 0, errors.New("regex return match but no matching string(s) found")
	}
	ff := f[0]
	ff = strings.TrimPrefix(ff, "%S")
	i, e := strconv.Atoi(ff)
	if e != nil {
		return 0, e
	}

	return i, nil
}

// getB gets the semi-semver/mod patch number
func getB() (string, bool) {
	t, exists := getLastTag()
	if !exists {
		return "", false
	}
	semvers := parseSemverFromTag(t)
	if len(semvers) != 3 {
		return "", false
	}
	p := semvers[2]

	c := getCommitCountFrom(t)

	pi, e := strconv.Atoi(p)
	if e != nil {
		log.Println(e)
		return "", false
	}
	pi = pi * 100

	ci, e := strconv.Atoi(c)
	if e != nil {
		log.Println(e)
		return "", false
	}

	bi := pi + ci
	b := strconv.Itoa(bi)

	return b, true
}

// GetVersion gets formatted git version
// It assumes tags are by semver standards
// format:
// %M, _M - major version
// %m, _m - minor version
// %P, _P - patch version
// %C, _C - commit count since last tag
// %S, _S - HEAD sha1
func GetVersion(format string) string {

	var (
		lastTag     string
		commitCount string = "0"
		semvers = []string{}
		sha         string
	)
	semvers = nil

	// Set default format.
	if format == "" {
		// v3.5.0+66-bbb06b1
		format = "v%M.%m.%P-%S"
	}

	// Need to get commit count
	lastTag, _ = getTagIfTagOnHEADCommit()
	// Is not 0
	if lastTag == "" {
		lastTag, _ = getLastTag()
		// Either from init (entire branch) or lastTag

	}

	commitCount = getCommitCountFrom(lastTag)
	if lastTag != "" {
		semvers = parseSemverFromTag(lastTag)
	}

	// Convention alert:
	// Want: when commit count is 0 (ie HEAD is on a tag), should yield only semver, eg v3.5.0
	//       when commit count is >0 (ie HEAD is above a tag), should yield full "nightly" version name, eg v3.5.0+14-adfe123
	// This syntax allows to signify tagged builds vs running builds.
	// -- The point of this is just to be able to shift some logic out of CI scripts.
	if format == "TAG_OR_NIGHTLY" {
		format = "v%M.%m.%P+%C-%S"
		if commitCount == "0" {
			format = "v%M.%m.%P"
		}
	}

	sha = getHEADHash(defaultHashLength)
	if strings.Index(format, "%S") >= 0 {
		l, e := parseHashLength(format)
		if e != nil {
			log.Println(e)
		}
		if l != defaultHashLength {
			cacheHEADHash = ""
			sha = getHEADHash(l)
		}
	}

	out := format

	if semvers != nil {
		// -1 to replace indefinitely. Allows maximum user-decision-making.
		out = strings.Replace(out, "%M", semvers[0], -1)
		out = strings.Replace(out, "_M", semvers[0], -1)
		out = strings.Replace(out, "%m", semvers[1], -1)
		out = strings.Replace(out, "_m", semvers[1], -1)
		out = strings.Replace(out, "%P", semvers[2], -1)
		out = strings.Replace(out, "_P", semvers[2], -1)
	} else {
		out = strings.Replace(out, "%M", "?", -1)
		out = strings.Replace(out, "_M", "?", -1)
		out = strings.Replace(out, "%m", "?", -1)
		out = strings.Replace(out, "_m", "?", -1)
		out = strings.Replace(out, "%P", "?", -1)
		out = strings.Replace(out, "_P", "?", -1)
	}

	out = strings.Replace(out, "%C", commitCount, -1)
	out = strings.Replace(out, "_C", commitCount, -1)

	re1 := regexp.MustCompile(`(%S(\d+|))`)
	re2 := regexp.MustCompile(`(_S(\d+|))`)
	out = re1.ReplaceAllLiteralString(out, sha)
	out = re2.ReplaceAllLiteralString(out, sha)

	b, ok := getB()
	if ok {
		out = strings.Replace(out, "%B", b, -1)
		out = strings.Replace(out, "_B", b, -1)
	} else {
		out = strings.Replace(out, "%B", "?", -1)
		out = strings.Replace(out, "_B", "?", -1)
	}

	return out
}
