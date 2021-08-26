GOBIN ?= $(shell go env GOPATH)/bin

.PHONY: all
all: test

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint: $(GOBIN)/staticcheck
	go vet ./...
	staticcheck ./...

$(GOBIN)/staticcheck:
	cd && go get honnef.co/go/tools/cmd/staticcheck

.PHONY: clean
clean:
	go clean
