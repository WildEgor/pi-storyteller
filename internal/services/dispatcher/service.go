package dispatcher

import (
	"errors"
	"github.com/google/uuid"
	"time"

	"github.com/robfig/cron/v3"
)

// Dispatcher maintains a pool for available workers
// and a job queue that workers will process
type Dispatcher struct {
	maxWorkers int
	maxQueue   int
	workers    []*Worker
	tickers    []*DispatchTicker
	crons      []*DispatchCron
	workerPool chan chan Job
	jobQueue   chan Job
	done       chan struct{}
	active     bool
}

// NewDispatcher creates a new dispatcher with the given
// number of workers and buffers the job queue based on maxQueue.
// It also initializes the channels for the worker pool and job queue
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		// TODO: move to config
		maxWorkers: 10,
		maxQueue:   1000,
		done:       make(chan struct{}),
		workers:    make([]*Worker, 0),
		tickers:    make([]*DispatchTicker, 0),
		crons:      make([]*DispatchCron, 0),
	}
}

// Start creates and starts workers, adding them to the worker pool.
// Then, it starts a select loop to wait for job to be dispatched
// to available workers
func (d *Dispatcher) Start() {
	d.workerPool = make(chan chan Job, d.maxWorkers)
	d.jobQueue = make(chan Job, d.maxQueue)

	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.workerPool)
		worker.Start()
		d.workers = append(d.workers, worker)
	}

	d.active = true

	go func() {
		for {
			select {
			case job := <-d.jobQueue:
				go func(job Job) {
					jobChannel := <-d.workerPool
					jobChannel <- job
				}(job)
			case <-d.done:
				return
			}
		}
	}()
}

// Stop ends execution for all workers/tickers and
// closes all channels, then removes all workers/tickers
func (d *Dispatcher) Stop() {
	if !d.active {
		return
	}

	d.active = false

	for i := range d.workers {
		d.workers[i].Stop()
	}

	for i := range d.tickers {
		d.tickers[i].Stop()
	}

	for i := range d.crons {
		d.crons[i].Stop()
	}

	d.workers = []*Worker{}
	d.tickers = []*DispatchTicker{}
	d.crons = []*DispatchCron{}
	d.done <- struct{}{}
}

// Dispatch pushes the given job into the job queue.
// The first available worker will perform the job
func (d *Dispatcher) Dispatch(run func()) (id string, err error) {
	if !d.active {
		return "", errors.New("dispatcher is not active")
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	d.jobQueue <- Job{
		ID:  newUUID.String(),
		Run: run,
	}
	return newUUID.String(), nil
}

// DispatchAt pushes the given job into the job queue
// at the given time
func (d *Dispatcher) DispatchAt(run func(), at time.Time) error {
	if !d.active {
		return errors.New("dispatcher is not active")
	}

	go func() {
		now := time.Now()
		diff := at.Sub(now)

		if diff < 0 {
			return
		}

		time.Sleep(diff)
		d.jobQueue <- Job{Run: run}
	}()

	return nil
}

// DispatchCron pushes the given job into the job queue
// each time the cron definition is met
func (d *Dispatcher) DispatchCron(run func(), cronStr string) (*DispatchCron, error) {
	if !d.active {
		return nil, errors.New("dispatcher is not active")
	}

	dc := &DispatchCron{cron: cron.New(cron.WithSeconds())}
	d.crons = append(d.crons, dc)

	_, err := dc.cron.AddFunc(cronStr, func() {
		d.jobQueue <- Job{Run: run}
	})

	if err != nil {
		return nil, errors.New("invalid cron definition")
	}

	dc.cron.Start()
	return dc, nil
}
