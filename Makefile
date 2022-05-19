.DEFAULT_GOAL := all
GO111MODULE    = on

SRC_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

init:
	go mod init

mod:
	go mod download
	go mod tidy
	go mod verify


fmt:
	go fmt ./...


fix:
	go fix ./...


vet:
	go vet ./...


build: fmt fix vet
	go build ./...


test: build
	go test ./...

clean:
	go clean -cache -modcache -testcache


all: build

.PHONY: init mod fmt fix vet build test clean all
