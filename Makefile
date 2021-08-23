# GO Parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) glean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

# Binary Name (edit if you don't want the default)
BINARY_NAME=$(shell basename $(CURDIR))

# Compiler flags
LD_FLAGS=-X 'main.version=$$(git describe --tags)' -X 'main.date=$$(date +"%Y.%m.%d_%H%M%S")' -X 'main.rev=$$(git rev-parse --short HEAD)' -X 'main.branch=$$(git rev-parse --abbrev-ref HEAD | tr -d '\040\011\012\015\n')'

# Tool Arguments
TAGS=json,yaml,xml

build: deps
    export GO111MODULE=on
    [ -d bin ] || mkdir bin
    GOOS=linux $(GOBUILD) -ldflags "$(LD_FLAGS)" -o bin/$(BINARY_NAME) -v .
    GOOS=windows $(OBUILD) -ldflags "$(LD_FLAGS)" -o bin/$(BINARY_NAME).exe -v .
    
clean:
    $(GOCLEAN)
    rm -rf bin

deps:
    export GOPRIVATE=github.com/bengrewell
    $(GOGET) -u ./...

install-tools:
    go install google.golang.org/protobuf/cmd/protoc-gen-go
    go get github.com/fatih/gomodifytags

run:
    $(GORUN) cmd/main.go
    
tags:
    gomodifytags -file $(FILE) -all -add-tags $(TAGS) -w
