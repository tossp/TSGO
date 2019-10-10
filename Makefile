.PHONY: build fmt fmt-check install test vet clean

-include .env
EXE:=
ifeq ($(LANG),)
EXE=.exe
endif
DIST_DIR=./dist/
BINARY:=$(DIST_DIR)$(PROJECTNAME)$(EXE)
GITTAG:=$(shell git describe --tags || echo 'unknown')
GITVERSION:=$(shell git rev-parse HEAD)
PACKAGES:=$(shell go list ./... | grep -v /vendor/)
VETPACKAGES=`go list ./... | grep -v /vendor/ | grep -v /examples/`
GOFILES=`find . -name "*.go" -type f -not -path "./vendor/*"`
BUILD_TIME=`date +%FT%T%z`
LDFLAGS=-ldflags "-s -w -X github.com/tossp/tsgo/pkg/setting.ProjectName=$(PROJECTNAME) -X github.com/tossp/tsgo/pkg/setting.GitTag=$(GITTAG) -X github.com/tossp/tsgo/pkg/setting.GitVersion=${GITVERSION} -X github.com/tossp/tsgo/pkg/setting.BuildTime=${BUILD_TIME} -X github.com/tossp/tsgo/pkg/setting.BuildVersion=${VERSION}"

all: run
build:
	@echo " > Building binary..."
	@go build -v ${LDFLAGS} -tags=jsoniter -o ${BINARY} ./cmd/app

run: build
	@echo " > exec..."
	@${BINARY}

list:
	@echo ${PACKAGES}
	@echo ${VETPACKAGES}
	@echo ${GOFILES}
fmt:
	@echo " > gofmt..."
	@goimports -w ${GOFILES}

check: fmt
	@golint "-set_exit_status" ${GOFILES}
	@go vet ${GOFILES}

fmt-check:
	@diff=$$(gofmt -s -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

install:
	@govendor sync -v

test:
	@go test -cpu=1,2,4 -v -tags integration ./...

vet:
	@go vet $(VETPACKAGES)


clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
