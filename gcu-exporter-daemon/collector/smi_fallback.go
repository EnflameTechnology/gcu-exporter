// Copyright (c) 2026 Enflame. All Rights Reserved.
//
// efsmi-json fallback for the daemon data path.
//
// The daemon collects device metrics from eflib by default. When eflib cannot
// produce a valid device snapshot (collectGCUMetrics returns an error or zero
// devices), the collector switches to querying `efsmi -q --json-format` and
// stays on it (sticky) for the rest of this process's lifetime; recovery back
// to the eflib path only happens when the daemon process restarts. If efsmi is
// missing, not executable, does not support --json-format, or the query/parse
// fails, the fallback returns an error and that scrape produces no device
// metrics.
//
// The field mapping is adapted from the abandoned gcu-exporter-smi work
// (gerrit 303568, TOPSCLOUD-12) which was validated on 150/217/34.

package collector

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"gcu-exporter/internal/eflib/efml"
)

const smiBinaryPath = "/usr/bin/efsmi"

// collectGCUMetricsFromSmi runs `efsmi -q --json-format`, parses it and maps
// the result to the same *Metrics the eflib path produces. Any failure — binary
// absent, json unsupported, query timeout, parse error, or zero devices —
// returns an error so the caller can surface it.
func collectGCUMetricsFromSmi() (*Metrics, error) {
	if _, err := os.Stat(smiBinaryPath); err != nil {
		return nil, fmt.Errorf("efsmi not found at %s: %w", smiBinaryPath, err)
	}
	out, err := runSmi([]string{"-q", "--json-format"}, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("efsmi -q --json-format failed (efsmi may not support json): %w", err)
	}
	var data smiQueryJSON
	if err := json.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("efsmi json unmarshal failed (efsmi may not support json): %w (output prefix: %.200s)", err, string(out))
	}
	if len(data.Devices) == 0 {
		return nil, fmt.Errorf("efsmi -q --json-format returned 0 devices")
	}

	host, _ := os.Hostname()
	slotByID := smiSlotByID()
	metrics := &Metrics{Count: uint32(len(data.Devices))}
	for _, ed := range data.Devices {
		d := mapSmiJSONToDevice(host, ed)
		// eflib reports the OAM/slot name, or "N/A" when empty; mirror that.
		d.Slot = slotByID[ed.ID]
		if d.Slot == "" {
			d.Slot = "N/A"
		}
		metrics.Devices = append(metrics.Devices, d)
	}

	// Pod labels come from the same kubelet pod-resources source as the eflib
	// path. efsmi -q does not expose vGCU, so cluster mapping is disabled.
	if IsInK8sCloud || NeedMonitPodInfo {
		if pods, err := getDevicePodInfo(socketPath, false); err != nil {
			fmt.Println("efsmi fallback getDevicePodInfo errored as ", err)
		} else {
			metrics.DeviceToPod = pods
		}
	}
	return metrics, nil
}

// runSmi runs efsmi with the given args, killing it after timeout.
func runSmi(args []string, timeout time.Duration) ([]byte, error) {
	cmd := exec.Command(smiBinaryPath, args...)
	cmd.Env = append(cmd.Env, "LANG=C", "LC_ALL=C")
	done := make(chan struct{}, 1)
	var out []byte
	var err error
	go func() {
		out, err = cmd.Output()
		done <- struct{}{}
	}()
	select {
	case <-time.After(timeout):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		return nil, fmt.Errorf("efsmi %v timed out after %v", args, timeout)
	case <-done:
		if err != nil {
			return nil, fmt.Errorf("exec efsmi %v: %w (output prefix: %.200s)", args, err, string(out))
		}
		return out, nil
	}
}

// smiQueryJSON is the subset of the `efsmi -q --json-format` schema we map.
type smiQueryJSON struct {
	Devices []smiDeviceJSON `json:"devices"`
}

