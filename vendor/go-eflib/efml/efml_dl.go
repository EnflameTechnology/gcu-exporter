/////////////////////////////////////////////////////////////////////////////
//  @brief API interface of Enflame Management Library
//
//  Enflame Tech, All Rights Reserved. 2024 Copyright (C)
/////////////////////////////////////////////////////////////////////////////

package efml

import (
	"unsafe"
)

/*
#include <dlfcn.h>
#include <stdbool.h>
#include <string.h>
#include "efml.h"
*/
import "C"

type dl_handle_ptr struct{ handles []unsafe.Pointer }

var dl dl_handle_ptr

func (dl *dl_handle_ptr) InitV2(noDriver bool) C.efmlReturn_t {
	openFlags := C.RTLD_LAZY | C.RTLD_GLOBAL
	searchPaths := []string{"", "/usr/lib/", "/usr/lib64/", "/usr/lib/x86_64-linux-gnu/"}
	searchNames := []string{"libefml.so", "libefml.so.1", "libefml.so.2", "libefml.so.2.1", "libefml.so.2.0", "libefml.so.1.0.0"}
	var handle unsafe.Pointer
	for _, path := range searchPaths {
		for _, name := range searchNames {
			handle = C.dlopen(C.CString(path+name), C.int(openFlags))
			if handle != C.NULL {
				dl.handles = append(dl.handles, handle)
				return C.EfmlInit(C.bool(noDriver))
			}
		}
	}
	return C.EFML_ERROR_LIBRARY_NOT_FOUND
}

func (dl *dl_handle_ptr) Shutdown() C.efmlReturn_t {
	C.EfmlShutdown()

	for _, handle := range dl.handles {
		err := C.dlclose(handle)
		if err != 0 {
			return C.EFML_ERROR_FUNCTION_NOT_FOUND
		}
	}

	return C.EFML_SUCCESS
}

func (dl *dl_handle_ptr) lookupSymbol(symbol string) C.efmlReturn_t {
	for _, handle := range dl.handles {
		C.dlerror()
		C.dlsym(handle, C.CString(symbol))
		if unsafe.Pointer(C.dlerror()) == C.NULL {
			return C.EFML_SUCCESS
		}
	}

	return C.EFML_ERROR_FUNCTION_NOT_FOUND
}
