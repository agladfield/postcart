// Package jdb is a simple json data keeping system for a project of this size
package jdb

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

func Close() {
	recipientsDB.save()
	sendersDB.save()
	blockedSendersDB.save()
	statsDB.save()
	queuedJobsDB.save()
}

const jdbErrFmtStr = "jdb err: %w"

func Load() error {
	dataDirErr := os.MkdirAll("./data/attachments", 0700)
	if dataDirErr != nil {
		return dataDirErr
	}
	recipientsDB = &jsondb[int8Map]{path: "./data/recip.json"}
	sendersDB = &jsondb[int8Map]{path: "./data/senders.json"}
	blockedSendersDB = &jsondb[intMap]{path: "./data/blocked.json"}
	statsDB = &jsondb[intMap]{path: "./data/stats.json"}
	queuedJobsDB = &jsondb[jobMap]{path: "./data/queue.json"}

	recipErr := recipientsDB.load()
	if recipErr != nil {
		return fmt.Errorf(jdbErrFmtStr, recipErr)
	}
	sendersErr := sendersDB.load()
	if sendersErr != nil {
		return fmt.Errorf(jdbErrFmtStr, sendersErr)
	}
	blockedErr := blockedSendersDB.load()
	if blockedErr != nil {
		return fmt.Errorf(jdbErrFmtStr, blockedErr)
	}
	statsErr := statsDB.load()
	if statsErr != nil {
		return fmt.Errorf(jdbErrFmtStr, statsErr)
	}
	queueErr := queuedJobsDB.load()
	if queueErr != nil {
		return fmt.Errorf(jdbErrFmtStr, queueErr)
	}

	return nil
}

type jsondb[T any] struct {
	data T
	path string
	mu   sync.RWMutex
}

type craftableMap[T any] interface {
	Make() any
}

func (db *jsondb[T]) load() error {
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		// Code to execute if the file does not exist
		var emptyV T
		newErr := json.Unmarshal([]byte("{}"), &emptyV)
		if newErr != nil {
			return fmt.Errorf(jdbErrFmtStr, newErr)
		}
		db.data = emptyV
	} else if err == nil {
		fmt.Printf("File '%s' exists\n", db.path)
		// Code to execute if the file exists
		bytes, bytesErr := os.ReadFile(db.path)
		if bytesErr != nil {
			return fmt.Errorf(jdbErrFmtStr, bytesErr)
		}
		var existV T
		newErr := json.Unmarshal(bytes, &existV)
		if newErr != nil {
			return fmt.Errorf(jdbErrFmtStr, newErr)
		}
		db.data = existV
	} else {
		return fmt.Errorf(jdbErrFmtStr, fmt.Errorf("error checking jdb file: %w", err))
	}
	return nil
}

func (db *jsondb[T]) save() error {
	bytes, jsonErr := json.Marshal(db.data)
	if jsonErr != nil {
		return fmt.Errorf(jdbErrFmtStr, jsonErr)
	}

	writeErr := os.WriteFile(db.path, bytes, 0700)
	if writeErr != nil {
		return fmt.Errorf(jdbErrFmtStr, writeErr)
	}

	return nil
}

type int8Map map[string]int8
type intMap map[string]int
type jobMap map[string]JobRecord

var (
	recipientsDB     *jsondb[int8Map]
	sendersDB        *jsondb[int8Map]
	blockedSendersDB *jsondb[intMap]
	statsDB          *jsondb[intMap]
	queuedJobsDB     *jsondb[jobMap]
	attachmentsDB    *jsondb[int8Map]
)

func BlockRecipient(email string) {
	recipientsDB.mu.Lock()
	defer recipientsDB.mu.Unlock()
	recipientsDB.data[email] = 1
}

func IsRecipientBlocked(email string) bool {
	recipientsDB.mu.RLock()
	defer recipientsDB.mu.RUnlock()
	val, exists := recipientsDB.data[email]
	if !exists {
		return false
	}

	return val > 0
}

func IncrementSender(email string) int8 {
	sendersDB.mu.Lock()
	defer sendersDB.mu.Unlock()
	sendersDB.data[email]++

	return sendersDB.data[email]
}

func BlockSender(email string, ruleID int) {
	blockedSendersDB.mu.Lock()
	defer blockedSendersDB.mu.Unlock()
	blockedSendersDB.data[email] = ruleID
}

