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
	cl := &http.Client{}
	a := bodyAnalyzerNew(nil)
	for {
		select {
		case j, ok := <-jobq:
			if !ok {
				return // Assume closed channel.
			}
			// Scan the link.
			j.Stat.StartTime = time.Now()
			resp, err := cl.Get(j.Stat.URL.String())
			j.Stat.EndTime = time.Now()
			if err != nil {
				j.Stat.StatusCode = -1 // We couldn't even get a HTTP status code.
			} else {
				j.Stat.StatusCode = resp.StatusCode
				if j.Stat.URLType == "html" {
					j.Body = resp.Body
					a.ScanJob = j
					a.analyzeBody()
					j.Body.Close()
				}
			}
			doneCh <- j
		default:
			time.Sleep(workerMaxSleep) // Sleep before peeking again.
		}
	}
}
