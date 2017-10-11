bin          := kafka-management-tool

.DEFAULT_GOAL := build

.PHONY: all
all: fmt imports clean build install

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: imports
imports:
	go get golang.org/x/tools/cmd/goimports
	goimports -l -w .

.PHONY: clean
clean:
	rm -f ./$(bin) ./$(bin)-*.zip $(test_path) build.log

.PHONY: build
build: clean
	go build

.PHONY: install
install:
	go install