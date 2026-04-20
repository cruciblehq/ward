package internal

import (
	"fmt"
	"runtime"
	"strings"
)

const (

	// String to indicate an undefined variable.
	defaultUndefined = "(undefined)"

	// String to indicate a local (non-pipeline) build.
	defaultLocalBuild = "(local)"

	// Main branch name used in version strings.
	mainBranch = "main"
)

var (
	version   = "" // Version number (e.g., "1.2.3").
	stage     = "" // Development stage or git branch (e.g., "staging", "main").
	gitCommit = "" // Git commit hash (e.g., "a1b2c3d4").

	rawQuiet = "false" // Whether to enable quiet mode.
	rawDebug = "false" // Whether to enable debug mode.
)

// Returns the semantic version of the ward binary.
//
// If the version was not set via ldflags at build time, returns "(undefined)".
// Any leading "v" or "V" prefix is stripped and the result is lowercased.
func Version() string {
	v := strings.TrimSpace(version)
	if v == "" {
		return defaultUndefined
	}

	v = strings.ToLower(v)
	v = strings.TrimPrefix(v, "v")

	return v
}

// Returns the build stage or git branch name.
//
// The stage is set via ldflags and typically matches the branch that produced
// the build (e.g. "main", "staging"). Returns "(undefined)" when not set.
func Stage() string {
	s := strings.TrimSpace(stage)
	if s == "" {
		return defaultUndefined
	}
	return strings.ToLower(s)
}

// Returns the abbreviated git commit hash baked into the binary.
//
// Set via ldflags at build time. Returns "(undefined)" when the binary was not
// built through the standard pipeline.
func GitCommit() string {
	c := strings.TrimSpace(gitCommit)
	if c == "" {
		return defaultUndefined
	}
	return c
}

// Returns the GOARCH of the running binary.
//
// Used in the version string to indicate the target architecture.
func Arch() string {
	return runtime.GOARCH
}

// Reports whether any build-time variable was left unset.
//
// A build is local when version, stage, or gitCommit is empty, which means it
// was not produced by the CI pipeline. This affects the output of VersionString.
func IsLocal() bool {
	return strings.TrimSpace(version) == "" ||
		strings.TrimSpace(gitCommit) == "" ||
		strings.TrimSpace(stage) == ""
}

// Formats a human-readable version string for display.
//
// Local builds return "(local)". Pipeline builds return a string in the
// form "<version>+<stage> <commit> [<arch>]", omitting the stage suffix
// when the build was produced from the main branch.
func VersionString() string {
	if IsLocal() {
		return defaultLocalBuild
	}

	s := Stage()
	if s == mainBranch {
		s = ""
	} else {
		s = "+" + s
	}

	return fmt.Sprintf("%s%s %s [%s]", Version(), s, GitCommit(), Arch())
}
