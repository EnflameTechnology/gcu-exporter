package collector

import (
	"testing"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewVMemoryUsedBytesCollector(t *testing.T) {
	logger := log.NewNopLogger()
	collector, err := NewVMemoryUsedBytesCollector(logger)
	if collector == nil || err != nil {
		t.Errorf("NewVMemoryUsedBytesCollector failed, got: %v, want: non-nil collector", collector)
	}
}

func TestVgcuMemoryUsedBytesUpdate(t *testing.T) {
	logger := log.NewNopLogger()
	collector, _ := NewVMemoryUsedBytesCollector(logger)
	ch := make(chan prometheus.Metric,1)
	metrics := &Metrics{
		vCount: 1,
		Devices: []*Device{
			{
				VMemoryUsed: []float64{512},
				Host: "localhost",
				Minor: "0",
				Uuid: "uuid",
				BusID: "busid",
				Slot: "slot",
				Name: "name",
				VIndexList: []uint{0},
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