package dispatcher

import (
	"context"
	"time"
)

// Status of Job
type Status string

const (
	// StatusStarted This is the initial state when a job is pushed onto the broker.
	StatusStarted Status = "queued"
	// StatusProcessing This is the state when a worker has received a job.
	StatusProcessing Status = "processing"
	// StatusFailed The state when a job completes, but returns an error (and all retries are over).
	StatusFailed Status = "failed"
	// StatusDone The state when a job completes without any error.
	StatusDone Status = "successful"
)

// Priority of Job
type Priority uint

var (
	// LowPriority the lowest priority (by default)
	LowPriority Priority
	// HighPriority the highest priority
	HighPriority Priority = 3
)

// JobCtx for callback
type JobCtx struct {
	context.Context

	Meta *JobMeta
}

// JobMeta some info
type JobMeta struct {
	ID        string
	OwnerID   string
	StartedAt time.Time
}

// handler ...
type handler func(JobCtx) error

// onHandler ...
type onHandler func(JobCtx)

// onHandler ...
type onFailHandler func(JobCtx, error)

// JobOpts settings
type JobOpts struct {
	Name         string
	At           time.Time
	Priority     Priority
	OwnerID      string
	OnSuccess    func(JobCtx)
	OnProcessing func(JobCtx)
	OnFail       func(JobCtx, error)
}

// NewDefaultOpts ...
func NewDefaultOpts() *JobOpts {
	return &JobOpts{
		OnSuccess: func(ctx JobCtx) {
		},
		OnProcessing: func(ctx JobCtx) {
		},
		OnFail: func(ctx JobCtx, err error) {
		},
	}
}

// Job represents a runnable process, where Start
// will be executed by a worker via the dispatch queue
type Job struct {
	ID      string
	Status  Status
	handler handler
	onStart onHandler
	onDone  onHandler
	onFail  onFailHandler
	opts    *JobOpts
}
