package dispatcher

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"log/slog"
	"time"

	"errors"
	"sync"

	"github.com/WildEgor/pi-storyteller/internal/adapters/monitor"
)

// Dispatcher maintains a pool for available workers
// and a job queue that workers will process
type Dispatcher struct {
	maxHighWorkers int
	maxLowWorkers  int
	maxQueueLen    int
	minQueueLen    int
	workers        []*Worker
	lowWorkerPool  chan chan Job
	highWorkerPool chan chan Job
	tickers        []*DispatchTicker
	crons          []*DispatchCron
	lowQueue       chan Job
	highQueue      chan Job
	done           chan struct{}
	active         bool
	mu             sync.Mutex
	inProgressMap  map[string][]string
	metrics        monitor.Monitor
}

// NewDispatcher creates a new dispatcher with the given
// number of workers and buffers the job queue based on maxQueue.
// It also initializes the channels for the worker pool and job queue
func NewDispatcher(metrics monitor.Monitor) *Dispatcher {
	return &Dispatcher{
		// TODO: move to config
		maxHighWorkers: 10,
		maxLowWorkers:  1,
		maxQueueLen:    1000,
		minQueueLen:    100,
		done:           make(chan struct{}),
		workers:        make([]*Worker, 0),
		inProgressMap:  make(map[string][]string),
		metrics:        metrics,
	}

}

// Start creates and starts workers, adding them to the worker pool.
// Then, it starts a select loop to wait for job to be dispatched
// to available workers
func (d *Dispatcher) Start() {
	d.tickers = []*DispatchTicker{}
	d.crons = []*DispatchCron{}
	d.lowWorkerPool = make(chan chan Job, d.maxLowWorkers)
	d.lowQueue = make(chan Job, d.minQueueLen)
	d.highWorkerPool = make(chan chan Job, d.maxHighWorkers)
	d.highQueue = make(chan Job, d.maxQueueLen)

	for i := 0; i < d.maxLowWorkers; i++ {
		worker := NewWorker(d.lowWorkerPool)
		worker.Start()
		d.workers = append(d.workers, worker)
	}

	for i := 0; i < d.maxHighWorkers; i++ {
		worker := NewWorker(d.highWorkerPool)
		worker.Start()
		d.workers = append(d.workers, worker)
	}

	d.active = true

	go func() {
		for {
			select {
			case job := <-d.highQueue:
				go func(job Job) {
					jobChannel := <-d.highWorkerPool
					jobChannel <- job
				}(job)

			case job := <-d.lowQueue:
				go func(job Job) {
					jobChannel := <-d.lowWorkerPool
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
func (d *Dispatcher) Dispatch(fn handler, opts *JobOpts) (id string, err error) {
	if !d.active {
		return "", errors.New("dispatcher is not active")
	}

	newUUID := d.uuid()

	onDone := func(ctx JobCtx) {
		d.dequeue(ctx.Meta)
		d.metrics.DecActiveJobsCounter()
		d.metrics.IncAllJobsCounter(ctx.Meta.OwnerID)
	}

	onStart := func(ctx JobCtx) {
		d.metrics.IncActiveJobsCounter()
	}

	onFail := func(ctx JobCtx, err error) {
		// TODO: err to kind mapping
		d.metrics.IncFailedJobsCounter(ctx.Meta.OwnerID, monitor.ProblemKindUnknown)
	}

	job := Job{
		ID:      newUUID,
		handler: fn,
		Status:  StatusStarted,
		onStart: onStart,
		onDone:  onDone,
		onFail:  onFail,
		opts:    opts,
	}

	if opts != nil && opts.Priority == LowPriority {
		d.lowQueue <- job
	} else {
		d.highQueue <- job
	}

	d.enqueue(&job)

	return newUUID, nil
}

// DispatchCron pushes the given job into the job queue
// each time the cron definition is met, using the given location
func (d *Dispatcher) DispatchCron(fn handler, cronStr string, loc *time.Location) (*DispatchCron, error) {
	if !d.active {
		return nil, errors.New("dispatcher is not active")
	}

	dc := &DispatchCron{cron: cron.New(cron.WithSeconds(), cron.WithLocation(loc))}
	d.crons = append(d.crons, dc)

	newUUID := d.uuid()

	onDone := func(ctx JobCtx) {
		d.dequeue(ctx.Meta)
		d.metrics.DecActiveJobsCounter()
		d.metrics.IncAllJobsCounter(ctx.Meta.OwnerID)
	}

	onStart := func(ctx JobCtx) {
		d.metrics.IncActiveJobsCounter()
	}

	onFail := func(ctx JobCtx, err error) {
		// TODO: err to kind mapping
		d.metrics.IncFailedJobsCounter(ctx.Meta.OwnerID, monitor.ProblemKindUnknown)
	}

	cronID, err := dc.cron.AddFunc(cronStr, func() {
		d.highQueue <- Job{
			ID:      newUUID,
			handler: fn,
			Status:  StatusStarted,
			onStart: onStart,
			onDone:  onDone,
			onFail:  onFail,
			opts:    NewDefaultOpts(),
		}
	})

	if err != nil {
		return nil, errors.New("invalid cron definition")
	}

	dc.cron.Start()

	slog.Debug(fmt.Sprintf("cron started %d", cronID))
	dc.cron.Run()

	return dc, nil
}

// DispatchAt pushes the given job into the job queue
// at the given time
func (d *Dispatcher) DispatchAt(fn handler, at time.Time) error {
	if !d.active {
		return errors.New("dispatcher is not active")
	}

	go func() {
		now := time.Now()
		diff := at.Sub(now)

		if diff < 0 {
			return
		}

		newUUID := d.uuid()

		onDone := func(ctx JobCtx) {
			d.dequeue(ctx.Meta)
			d.metrics.DecActiveJobsCounter()
			d.metrics.IncAllJobsCounter(ctx.Meta.OwnerID)
		}

		onStart := func(ctx JobCtx) {
			d.metrics.IncActiveJobsCounter()
		}

		onFail := func(ctx JobCtx, err error) {
			// TODO: err to kind mapping
			d.metrics.IncFailedJobsCounter(ctx.Meta.OwnerID, monitor.ProblemKindUnknown)
		}

		time.Sleep(diff)
		d.highQueue <- Job{
			ID:      newUUID,
			handler: fn,
			Status:  StatusStarted,
			onStart: onStart,
			onDone:  onDone,
			onFail:  onFail,
			opts:    NewDefaultOpts(),
		}
	}()

	return nil
}

// CountActiveJobs ...
func (d *Dispatcher) CountActiveJobs(ownerId string) int {
	d.mu.Lock()
	defer d.mu.Unlock()

	if v, ok := d.inProgressMap[ownerId]; ok {
		return len(v)
	}

	return 0
}

// uuid ...
func (d *Dispatcher) uuid() string {
	//nolint
	newUUID, _ := uuid.NewUUID()
	return newUUID.String()
}

// enqueue ...
func (d *Dispatcher) enqueue(job *Job) {
	d.mu.Lock()
	defer d.mu.Unlock()

	//nolint
	v, _ := d.inProgressMap[job.opts.OwnerID]
	v = append(v, job.ID)
	d.inProgressMap[job.opts.OwnerID] = lo.Uniq(v)
}

// dequeue ...
func (d *Dispatcher) dequeue(meta *JobMeta) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if v, ok := d.inProgressMap[meta.OwnerID]; ok {
		v = lo.Filter(v, func(item string, _ int) bool {

			return item != meta.ID

		})
		d.inProgressMap[meta.OwnerID] = v
	}
}
