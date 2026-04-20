package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/cruciblehq/crex"

	"ward"
	"ward/internal/parser"
)

const (

	// Marker for stdin input. When the sole source argument is "-", stdin is
	// used as the source reader.
	stdinMarker = "-"

	// Marker for stdout output. When the Output argument is empty, stdout is
	// used as the output writer.
	stdoutMarker = ""
)

// Subcommand that builds .w sources into a single .wo binary object.
//
// Each source is read from a file path. When a single source of "-" is given,
// stdin is used instead. Mixing "-" with explicit file paths is not allowed.
// Each source is built independently and the results are linked into one
// binary. Output goes to stdout by default, or to the path given by Output.
// When DryRun is set no output is produced.
type BuildCmd struct {
	DryRun  bool     `short:"n" help:"Validate without producing output."`
	Output  string   `short:"o" help:"Output file path." type:"path"`
	Sources []string `arg:"" help:"Affordance source files (.w), or - for stdin."`
}

// Runs the build pipeline.
//
// Each source is built into a binary object. The objects are linked together
// and written to the output. When DryRun is set, the sources are built, but
// the output is discarded.
func (c *BuildCmd) Run() error {
	if err := validateSources(c.Sources); err != nil {
		return crex.Wrap(ErrBuild, err)
	}

	objects := make([]io.Reader, 0, len(c.Sources))
	for _, src := range c.Sources {
		r, closer, err := input(src)
		if err != nil {
			return crex.Wrap(ErrParse, err)
		}

		var buf bytes.Buffer
		err = build(r, &buf)
		if closer != nil {
			closer.Close()
		}
		if err != nil {
			return crex.Wrap(ErrBuild, err)
		}

		objects = append(objects, &buf)
	}

	if c.DryRun {
		return nil
	}

	writer, closer, err := output(c.Output)
	if err != nil {
		return crex.Wrap(ErrOutput, err)
	}
	if closer != nil {
		defer closer.Close()
	}

	if err := ward.Link(writer, objects...); err != nil {
		return crex.Wrap(ErrBuild, err)
	}

	return nil
}

// Validates that sources are not mixed between stdin and file paths.
//
// When no sources are given, an error is returned. When multiple sources are
// given, an error is returned if any of them is the stdin marker.
func validateSources(sources []string) error {
	if len(sources) == 0 {
		return errors.New("no source files specified")
	}
	if len(sources) > 1 {
		if slices.Contains(sources, stdinMarker) {
			return fmt.Errorf("stdin (%s) cannot be mixed with file paths", stdinMarker)
		}
	}
	return nil
}

// Opens a source for reading.
//
// When path is "-", stdin is returned with no closer. Otherwise the file at
// path is opened and the caller is responsible for closing it.
func input(path string) (io.Reader, io.Closer, error) {
	if path == stdinMarker {
		return os.Stdin, nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return f, f, nil
}

// Opens the output for writing.
//
// When path is empty, stdout is returned with no closer. Otherwise the file at
// path is created or truncated and the caller is responsible for closing it.
func output(path string) (io.Writer, io.Closer, error) {
	if path == stdoutMarker {
		return os.Stdout, nil, nil
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, nil, err
	}
	return f, f, nil
}

// Builds .w source from a reader into a binary .wo object written to w.
func build(r io.Reader, w io.Writer) error {
	b := ward.NewBuilder()

	version, err := parser.Parse(r, b)
	if err != nil {
		return err
	}

	return b.Build(version, w)
}
