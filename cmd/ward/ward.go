package main

import (
	"log/slog"
	"os"

	"github.com/cruciblehq/crex"

	"ward/cmd/ward/internal"
	"ward/cmd/ward/internal/cli"
)

// Entry point for the ward CLI.
//
// Initialises structured logging, then delegates to the CLI dispatcher. Exits
// with code 1 if any stage returns an error.
func main() {
	setUpLogger()

	slog.Debug("build",
		"version", internal.VersionString(),
	)

	slog.Debug("ward is running",
		"pid", os.Getpid(),
		"cwd", cwd(),
		"args", os.Args,
	)

	if err := cli.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

// Sets up a buffered logger.
//
// The initial log level comes from build-time linker flags (debug or quiet).
// After Kong finishes parsing, the CLI reconfigures and flushes the handler
// so that user-provided flags take effect.
func setUpLogger() {
	handler := crex.NewHandler()
	handler.SetLevel(logLevel())
	logger := slog.New(handler.WithGroup(internal.Name))
	slog.SetDefault(logger)
}

// Returns the log level derived from build-time linker flags.
//
// Debug builds default to debug level, quiet builds default to warn level, and
// all others default to info level.
func logLevel() slog.Level {
	if internal.IsDebug() {
		return slog.LevelDebug
	}
	if internal.IsQuiet() {
		return slog.LevelWarn
	}
	return slog.LevelInfo
}

// Returns the current working directory or "(unknown)".
func cwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return "(unknown)"
	}
	return dir
}
