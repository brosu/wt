package color

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// ANSI escape code constants.
const (
	Reset   = "\033[0m"
	Bold    = "1"
	Dim     = "2"
	Red     = "31"
	Green   = "32"
	Yellow  = "33"
	Cyan    = "36"
	CyanRaw = "36" // for combining with other codes
)

// Colorize wraps s in ANSI escape codes. The code parameter can be a single
// code (e.g. "31" for red) or combined codes separated by semicolons
// (e.g. "1;36" for bold cyan).
func Colorize(s string, code string) string {
	return fmt.Sprintf("\033[%sm%s%s", code, s, Reset)
}

// IsEnabled returns true if color output should be used.
// Color is disabled when:
//   - the NO_COLOR environment variable is set (any value, per https://no-color.org/)
//   - stdout is not a terminal (i.e., output is piped)
func IsEnabled() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	return term.IsTerminal(int(os.Stdout.Fd()))
}
