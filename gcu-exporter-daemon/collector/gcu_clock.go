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
)

type clockCollector struct {
	metric []typedDesc
	logger log.Logger
}

func init() {
	registerCollector("gcu_clock", defaultEnabled, NewClockCollector)
}

func NewClockCollector(logger log.Logger) (Collector, error) {
	return &clockCollector{
		metric: []typedDesc{
			{prometheus.NewDesc(namespace+"_clock", "Gcu clock as reported by the device, -1 means not supported", []string{"host", "minor_number", "uuid", "busid", "slot", "name", "power_mode", "pod_name", "pod_namespace", "container_name"}, nil), prometheus.GaugeValue},
		},
		logger: logger,
	}, nil
}

func (c *clockCollector) Update(ch chan<- prometheus.Metric, metrics *Metrics) error {
	if metrics.vCount > 0 {
		return nil
	}

	for _, device := range metrics.Devices {
		if device.ClockVisible == true {
		   ch <- c.metric[0].mustNewConstMetric(device.GcuClock, device.Host, device.Minor, device.Uuid, device.BusID, device.Slot, device.Name, device.PowerMode, metrics.DeviceToPod[device.Minor].name, metrics.DeviceToPod[device.Minor].namespace, metrics.DeviceToPod[device.Minor].container)
	        }
	}
	return nil
}
