# Copyright (C) 2025 Enflame Technologies, Inc.

.PHONY: all  clean

VERSION_TXT ?= $(shell if [ -f "../../../g_version.txt" ]; then echo "../../../g_version.txt"; fi)
RELEASE_VERSION ?= $(shell if [ -n "$(VERSION_TXT)" ]; then grep -w 'GCU_EXPORTER_VERSION' $(VERSION_TXT) | awk -F '=' '{print $$2}'; fi)
GCU_EXPORTER_VERSION = $(RELEASE_VERSION)
ifeq ($(GCU_EXPORTER_VERSION),)
    GCU_EXPORTER_VERSION := 0.0.0
endif

all:
	@echo "Version: $(GCU_EXPORTER_VERSION)"
	@go build  -ldflags "-X main.Version=$(GCU_EXPORTER_VERSION)" -v -o gcu-exporter

clean:
	@if [ -f ./gcu-exporter ]; then rm gcu-exporter; else echo " "; fi
