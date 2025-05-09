package collector

import (
	"testing"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewVPowerCapabilityCollector(t *testing.T) {
	logger := log.NewNopLogger()
	collector, err := NewVPowerCapabilityCollector(logger)
	if collector == nil || err != nil {
		t.Errorf("NewVPowerCapabilityCollector failed, got: %v, want: non-nil collector", collector)
	}
}

func TestVgcuPowerCapabilityUpdate(t *testing.T) {
	logger := log.NewNopLogger()
	collector, _ := NewVPowerCapabilityCollector(logger)
	ch := make(chan prometheus.Metric,1)
	metrics := &Metrics{
		vCount: 1,
		Devices: []*Device{
			{
				PowerCapability: 70,
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