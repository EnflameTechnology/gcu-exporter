package collector

import (
	"testing"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewMemoryUsageCollector(t *testing.T) {
	logger := log.NewNopLogger()
	collector, err := NewMemoryUsageCollector(logger)
	if collector == nil || err != nil {
		t.Errorf("NewMemoryUsageCollector failed, got: %v, want: non-nil collector", collector)
	}
}

func TestGcuMemoryUsageUpdate(t *testing.T) {
	logger := log.NewNopLogger()
	collector, _ := NewMemoryUsageCollector(logger)
	ch := make(chan prometheus.Metric,1)
	metrics := &Metrics{
		vCount: 0,
		Devices: []*Device{
			{
				MemoryUsage: 500,
				Host: "localhost",
				Minor: "0",
				Uuid: "uuid",
				BusID: "busid",
				Slot: "slot",
				Name: "name",
			},
		},
		DeviceToPod: map[string]devicePodInfo{
			"enflame": {
				name: "pod_name",
				namespace: "pod_namespace",
				container: "container_name",
			},
		},
	}
	err := collector.Update(ch, metrics)
	if err != nil {
		t.Errorf("Update failed, got: %v, want: nil", err)
	}
}