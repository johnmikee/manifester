all: build

.PHONY: build

ifeq ($(GOPATH),)
	PATH := $(HOME)/go/bin:$(PATH)
else
	PATH := $(GOPATH)/bin:$(PATH)
endif

export GO111MODULE=on

BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
REVISION = $(shell git rev-parse HEAD)
REVSHORT = $(shell git rev-parse --short HEAD)
USER = $(shell whoami)
GOVERSION = $(shell go version | awk '{print $$3}')
NOW	= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
SHELL = /bin/sh
VERSION = $(shell git describe --tags --always)
WHEN = $(shell date +%Y-%m-%d)


ifneq ($(OS), Windows_NT)
	CURRENT_PLATFORM = linux
	ifeq ($(shell uname), Darwin)
		SHELL := /bin/sh
		CURRENT_PLATFORM = darwin
	endif
else
	CURRENT_PLATFORM = windows
endif

.pre-build:
	mkdir -p build/darwin
	mkdir -p build/linux

APP_NAME = manifester

build: manifester

clean:
	rm -rf build/
	rm -f *.zip

deps:
	go mod download

fmt:
	gofumpt -l -w .

static-check:
	staticcheck ./...

test:
	go test -cover ./...

tidy:
	go mod tidy

vet:
	go vet ./...

manifester: .pre-build
	go build -o build/$(CURRENT_PLATFORM)/$(APP_NAME) -ldflags="-X 'main.Version=$(VERSION)' -X 'app/build.User=$(USER)' -X 'app/build.Time=$(WHEN)'" main.go
