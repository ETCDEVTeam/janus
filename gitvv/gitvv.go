package gitvv

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"strconv"
	"errors"
)

type gitDescription int

const (
	gitTag             gitDescription = iota // v3.5.0
	gitHash                                  // bbb06b1 or gbbb06b1
	gitFullDescription                       // v3.5.0-66-gbbb06b1
)

const defaultHashLength = 7

var (
	cacheDescription string
	cacheCommitCount string
	cacheHEADHash    string
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

func getCommitCountFromDescription(s string) string {
	switch getTypeofDescription(s) {
	case gitTag:
		return "0"
	case gitFullDescription:
		return strings.Split(s, "-")[1]
	}
	return getCommitCount()
}

// Assumes using semver format for tags, eg v3.5.0 or 3.4.0
func getSemverFromDescription(s string) []string {
	if getTypeofDescription(s) == gitHash {
		return []string{"0", "0", "0"}
	}
	tag := strings.Split(s, "-")[0]
	tag = strings.TrimPrefix(tag, "v")
	return strings.Split(tag, ".")
}

// getTypeofDescription returns the kind of log/tag as iota const
func getTypeofDescription(s string) gitDescription {
	ss := strings.Split(s, "-")
	if len(ss) == 1 {
		if isHash(ss[0]) {
			return gitHash
		}
		return gitTag
	}
	return gitFullDescription
}

func getCommitCount() string {
	if cacheCommitCount != "" {
		return cacheCommitCount
	}

	c, e := exec.Command("git", "rev-list", "HEAD", "--count").Output()
	if e != nil {
		// TODO: handle error better
		fmt.Println(e)
		return "0"
	}

	cacheCommitCount = strings.TrimSpace(string(c))
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

	return cacheHEADHash[:length]
}

func getDescription() (string, error) {
	if cacheDescription != "" {
		return cacheDescription, nil
	}

	vOut, verErr := exec.Command("git", "describe", "--tags", "--always").Output()
	if verErr != nil {
		return "???", verErr
	}

	cacheDescription = strings.TrimSpace(string(vOut))
	return cacheDescription, nil
}

func getSHA1FromDescription(s string, hashLength int) string {
	switch getTypeofDescription(s) {
	case gitFullDescription:
		ss := strings.Split(s, "-")
		// safety measure
		if len(ss) == 3 {
			sha := ss[2]
			if isHash(sha) {
				return strings.TrimPrefix(sha, "g")
			}
		}
	case gitHash:
		return strings.TrimPrefix(s, "g")
	}
	return getHEADHash(hashLength)
}

// parseHashLength parses desired hash length output with default for none set
// eg.
// %S8 -> 8
// %S123 -> 123
// %S -> defaultLen
// NOTE: only compatible with single use #TODO?
func parseHashLength(s string, defaultLen int) (int, error) {
	re := regexp.MustCompile(`\%S(\d+)`)
	m := re.MatchString(s)
	// no digits following %S, use default
	if !m {
		return defaultLen, nil
	}
	f := re.FindAllString(s, 1)
	if f == nil || len(f) == 0 {
		return defaultLen, errors.New("regex return match but no matching string(s) found")
	}
	ff := f[0]
	ff = strings.TrimPrefix(ff, "%S")
	i, e := strconv.Atoi(ff)
	if e != nil {
		return defaultLen, e
	}

	return i, nil
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
	if format == "" {
		// v3.5.0+66-bbb06b1
		format = "v%M.%m.%P-%S"
	}

	d, e := getDescription()
	if e != nil {
		fmt.Println(e)
		return ""
	}
	out := format

	commitCount := getCommitCountFromDescription(d)

	// Convention alert:
	// Want: when commit count is 0 (ie HEAD is on a tag), should yield only semver, eg v3.5.0
	//       when commit count is >0 (ie HEAD is above a tag), should yield full "nightly" version name, eg v3.5.0+14-adfe123
	// This syntax allows to signify tagged builds vs running builds.
	if format == "TAG_OR_NIGHTLY" {
		out = "v%M.%m.%P+%C-%S"
		if commitCount == "0" {
			out = "v%M.%m.%P"
		}
	}

	sha := getHEADHash(defaultHashLength)
	if strings.Index(format, "%S") >= 0 {
		l, e := parseHashLength(format, defaultHashLength)
		if e != nil {
			log.Println(e)
		}
		sha = getSHA1FromDescription(d, l)
	}

	semvers := getSemverFromDescription(d)

	// -1 to replace indefinitely. Allows maximum user-decision-making.
	out = strings.Replace(out, "%M", semvers[0], -1)
	out = strings.Replace(out, "%m", semvers[1], -1)
	out = strings.Replace(out, "%P", semvers[2], -1)
	out = strings.Replace(out, "%C", commitCount, -1)
	out = strings.Replace(out, "%S", sha, -1)

	// Problem: escaping %'s in windows (ps or batch) is a pain in the ass.
	// Solution: offer an extra set of escape vars for use with AppVeyor.
	out = strings.Replace(out, "_M", semvers[0], -1)
	out = strings.Replace(out, "_m", semvers[1], -1)
	out = strings.Replace(out, "_P", semvers[2], -1)
	out = strings.Replace(out, "_C", commitCount, -1)
	out = strings.Replace(out, "_S", sha, -1)

	return out
}
