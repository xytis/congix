PACKAGE  = congix
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || cat $(CURDIR)/.version 2> /dev/null || echo v0)

GO      = go
GODOC   = godoc
GOFMT   = gofmt
GODEP   = dep
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: run
run:
	$(GO) build \
		-tags release \
		-ldflags '-X $(PACKAGE)/main.Version=$(VERSION) -X $(PACKAGE)/main.BuildDate=$(DATE)' \
		-o bin/$(PACKAGE)
	CONSUL_HTTP_ADDR=consul-ui.hosts.lamoda.ru:8500 bin/$(PACKAGE) agent
