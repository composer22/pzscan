package scanner

import (
	"fmt"
	"os"
)

const usageStr = `
Description: scan a host for embedded anchor links.

Usage: pzscan [options...]

Server options:
    -H, --hostname HOSTNAME          HOSTNAME to scan (default: example.com).
    -X, --procs MAX                  MAX processor cores to use from the
	                                 machine (default 1).
    -m, --minutes MAX                MAX minutes to live (default: 5).
    -W, --workers MAX                MAX running workers allowed (default: 4).

Common options:
    -h, --help                       Show this message

Example:

    # Scan craigslist; 1 processor; 2 min max; 10 worker go routines.

    ./pzscan -H "craigslist.com" -X 1 -m 2 -W 10
`

// end help text

// PrintUsageAndExit is used to print out command line options.
func PrintUsageAndExit() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}
