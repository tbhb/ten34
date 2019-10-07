export GOOS=$(shell go env GOOS)
export GO_BUILD=env GO11MODULE=on go build -ldflags="-s -w"
export GO_INSTALL=env GO11MODULE=on go install
export GO_TEST=env GOTRACEBACK=all GO11MODULE=on go test
export GO_VET=env GO11MODULE=on go vet
export GO_RUN=env GO11MODULE=on go run
export PATH := $(PWD)/bin/$(GOOS):$(PATH)

SOURCES := $(shell find . -name '*.go' -not -name '*_test.go') go.mod go.sum
SOURCES_NO_VENDOR := $(shell find . -path ./vendor -prune -o -name "*.go" -not -name '*_test.go' -print)

all: clean vet test build

bench:
	$(GO_TEST) -bench=. -run=^$$ ./...

build: $(SOURCES)
	$(GO_BUILD) -o bin/ten34 cmd/ten34/main.go

clean:
	$(RM) -r bin
	$(RM) -r dist

fmt: $(SOURCES_NO_VENDOR)
	gofmt -w -s $^

test:
	$(GO_TEST) ./...

tidy:
	GO11MODULE=on go mod tidy

vet:
	$(GO_VET) -v ./...

.PHONY: all bench clean fmt test tidy