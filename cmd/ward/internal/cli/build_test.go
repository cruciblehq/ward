package cli

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInputStdin(t *testing.T) {
	r, closer, err := input(stdinMarker)
	if err != nil {
		t.Fatal(err)
	}
	if closer != nil {
		t.Error("expected nil closer for stdin")
	}
	if r != os.Stdin {
		t.Error("expected os.Stdin reader")
	}
}

func TestInputFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.w")
	if err := os.WriteFile(path, []byte("version 1\n.expose /tmp\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	r, closer, err := input(path)
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("expected non-nil reader")
	}
	if closer == nil {
		t.Fatal("expected non-nil closer")
	}
	closer.Close()
}

func TestInputMissingFile(t *testing.T) {
	_, _, err := input("/nonexistent/path.w")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestOutputStdout(t *testing.T) {
	w, closer, err := output(stdoutMarker)
	if err != nil {
		t.Fatal(err)
	}
	if closer != nil {
		t.Error("expected nil closer for stdout")
	}
	if w != os.Stdout {
		t.Error("expected os.Stdout writer")
	}
}

func TestOutputFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.wo")

	w, closer, err := output(path)
	if err != nil {
		t.Fatal(err)
	}
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
	if closer == nil {
		t.Fatal("expected non-nil closer")
	}

	if _, err := w.Write([]byte("data")); err != nil {
		t.Fatal(err)
	}
	closer.Close()

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "data" {
		t.Errorf("got %q, want %q", got, "data")
	}
}

func TestOutputBadPath(t *testing.T) {
	_, _, err := output("/nonexistent/dir/out.wo")
	if err == nil {
		t.Error("expected error for bad output path")
	}
}

func TestOutputTruncates(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.wo")
	if err := os.WriteFile(path, []byte("old content that is long"), 0o644); err != nil {
		t.Fatal(err)
	}

	w, closer, err := output(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("new")); err != nil {
		t.Fatal(err)
	}
	closer.Close()

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "new" {
		t.Errorf("got %q, want %q", got, "new")
	}
}

func TestValidateSourcesEmpty(t *testing.T) {
	if err := validateSources(nil); err == nil {
		t.Error("expected error for no sources")
	}
}

func TestValidateSourcesSingle(t *testing.T) {
	if err := validateSources([]string{"a.w"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSourcesStdin(t *testing.T) {
	if err := validateSources([]string{"-"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSourcesMultipleFiles(t *testing.T) {
	if err := validateSources([]string{"a.w", "b.w"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSourcesStdinMixed(t *testing.T) {
	if err := validateSources([]string{"a.w", "-"}); err == nil {
		t.Error("expected error for stdin mixed with files")
	}
}

func TestValidateSourcesStdinFirst(t *testing.T) {
	if err := validateSources([]string{"-", "a.w"}); err == nil {
		t.Error("expected error for stdin mixed with files")
	}
}

func TestBuildValid(t *testing.T) {
	src := "version 1\n.expose /tmp\n"
	var buf bytes.Buffer
	if err := build(strings.NewReader(src), &buf); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestBuildCommentsAndBlanks(t *testing.T) {
	src := "version 1\n# header\n\n-- comment\n.expose /tmp\n\n"
	var buf bytes.Buffer
	if err := build(strings.NewReader(src), &buf); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestBuildInvalidRule(t *testing.T) {
	src := "not a valid rule\n"
	if err := build(strings.NewReader(src), io.Discard); err == nil {
		t.Error("expected error for invalid rule")
	}
}

func TestBuildEmpty(t *testing.T) {
	if err := build(strings.NewReader(""), io.Discard); err == nil {
		t.Error("expected error for missing version")
	}
}

func TestBuildBadVersion(t *testing.T) {
	src := "version abc\n"
	if err := build(strings.NewReader(src), io.Discard); err == nil {
		t.Error("expected error for bad version")
	}
}

func TestRunSingleSource(t *testing.T) {
	srcPath := filepath.Join(t.TempDir(), "test.w")
	outPath := filepath.Join(t.TempDir(), "out.wo")

	if err := os.WriteFile(srcPath, []byte("version 1\n.expose /tmp\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := &BuildCmd{Sources: []string{srcPath}, Output: outPath}
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty output file")
	}
}

func TestRunMultipleSources(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.w")
	b := filepath.Join(dir, "b.w")
	out := filepath.Join(dir, "out.wo")

	if err := os.WriteFile(a, []byte("version 1\n.expose /tmp\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte("version 1\n.expose /proc\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := &BuildCmd{Sources: []string{a, b}, Output: out}
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty linked output")
	}
}

func TestRunDryRun(t *testing.T) {
	srcPath := filepath.Join(t.TempDir(), "test.w")
	outPath := filepath.Join(t.TempDir(), "out.wo")

	if err := os.WriteFile(srcPath, []byte("version 1\n.expose /tmp\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := &BuildCmd{Sources: []string{srcPath}, Output: outPath, DryRun: true}
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(outPath); !errors.Is(err, os.ErrNotExist) {
		t.Error("expected no output file in dry-run mode")
	}
}

func TestRunNoSources(t *testing.T) {
	cmd := &BuildCmd{}
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected error for no sources")
	}
	if !errors.Is(err, ErrBuild) {
		t.Errorf("got %v, want ErrBuild wrapper", err)
	}
}

func TestRunMissingSource(t *testing.T) {
	cmd := &BuildCmd{Sources: []string{"/nonexistent/path.w"}}
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected error for missing source")
	}
	if !errors.Is(err, ErrParse) {
		t.Errorf("got %v, want ErrParse wrapper", err)
	}
}

func TestRunInvalidSource(t *testing.T) {
	srcPath := filepath.Join(t.TempDir(), "bad.w")
	if err := os.WriteFile(srcPath, []byte("not a valid rule\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := &BuildCmd{Sources: []string{srcPath}}
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected error for invalid source")
	}
	if !errors.Is(err, ErrBuild) {
		t.Errorf("got %v, want ErrBuild wrapper", err)
	}
}

func TestRunBadOutputPath(t *testing.T) {
	srcPath := filepath.Join(t.TempDir(), "test.w")
	if err := os.WriteFile(srcPath, []byte("version 1\n.expose /tmp\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := &BuildCmd{Sources: []string{srcPath}, Output: "/nonexistent/dir/out.wo"}
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected error for bad output path")
	}
	if !errors.Is(err, ErrOutput) {
		t.Errorf("got %v, want ErrOutput wrapper", err)
	}
}

func TestRunStdinMixed(t *testing.T) {
	srcPath := filepath.Join(t.TempDir(), "test.w")
	if err := os.WriteFile(srcPath, []byte("version 1\n.expose /tmp\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := &BuildCmd{Sources: []string{srcPath, "-"}}
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected error for mixed stdin and files")
	}
	if !errors.Is(err, ErrBuild) {
		t.Errorf("got %v, want ErrBuild wrapper", err)
	}
}
