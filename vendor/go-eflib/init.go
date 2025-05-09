// Copyright (c) 2022 Enflame. All Rights Reserved.

package eflib

import (
	"go-eflib/efml"
)

func Init(with_driver bool) error {
	if (with_driver) {
		// init with driver
		return efml.InitV2(false)
	} else {
		// init without driver
		return efml.InitV2(true)
	}
}

func Shutdown() error {
	return efml.Shutdown()
}
