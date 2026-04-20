// Package cli implements the ward command-line interface.
//
// The CLI uses Kong for argument parsing. The build subcommand compiles
// one or more .w affordance sources into .wo binary objects, links them
// into a single output, and writes the result to a file or stdout. The
// -n flag validates without producing output. Stdin is used when the sole
// source argument is "-".
package cli
