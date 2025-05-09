package collector

import (
	"testing"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewUsageCollector(t *testing.T) {
	logger := log.NewNopLogger()
	collector, err := NewUsageCollector(logger)
	if collector == nil || err != nil {
		t.Errorf("NewUsageCollector failed, got: %v, want: non-nil collector", collector)
	}
}

func TestGcuUsageUpdate(t *testing.T) {
	logger := log.NewNopLogger()
	collector, _ := NewUsageCollector(logger)
	ch := make(chan prometheus.Metric,1)
	metrics := &Metrics{
		vCount: 0,
		Devices: []*Device{
			{
				GcuUsage: 500,
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