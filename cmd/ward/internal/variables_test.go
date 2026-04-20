package internal

import (
	"runtime"
	"testing"
)

func setVars(v, s, g string) func() {
	ov, os, og := version, stage, gitCommit
	version, stage, gitCommit = v, s, g
	return func() { version, stage, gitCommit = ov, os, og }
}

func TestVersionTrimsPrefix(t *testing.T) {
	defer setVars("v1.2.3", "", "")()
	if got := Version(); got != "1.2.3" {
		t.Errorf("got %q, want %q", got, "1.2.3")
	}
}

func TestVersionUpperPrefix(t *testing.T) {
	defer setVars("V2.0.0", "", "")()
	if got := Version(); got != "2.0.0" {
		t.Errorf("got %q, want %q", got, "2.0.0")
	}
}

func TestVersionUndefined(t *testing.T) {
	defer setVars("", "", "")()
	if got := Version(); got != defaultUndefined {
		t.Errorf("got %q, want %q", got, defaultUndefined)
	}
}

func TestVersionWhitespace(t *testing.T) {
	defer setVars("  v3.0.0  ", "", "")()
	if got := Version(); got != "3.0.0" {
		t.Errorf("got %q, want %q", got, "3.0.0")
	}
}

func TestStage(t *testing.T) {
	defer setVars("", "Staging", "")()
	if got := Stage(); got != "staging" {
		t.Errorf("got %q, want %q", got, "staging")
	}
}

func TestStageUndefined(t *testing.T) {
	defer setVars("", "", "")()
	if got := Stage(); got != defaultUndefined {
		t.Errorf("got %q, want %q", got, defaultUndefined)
	}
}

func TestGitCommit(t *testing.T) {
	defer setVars("", "", "abc1234")()
	if got := GitCommit(); got != "abc1234" {
		t.Errorf("got %q, want %q", got, "abc1234")
	}
}

func TestGitCommitUndefined(t *testing.T) {
	defer setVars("", "", "")()
	if got := GitCommit(); got != defaultUndefined {
		t.Errorf("got %q, want %q", got, defaultUndefined)
	}
}

func TestArch(t *testing.T) {
	if got := Arch(); got != runtime.GOARCH {
		t.Errorf("got %q, want %q", got, runtime.GOARCH)
	}
}

func TestIsLocalAllSet(t *testing.T) {
	defer setVars("1.0.0", "main", "abc")()
	if IsLocal() {
		t.Error("expected false when all vars set")
	}
}

func TestIsLocalMissingVersion(t *testing.T) {
	defer setVars("", "main", "abc")()
	if !IsLocal() {
		t.Error("expected true when version is empty")
	}
}

func TestIsLocalMissingStage(t *testing.T) {
	defer setVars("1.0.0", "", "abc")()
	if !IsLocal() {
		t.Error("expected true when stage is empty")
	}
}

func TestIsLocalMissingCommit(t *testing.T) {
	defer setVars("1.0.0", "main", "")()
	if !IsLocal() {
		t.Error("expected true when gitCommit is empty")
	}
}

func TestVersionStringLocal(t *testing.T) {
	defer setVars("", "", "")()
	if got := VersionString(); got != defaultLocalBuild {
		t.Errorf("got %q, want %q", got, defaultLocalBuild)
	}
}

func TestVersionStringMain(t *testing.T) {
	defer setVars("1.0.0", "main", "abc1234")()
	want := "1.0.0 abc1234 [" + runtime.GOARCH + "]"
	if got := VersionString(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestVersionStringNonMain(t *testing.T) {
	defer setVars("1.0.0", "staging", "abc1234")()
	want := "1.0.0+staging abc1234 [" + runtime.GOARCH + "]"
	if got := VersionString(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
