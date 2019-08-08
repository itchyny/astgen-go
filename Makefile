GOBIN ?= $(shell go env GOPATH)/bin
export GO111MODULE=on

.PHONY: all
all: clean test

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint: $(GOBIN)/golint
	go vet ./...
	golint -set_exit_status ./...

$(GOBIN)/golint:
	GO111MODULE=off go get golang.org/x/lint/golint

.PHONY: clean
clean:
	go clean
