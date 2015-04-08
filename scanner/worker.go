package scanner

import (
	"net/http"
	"sync"
	"time"
)

var (
	workerMaxSleep = 250 * time.Millisecond // How long should a worker sleep between jobq peeks.
)

// scanWorker is used as a go routine wrapper to handle URL scan jobs.
func scanWorker(jobq chan *scanJob, doneCh chan *scanJob, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}
	analyzer := bodyAnalyzerNew(nil)
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
				if job.Stat.URLType == "html" {
					job.Body = resp.Body
					analyzer.ScanJob = job
					analyzer.analyzeBody()
					job.Body.Close()
				}
			}
			doneCh <- job
		default:
			time.Sleep(workerMaxSleep) // Sleep before peeking again.
		}
	}
}
