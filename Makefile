bin          := kafka-management-tool

.PHONY: all
all: clean build install

.PHONY: fmt
fmt:
	go fmt ./...
	goimports ./..

.PHONY: test-unit
test-unit:
	go test ./... -run 'Unit'

.PHONY: deps
deps:
	go get github.com/companieshouse/$(bin)

.PHONY: clean
clean:
	rm -f ./$(bin) ./$(bin)-*.zip $(test_path) build.log

.PHONY: build
build: deps fmt
	go build -o ./$(bin)

.PHONY: install
install:
	go install