// pzscan is a simple web crawler and link tester.
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
	flag.StringVar(&hostname, "h", "prezi.com", "Hostname of the host to scan (default: example.com).")
	flag.StringVar(&hostname, "--hostname", "prezi.com", "Hostname of the host to scan (default: example.com).")
	flag.IntVar(&procs, "X", 1, "Maximum processor cores to use (default: 1).")
	flag.IntVar(&procs, "--procs", 1, "Maximum processor cores to use (default: 1).")
	flag.IntVar(&maxRunMin, "m", 5, "Maximum minutes you want to run this routine (default: 5).")
	flag.IntVar(&maxRunMin, "--minutes", 5, "Maximum minutes you want to run this routine (default: 5).")
	flag.IntVar(&maxWorkers, "W", 4, "Maximum Job Workers (default: 4).")
	flag.IntVar(&maxWorkers, "--workers", 4, "Maximum Job Workers (default: 4).")
	flag.Parse()
	runtime.GOMAXPROCS(procs)
	s := scanner.New(hostname, maxRunMin, maxWorkers)
	s.Run()
	//	s.Dump() // We will post it as we test it.
}
