package cli

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/cruciblehq/crex"

	"ward/cmd/ward/internal"
)

// Top-level command structure parsed by Kong.
//
// Global flags (Quiet, Verbose, Debug) control log verbosity. Subcommands
// are embedded as struct fields and dispatched automatically by Kong.
var RootCmd struct {
	Quiet   bool     `short:"q" help:"Suppress informational output."`
	Verbose bool     `short:"v" help:"Enable verbose output."`
	Debug   bool     `short:"d" help:"Enable debug output."`
	Build   BuildCmd `cmd:"" help:"Build .w affordance files into a .wo object."`
}

// Executes the CLI.
//
// Parses the command-line arguments and dispatches to the selected subcommand.
// Configures the global logger based on CLI flags. Returns any error returned
// by the subcommand.
func Execute() error {
	kctx := kong.Parse(&RootCmd,
		kong.Name(internal.Name),
		kong.UsageOnError(),
		kong.Vars{
			"version": internal.VersionString(),
		},
	)

	configureLogger()

	return kctx.Run()
}

// Configures the global logger based on CLI flags.
//
// The log level is set to debug if the Debug flag is set, warn if the Quiet
// flag is set, and info otherwise. The formatter is set to a pretty human-
// friendly format when outputting to a terminal, and a compact JSON format
// otherwise. When the Verbose flag is set, the formatter includes caller and
// timestamp information.
func configureLogger() {
	handler, ok := slog.Default().Handler().(crex.Handler)
	if !ok {
		return
	}

	formatter := crex.NewPrettyFormatter(istty(os.Stderr))
	formatter.SetVerbose(RootCmd.Verbose)

	if RootCmd.Debug {
		handler.SetLevel(slog.LevelDebug)
	} else if RootCmd.Quiet {
		handler.SetLevel(slog.LevelWarn)
	} else {
		handler.SetLevel(slog.LevelInfo)
	}

	handler.SetFormatter(formatter)
	handler.SetStream(os.Stderr)
	handler.Flush()
}

// Whether the given file is an interactive terminal.
func istty(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
