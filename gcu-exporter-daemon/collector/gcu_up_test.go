package collector

import (
	"testing"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewUpCollector(t *testing.T) {
	logger := log.NewNopLogger()
	collector, err := NewUpCollector(logger)
	if collector == nil || err != nil {
		t.Errorf("NewUpCollector failed, got: %v, want: non-nil collector", collector)
	}
}

func TestGcuUpUpdate(t *testing.T) {
	logger := log.NewNopLogger()
	collector, _ := NewUpCollector(logger)
	ch := make(chan prometheus.Metric,1)
	metrics := &Metrics{
		vCount: 0,
		Devices: []*Device{
			{
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