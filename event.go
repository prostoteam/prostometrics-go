package prostometrics

import (
	"sync"
	"time"
)

type metricType uint8

const (
	metricTypeCounter metricType = iota + 1
	metricTypeValue
	metricTypeValueSparse
	// metricTypeTotal is internal-only and converted to a counter delta before sending.
	metricTypeTotal
	metricTypeUnique
)

type event struct {
	typ      metricType
	name     string
	value    float64
	uniqueID string
	labels   []string
	ts       time.Time
}

var eventPool = sync.Pool{
	New: func() any {
		return &event{}
	},
}

func borrowEvent(typ metricType, name string, value float64, labels []string) *event {
	e := eventPool.Get().(*event)
	e.typ = typ
	e.name = name
	e.value = value
	e.labels = e.labels[:0]
	e.labels = append(e.labels, labels...)
	e.ts = time.Now()
	return e
}

func borrowUniqueEvent(name string, uniqueID string, labels []string) *event {
	e := borrowEvent(metricTypeUnique, name, 0, labels)
	e.uniqueID = uniqueID
	return e
}

func releaseEvent(e *event) {
	if e == nil {
		return
	}
	e.typ = 0
	e.name = ""
	e.value = 0
	e.uniqueID = ""
	if cap(e.labels) > maxLabelsPerSeries {
		e.labels = nil
	} else {
		e.labels = e.labels[:0]
	}
	eventPool.Put(e)
}
