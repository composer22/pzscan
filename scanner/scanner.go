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
	mu         sync.Mutex        // For locking access.
	log        *logger.Logger    // Logger for writing final results.
	RootURL    *url.URL          // The original URL that we started the scan from.
	Tests      map[string]*Stats // URL test results go in here.
	MaxRunMin  int               // The Maximum number of minutes we want the scanner to run.
	MaxWorkers int               // The maximumm job workers we want in the pool.
	StartTime  time.Time         // When the scanner started runnning.
	ExpireTime time.Time         // The expire time: when the scanner should stop running.
	EndTime    time.Time         // When the scanner ended.
	jobq       chan *scanJob     // Channel to send jobs.
	doneCh     chan *scanJob     // Channel to receive done jobs.
	wg         sync.WaitGroup    // Synchronize close() of job channel.
	stopOnce   sync.Once         // Used to close down the system once and once only.
}

// New is a factory function that creates a new Scanner instance.
func New(hostname string, maxRunMin int, maxWorkers int) *Scanner {
	u, _ := url.Parse(fmt.Sprintf("http://%s", hostname))
	return &Scanner{
		log:        logger.New(logger.UseDefault, false),
		RootURL:    u,
		Tests:      make(map[string]*Stats),
		MaxRunMin:  maxRunMin,
		MaxWorkers: maxWorkers,
		jobq:       make(chan *scanJob, maxJobs),
		doneCh:     make(chan *scanJob, maxJobs),
	}
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
	s.jobq <- scanJobNew(s.RootURL) // Create first job
	for {
		select {
		case doneJob, ok := <-s.doneCh:
			if !ok {
				s.Stop()
			}
			s.evaluate(doneJob)
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
	// Store the result of this scan.
	if _, ok := s.Tests[job.Stat.URL.Path]; !ok {
		s.Tests[job.Stat.URL.String()] = job.Stat
		s.log.Infof(fmt.Sprint(job.Stat))
	}

	// Check for valid chidren under this host, and if any is not already scanned,
	// then submit new a job for this additional URL.
	for _, c := range job.Childen {
		if strings.Contains(c.Host, s.RootURL.Host) {
			if _, ok := s.Tests[c.String()]; !ok {
				s.jobq <- scanJobNew(c)
			}
		}
	}
}

// Dump prints the results to the log output as json.
func (s *Scanner) Dump() {
	for _, v := range s.Tests {
		s.log.Infof(fmt.Sprint(v))
	}
}
