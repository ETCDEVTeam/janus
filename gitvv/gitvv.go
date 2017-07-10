package gitvv

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type gitDescription int

const (
	gitTag             gitDescription = iota // v3.5.0
	gitHash                                  // bbb06b1 or gbbb06b1
	gitFullDescription                       // v3.5.0-66-gbbb06b1
)

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
	r := regexp.MustCompile(`\b([^g-z]\w+)\b`)
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
		log.Println(e)
		return "0"
	}

	cacheCommitCount = strings.TrimSpace(string(c))
	return cacheCommitCount
}

func getHEADHash() string {
	if cacheHEADHash != "" {
		return cacheHEADHash
	}

	c, e := exec.Command("git", "rev-list", "HEAD", "--max-count=1").Output()
	if e != nil {
		log.Println(e)
		return "???????"
	}

	sha1 := strings.TrimSpace(string(c)[:7])

	cacheHEADHash = sha1
	return cacheHEADHash
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

func getSHA1FromDescription(s string) string {
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
	return getHEADHash()
}

// GetVersion gets formatted git version
// It assumes tags are by semver standards
// format:
// %M - major version
// %m - minor version
// %P - patch version
// %C - commit count since last tag
// %S - HEAD sha1
func GetVersion(format string) string {
	if format == "" {
		// v3.5.0+66-bbb06b1
		format = "v%M.%m.%P+%C-%S"
	}
	d, e := getDescription()
	if e != nil {
		log.Println(e)
		return ""
	}
	out := format
	commitCount := getCommitCountFromDescription(d)
	sha := getSHA1FromDescription(d)
	semvers := getSemverFromDescription(d)

	// -1 to replace indefinitely. Allows maximum user-decision-making.
	out = strings.Replace(out, "%M", semvers[0], -1)
	out = strings.Replace(out, "%m", semvers[1], -1)
	out = strings.Replace(out, "%P", semvers[2], -1)
	out = strings.Replace(out, "%C", commitCount, -1)
	out = strings.Replace(out, "%S", sha, -1)
	return out
}
