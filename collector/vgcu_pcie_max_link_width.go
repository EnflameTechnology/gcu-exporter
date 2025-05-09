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

type vtemperatureCollector struct {
	metric []typedDesc
	logger log.Logger
}

func init() {
	registerCollector("vgcu_temperatures", defaultEnabled, NewVTemperatureCollector)
}

// NewLoadavgCollector returns a new Collector exposing load average stats.
func NewVTemperatureCollector(logger log.Logger) (Collector, error) {
	return &vtemperatureCollector{
		metric: []typedDesc{
			{prometheus.NewDesc(vgcu_namespace+"_temperatures", "Temperature as reported by the device, -1 means not supported", []string{"vindex", "host", "minor_number", "uuid", "busid", "slot", "name", "pod_name", "pod_namespace", "container_name"}, nil), prometheus.GaugeValue},
		},
		logger: logger,
	}, nil
}

func (c *vtemperatureCollector) Update(ch chan<- prometheus.Metric, metrics *Metrics) error {
	if metrics.vCount == 0 {
		return nil
	}
	for _, device := range metrics.Devices {
		for _, vindex := range device.VIndexList {
			value, ok := metrics.DeviceToPod[strconv.Itoa(int(vindex))]
			if ok {
				ch <- c.metric[0].mustNewConstMetric(device.Temperature, strconv.Itoa(int(vindex)), device.Host, device.Minor, device.Uuid, device.BusID, device.Slot, device.Name, value.name, value.namespace, value.container)
			} else {
				ch <- c.metric[0].mustNewConstMetric(device.Temperature, strconv.Itoa(int(vindex)), device.Host, device.Minor, device.Uuid, device.BusID, device.Slot, device.Name, "", "", "")
			}
		}
	}
	return nil
}
