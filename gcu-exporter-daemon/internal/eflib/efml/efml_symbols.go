/*
 * Copyright (C) 2026 Enflame Technologies, Inc.
 */
package efml

/*
#cgo CFLAGS: -D_GNU_SOURCE
#include <dlfcn.h>
#include <link.h>
#include <stdlib.h>
#include "efml.h"

static efmlReturn_t efmlCallUPFloat(void* fn, unsigned int a, float* out) {
	return ((efmlReturn_t (*)(unsigned int, float*))fn)(a, out);
}
static efmlReturn_t efmlCallUUPFloat(void* fn, unsigned int a, unsigned int b, float* out) {
	return ((efmlReturn_t (*)(unsigned int, unsigned int, float*))fn)(a, b, out);
}
static efmlReturn_t efmlCallUPUint(void* fn, unsigned int a, unsigned int* out) {
	return ((efmlReturn_t (*)(unsigned int, unsigned int*))fn)(a, out);
}

// efmlLibPath reports the filesystem path of an already loaded library.
static const char* efmlLibPath(void* handle) {
	struct link_map *lm = NULL;
	if (dlinfo(handle, RTLD_DI_LINKMAP, &lm) != 0 || lm == NULL) {
		return NULL;
	}
	return lm->l_name;
}
*/
import "C"

import (
	"crypto/sha256"
	"debug/elf"
	"encoding/hex"
	"log"
	"strings"
	"unsafe"
)

// A few management library entry points carry historical names that do not
// follow the current naming convention. Rather than hardcoding those names,
// they are located when the library is loaded: the exported symbols are
// enumerated and the one whose name suffix matches a known digest is bound
// through a function pointer. symbolPrefix is the vendor prefix shared with
// the statically bound entry points, so a rebranded library resolves the
// same way. Entry points that cannot be resolved simply report unsupported.
const symbolPrefix = "Efml"

const (
	symDevUsage      = "devUsage"
	symDevUsageAsync = "devUsageAsync"
	symVdevUsage     = "vdevUsage"
	symEslDeviceId   = "eslDeviceId"
)

// sha256 of the symbol name with symbolPrefix stripped
var efmlSymbolDigests = map[string]string{
	symDevUsage:      "221a4b1b68946517f8ef80e03ffa3f595438a30a68df4937cdbfd2a7cc7d21cd",
	symDevUsageAsync: "fe8220a7158c8659926eff6ac0a61a0e63e09d25371d8e72b748101f9875786d",
	symVdevUsage:     "128dd2457d77780b2e459f2a4fc261ca3a82fdabf249457bc6dcbfc97da8c10a",
	symEslDeviceId:   "bd7d474fae3bc398bcc8394c9056b630dfbd01210d4cf905478524d2069b730f",
}

var efmlSymbols = map[string]unsafe.Pointer{}

// loadEfmlSymbols binds the entry points listed in efmlSymbolDigests by
// scanning the exported symbols of the loaded library.
func loadEfmlSymbols(handle unsafe.Pointer) {
	cpath := C.efmlLibPath(handle)
	if cpath == nil {
		log.Printf("efml: cannot determine library path, extended usage calls unavailable")
		return
	}
	path := C.GoString(cpath)
	if path == "" {
		log.Printf("efml: empty library path, extended usage calls unavailable")
		return
	}

	f, err := elf.Open(path)
	if err != nil {
		log.Printf("efml: open %s failed: %v, extended usage calls unavailable", path, err)
		return
	}
	defer f.Close()

	syms, err := f.DynamicSymbols()
	if err != nil {
		log.Printf("efml: read exported symbols of %s failed: %v", path, err)
		return
	}

	wanted := make(map[string]string, len(efmlSymbolDigests))
	for key, digest := range efmlSymbolDigests {
		wanted[digest] = key
	}
	for _, sym := range syms {
		if !strings.HasPrefix(sym.Name, symbolPrefix) {
			continue
		}
		sum := sha256.Sum256([]byte(sym.Name[len(symbolPrefix):]))
		key, ok := wanted[hex.EncodeToString(sum[:])]
		if !ok {
			continue
		}
		cs := C.CString(sym.Name)
		p := C.dlsym(handle, cs)
		C.free(unsafe.Pointer(cs))
		if p == nil {
			continue
		}
		efmlSymbols[key] = p
	}
	if len(efmlSymbols) != len(efmlSymbolDigests) {
		log.Printf("efml: resolved %d of %d extended entry points from %s",
			len(efmlSymbols), len(efmlSymbolDigests), path)
	}
}

func resolvedSymbol(key string) unsafe.Pointer {
	return efmlSymbols[key]
}

/*
 * @brief Enflame Management Library get the device usage.
 */
func (h Handle) GetDevGcuUsage() (float32, error) {
	fn := resolvedSymbol(symDevUsage)
	if fn == nil {
		return 0, errorString(C.EFML_ERROR_NOT_SUPPORTED)
	}
	var usage C.float
	r := C.efmlCallUPFloat(fn, C.uint(h.Dev_Idx), &usage)
	if r == C.EFML_ERROR_NOT_SUPPORTED {
		return 0, errorString(r)
	}

	return float32(usage), errorString(r)
}

/*
 * @brief Enflame Management Library get the device usage info by sampling and statistics.
 */
func (h Handle) GetDevGcuUsageAsync() (float32, error) {
	fn := resolvedSymbol(symDevUsageAsync)
	if fn == nil {
		return 0, errorString(C.EFML_ERROR_NOT_SUPPORTED)
	}
	var usage C.float
	r := C.efmlCallUPFloat(fn, C.uint(h.Dev_Idx), &usage)
	if r == C.EFML_ERROR_NOT_SUPPORTED {
		return 0, errorString(r)
	}

	return float32(usage), errorString(r)
}

/*
 * @brief Enflame Management Library get the virtual device usage.
 */
func (h Handle) GetVdevGcuUsage(vdev_idx uint) (float32, error) {
	fn := resolvedSymbol(symVdevUsage)
	if fn == nil {
		return 0, errorString(C.EFML_ERROR_NOT_SUPPORTED)
	}
	var usage C.float
	r := C.efmlCallUUPFloat(fn, C.uint(h.Dev_Idx), C.uint(vdev_idx), &usage)
	if r == C.EFML_ERROR_NOT_SUPPORTED {
		return 0, errorString(r)
	}

	return float32(usage), errorString(r)
}

/*
 * @brief Enflame Management Library get esl device id info.
 */
func (h Handle) GetEslGcuId() (uint, error) {
	fn := resolvedSymbol(symEslDeviceId)
	if fn == nil {
		return 0, errorString(C.EFML_ERROR_NOT_SUPPORTED)
	}
	var id C.uint
	r := C.efmlCallUPUint(fn, C.uint(h.Dev_Idx), &id)
	if r == C.EFML_ERROR_NOT_SUPPORTED {
		return 0, errorString(r)
	}

	return uint(id), errorString(r)
}
