package prostometrics

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNormalizeEventInputCanonicalizesLabels(t *testing.T) {
	labels, ok := normalizeEventInput(metricTypeValue, " latency_ms ", 12.5, []string{"region=eu", "host=a"})
	if !ok {
		t.Fatal("valid input rejected")
	}
	if got := labels[0] + "," + labels[1]; got != "host=a,region=eu" {
		t.Fatalf("labels = %q", got)
	}
	if _, ok := normalizeEventInput(metricTypeValue, "latency_ms", 1, []string{"host=a", "host=b"}); ok {
		t.Fatal("duplicate label names accepted")
	}
	if _, ok := normalizeEventInput(metricTypeValue, "bad|metric", 1, nil); ok {
		t.Fatal("protocol delimiter accepted in metric")
	}
}

type blockingTransport struct {
	started chan struct{}
	release chan struct{}
	once    sync.Once
}

func (t *blockingTransport) Send(context.Context, *Payload) error {
	t.once.Do(func() { close(t.started) })
	<-t.release
	return nil
}

func TestCloseCanWaitAgainAfterTimeout(t *testing.T) {
	transport := &blockingTransport{started: make(chan struct{}), release: make(chan struct{})}
	client, err := NewClient("api", Config{Transport: transport})
	if err != nil {
		t.Fatal(err)
	}
	client.Count("requests", 1)
	select {
	case <-transport.started:
	case <-time.After(time.Second):
		t.Fatal("transport did not start")
	}
	first, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if err := client.Close(first); err == nil {
		t.Fatal("first Close succeeded while transport was blocked")
	}
	close(transport.release)
	second, secondCancel := context.WithTimeout(context.Background(), time.Second)
	defer secondCancel()
	if err := client.Close(second); err != nil {
		t.Fatalf("second Close() = %v", err)
	}
}

func TestSparseValueUsesDedicatedWireType(t *testing.T) {
	payload := &Payload{Values: []ValueEvent{{Metric: "capacity", Value: 42, Sparse: true, Timestamp: 1}}}
	body := encodeLinePayloadV5(payload, newDictionaryState(), false)
	if !containsLinePrefix(body, "s|") {
		t.Fatalf("sparse payload = %q", body)
	}
}

func containsLinePrefix(body []byte, prefix string) bool {
	for start := 0; start < len(body); {
		end := start
		for end < len(body) && body[end] != '\n' {
			end++
		}
		if end-start >= len(prefix) && string(body[start:start+len(prefix)]) == prefix {
			return true
		}
		start = end + 1
	}
	return false
}
