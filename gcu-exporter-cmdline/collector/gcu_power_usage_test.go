package collector

import (
	"testing"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewPowerUsageCollector(t *testing.T) {
	logger := log.NewNopLogger()
	collector, err := NewPowerUsageCollector(logger)
	if collector == nil || err != nil {
		t.Errorf("NewPowerUsageCollector failed, got: %v, want: non-nil collector", collector)
	}
}

func TestGcuPowerUsageUpdate(t *testing.T) {
	logger := log.NewNopLogger()
	collector, _ := NewPowerUsageCollector(logger)
	ch := make(chan prometheus.Metric,1)
	metrics := &Metrics{
		vCount: 0,
		Devices: []*Device{
			{
				PowerUsage: 500,
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