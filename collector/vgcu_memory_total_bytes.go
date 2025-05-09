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

type vmemoryTotalBytesCollector struct {
	metric []typedDesc
	logger log.Logger
}

func init() {
	registerCollector("vgcu_memory_total_bytes", defaultEnabled, NewVMemoryTotalBytesCollector)
}

func NewVMemoryTotalBytesCollector(logger log.Logger) (Collector, error) {
	return &vmemoryTotalBytesCollector{
		metric: []typedDesc{
			{prometheus.NewDesc(vgcu_namespace+"_memory_total_bytes", "Total memory size as reorted by the vdevices", []string{"vindex", "host", "minor_number", "uuid", "busid", "slot", "name", "pod_name", "pod_namespace", "container_name"}, nil), prometheus.GaugeValue},
		},
		logger: logger,
	}, nil
}

func (c *vmemoryTotalBytesCollector) Update(ch chan<- prometheus.Metric, metrics *Metrics) error {
	if metrics.vCount == 0 {
		return nil
	}
	for _, device := range metrics.Devices {
		for i, vindex := range device.VIndexList {
			vmemSize := device.VMemorySize[i]
			value, ok := metrics.DeviceToPod[strconv.Itoa(int(vindex))]
			if ok {
				ch <- c.metric[0].mustNewConstMetric(vmemSize, strconv.Itoa(int(vindex)), device.Host, device.Minor, device.Uuid, device.BusID, device.Slot, device.Name, value.name, value.namespace, value.container)
			} else {
				ch <- c.metric[0].mustNewConstMetric(vmemSize, strconv.Itoa(int(vindex)), device.Host, device.Minor, device.Uuid, device.BusID, device.Slot, device.Name, "", "", "")

			}
		}
	}
	return nil
}
