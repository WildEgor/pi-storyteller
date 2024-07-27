package dispatcher

import "time"

// Dispatcher maintains a pool for available workers
// and a job queue that workers will process
type Dispatcher interface {
	Start()
	Stop()
	Dispatch(fn handler, opts *JobOpts) (id string, err error)
	DispatchCron(fn handler, cronStr string, loc *time.Location) (*DispatchCron, error)
	DispatchAt(fn handler, at time.Time) error
	CountActiveJobs(ownerId string) int
}
