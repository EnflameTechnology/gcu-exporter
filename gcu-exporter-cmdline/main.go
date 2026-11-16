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

package main

import (
	"encoding/json"
	"fmt"
	"gcu-exporter/internal/eflib"
	"os"
	"os/user"
	"runtime"
	"sort"
	"time"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"gcu-exporter/collector"
)

var Version string

// MetricJSON represents a single metric in JSON format
type MetricJSON struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Help        string            `json:"help,omitempty"`
	Unit        string            `json:"unit,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Value       float64           `json:"value,omitempty"`
	StringValue string            `json:"string_value,omitempty"`
	Timestamp   int64             `json:"timestamp,omitempty"`
}

// MetricsJSON represents the complete JSON output structure
type MetricsJSON struct {
	Timestamp int64        `json:"timestamp"`
	Metrics   []MetricJSON `json:"metrics"`
}

// convertToJSON converts Prometheus MetricFamily to JSON format
func convertToJSON(metricFamilies []*dto.MetricFamily) (*MetricsJSON, error) {
	result := &MetricsJSON{
		Timestamp: time.Now().Unix(),
		Metrics:   []MetricJSON{},
	}

	for _, mf := range metricFamilies {
		metricType := mf.GetType().String()
		help := mf.GetHelp()

		for _, metric := range mf.GetMetric() {
			labels := make(map[string]string)
			for _, labelPair := range metric.GetLabel() {
				labels[labelPair.GetName()] = labelPair.GetValue()
			}

			metricJSON := MetricJSON{
				Name:   mf.GetName(),
				Type:   metricType,
				Help:   help,
				Labels: labels,
			}

			// Extract value based on metric type
			if metric.Gauge != nil {
				metricJSON.Value = metric.GetGauge().GetValue()
			} else if metric.Counter != nil {
				metricJSON.Value = metric.GetCounter().GetValue()
			} else if metric.Histogram != nil {
				// For histogram, we'll use the sample count as the value
				metricJSON.Value = float64(metric.GetHistogram().GetSampleCount())
			} else if metric.Summary != nil {
				// For summary, we'll use the sample count as the value
				metricJSON.Value = float64(metric.GetSummary().GetSampleCount())
			} else if metric.Untyped != nil {
				metricJSON.Value = metric.GetUntyped().GetValue()
			}

			if metric.TimestampMs != nil {
				metricJSON.Timestamp = metric.GetTimestampMs()
			}

			result.Metrics = append(result.Metrics, metricJSON)
		}
	}

	return result, nil
}

func main() {
	var (
		disableExporterMetrics = kingpin.Flag(
			"web.disable-exporter-metrics",
			"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
		).Default("false").Bool()
		disableDefaultCollectors = kingpin.Flag(
			"collector.disable-defaults",
			"Set all collectors to disabled by default.",
		).Default("false").Bool()
		outputFile = kingpin.Flag(
			"output.file",
			"Path to output file. If not specified, metrics will be written to stdout.",
		).Default("").String()
		outputFormat = kingpin.Flag(
			"output.format",
			"Output format: 'json' or 'prometheus' (default: 'prometheus').",
		).Default("prometheus").Enum("json", "prometheus")
		prettyPrint = kingpin.Flag(
			"output.pretty",
			"Pretty print JSON output with indentation (only applies to JSON format).",
		).Default("false").Bool()
		warmup = kingpin.Flag(
			"warmup",
			"Do a warm-up scrape before final output. Helps avoid zero usage on short-lived one-shot runs.",
		).Default("true").Bool()
		warmupIntervalMs = kingpin.Flag(
			"warmup.interval-ms",
			"Sleep interval in milliseconds between warm-up scrape and final scrape. 0 disables waiting.",
		).Default("500").Int()
	)

	if err := eflib.Init(true); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't initialize efml: %v. Make sure EFML is in the shared library search path, please mount EFML Lib from host, for details, read user guide\n", err)
		os.Exit(1)
	}

	defer eflib.Shutdown()

	runtime.GOMAXPROCS(1)

	collector.IsInK8sCloud = false
	collector.NeedMonitPodInfo = false

	kubernetes_service := os.Getenv("KUBERNETES_SERVICE_PORT")
	if kubernetes_service != "" {
		collector.IsInK8sCloud = true
	}

	monitPodInfo := os.Getenv("MONIT_PODS")
	if monitPodInfo == "true" {
		collector.NeedMonitPodInfo = true
	}

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(Version)
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	if *warmupIntervalMs < 0 {
		level.Warn(logger).Log("msg", "Invalid warmup interval, falling back to 1000ms", "warmup.interval-ms", *warmupIntervalMs)
		*warmupIntervalMs = 500
	}

	if *disableDefaultCollectors {
		collector.DisableDefaultCollectors()
	}
	level.Info(logger).Log("msg", "Starting gcu-exporter-cmdline", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())
	if user, err := user.Current(); err == nil && user.Uid == "0" {
		level.Warn(logger).Log("msg", "Gcu Exporter is running as root user. This exporter is designed to run as unprivileged user, root is not required.")
	}

	// Create collector
	nc, err := collector.NewGcuCollector(logger)
	if err != nil {
		level.Error(logger).Log("msg", "Couldn't create collector", "err", err)
		os.Exit(1)
	}

	// Log enabled collectors
	level.Info(logger).Log("msg", "Enabled collectors")
	collectors := []string{}
	for n := range nc.Collectors {
		collectors = append(collectors, n)
	}
	sort.Strings(collectors)
	for _, c := range collectors {
		level.Info(logger).Log("collector", c)
	}

	// Create registry and register collectors
	registry := prometheus.NewRegistry()
	registry.MustRegister(version.NewCollector("gcu_exporter"))
	if err := registry.Register(nc); err != nil {
		level.Error(logger).Log("msg", "Couldn't register collector", "err", err)
		os.Exit(1)
	}

	// Optionally register exporter metrics
	if !*disableExporterMetrics {
		registry.MustRegister(
			promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
			promcollectors.NewGoCollector(),
		)
	}

	// Gather metrics. In one-shot mode, usage values can stay at 0 if sampled
	// immediately after process startup. Optionally do a warm-up scrape first.
	if *warmup {
		level.Info(logger).Log("msg", "Running warm-up scrape before final collection")
		if _, err := registry.Gather(); err != nil {
			level.Warn(logger).Log("msg", "Warm-up scrape failed, continuing with final collection", "err", err)
		}
		if *warmupIntervalMs > 0 {
			time.Sleep(time.Duration(*warmupIntervalMs) * time.Millisecond)
		}
	}

	// Gather metrics
	level.Info(logger).Log("msg", "Collecting metrics...")
	metricFamilies, err := registry.Gather()
	if err != nil {
		level.Error(logger).Log("msg", "Error gathering metrics", "err", err)
		os.Exit(1)
	}

	// Determine output destination
	var output *os.File
	if *outputFile != "" {
		output, err = os.Create(*outputFile)
		if err != nil {
			level.Error(logger).Log("msg", "Error creating output file", "file", *outputFile, "err", err)
			os.Exit(1)
		}
		defer output.Close()
		level.Info(logger).Log("msg", "Writing metrics to file", "file", *outputFile, "format", *outputFormat)
	} else {
		output = os.Stdout
		level.Info(logger).Log("msg", "Writing metrics to stdout", "format", *outputFormat)
	}

	// Output in the selected format
	switch *outputFormat {
	case "json":
		// Convert to JSON
		metricsJSON, err := convertToJSON(metricFamilies)
		if err != nil {
			level.Error(logger).Log("msg", "Error converting metrics to JSON", "err", err)
			os.Exit(1)
		}

		// Encode and write JSON
		encoder := json.NewEncoder(output)
		if *prettyPrint {
			encoder.SetIndent("", "  ")
		}
		if err := encoder.Encode(metricsJSON); err != nil {
			level.Error(logger).Log("msg", "Error encoding JSON", "err", err)
			os.Exit(1)
		}
		level.Info(logger).Log("msg", "Metrics collection completed successfully", "format", "json", "count", len(metricsJSON.Metrics))

	case "prometheus":
		// Write Prometheus text format
		encoder := expfmt.NewEncoder(output, expfmt.FmtText)
		for _, mf := range metricFamilies {
			if err := encoder.Encode(mf); err != nil {
				level.Error(logger).Log("msg", "Error encoding metric family", "err", err)
				os.Exit(1)
			}
		}
		// Close encoder if it implements Closer interface
		if closer, ok := encoder.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				level.Error(logger).Log("msg", "Error closing encoder", "err", err)
				os.Exit(1)
			}
		}
		level.Info(logger).Log("msg", "Metrics collection completed successfully", "format", "prometheus", "count", len(metricFamilies))

	default:
		level.Error(logger).Log("msg", "Unknown output format", "format", *outputFormat)
		os.Exit(1)
	}
}
