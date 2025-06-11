package cards

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/agladfield/postcart/pkg/jdb"
	"github.com/agladfield/postcart/pkg/postmark"
	"github.com/agladfield/postcart/pkg/shared/enum"
	"github.com/agladfield/postcart/pkg/shared/env"
	"golang.org/x/time/rate"
)

const (
	queueSize     = 1_000
	genAIJobLimit = 20 // vertex ai rate limit is 20 requests/minute
)

var (
	numWorkers    = runtime.NumCPU()
	jobsPerMinute = 120 // one job every 500ms
	jobCh         chan Params
)

const cardsQueueErrFmtStr = "cards queue err: %w"

func AddToQueue(job Params) error {
	jdbRecord := job.toJDBJobRecord()
	queueSaveErr := jdb.RecordQueuedJob(jdbRecord)
	if queueSaveErr != nil {
		return fmt.Errorf(cardsQueueErrFmtStr, fmt.Errorf("failed to record job to queue: %w", queueSaveErr))
	}
	if job.Attachment != nil {
		writeErr := os.WriteFile(fmt.Sprintf(jdb.AttachmentFmtStr, job.ID), []byte(job.Attachment.Content), 0400)
		if writeErr != nil {
			return fmt.Errorf(cardsQueueErrFmtStr, writeErr)
		}
	}
	jdb.IncrementQueueSize()
	jobCh <- job
	return nil
}

func queueJobRecordToJob(record *jdb.JobRecord) (*Params, error) {
	var attachment *postmark.EmailAttachment
	if record.AttachmentType != "" {
		readBytes, readErr := os.ReadFile(fmt.Sprintf(jdb.AttachmentFmtStr, record.ID))
		if readErr != nil {
			return nil, readErr
		}
		attachment = &postmark.EmailAttachment{
			Content:     string(readBytes),
			ContentType: record.AttachmentType,
		}
	}

	params := Params{ID: record.ID,
		To: Person{
			Name:  record.ToName,
			Email: record.ToEmail,
		},
		From: Person{
			Name:  record.FromName,
			Email: record.FromEmail,
		},
		Artwork:    enum.ArtworkEnum(record.Artwork),
		Style:      enum.StyleEnum(record.Style),
		Font:       enum.FontEnum(record.Font),
		Border:     enum.BorderEnum(record.Border),
		StampShape: enum.StampShapeEnum(record.StampShape),
		Textured:   enum.TexturedEnum(record.Textured),
		Country:    record.Country,
		Subject:    record.Subject,
		Message:    record.Message,
		Attachment: attachment,
	}
	return &params, nil
}

func createQueue(ctx context.Context, wg *sync.WaitGroup) {
	jobCh = make(chan Params, queueSize)

	// for queued up/missed jobs
	queuedJobs := jdb.GetUncompletedQueuedJobs()
	for _, qJobRecord := range queuedJobs {
		qJob, qJobErr := queueJobRecordToJob(&qJobRecord)
		if qJobErr != nil {
			// log error
			fmt.Println(fmt.Errorf(cardsQueueErrFmtStr, qJobErr))
			continue
		}
		jobCh <- *qJob
	}
	jdb.RecordQueueSize(len(jobCh))

	if env.UseAI() {
		jobsPerMinute = genAIJobLimit
	}

	burstSize := numWorkers
	if jobsPerMinute < burstSize {
		burstSize = jobsPerMinute
	}
	limiter := rate.NewLimiter(rate.Limit(float64(jobsPerMinute)/60), burstSize)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()
			for {
				if err := limiter.Wait(ctx); err != nil {
					fmt.Printf("rate limiter error occured [thread %d]: %v\n", fmt.Errorf(cardsErrFmtStr, err), i+1)
					return
				}
				select {
				case job, ok := <-jobCh:
					if !ok {
						return
					}
					fmt.Printf("processing email %q [thread %d]\n", job.ID, i+1)

					jobCtx := context.WithValue(ctx, "job-id", job.ID)
					jobErr := processJob(jobCtx, job)
					if jobErr != nil {
						// record job error
						jdb.RecordError()
						fmt.Printf("failed email %q job [thread %d] with error: %v\n", job.ID, i+1, fmt.Errorf(cardsErrFmtStr, jobErr))
					} else {
						fmt.Printf("finished email %q [thread %d]\n", job.ID, i+1)

					}
					jdb.DecrementQueueSize()
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}
}

// © Arthur Gladfield
