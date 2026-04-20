package internal

import (
	"strconv"
	"sync/atomic"
)

var (

	// Indicates whether quiet mode is enabled.
	quietMode atomic.Bool

	// Indicates whether debug logging is enabled.
	debugMode atomic.Bool
)

// Parses the linker flags into usable runtime variables.
//
// The rawQuiet and rawDebug variables should be set via ldflags during the
// build process. If not set, they default to "false".
func init() {
	if v, err := strconv.ParseBool(rawQuiet); err == nil {
		quietMode.Store(v)
	}
	if v, err := strconv.ParseBool(rawDebug); err == nil {
		debugMode.Store(v)
	}
}

// Reports whether quiet mode is enabled.
//
// Quiet mode suppresses informational output and raises the minimum log
// level to warn. It is set via the rawQuiet linker flag.
func IsQuiet() bool {
	return quietMode.Load()
}

// Reports whether debug mode is enabled.
//
// Debug mode lowers the minimum log level to debug. It is set via the
// rawDebug linker flag.
func IsDebug() bool {
	return debugMode.Load()
}