type smiDeviceJSON struct {
	ID         int `json:"id"`
	DeviceInfo struct {
		DevName string `json:"Dev Name"`
		DevUUID string `json:"Dev UUID"`
		DevSN   string `json:"Dev SN"`
		Health  bool   `json:"Health"`
	} `json:"Device Info"`
	Firmware map[string]interface{} `json:"Firmware Info"`
	PCIe     struct {
		Domain   string `json:"Domain"`
		Bus      string `json:"Bus"`
		Dev      string `json:"Dev"`
		Func     string `json:"Func"`
		LinkInfo struct {
			MaxLinkSpeed string `json:"Max Link Speed"`
			MaxLinkWidth string `json:"Max Link Width"`
			CurLinkSpeed string `json:"Cur Link Speed"`
			CurLinkWidth string `json:"Cur Link Width"`
		} `json:"Link Info"`
	} `json:"PCIe Info"`
	Clock struct {
		GcuClk float64 `json:"GCU CLK"`
	} `json:"Clock Info"`
	Power struct {
		Usage        float64 `json:"Power Usage"`
		CurrentLimit float64 `json:"Current Power Limit"`
		DefaultLimit float64 `json:"Default Power Limit"`
		MaxLimit     float64 `json:"Max Power Limit"`
	} `json:"GCU Power Readings"`
	LowPower int `json:"Low Power Mode"`
	DevMem   struct {
		TotalSize    float64 `json:"Total Size"`
		ReservedSize float64 `json:"Reserved Size"`
		UsedSize     float64 `json:"Used Size"`
	} `json:"Device Mem Info"`
	Temp struct {
		GcuTemp float64 `json:"GCU Temp"`
	} `json:"Temperature Info"`
	Usage struct {
		GCUUsage float64 `json:"GCU Usage"`
	} `json:"Device Usage Info"`
	EccErrors struct {
		Aggregate struct {
			L1CECC     uint64 `json:"L1 Cache CECC"`
			L1UECC     uint64 `json:"L1 Cache UECC/Parity"`
			L2CECC     uint64 `json:"L2 Cache CECC"`
			L2UECC     uint64 `json:"L2 Cache UECC"`
			DevMemCECC uint64 `json:"Device Memory CECC"`
			DevMemUECC uint64 `json:"Device Memory UECC"`
		} `json:"Aggregate"`
	} `json:"ECC Errors"`
}

var (
	smiLinkSpeedRe = regexp.MustCompile(`(?i)gen(\d+)`)
	smiLinkWidthRe = regexp.MustCompile(`(?i)x(\d+)`)
)

