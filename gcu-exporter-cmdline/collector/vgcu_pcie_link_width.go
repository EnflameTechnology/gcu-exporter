// Copyright 2015 The Prometheus Authors
// Copyright (c) 2022 Enflame. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type vpcieLinkWidthCollector struct {
	metric []typedDesc
	logger log.Logger
}

func init() {
	registerCollector("vgcu_pcie_link_width", defaultEnabled, NewVPcieLinkWidthCollector)
}

func NewVPcieLinkWidthCollector(logger log.Logger) (Collector, error) {
	return &vpcieLinkWidthCollector{
		metric: []typedDesc{
			{prometheus.NewDesc(vgcu_namespace+"_pcie_link_width", "Gcu pcie link width as reported by the device, -1 means not supported", []string{"vindex", "host", "minor_number", "uuid", "busid", "slot", "name", "pod_name", "pod_namespace","container_name", "metrics"}, nil), prometheus.GaugeValue},
		},
		logger: logger,
	}, nil
}

func (c *vpcieLinkWidthCollector) Update(ch chan<- prometheus.Metric, metrics *Metrics) error {
	if metrics.vCount == 0 {
		return nil
	}
	for _, device := range metrics.Devices {
		for _, vindex := range device.VIndexList {
			value, ok := metrics.DeviceToPod[strconv.Itoa(int(vindex))]
			if ok {
				ch <- c.metric[0].mustNewConstMetric(float64(device.PcieLink.Max_Link_Width), strconv.Itoa(int(vindex)), device.Host, device.Minor, device.Uuid, device.BusID, device.Slot, device.Name, value.name, value.namespace, value.container, "vgcu_pcie_link_width")
			} else {
				ch <- c.metric[0].mustNewConstMetric(float64(device.PcieLink.Max_Link_Width), strconv.Itoa(int(vindex)), device.Host, device.Minor, device.Uuid, device.BusID, device.Slot, device.Name, "", "", "", "vgcu_pcie_link_width")
			}
		}
	}
	return nil
}
