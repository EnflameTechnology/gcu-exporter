// Copyright (c) 2022 Enflame. All Rights Reserved.

package eflib

import (
	"log"
	"time"

	"go-eflib/efml"
)


var oldSsmCount = [GCU_DEV_MAX]string{}

func GetPcieLinkWidthHealthy(h efml.Handle) (bool) {
	linkInfo, _ := h.GetPcieLinkInfo()
	if linkInfo.Link_Width < linkInfo.Max_Link_Width {
		log.Printf("Device pcie link width unhealthy.\n")
		return false
	}
	return true
}

func GetDeviceTemperatureHealthy(h efml.Handle) (bool) {
	var count uint32 = 0
	thermalInfo, _ := h.GetDevTempV2()
	if thermalInfo.Cur_Asic_Temp > 95 {
		// waiting for 15 seconds
		for idx := 0; idx < 15; idx++ {
			time.Sleep(time.Duration(1000) * time.Millisecond)
			thermalInfo, _ := h.GetDevTempV2()
			if thermalInfo.Cur_Asic_Temp > 95 {
				count++
			}
		}
	}
	// > 10 seconds
	if count > 10 {
		log.Printf("The device Temperature > 95°\n")
		return false
	}

	return true
}

func GetDeviceEccDbeHealthy(h efml.Handle) (bool) {
	eccStatus, _ := h.GetDevEccStatus()
	if eccStatus.Ecnt_db > 0 {
		log.Printf("The device ecc DBE error\n")
		return false
	}

	return true
}

func GetDeviceRmaFlags(h efml.Handle) (bool) {
/*
	rmaStatus, _ := h.GetDevRmaStatus()
	if rmaStatus.Flags {
		log.Printf("The device rma flasg == true\n")
		return true
	}
*/
	return false
}

func GetDeviceHealthState(id uint32, h efml.Handle) (bool, string) {
	var healthmsg string = ""
	var errcount uint = 0
/*
	eccdbe := GetDeviceEccDbeHealthy(h)
	if eccdbe == false {
		errcount++
		healthmsg = healthmsg + " ecc dbe > 0"
	}
	pciewidth := GetPcieLinkWidthHealthy(h)
	if pciewidth == false {
		errcount++
		healthmsg = healthmsg + " pcie link width < max link width"
	}
	heartbeat := GetDeviceSSmHeartbeatHealthy(id)
	if heartbeat == false {
		errcount++
		healthmsg = healthmsg + " SSM no heartbeat"
	}
	temperature := GetDeviceTemperatureHealthy(h)
	if temperature == false {
		errcount++
		healthmsg = healthmsg + " temperature > 95°"
	}
*/
	rmaflags := GetDeviceRmaFlags(h)
	if rmaflags == true {
		errcount++
		healthmsg = healthmsg + " RMA true"
	}

	if errcount > 0 {
		healthmsg = "Unhealthy: " + healthmsg
		return false, healthmsg
	} else {
		healthmsg = "Healthy"
		return true, healthmsg
	}
}
