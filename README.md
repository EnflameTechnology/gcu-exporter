# Enflame GCU Exporter (gcu-exporter)

> **Language**: English | [中文](README_zh.md)

## Table of Contents

- [Overview](#overview)
- [Project Structure](#project-structure)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Build from Source](#build-from-source)
  - [Build the Container Image](#build-the-container-image)
- [Quick Start](#quick-start)
  - [Run in Docker](#run-in-docker)
  - [Deploy on Kubernetes](#deploy-on-kubernetes)
  - [Scrape Metrics](#scrape-metrics)
- [Key Features](#key-features)
- [Exported Metrics](#exported-metrics)
- [Version Compatibility](#version-compatibility)
- [Hardware Support](#hardware-support)
- [User Manual](#user-manual)
- [Trouble Shooting](#trouble-shooting)
- [License](#license)

## Overview

`gcu-exporter` is a [Prometheus](https://prometheus.io/) exporter that exposes
operational metrics for Enflame GCU (General Compute Unit) devices. It collects
per-device health, utilization, memory, temperature and power information from
the host and publishes it on an HTTP `/metrics` endpoint so that it can be
scraped by Prometheus and visualized with Grafana.

The exporter delivers:

- **Rich GCU metrics** — utilization, memory usage, temperature, power, clock and health status for every GCU on the node.
- **Multiple collection backends** — reads device data through `efsmi` and the EFML management library ([go-eflib](../go-eflib)), with a graceful fallback path that keeps working on older drivers.
- **Kubernetes-native** — ships DaemonSet manifests (privileged / non-privileged, x86 / arm64, optional Basic Auth) plus Prometheus and Grafana examples.
- **Container & bare-metal** — can run directly on a host, in Docker, or as a Kubernetes DaemonSet.

## Project Structure

```
gcu-exporter/
├── gcu-exporter-daemon/     # Long-running daemon collector (efsmi-first, EFML fallback)
├── gcu-exporter-cmdline/    # Command-line collector variant
├── gcu-exporter-server/     # Lightweight efsmi-based HTTP metrics server
├── deployments/             # Docker / Kubernetes manifests and image build files
│   ├── yaml/                # DaemonSet variants (privileged, non-privileged, arm, basic-auth)
│   ├── examples/            # docker-compose + k8s + Prometheus/Grafana examples
│   └── dockerfiles/         # Dockerfile.ubuntu / .tlinux / .openeuler
├── Makefile                 # Builds all binaries
└── build.sh                 # Build + package entry point
```

## Installation

### Prerequisites

- An Enflame GCU device with the Enflame GCU driver installed.
- `efsmi` and the EFML runtime library (`libefml.so`) available on the host.
- [Go](https://go.dev/) 1.22+ and a C toolchain (the collectors use cgo) for building from source.
- Docker and/or a Kubernetes cluster for containerized deployment.

### Build from Source

```bash
# Build all collector binaries and produce a release package under dist/
./build.sh all

# Clean build artifacts
./build.sh clean
```

The binaries (`gcu-exporter-daemon`, `gcu-exporter-cmdline`, `gcu-exporter-server`)
are produced under `bin/` and packaged into `dist/gcu-exporter_<version>/`.

### Build the Container Image

```bash
cd deployments
# Configure the target base image / registry in build-image.conf, then:
./build-image.sh
```

Dockerfiles for Ubuntu, TencentOS (tlinux) and openEuler are provided under
`deployments/dockerfiles/`.

## Quick Start

### Run in Docker

```bash
cd deployments/examples/docker
./docker-compose.sh          # brings up gcu-exporter + Prometheus
```

### Deploy on Kubernetes

```bash
# Pick the DaemonSet manifest that matches your environment
kubectl apply -f deployments/yaml/gcu-exporter.yaml

# Non-privileged / arm64 / Basic-Auth variants are available in the same directory, e.g.
# deployments/yaml/gcu-exporter-non-privileged.yaml
# deployments/yaml/gcu-exporter-for-arm.yaml
# deployments/yaml/gcu-exporter-with-basic-auth.yaml
```

### Scrape Metrics

The exporter listens on port **9400** by default:

```bash
curl http://<node-ip>:9400/metrics
```

Point Prometheus at the endpoint (the DaemonSet Service is annotated with
`prometheus.io/port: "9400"`), and import the Grafana dashboard from
`deployments/examples/k8s/yaml/grafana.yaml`.

## Key Features

| Feature | Description |
| --- | --- |
| Prometheus format | Exposes GCU metrics on `/metrics` in the standard Prometheus text format |
| Health monitoring | Reports per-device health status for alerting |
| Utilization & memory | GCU compute utilization and memory usage per device |
| Thermal & power | Temperature and power draw per device |
| EFML + efsmi backends | Collects via EFML/`efsmi` with a fallback path for older drivers |
| Flexible deployment | Docker, bare-metal or Kubernetes DaemonSet (privileged / non-privileged) |
| Multi-arch | x86_64 and arm64 manifests and images |

## Exported Metrics

Metrics cover, per GCU device, dimensions such as device health, compute
utilization, memory total/used, temperature, power usage and clock. Each series
is labelled with the device index and node so that dashboards and alerts can be
built per card and per node. Run `curl http://<node-ip>:9400/metrics` to see the
full, up-to-date set for your driver version.

## Version Compatibility

`gcu-exporter` reads device data through `efsmi` / EFML, so its capabilities
track the installed Enflame GCU driver. Newer drivers expose more metrics; the
exporter degrades gracefully on older drivers via its fallback collection path.
Always match the exporter release to the GCU driver deployed on your nodes.

## Hardware Support

Supports Enflame GCU devices that provide the Enflame GCU driver together with
`efsmi` and the EFML runtime library. The same binary works across GCU product
generations; the exact metric set available depends on the driver.

## User Manual

| Document | Description |
| --- | --- |
| [deployments/README.md](deployments/README.md) | Deployment and image build guide |
| [deployments/examples/README.md](deployments/examples/README.md) | Docker / Kubernetes / Prometheus / Grafana examples |

## Trouble Shooting

- **`/metrics` returns no GCU series** — verify the driver is loaded and `efsmi`
  runs on the host; check that the container is privileged (or that the required
  devices `/dev/gcu*` and `libefml.so` are mounted for the non-privileged variant).
- **`libefml.so` not found** — ensure the Enflame GCU driver / management library
  is installed and on the library search path.
- **Prometheus does not scrape the target** — confirm the Service annotations and
  that port `9400` is reachable from the Prometheus pod/network.

## License

`gcu-exporter` is licensed under the **Apache License 2.0**. Copyright (c)
Enflame Technology. All Rights Reserved. See [LICENSE](LICENSE).
