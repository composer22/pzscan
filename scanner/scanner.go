package scanner

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/composer22/pzscan/logger"
)

const (
	maxJobs = 10000 // The jobq maximum number of jobs to hold. We need something for non blocking.
)

var (
	maxIdleDuration = 15 * time.Second       // How long we should wait on empty queues before auto quitting.
	maxScannerSleep = 100 * time.Millisecond // How long should the scanner sleep before checking for results.
)

// Scanner is a manager of scanning jobs and evaluates the results of the workers.
type Scanner struct {
	RootURL    *url.URL                     // The original URL that we started the scan from.
	Tests      map[string]map[string]*Stats // URL test results go in here.
	MaxRunMin  int                          // The Maximum number of minutes we want the scanner to run.
	MaxWorkers int                          // The maximumm job workers we want in the pool.
	StartTime  time.Time                    // When the scanner started runnning.
	ExpireTime time.Time                    // The expire time: when the scanner should stop running.
	EndTime    time.Time                    // When the scanner ended.
	mu         sync.Mutex                   // For locking access.
	wg         sync.WaitGroup               // Synchronize close() of job channel.
	stopOnce   sync.Once                    // Used to close down the system once and once only.
	log        *logger.Logger               // Logger for writing final results.
	jobq       chan *scanJob                // Channel to send jobs.
	doneCh     chan *scanJob                // Channel to receive done jobs.
}

// New is a factory function that creates a new Scanner instance.
func New(hostname string, maxRunMin int, maxWorkers int) *Scanner {
	u, _ := url.Parse(fmt.Sprintf("http://%s", hostname))
	return &Scanner{
		RootURL:    u,
		Tests:      make(map[string]map[string]*Stats),
		MaxRunMin:  maxRunMin,
		MaxWorkers: maxWorkers,
		log:        logger.New(logger.UseDefault, false),
		jobq:       make(chan *scanJob, maxJobs),
		doneCh:     make(chan *scanJob, maxJobs),
	}
}

// PrintVersionAndExit prints the version of the scanner then exits.
func PrintVersionAndExit() {
	fmt.Printf("pzscan version %s\n", version)
	os.Exit(0)
}

// Run starts the scanner and manages the jobs.
func (s *Scanner) Run() {
	// Trap all signals to quit.
	s.handleSignals()

	s.mu.Lock()

	// Spin up the workers
	for i := 0; i < s.MaxWorkers; i++ {
		s.wg.Add(1)
		go scanWorker(s.jobq, s.doneCh, &s.wg)
	}

	s.StartTime = time.Now()
	s.ExpireTime = s.StartTime.Add(time.Duration(s.MaxRunMin) * time.Minute)
	s.mu.Unlock()

	// Main event loop.
	p, _ := url.Parse(fmt.Sprintf("http://%s", ""))
	s.jobq <- scanJobNew(s.RootURL, "html", p) // Create first job.  Assume its a page.
	for {
		select {
		case j, ok := <-s.doneCh:
			if !ok {
				s.Stop()
				return
			}
			s.evaluate(j)
		default:
			// Drop dead time reached?
			if time.Now().After(s.ExpireTime) {
				s.Stop()
				return
			}
			// Test a reasonable time for all jobs to be cleared, and then die if no more
			// jobs need to be done or are returned.
			if len(s.jobq) == 0 && len(s.doneCh) == 0 {
				time.Sleep(maxIdleDuration)
				if len(s.doneCh) == 0 {
					s.Stop()
					return
				}
			}
			time.Sleep(maxScannerSleep) // Sleep a while.
		}
	}
}

// Stop performs close out procedures.
func (s *Scanner) Stop() {
	s.stopOnce.Do(func() {
		s.EndTime = time.Now()
		close(s.doneCh)
		close(s.jobq)
		s.wg.Wait()
	})
}

// handleSignals responds to operating system interrupts such as application kills.
func (s *Scanner) handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			s.Stop()
			os.Exit(0)
		}
	}()
}

// evaluate examines the result of the job and launches new jobs if site children are found.
func (s *Scanner) evaluate(job *scanJob) {
	pURL := job.Stat.ParentURL.String()
	cURL := job.Stat.URL.String()
	// Initialize result slot
	if _, ok := s.Tests[cURL]; !ok {
		s.Tests[cURL] = make(map[string]*Stats)
	}

	// Store the result of this scan and print it to the log.
	if _, ok := s.Tests[cURL][pURL]; !ok {
		s.Tests[cURL][pURL] = job.Stat
		s.log.Infof(fmt.Sprint(job.Stat))
	}
	// Check for any URL's returned and create new jobs.
	for _, c := range job.Children {
		// No Scheme?  Assume http:
		if c.URL.Scheme == "" {
			c.URL.Scheme = "http"
		}
		// No host?  Assume its us.
		if c.URL.Host == "" {
			c.URL.Host = s.RootURL.Host
		}
		switch c.URLType {
		case "html":
			// Don't scan foreign pages.
			if !strings.Contains(c.URL.Host, s.RootURL.Host) {
				continue
			}
			// If we haven't scanned this url, do it. [new][sourcepage]
			if _, ok := s.Tests[c.URL.String()][cURL]; !ok {
				s.jobq <- scanJobNew(c.URL, c.URLType, job.Stat.URL)
			}
		default:
			// If it is a site asset
			if strings.Contains(c.URL.Host, s.RootURL.Host) {
				// If we haven't scanned this asset, do it.
				if _, ok := s.Tests[c.URL.String()]; !ok {
					s.jobq <- scanJobNew(c.URL, c.URLType, job.Stat.URL)
				}
			} else { // Foreign asset
				// If we haven't scanned this url, do it. [new][sourcepage]
				if _, ok := s.Tests[c.URL.String()][cURL]; !ok {
					s.jobq <- scanJobNew(c.URL, c.URLType, job.Stat.URL)
				}
			}
		}
	}
}