// smiParseLinkSpeed turns "Gen5" into 5.0; returns 0 on parse failure.
func smiParseLinkSpeed(s string) float64 {
	m := smiLinkSpeedRe.FindStringSubmatch(strings.TrimSpace(s))
	if len(m) < 2 {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
	return float64(n)
}

// smiParseLinkWidth turns "X16" into 16.0; returns 0 on parse failure.
func smiParseLinkWidth(s string) float64 {
	m := smiLinkWidthRe.FindStringSubmatch(strings.TrimSpace(s))
	if len(m) < 2 {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
	return float64(n)
}

// smiSlotByID runs `efsmi -L` and returns device id -> slot/OAM name so the
// efsmi path can report the same slot label as eflib's GetDeviceSlotNumber.
// Returns nil on failure; callers default a missing slot to "N/A" (eflib's
// value for an empty slot).
var (
	slotCacheMu     sync.Mutex
	slotCache       map[int]string
	slotCacheLoaded bool
)

func smiSlotByID() map[int]string {
	slotCacheMu.Lock()
	defer slotCacheMu.Unlock()
	if slotCacheLoaded {
		return slotCache
	}
	out, err := runSmi([]string{"-L"}, 10*time.Second)
	if err != nil {
		return nil // leave uncached so a later scrape retries
	}
	res := map[int]string{}
	lines := strings.Split(string(out), "\n")
	// The Slot column value can contain a space (e.g. "Slot 8"), so slice by
	// the header column positions instead of splitting on whitespace.
	slotStart, idStart := -1, -1
	for _, line := range lines {
		if strings.Contains(line, "DEV") && strings.Contains(line, "Slot") && strings.Contains(line, "UUID") {
			slotStart = strings.Index(line, "Slot")
			idStart = strings.Index(line, "ID")
			break
		}
	}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		id, err := strconv.Atoi(fields[0])
		if err != nil {
			continue // header / separator lines
		}
		slot := ""
		if slotStart >= 0 && idStart > slotStart && len(line) >= idStart {
			slot = strings.TrimSpace(line[slotStart:idStart])
		} else if len(fields) >= 3 {
			slot = fields[2]
		}
		res[id] = slot
	}
	// slot/OAM topology is static; cache so efsmi -L runs only once.
	slotCache = res
	slotCacheLoaded = true
	return slotCache
}

// mapSmiJSONToDevice maps one efsmi -q JSON device to the daemon Device.
// Fields efsmi does not expose (Sip / Cluster / PG / vGCU / ESL / Slot) are
// left at their eflib-compatible zero/sentinel values.
func mapSmiJSONToDevice(host string, ed smiDeviceJSON) *Device {
	// busid: keep the sysfs-style path the eflib path produces.
	busid := fmt.Sprintf("/sys/module/enflame/drivers/pci:enflame/%s:%s:%s.%s",
		ed.PCIe.Domain, ed.PCIe.Bus, ed.PCIe.Dev, ed.PCIe.Func)

	// uuid: prefix "GCU-" to match eflib.
	uuid := ed.DeviceInfo.DevUUID
	if uuid != "" && !strings.HasPrefix(uuid, "GCU-") {
		uuid = "GCU-" + uuid
	}

	// Health: efsmi bool -> eflib semantics 1=unhealthy 2=healthy.
	health := 1.0
	healthMsg := "Unhealthy"
	if ed.DeviceInfo.Health {
		health = 2.0
		healthMsg = "Healthy"
	}

	// efsmi mem sizes are MiB; convert to bytes. eflib's memory_used includes
	// the driver-reserved memory, so add Reserved to Used to match it.
	memTotal := ed.DevMem.TotalSize * 1024 * 1024
	memUsed := (ed.DevMem.UsedSize + ed.DevMem.ReservedSize) * 1024 * 1024
	var memUsage float64
	if memTotal > 0 {
		memUsage = memUsed / memTotal
	}

	// FwVersion: eflib's GetFwVersion reports the SSM Boot FW Version (verified
	// on driver 1.9.7: eflib "31.5.3" == efsmi "SSM Boot FW Version"), so read
	// the same field to keep fw_version identical across data sources.
	fwVer := ""
	for _, k := range []string{"SSM Boot FW Version", "SSM RT FW Version", "SP FW Version"} {
		if v, ok := ed.Firmware[k].(string); ok && v != "" {
			fwVer = v
			break
		}
	}

	// power_capability: Current -> Max -> Default, matching the eflib path
	// (the currently effective cap, not the hardware TDP ceiling).
	powerCap := ed.Power.CurrentLimit
	if powerCap == 0 {
		powerCap = ed.Power.MaxLimit
	}
	if powerCap == 0 {
		powerCap = ed.Power.DefaultLimit
	}
	// power_usage: efsmi gives absolute W; eflib exposes a percentage. Convert.
	var powerUsagePct float64
	if powerCap > 0 {
		powerUsagePct = ed.Power.Usage / powerCap * 100
	}

	// name: keep the "Enflame " prefix so the value matches the eflib naming convention.
	name := strings.TrimSpace(ed.DeviceInfo.DevName)
	if name != "" && !strings.HasPrefix(name, "Enflame ") {
		name = "Enflame " + name
	}

	idStr := strconv.Itoa(ed.ID)

	return &Device{
		Host:             host,
		Index:            idStr,
		Minor:            idStr, // efsmi has no /dev minor; use the efsmi id (documented divergence)
		Name:             name,
		Uuid:             uuid,
		BusID:            busid,
		Slot:             "", // set by the caller from efsmi -L
		FwVersion:        fwVer,
		DevSn:            ed.DeviceInfo.DevSN,
		Health:           health,
		HealthMsg:        healthMsg,
		PowerUsage:       powerUsagePct,
		PowerConsumption: ed.Power.Usage,
		PowerCapability:  powerCap,
		MemorySize:       memTotal,
		MemoryUsed:       memUsed,
		MemoryUsage:      memUsage,
		GcuUsage:         ed.Usage.GCUUsage,
		SipUsage:         0,             // efsmi -q does not expose sip usage
		ClusterUsage:     []float64{-1}, // eflib S60 sentinel (-1 = not supported)
		PGUsage:          nil,           // eflib S60 emits none
		Temperature:      ed.Temp.GcuTemp,
		GcuClock:         ed.Clock.GcuClk,
		ClockVisible:     true,
		PowerMode:        fmt.Sprintf("LP%d", ed.LowPower),
		// EccStatus / PcieLink must be non-nil (collectors dereference them).
		EccStatus: &efml.DevEccStatus{
			Enabled: true,
			Ecnt_sb: uint(ed.EccErrors.Aggregate.L1CECC + ed.EccErrors.Aggregate.L2CECC + ed.EccErrors.Aggregate.DevMemCECC),
			Ecnt_db: uint(ed.EccErrors.Aggregate.L1UECC + ed.EccErrors.Aggregate.L2UECC + ed.EccErrors.Aggregate.DevMemUECC),
		},
		EslThroughput: nil,
		EslLink:       nil,
		PcieLinkSpeed: smiParseLinkSpeed(ed.PCIe.LinkInfo.CurLinkSpeed),
		PcieLinkWidth: smiParseLinkWidth(ed.PCIe.LinkInfo.CurLinkWidth),
		PcieLink: &efml.LinkInfo{
			Link_Speed:     int64(smiParseLinkSpeed(ed.PCIe.LinkInfo.CurLinkSpeed)),
			Max_Link_Speed: int64(smiParseLinkSpeed(ed.PCIe.LinkInfo.MaxLinkSpeed)),
			Link_Width:     int64(smiParseLinkWidth(ed.PCIe.LinkInfo.CurLinkWidth)),
			Max_Link_Width: int64(smiParseLinkWidth(ed.PCIe.LinkInfo.MaxLinkWidth)),
		},
		VirtMode: 0,
	}
}
