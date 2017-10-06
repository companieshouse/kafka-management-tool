.PHONY: all
all: fmt imports build

.PHONY: fmt
fmt:
	gofmt .

.PHONY: imports
imports:
	go get golang.org/x/tools/cmd/goimports
	goimports -l -w .

.PHONY: install
install:
	go install

.PHONY: build
build:
	 rm kafka-management-tool
	go build

.PHONY: all install