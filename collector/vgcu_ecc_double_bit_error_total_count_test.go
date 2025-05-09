package collector

import (
	"testing"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"go-eflib/efml"
)

func TestNewVEccDoubleBitCountCollector(t *testing.T) {
	logger := log.NewNopLogger()
	collector, err := NewVEccDoubleBitCountCollector(logger)
	if collector == nil || err != nil {
		t.Errorf("NewVEccDoubleBitCountCollector failed, got: %v, want: non-nil collector", collector)
	}
}

func TestVgcuEccDoubleBitErrorTotalCountUpdate(t *testing.T) {
	logger := log.NewNopLogger()
	collector, _ := NewVEccDoubleBitCountCollector(logger)
	ch := make(chan prometheus.Metric,1)
	metrics := &Metrics{
		vCount: 1,
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