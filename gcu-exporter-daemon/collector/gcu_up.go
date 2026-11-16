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

type upCollector struct {
	metric []typedDesc
	logger log.Logger
}

func init() {
	registerCollector("gcu_up", defaultEnabled, NewUpCollector)
}

func NewUpCollector(logger log.Logger) (Collector, error) {
	return &upCollector{
		metric: []typedDesc{
			{prometheus.NewDesc(namespace+"_exporter_up", "Enflame GCU Metric Collection Services", []string{"host"}, nil), prometheus.GaugeValue},
		},
		logger: logger,
	}, nil
}

func (c *upCollector) Update(ch chan<- prometheus.Metric, metrics *Metrics) error {
	for _, device := range metrics.Devices {
		ch <- c.metric[0].mustNewConstMetric(float64(1), device.Host)
		break

	}
	return nil
}