func IsSenderBlocked(email string) bool {
	blockedSendersDB.mu.RLock()
	defer blockedSendersDB.mu.RUnlock()
	val, exists := blockedSendersDB.data[email]
	if !exists {
		return false
	}

	return val > 0
}

func RecordInbound() {
	statsDB.mu.Lock()
	defer statsDB.mu.Unlock()
	statsDB.data["inbound"]++
}

func RecordSent() {
	statsDB.mu.Lock()
	defer statsDB.mu.Unlock()
	statsDB.data["sent"]++
}

func RecordError() {
	statsDB.mu.Lock()
	defer statsDB.mu.Unlock()
	statsDB.data["errors"]++
}

func RecordRejection() {
	statsDB.mu.Lock()
	defer statsDB.mu.Unlock()
	statsDB.data["rejections"]++
}

func RecordRetry() {
	statsDB.mu.Lock()
	defer statsDB.mu.Unlock()
	statsDB.data["retries"]++
}

func RecordBounce() {
	statsDB.mu.Lock()
	defer statsDB.mu.Unlock()
	statsDB.data["bounces"]++
}

func RecordSpamComplaint() {
	statsDB.mu.Lock()
	defer statsDB.mu.Unlock()
	statsDB.data["spamc"]++
}

func RecordDelivery() {
	statsDB.mu.Lock()
	defer statsDB.mu.Unlock()
	statsDB.data["deliveries"]++
}

func RecordBlockedSender() {
	statsDB.mu.Lock()
	defer statsDB.mu.Unlock()
	statsDB.data["blocked_senders"]++
}

func RecordBlockedRecipient() {
	statsDB.mu.Lock()
	defer statsDB.mu.Unlock()
	statsDB.data["blocked_recipients"]++
}

func RecordQueueSize(q int) {
	queuedJobsDB.mu.RLock()
	defer queuedJobsDB.mu.RUnlock()
	statsDB.mu.Lock()
	defer statsDB.mu.Unlock()
	statsDB.data["queue"] = q
}

func IncrementQueueSize() {
	statsDB.mu.Lock()
	defer statsDB.mu.Unlock()
	statsDB.data["queue"]++
}
func DecrementQueueSize() {
	statsDB.mu.Lock()
	defer statsDB.mu.Unlock()
	statsDB.data["queue"]--
}

func GetStats() map[string]int {
	statsDB.mu.RLock()
	defer statsDB.mu.RUnlock()
	return statsDB.data
}

type JobRecord struct {
	ID             string `json:"id"`
	ToEmail        string `json:"to_email"`
	ToName         string `json:"to_name,omitempty"`
	FromEmail      string `json:"from_email"`
	FromName       string `json:"from_name,omitempty"`
	Artwork        int8   `json:"artwork"`
	Style          int8   `json:"style"`
	Font           int8   `json:"font"`
	Border         int8   `json:"border"`
	StampShape     int8   `json:"stamp"`
	Textured       int8   `json:"textured"`
	Country        string `json:"country"`
	Subject        string `json:"subject"`
	Message        string `json:"message"`
	AttachmentType string `json:"attachment_type,omitempty"`
}

func RecordQueuedJob(job JobRecord) error {
	queuedJobsDB.mu.Lock()
	defer queuedJobsDB.mu.Unlock()
	queuedJobsDB.data[job.ID] = job
	saveErr := queuedJobsDB.save()
	if saveErr != nil {
		return fmt.Errorf(jdbErrFmtStr, saveErr)
	}
	return nil
}

func RemoveJobFromRecords(jobID string) {
	queuedJobsDB.mu.Lock()
	defer queuedJobsDB.mu.Unlock()
	delete(queuedJobsDB.data, jobID)
}

func GetUncompletedQueuedJobs() map[string]JobRecord {
	queuedJobsDB.mu.RLock()
	defer queuedJobsDB.mu.RUnlock()
	return queuedJobsDB.data
}

const AttachmentFmtStr = "./data/attachments/attachment-%s"

func RemoveAttachmentForJob(jobID string) error {
	remErr := os.Remove(fmt.Sprintf(AttachmentFmtStr, jobID))
	if remErr != nil {
		return fmt.Errorf(jdbErrFmtStr, remErr)
	}
	return nil
}

// © Arthur Gladfield
