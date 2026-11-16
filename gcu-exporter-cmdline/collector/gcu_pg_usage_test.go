package collector

import (
	"testing"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewPGUsageCollector(t *testing.T) {
	logger := log.NewNopLogger()
	collector, err := NewPGUsageCollector(logger)
	if collector == nil || err != nil {
		t.Errorf("NewPGUsageCollector failed, got: %v, want: non-nil collector", collector)
	}
}

func TestGcuPgUsageUpdate(t *testing.T) {
	logger := log.NewNopLogger()
	collector, _ := NewPGUsageCollector(logger)
	ch := make(chan prometheus.Metric,10)
	metrics := &Metrics{
		vCount: 0,
		Devices: []*Device{
			{
				PGUsage: []float64{0.5, 0.6},
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
