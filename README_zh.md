# Enflame GCU Exporter (gcu-exporter)

> **语言**: [English](README.md) | 中文

## 目录

- [概述](#概述)
- [项目结构](#项目结构)
- [安装](#安装)
  - [前置条件](#前置条件)
  - [从源码构建](#从源码构建)
  - [构建容器镜像](#构建容器镜像)
- [快速开始](#快速开始)
  - [在 Docker 中运行](#在-docker-中运行)
  - [在 Kubernetes 上部署](#在-kubernetes-上部署)
  - [抓取指标](#抓取指标)
- [核心特性](#核心特性)
- [导出的指标](#导出的指标)
- [版本兼容性](#版本兼容性)
- [硬件支持](#硬件支持)
- [用户手册](#用户手册)
- [故障排查](#故障排查)
- [许可证](#许可证)

## 概述

`gcu-exporter` 是一个 [Prometheus](https://prometheus.io/) exporter，用于暴露
Enflame GCU（General Compute Unit，通用计算单元）设备的运行指标。它从主机上采集
每个设备的健康、利用率、显存、温度和功耗信息，并通过 HTTP `/metrics` 端点发布，
以便被 Prometheus 抓取并用 Grafana 可视化。

该 exporter 提供：

- **丰富的 GCU 指标** — 涵盖节点上每张 GCU 的利用率、显存使用、温度、功耗、时钟频率和健康状态。
- **多种采集后端** — 通过 `efsmi` 和 EFML 管理库（[go-eflib](../go-eflib)）读取设备数据，并具备优雅的兜底路径，在较旧驱动上仍可正常工作。
- **Kubernetes 原生** — 提供 DaemonSet 清单（特权 / 非特权、x86 / arm64、可选 Basic Auth），以及 Prometheus 和 Grafana 示例。
- **容器与裸金属** — 可直接运行在主机上、Docker 中，或作为 Kubernetes DaemonSet 运行。

## 项目结构

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

## 安装

### 前置条件

- 一台已安装 Enflame GCU 驱动的 Enflame GCU 设备。
- 主机上可用的 `efsmi` 和 EFML 运行时库（`libefml.so`）。
- 用于从源码构建的 [Go](https://go.dev/) 1.22+ 和 C 工具链（采集器使用 cgo）。
- 用于容器化部署的 Docker 和/或 Kubernetes 集群。

### 从源码构建

```bash
# Build all collector binaries and produce a release package under dist/
./build.sh all

# Clean build artifacts
./build.sh clean
```

各二进制文件（`gcu-exporter-daemon`、`gcu-exporter-cmdline`、`gcu-exporter-server`）
生成在 `bin/` 目录下，并打包到 `dist/gcu-exporter_<version>/`。

### 构建容器镜像

```bash
cd deployments
# Configure the target base image / registry in build-image.conf, then:
./build-image.sh
```

`deployments/dockerfiles/` 下提供了适用于 Ubuntu、TencentOS（tlinux）和 openEuler 的
Dockerfile。

## 快速开始

### 在 Docker 中运行

```bash
cd deployments/examples/docker
./docker-compose.sh          # brings up gcu-exporter + Prometheus
```

### 在 Kubernetes 上部署

```bash
# Pick the DaemonSet manifest that matches your environment
kubectl apply -f deployments/yaml/gcu-exporter.yaml

# Non-privileged / arm64 / Basic-Auth variants are available in the same directory, e.g.
# deployments/yaml/gcu-exporter-non-privileged.yaml
# deployments/yaml/gcu-exporter-for-arm.yaml
# deployments/yaml/gcu-exporter-with-basic-auth.yaml
```

### 抓取指标

该 exporter 默认监听 **9400** 端口：

```bash
curl http://<node-ip>:9400/metrics
```

将 Prometheus 指向该端点（DaemonSet 的 Service 已通过
`prometheus.io/port: "9400"` 注解标注），并从
`deployments/examples/k8s/yaml/grafana.yaml` 导入 Grafana 仪表盘。

## 核心特性

| 特性 | 描述 |
| --- | --- |
| Prometheus 格式 | 以标准 Prometheus 文本格式在 `/metrics` 上暴露 GCU 指标 |
| 健康监控 | 上报每个设备的健康状态以供告警 |
| 利用率与显存 | 每个设备的 GCU 计算利用率和显存使用 |
| 温度与功耗 | 每个设备的温度和功耗 |
| EFML + efsmi 后端 | 通过 EFML/`efsmi` 采集，并为较旧驱动提供兜底路径 |
| 灵活部署 | Docker、裸金属或 Kubernetes DaemonSet（特权 / 非特权） |
| 多架构 | x86_64 和 arm64 清单与镜像 |

## 导出的指标

指标按 GCU 设备覆盖多个维度，例如设备健康、计算利用率、显存总量/已用、温度、
功耗和时钟频率。每条时间序列都带有设备索引和节点标签，从而可以按卡、按节点构建
仪表盘和告警。运行 `curl http://<node-ip>:9400/metrics` 可查看你的驱动版本下
完整、最新的指标集合。

## 版本兼容性

`gcu-exporter` 通过 `efsmi` / EFML 读取设备数据，因此其能力随所安装的 Enflame GCU
驱动而变化。较新的驱动会暴露更多指标；在较旧驱动上，exporter 会通过其兜底采集路径
优雅降级。请始终让 exporter 版本与节点上部署的 GCU 驱动相匹配。

## 硬件支持

支持提供了 Enflame GCU 驱动以及 `efsmi` 和 EFML 运行时库的 Enflame GCU 设备。同一份
二进制文件可跨 GCU 产品代际工作；具体可用的指标集取决于驱动。

## 用户手册

| 文档 | 描述 |
| --- | --- |
| [deployments/README.md](deployments/README.md) | 部署与镜像构建指南 |
| [deployments/examples/README.md](deployments/examples/README.md) | Docker / Kubernetes / Prometheus / Grafana 示例 |

## 故障排查

- **`/metrics` 未返回任何 GCU 序列** — 确认驱动已加载且 `efsmi` 能在主机上运行；
  检查容器是否为特权模式（或对于非特权变体，是否已挂载所需的设备 `/dev/gcu*` 和
  `libefml.so`）。
- **找不到 `libefml.so`** — 确保 Enflame GCU 驱动 / 管理库已安装并位于库搜索路径中。
- **Prometheus 未抓取到目标** — 确认 Service 注解正确，且 Prometheus pod/网络能够
  访问 `9400` 端口。

## 许可证

`gcu-exporter` 采用 **Apache License 2.0** 许可证。版权所有 (c) Enflame Technology，
保留所有权利。详见 [LICENSE](LICENSE)。
