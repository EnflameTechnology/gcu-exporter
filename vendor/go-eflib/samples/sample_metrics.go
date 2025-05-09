/**
 * copyright (C) 2025 Enflame Technologies, Inc.
 */

package metrics

import (
	"strconv"

	"go-eflib"
	"go-eflib/efml"
)

type Metrics struct {
	Count   uint32
	vCount  uint32
	Devices []*Device
}

type Device struct {
	Host              string
	Index             string
	Minor             string
	Name              string
	UUID              string
	State             string
	Health            bool
	HealthMsg         string
	Heartbeat         bool
	PowerUsage        float64
	MemorySize        float64
	MemoryUsed        float64
	MemoryUsage       float64
	GcuUsage          float64
	ClusterUsage      []float64
	Temperature       float64
	GcuClock          float64
	EccStatus         *efml.DevEccStatus
	EslThroughput     []*efml.ThroughputInfo
	EslLink           []*efml.LinkInfo
	PcieLinkSpeed     float64
	PcieLinkWidth     float64
	PcieLink          *efml.LinkInfo
	CurrentResetCount float64
	TotalResetCount   float64
}

func collectMetrics() (*Metrics, error) {
	host := eflib.GetHostName()

	count, err := eflib.GetDeviceCount()
	if err != nil {
		return nil, err
	}

	vcount, err := eflib.GetVDeviceCount()
	if err != nil {
		return nil, err
	}

	metrics := &Metrics{
		Count:  count,
		vCount: vcount,
	}

	for index := uint32(0); index < uint32(count); index++ {
		h, err := eflib.GetDeviceHandle(index)
		if err != nil {
			continue
		}

		uuid, err := eflib.GetDeviceUUID(h)
		if err != nil {
			continue
		}

		name, err := eflib.GetDeviceName(h)
		if err != nil {
			continue
		}

		minor, err := eflib.GetDeviceMinor(h)
		if err != nil {
			continue
		}

		temperature, err := eflib.GetDeviceTemperature(h)
		if err != nil {
			continue
		}

		powerUsage, _, _, err := eflib.GetDevicePowerInfo(h)
		if err != nil {
			continue
		}

		memSize, err := eflib.GetDeviceMemoryTotalGBSize(h)
		if err != nil {
			continue
		}

		memUsed, err := eflib.GetDeviceMemoryUsedMBSize(h)
		if err != nil {
			continue
		}

		memUsage, err := eflib.GetDeviceMemoryUsage(h)
		if err != nil {
			continue
		}

		gcuUsage, err := eflib.GetDeviceGcuUsage(h)
		if err != nil {
			continue
		}

		clusterUsage, err := eflib.GetDeviceClusterUsage(h)
		if err != nil {
			continue
		}

		var eslThroughput []*efml.ThroughputInfo
		var eslLink []*efml.LinkInfo

		if hasEsl, _ := eflib.HasEslLink(h); hasEsl {
			eslThroughput, err = eflib.GetDeviceEslThroughput(h)
			if err != nil {
				continue
			}

			eslLink, err = eflib.GetDeviceEslLinkInfo(h)
			if err != nil {
				continue
			}
		}

		state, _, err := eflib.GetDeviceState(minor)
		if err != nil {
			continue
		}

		clock, err := eflib.GetDeviceClock(h)
		if err != nil {
			continue
		}

		heartbeat := eflib.GetDeviceSSmHeartbeatHealthy(minor)

		eccStatus, err := eflib.GetDeviceEccStatus(h)
		if err != nil {
			continue
		}

		pcieLinkSpeed, err := eflib.GetDevicePcieLinkSpeed(h)
		if err != nil {
			continue
		}

		pcieLinkWidth, err := eflib.GetDevicePcieLinkWidth(h)
		if err != nil {
			continue
		}

		pcieLink, err := eflib.GetDevicePcieLinkInfo(h)
		if err != nil {
			continue
		}

		currentResetCount, err := eflib.GetDeviceCurrentResetCount(minor)
		if err != nil {
			continue
		}

		totalResetCount, err := eflib.GetDeviceTotalResetCount(minor)
		if err != nil {
			continue
		}

		health, healthmsg := eflib.GetDeviceHealthState(minor, h)

		metrics.Devices = append(metrics.Devices,
			&Device{
				Host:              host,
				Index:             strconv.Itoa(int(index)),
				Minor:             strconv.Itoa(int(minor)),
				Name:              name,
				UUID:              uuid,
				State:             state,
				Health:            health,
				HealthMsg:         healthmsg,
				Heartbeat:         heartbeat,
				PowerUsage:        float64(powerUsage),
				MemorySize:        float64(memSize),
				MemoryUsed:        float64(memUsed),
				MemoryUsage:       float64(memUsage),
				GcuUsage:          float64(gcuUsage),
				ClusterUsage:      clusterUsage,
				Temperature:       float64(temperature),
				GcuClock:          float64(clock),
				EccStatus:         eccStatus,
				EslThroughput:     eslThroughput,
				EslLink:           eslLink,
				PcieLinkSpeed:     float64(pcieLinkSpeed),
				PcieLinkWidth:     float64(pcieLinkWidth),
				PcieLink:          pcieLink,
				CurrentResetCount: float64(currentResetCount),
				TotalResetCount:   float64(totalResetCount),
			})
	}

	return metrics, nil
}
