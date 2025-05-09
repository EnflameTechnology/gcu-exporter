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

type vhealthCollector struct {
	metric []typedDesc
	logger log.Logger
}

func init() {
	registerCollector("vgcu_health", defaultEnabled, NewVHealthCollector)
}

func NewVHealthCollector(logger log.Logger) (Collector, error) {
	return &vhealthCollector{
		metric: []typedDesc{
			{prometheus.NewDesc(vgcu_namespace+"_health", "Gcu health as reported by the device (2:healthy,1:unhealthy,0:unknown)", []string{"vindex", "host", "minor_number", "uuid", "busid", "slot", "name", "healthmsg"}, nil), prometheus.GaugeValue},
		},
		logger: logger,
	}, nil
}

func (c *vhealthCollector) Update(ch chan<- prometheus.Metric, metrics *Metrics) error {
	if metrics.vCount == 0 {
		return nil
	}
	for _, device := range metrics.Devices {
		for _, vindex := range device.VIndexList {
			ch <- c.metric[0].mustNewConstMetric(device.Health, strconv.Itoa(int(vindex)), device.Host, device.Minor, device.Uuid, device.BusID, device.Slot, device.Name, device.HealthMsg)
		}
	}
	return nil
}
