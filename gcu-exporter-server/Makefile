# Copyright (C) 2022 Enflame Technologies, Inc.

.PHONY: all build clean

VERSION_TXT ?= $(shell if [ -f "../../../g_version.txt" ]; then echo "../../../g_version.txt"; fi)
RELEASE_VERSION ?= $(shell if [ -n "$(VERSION_TXT)" ]; then grep -w 'GCU_EXPORTER_VERSION' $(VERSION_TXT) | awk -F '=' '{print $$2}'; fi)
GCU_EXPORTER_VERSION = $(RELEASE_VERSION)
ifeq ($(GCU_EXPORTER_VERSION),)
    GCU_EXPORTER_VERSION := 0.0.0
endif

all: build

build:
	@echo "Version: $(GCU_EXPORTER_VERSION)"
	@go build -o gcu-exporter-server -ldflags "-X main.Version=$(GCU_EXPORTER_VERSION)" .

clean:
	rm -f gcu-exporter-server

install: build
	cp gcu-exporter-server /usr/local/bin/

