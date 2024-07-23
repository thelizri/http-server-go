.ONESHELL:
	SHELL := /bin/bash

all: server

server:
	@set -e
	@tmpFile=$$(mktemp)
	@go build -o "$$tmpFile" cmd/server/*.go
	@exec "$$tmpFile" $(ARGS)

test:
	@go test ./... -v

ARGS := $(filter-out --,$(MAKEFLAGS) $(MAKECMDGOALS))

.PHONY: all server test
