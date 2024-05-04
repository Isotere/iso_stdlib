SHELL   := /bin/bash -euo pipefail
TIMEOUT := 10s
GOFLAGS := -mod=vendor
BINPATH = $(PWD)/bin

.PHONY: deps
deps:
	@go mod tidy && go mod vendor && go mod verify