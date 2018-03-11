# Go parameters
GOCMD?=go
GOBUILD=$(GOCMD) build
GOVET=$(GOCMD) vet
GOTEST=$(GOCMD) test
GOBIN=${GOPATH}/bin
GODEP=$(GOTEST) -i
GOTESTFLAGS?=-coverprofile=cover.out
GOBUILDFLAGS?=

GOOS?=`go env GOOS`
GOARCH?=`go env GOARCH`

# project
VERSION?=`cat VERSION`
BUILD?=`git rev-parse HEAD`

TEST_PATHS=. ./pkg/...

LDFLAGS=-ldflags "-X main.AppVersion=$(VERSION) -X main.AppBuildID=$(BUILD) -s -w"

.DEFAULT_GOAL := test

.PHONY: test build clean

test:
	${GOVET} ${TEST_PATHS}
	${GOTEST} ${GOTESTFLAGS} ${TEST_PATHS}

clean:
	rm -f chuckssim || true

build:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} ${GOBUILD} ${GOBUILDFLAGS} ${LDFLAGS} 

