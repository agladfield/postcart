package cards

import (
	"context"
	"fmt"
	"sync"

	"github.com/agladfield/postcart/pkg/shared/env"
	"golang.org/x/time/rate"
)

const (
	numWorkers    = 10
	jobsPerMinute = 20
	queueSize     = 1_000
)

var jobCh chan EmailParams

func AddToQueue(job EmailParams) {
	jobCh <- job
}

func createQueue(ctx context.Context, wg *sync.WaitGroup) {
	jobCh = make(chan EmailParams, queueSize)

	// for queued up/missed jobs

	var limiter *rate.Limiter
	if env.UseAI() {
		limiter = rate.NewLimiter(rate.Limit(float64(jobsPerMinute)/60), numWorkers)
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case job, ok := <-jobCh:
					if !ok {
						return
					}
					wg.Add(1)
					defer wg.Done()

					if env.UseAI() {
						limitErr := limiter.Wait(ctx)
						if limitErr != nil {
							// log or some shit
							limitErr = fmt.Errorf(cardsErrFmtStr, limitErr)
							return
						}
					}

					jobCtx := context.WithValue(ctx, "job-id", job.ID)
					jobErr := processJob(jobCtx, job)
					if jobErr != nil {
						jobErr = fmt.Errorf(cardsErrFmtStr, jobErr)
						// record job error
						fmt.Printf("ERROR OCCURED: %v\n", jobErr)
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}
