package monitor

type Monitor interface {
	IncActiveJobsCounter()
	DecActiveJobsCounter()
	IncAllJobsCounter(username string)
	IncFailedJobsCounter(username, kind string)
}