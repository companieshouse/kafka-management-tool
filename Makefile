CHS_ENV_HOME ?= $(HOME)/.chs_env
TESTS        ?= ./...

bin          := kafka-management-tool
xunit_output := test.xml
lint_output  := lint.txt

commit       := $(shell git rev-parse --short HEAD)
tag          := $(shell git tag -l 'v*-rc*' --points-at HEAD)
version      := $(shell if [[ -n "$(tag)" ]]; then echo $(tag) | sed 's/^v//'; else echo $(commit); fi)

.PHONY: all
all: clean build install

.PHONY: fmt
fmt:
	go fmt ./...

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


.PHONY: test-deps
test-deps: deps
	go get github.com/smartystreets/goconvey

.PHONY: xunit-tests
xunit-tests: test-deps
	go get github.com/tebeka/go2xunit
	go test -v $(TESTS) -run 'Unit' | go2xunit -output $(xunit_output)

.PHONY: lint
lint:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
	gometalinter ./... > $(lint_output); true