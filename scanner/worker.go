package scanner

import (
	"net/http"
	"sync"
	"time"
)

var (
	workerMaxSleep = 100 * time.Millisecond // How long should a worker sleep between jobs peeks.
)

// scanWorker is used as a go routine wrapper to handle URL scan jobs.
func scanWorker(jobq chan *scanJob, doneCh chan *scanJob, wg *sync.WaitGroup) {
	defer wg.Done()
	// prepare client
	client := &http.Client{}
	for {
		select {
		case job, ok := <-jobq:
			if !ok {
				return // Assume closed channel.
			}
			// Scan the link.
			job.Stat.StartTime = time.Now()
			resp, err := client.Get(job.Stat.URL.String())
			job.Stat.EndTime = time.Now()
			if err != nil {
				job.Stat.StatusCode = -1 // We couldn't even get a HTTP status code.
			} else {
				job.Stat.StatusCode = resp.StatusCode
				job.Childen = AnalyzePage(resp.Body, job.Stat)
			}
			doneCh <- job
		default:
			time.Sleep(workerMaxSleep) // Sleep before peeking again.
		}
	}
}
