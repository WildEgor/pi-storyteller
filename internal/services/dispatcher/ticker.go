package dispatcher

import "time"

// DispatchTicker represents a dispatched job ticker
// that executes on a given interval. This provides
// a means for stopping the execution cycle from continuing.
type DispatchTicker struct {
	ticker *time.Ticker
	quit   chan bool
}

// Stop ends the execution cycle for the given ticker.
func (dt *DispatchTicker) Stop() {
	dt.ticker.Stop()
	dt.quit <- true
}
