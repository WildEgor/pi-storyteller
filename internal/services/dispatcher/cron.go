package dispatcher

import "github.com/robfig/cron/v3"

// DispatchCron represents a dispatched cron job
// that executes using cron expression formats.
type DispatchCron struct {
	cron *cron.Cron
}

// Stops ends the execution cycle for the given cron.
func (c *DispatchCron) Stop() {
	c.cron.Stop()
}
