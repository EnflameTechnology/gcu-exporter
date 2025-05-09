/////////////////////////////////////////////////////////////////////////////
//  @brief API interface of Enflame Managerment Library
//
//  Enflame Tech, All Rights Reserved. 2022 Copyright (C)
/////////////////////////////////////////////////////////////////////////////

package efml

// #cgo LDFLAGS: -ldl  -Wl,--unresolved-symbols=ignore-in-object-files
// #include "stdbool.h"
// #include "efml.h"
import "C"

/*
 * @brief Enfalme Management Library get the device low power mode.
 */
func (h Handle) GetDevIsLowPowerMode() (bool, error) {
	var is_low_power_mode C.bool
	r := C.EfmlGetDevIsLowPowerMode(C.uint(h.Dev_Idx), &is_low_power_mode)
	if r == C.EFML_ERROR_NOT_SUPPORTED {
		return false, nil
	}

	return bool(is_low_power_mode), errorString(r)
}

/*
 * @brief Enfalme Management Library get the device power mode.
 */
func (h Handle) GetDevPowerMode() (int, error) {
	var is_low_power_mode C.bool
	var power_mode C.efmlPowerMode_t
	var r C.efmlReturn_t
	if GetPowerModeEnabled == false {
		r := C.EfmlGetDevIsLowPowerMode(C.uint(h.Dev_Idx), &is_low_power_mode)
		if r == C.EFML_ERROR_NOT_SUPPORTED {
			return -1, nil
		}

		if bool(is_low_power_mode) == true {
			return int(C.EFML_SLEEP), nil
		} else {
			return int(C.EFML_ACTIVE), nil
		}

	} else {
		r := C.EfmlGetPowerMode(C.uint(h.Dev_Idx), &power_mode)
		if r == C.EFML_ERROR_NOT_SUPPORTED {
			return -1, nil
		}
	}

	return int(power_mode), errorString(r)
}

/*
 * @brief Enflame Management Library get the device sip usage info.
 */
func (h Handle) GetDevSipUsage() (float32, error) {
	var usage C.float
	if GetDevSipUsageEnabled == true {
		r := C.EfmlGetDevSipUsage(C.uint(h.Dev_Idx), &usage)
		if r == C.EFML_ERROR_NOT_SUPPORTED {
			return -1, nil
		}
		return float32(usage), errorString(r)
	} else {
		return -1, nil
	}
}
