# Prostometrics Go SDK

Go SDK for sending application metrics to Prostometrics.

## Install

```bash
go get github.com/prostoteam/prostometrics-go@latest
```

## Quick start

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/prostoteam/prostometrics-go"
)

func main() {
	client, err := prostometrics.NewClient("payments-api", prostometrics.Config{
		APIKey: "your-api-key",
	})
	if err != nil {
		log.Fatal(err)
	}

	client.Count("requests", 1, prostometrics.Label("method", "GET"))
	client.CountUnique(uint64(42), "daily_active_users", "plan=pro")
	client.Total("bytes_sent_kb", 2048, "interface=eth0")
	client.Value("latency_ms", 123.4, "route=/login")
	client.ValueSparse("capacity_kb", 1024*1024, "mount=/")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = client.Close(ctx)
}
```

## Metric methods

- `Count` increases a counter.
- `CountUnique` counts distinct identifiers approximately.
- `Total` records the current value of a monotonically increasing total.
- `Value` records a numeric observation.
- `ValueSparse` records a value whose last observation should carry across missing time buckets.

Labels use `name=value` syntax. Use `prostometrics.Label(name, value)` when label values are assembled dynamically.

## Configuration

| Field | Default | Purpose |
|---|---|---|
| `APIKey` | empty | Prostometrics API key |
| `Logger` | no-op | Optional client logger |
| `Verbose` | false | Enables client logs |
