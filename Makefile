# Set VERSION to the latest version tag name. Assuming version tags are formatted 'v*'
VERSION := $(shell git describe --always --abbrev=0 --tags --match "v*" $(git rev-list --tags --max-count=1))
BUILD := $(shell git rev-parse $(VERSION))
PROJECTNAME := "echo-server"
# We pass that to the main module to generate the correct help text
PROGRAMNAME := $(PROJECTNAME)

# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/bin
GOBUILD := $(GOBASE)/build
GOFILES := $(shell find . -type f -name '*.go' -not -path './vendor/*')
GOOS_DARWIN := "darwin"
GOOS_LINUX := "linux"
GOARCH := "amd64"

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -X=main.ProgramName=$(PROGRAMNAME)"

# Redirect error output to a file, so we can show it in development mode.
STDERR := $(GOBUILD)/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID := $(GOBUILD)/.$(PROJECTNAME).pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

default: install format lint test compile

## install: Install missing dependencies.
install: go-get

format: go-format

lint: go-lint

## compile: Compile the binary.
compile:
	@[ -d $(GOBUILD) ] || mkdir -p $(GOBUILD)
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s go-compile 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2


test: go-test

## clean: Clean build files. Runs `go clean` internally.
clean:
	@-rm $(GOBIN)/$(PROGRAMNAME)* 2> /dev/null
	@-$(MAKE) go-clean

go-lint:
	@echo "  >  Linting source files..."
	 gofmt -d $(GOFILES)

go-format:
	@echo "  >  Formating source files..."
	 gofmt -s -w $(GOFILES)

go-compile: go-get go-build-linux go-build-darwin

go-test:
	go test -v `go list ./...`

go-build-linux:
	@echo "  >  Building binary linux binaries..."
	@GOPATH=$(GOPATH) GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH) GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/$(PROGRAMNAME)-$(GOOS_LINUX)-$(GOARCH) ./cmd/echoserver

go-build-darwin:
	@echo "  >  Building binary darwin binaries..."
	@GOPATH=$(GOPATH) GOOS=$(GOOS_DARWIN) GOARCH=$(GOARCH) GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/$(PROGRAMNAME)-$(GOOS_DARWIN)-$(GOARCH) ./cmd/echoserver

go-generate:
	@echo "  >  Generating dependency files..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go generate $(generate)

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get ./...

go-install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean ./cmd/echoserver

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo