# Copyright (C) 2022 Enflame Technologies, Inc.

.PHONY: all clean cmdline server daemon

VERSION_TXT ?= $(shell if [ -f "../../g_version.txt" ]; then echo "../../g_version.txt"; fi)
RELEASE_VERSION ?= $(shell if [ -n "$(VERSION_TXT)" ]; then grep -w 'GCU_EXPORTER_VERSION' $(VERSION_TXT) | awk -F '=' '{print $$2}'; fi)
GCU_EXPORTER_VERSION = $(RELEASE_VERSION)
ifeq ($(GCU_EXPORTER_VERSION),)
    GCU_EXPORTER_VERSION := 0.0.0
endif

all: cmdline server daemon
	@echo "Build completed: gcu-exporter-cmdline, gcu-exporter-server, gcu-exporter-daemon"

cmdline:
	@echo "Building gcu-exporter-cmdline..."
	@mkdir -p bin
	@$(MAKE) -C gcu-exporter-cmdline all
	@if [ -f gcu-exporter-cmdline/gcu-exporter-cmdline ]; then \
		cp gcu-exporter-cmdline/gcu-exporter-cmdline bin/gcu-exporter-cmdline; \
		echo "gcu-exporter-cmdline built successfully -> bin/gcu-exporter-cmdline"; \
	else \
		echo "Error: gcu-exporter-cmdline binary not found"; \
		exit 1; \
	fi

server:
	@echo "Building gcu-exporter-server..."
	@mkdir -p bin
	@$(MAKE) -C gcu-exporter-server build
	@if [ -f gcu-exporter-server/gcu-exporter-server ]; then \
		cp gcu-exporter-server/gcu-exporter-server bin/gcu-exporter-server; \
		echo "gcu-exporter-server built successfully -> bin/gcu-exporter-server"; \
	else \
		echo "Error: gcu-exporter-server binary not found"; \
		exit 1; \
	fi

daemon:
	@echo "Building gcu-exporter-daemon..."
	@mkdir -p bin
	@$(MAKE) -C gcu-exporter-daemon all
	@if [ -f gcu-exporter-daemon/gcu-exporter-daemon ]; then \
		cp gcu-exporter-daemon/gcu-exporter-daemon bin/gcu-exporter-daemon; \
		echo "gcu-exporter-daemon built successfully -> bin/gcu-exporter-daemon"; \
	else \
		echo "Error: gcu-exporter-daemon binary not found"; \
		exit 1; \
	fi

clean:
	@if [ -f ./bin/gcu-exporter-cmdline ]; then rm -f ./bin/gcu-exporter-cmdline; fi
	@if [ -f ./bin/gcu-exporter-server ]; then rm -f ./bin/gcu-exporter-server; fi
	@if [ -f ./bin/gcu-exporter-daemon ]; then rm -f ./bin/gcu-exporter-daemon; fi
	@if [ -d ./bin ] && [ -z "$$(ls -A ./bin 2>/dev/null)" ]; then rmdir ./bin; fi
	@cd gcu-exporter-cmdline && $(MAKE) clean
	@cd gcu-exporter-server && $(MAKE) clean
	@cd gcu-exporter-daemon && $(MAKE) clean
	@echo "Clean completed"
