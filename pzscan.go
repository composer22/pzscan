// pzscan is a simple SEO web crawler and link tester.
package main

import (
	"flag"
	"runtime"

	"github.com/composer22/pzscan/scanner"
)

// main is the main entry point for the application.
func main() {
	var hostname string
	var procs int
	var maxRunMin int
	var maxWorkers int
	var showVersion bool

	flag.StringVar(&hostname, "H", scanner.DefaultHostname, "Hostname to scan.")
	flag.StringVar(&hostname, "--hostname", scanner.DefaultHostname, "Hostname to scan.")
	flag.IntVar(&procs, "X", scanner.DefaultMaxProcs, "Maximum processor cores to use.")
	flag.IntVar(&procs, "--procs", scanner.DefaultMaxProcs, "Maximum processor cores to use.")
	flag.IntVar(&maxRunMin, "m", scanner.DefaultMaxMin, "Maximum minutes you want to run this routine.")
	flag.IntVar(&maxRunMin, "--minutes", scanner.DefaultMaxMin, "Maximum minutes you want to run this routine.")
	flag.IntVar(&maxWorkers, "W", scanner.DefaultMaxWorkers, "Maximum Job Workers.")
	flag.IntVar(&maxWorkers, "--workers", scanner.DefaultMaxWorkers, "Maximum Job Workers.")
	flag.BoolVar(&showVersion, "V", false, "Show version")
	flag.BoolVar(&showVersion, "--version", false, "Show version")
	flag.Usage = scanner.PrintUsageAndExit
	flag.Parse()

	// Version flag request?
	if showVersion {
		scanner.PrintVersionAndExit()
	}

	runtime.GOMAXPROCS(procs)
	s := scanner.New(hostname, maxRunMin, maxWorkers)
	s.Run()
}
