// Copyright (c) 2022 Enflame. All Rights Reserved.

package eflib

import (
	"log"
	"os"
)

func FileIsExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

func GetHostName() string {
	name, err := os.Hostname()
	if err != nil {
		log.Printf("Failed to get host name\n")
		return "Unknown Host"
	}
	return name
}
