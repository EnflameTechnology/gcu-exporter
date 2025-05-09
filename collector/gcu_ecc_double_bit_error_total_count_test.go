package collector

import (
	"testing"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"go-eflib/efml"
)

func TestNewEccDoubleBitCountCollector(t *testing.T) {
	logger := log.NewNopLogger()
	collector, err := NewEccDoubleBitCountCollector(logger)
	if collector == nil || err != nil {
		t.Errorf("NewEccDoubleBitCountCollector failed, got: %v, want: non-nil collector", collector)
	}
}

func TestGcuEccDoubleBitErrorTotalCountUpdate(t *testing.T) {
	logger := log.NewNopLogger()
	collector, _ := NewEccDoubleBitCountCollector(logger)
	ch := make(chan prometheus.Metric,1)
	metrics := &Metrics{
		vCount: 0,
		Devices: []*Device{
			{
				EccStatus: &efml.DevEccStatus{
					Ecnt_db: 1,
				},
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