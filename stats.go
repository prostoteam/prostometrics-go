package prostometrics

// Stats is a lock-free snapshot of client-side delivery health.
type Stats struct {
	QueueDropped       uint64
	InvalidDropped     uint64
	RateLimitedDropped uint64
	RetryDropped       uint64
	QueueDepth         int
	IngestDisabled     bool
}

// Stats returns current client-side delivery statistics.
func (c *Client) Stats() Stats {
	if c == nil {
		return Stats{}
	}
	return Stats{
		QueueDropped:       c.dropped.Load(),
		InvalidDropped:     c.invalidDropped.Load(),
		RateLimitedDropped: c.rateLimitedDropped.Load(),
		RetryDropped:       c.retryDropped.Load(),
		QueueDepth:         len(c.queue),
		IngestDisabled:     c.stopSending.Load(),
	}
}
