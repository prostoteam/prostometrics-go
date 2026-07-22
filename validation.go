package prostometrics

import (
	"math"
	"sort"
	"strings"
)

func normalizeEventInput(typ metricType, metric string, value float64, labels []string) ([]string, bool) {
	metric = strings.TrimSpace(metric)
	if metric == "" || len(metric) > maxMetricBytes || strings.ContainsAny(metric, "|\r\n") {
		return nil, false
	}
	if typ != metricTypeUnique {
		if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
			return nil, false
		}
		limit := maxCounterValue
		if typ == metricTypeValue || typ == metricTypeValueSparse {
			limit = maxSampleValue
		}
		if value > limit {
			return nil, false
		}
	}
	if len(labels) > maxLabelsPerSeries {
		return nil, false
	}
	normalized := make([]string, 0, len(labels))
	seenNames := make(map[string]struct{}, len(labels))
	for _, raw := range labels {
		label := strings.TrimSpace(raw)
		if label == "" || len(label) > maxLabelBytes || strings.ContainsAny(label, "|\r\n") {
			return nil, false
		}
		name, _, ok := strings.Cut(label, "=")
		name = strings.TrimSpace(name)
		if !ok || name == "" || name == "workload" {
			return nil, false
		}
		if _, exists := seenNames[name]; exists {
			return nil, false
		}
		seenNames[name] = struct{}{}
		normalized = append(normalized, label)
	}
	sort.Strings(normalized)
	return normalized, true
}
