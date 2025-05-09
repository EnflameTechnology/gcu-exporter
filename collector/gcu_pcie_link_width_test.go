package collector

import (
	"testing"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"go-eflib/efml"
)

func TestNewPcieLinkWidthCollector(t *testing.T) {
	logger := log.NewNopLogger()
	collector, err := NewPcieLinkWidthCollector(logger)
	if collector == nil || err != nil {
		t.Errorf("NewPcieLinkWidthCollector failed, got: %v, want: non-nil collector", collector)
	}
}

func TestGcuPcieLinkWidthUpdate(t *testing.T) {
	logger := log.NewNopLogger()
	collector, _ := NewPcieLinkWidthCollector(logger)
	ch := make(chan prometheus.Metric,1)
	metrics := &Metrics{
		vCount: 0,
		Devices: []*Device{
			{
				PcieLink: &efml.LinkInfo{
					Max_Link_Width: 16,
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