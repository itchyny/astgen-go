BIN := goastgen
export GO111MODULE=on

.PHONY: all
all: clean test

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint: lintdeps
	go vet ./...
	golint -set_exit_status ./...

.PHONY: lintdeps
lintdeps:
	GO111MODULE=off go get golang.org/x/lint/golint

.PHONY: clean
clean:
	go clean
