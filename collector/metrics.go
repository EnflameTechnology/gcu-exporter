// Copyright (c) 2022 Enflame. All Rights Reserved.

package collector

import (
	"fmt"
	"go-eflib"
	"go-eflib/efml"
	"strconv"
)

const (
	socketPath = "/var/lib/kubelet/pod-resources/kubelet.sock"
)

type Metrics struct {
	Count       uint32
	vCount      uint32
	Devices     []*Device
	DeviceToPod map[string]devicePodInfo
	//ClusterToPod  map[string]devicePodInfo
}

var IsInK8sCloud bool
var NeedMonitPodInfo bool

type Device struct {
	Host          string
	Index         string
	Minor         string
	Name          string
	Uuid          string
	BusID         string
	Slot          string
	FwVersion     string
	DevSn         string
	Health        float64
	HealthMsg     string
	PowerUsage    float64
	PowerConsumption    float64
	PowerCapability     float64
	MemorySize    float64
	MemoryUsed    float64
	MemoryUsage   float64
	GcuUsage      float64
	SipUsage      float64
	ClusterUsage  []float64
	PGUsage       []float64
	VGcuUsage     []float64
	VMemorySize   []float64
	VMemoryUsed   []float64
	VMemoryUsage  []float64
	VIndexList    []uint
	Temperature   float64
	GcuClock      float64
	ClockVisible  bool
	PowerMode     string
	EccStatus     *efml.DevEccStatus
	EslThroughput []*efml.ThroughputInfo
	EslLink       []*efml.LinkInfo
	PcieLinkSpeed float64
	PcieLinkWidth float64
	PcieLink      *efml.LinkInfo
	VirtMode      float64
}

func collectGCUMetrics() (*Metrics, error) {
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

	if IsInK8sCloud || NeedMonitPodInfo {
		clusterMap := false

		if vcount != 0 {
			clusterMap = true
		}

		DeviceToPod, err := getDevicePodInfo(socketPath, clusterMap)
		/*for key, item := range DeviceToPod {
			fmt.Println("Here get DeviceToPod", key, item.name, item.namespace)
		}*/

		if err != nil {
			fmt.Println("getDevicePodInfo errored as ", err)
		} else {
			//if clusterMap {
			//		metrics.ClusterToPod = DeviceToPod
			//	} else {
			metrics.DeviceToPod = DeviceToPod
			//	}
		}
	}

	if vcount == 0 {

		for index := uint32(0); index < uint32(count); index++ {
			h, err := eflib.GetDeviceHandle(index)
			if err != nil {
				continue
			}

			uuid, err := eflib.GetDeviceUUID(h)
			if err != nil {
				continue
			}

			//name, err := eflib.GetDeviceName(h)
			name, err := eflib.GetDeviceType(h)
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

			powerUsage, powerConsumption, powerCapability, err := eflib.GetDevicePowerInfo(h)
			if err != nil {
				continue
			}

			memUsage, memSize, memUsed, err := eflib.GetDeviceMemoryInfo(h)
			if err != nil {
				continue
			}

			gcuUsage, err := eflib.GetDeviceGcuUsage(h)
			if err != nil {
				continue
			}

			sipUsage, err := eflib.GetDeviceSipUsage(h)
			if err != nil {
				continue
			}

			clusterUsage, err := eflib.GetDeviceClusterUsage(h)
			if err != nil {
				continue
			}

			PGUsage, err := eflib.GetDevicePGUsage(h)
			if err != nil {
				continue
			}

			var eslThroughput []*efml.ThroughputInfo
			var eslLink []*efml.LinkInfo

			if hasEsl, _ := eflib.HasEslLink(h); hasEsl {
				eslThroughput, eslLink, err = eflib.GetDeviceEslInfo(h)
				if err != nil {
					continue
				}
			}

			gcuClock, clockVisible, err := eflib.GetDeviceClock(h)
			if err != nil {
				continue
			}

			powerMode, err := eflib.GetDevicePowerMode(h)
			if err != nil {
				continue
			}

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

			virtMode, err := eflib.GetDeviceVirtMode(h)
			if err != nil {
				continue
			}

			busid, err := eflib.GetDeviceBusID(h)
			if err != nil {
				continue
			}

			slot, err := eflib.GetDeviceSlotNumber(h)
			if err != nil {
				continue
			}

			fwVersion, err := eflib.GetFwVersion(h)
			if err != nil {
				continue
			}

			devSn, err := eflib.GetDevSn(h)
			if err != nil {
				continue
			}

			health, healthmsg := eflib.GetDeviceHealthState(minor, h)
			health2 := 1

			if health == true {
				health2 = 2
			}

			metrics.Devices = append(metrics.Devices,
				&Device{
					Host:          host,
					Index:         strconv.Itoa(int(index)),
					Minor:         strconv.Itoa(int(minor)),
					Name:          name,
					Uuid:          uuid,
					BusID:         busid,
					Slot:          slot,
					Health:        float64(health2),
					HealthMsg:     healthmsg,
					PowerUsage:    float64(powerUsage),
					PowerConsumption:   float64(powerConsumption),
					PowerCapability:    float64(powerCapability),
					MemorySize:    float64(memSize),
					MemoryUsed:    float64(memUsed),
					MemoryUsage:   float64(memUsage),
					GcuUsage:      float64(gcuUsage),
					SipUsage:      float64(sipUsage),
					ClusterUsage:  clusterUsage,
					PGUsage:       PGUsage,
					Temperature:   float64(temperature),
					GcuClock:      float64(gcuClock),
					ClockVisible:  clockVisible,
					PowerMode:     powerMode,
					EccStatus:     eccStatus,
					EslThroughput: eslThroughput,
					EslLink:       eslLink,
					PcieLinkSpeed: float64(pcieLinkSpeed),
					PcieLinkWidth: float64(pcieLinkWidth),
					PcieLink:      pcieLink,
					VirtMode:      virtMode,
					FwVersion:     fwVersion,
					DevSn:         devSn,
				})
		}
	} else {

		for index := uint32(0); index < uint32(count); index++ {
				h, err := eflib.GetDeviceHandle(index)
				if err != nil {
					continue
				}

				uuid, err := eflib.GetDeviceUUID(h)
                                if err != nil {
					continue
				}

				name, err := eflib.GetDeviceType(h)
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

				powerUsage, powerConsumption, powerCapability, err := eflib.GetDevicePowerInfo(h)
				if err != nil {
					continue
				}

				memUsage, memSize, memUsed, err := eflib.GetDeviceMemoryInfo(h)
				if err != nil {
					continue
				}

				gcuUsage, err := eflib.GetDeviceGcuUsage(h)
				if err != nil {
					continue
				}

				sipUsage, err := eflib.GetDeviceSipUsage(h)
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
					eslThroughput, eslLink, err = eflib.GetDeviceEslInfo(h)
					if err != nil {
						continue
					}
				}

				gcuClock, clockVisible, err := eflib.GetDeviceClock(h)
				if err != nil {
					continue
				}

				powerMode, err := eflib.GetDevicePowerMode(h)
				if err != nil {
					continue
				}

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

				busid, err := eflib.GetDeviceBusID(h)
				if err != nil {
					continue
				}

				slot, err := eflib.GetDeviceSlotNumber(h)
				if err != nil {
					continue
				}

				fwVersion, err := eflib.GetFwVersion(h)
				if err != nil {
					continue
				}

				devSn, err := eflib.GetDevSn(h)
				if err != nil {
					continue
				}
				health, healthmsg := eflib.GetDeviceHealthState(minor, h)
				health2 := 1

				if health == true {
					health2 = 2
				}
                                vindexList, err := eflib.GetVIndexList(h)
				if err != nil {
					continue
				}

				vusages, err := eflib.GetDeviceVUsage(h, vindexList)
				if err != nil {
					continue
				}

				vmemUsed, vmemTotal, vmemUsage, err := eflib.GetDeviceVMem(
					h,
					vindexList,
				)
				if err != nil {
					continue
				}


				metrics.Devices = append(metrics.Devices,
					&Device{
						Host:          host,
						Index:         strconv.Itoa(int(index)),
						Minor:         strconv.Itoa(int(minor)),
						Name:          name,
						Uuid:          uuid,
						BusID:         busid,
						Slot:          slot,
						Health:        float64(health2),
						HealthMsg:     healthmsg,
						PowerUsage:    float64(powerUsage),
						PowerConsumption:   float64(powerConsumption),
						PowerCapability:    float64(powerCapability),
						MemorySize:    float64(memSize),
						MemoryUsed:    float64(memUsed),
						MemoryUsage:   float64(memUsage),
						GcuUsage:      float64(gcuUsage),
						SipUsage:      float64(sipUsage),
						VGcuUsage:     vusages,
						VMemorySize:   vmemTotal,
						VMemoryUsed:   vmemUsed,
						VMemoryUsage:  vmemUsage,
						VIndexList:    vindexList,
						ClusterUsage:  clusterUsage,
						Temperature:   float64(temperature),
						GcuClock:      float64(gcuClock),
						ClockVisible:  clockVisible,
						PowerMode:     powerMode,
						EccStatus:     eccStatus,
						EslThroughput: eslThroughput,
						EslLink:       eslLink,
						PcieLinkSpeed: float64(pcieLinkSpeed),
						PcieLinkWidth: float64(pcieLinkWidth),
						PcieLink:      pcieLink,
						FwVersion:     fwVersion,
						DevSn:         devSn,
					})
		}
	}
	return metrics, nil
}
