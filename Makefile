GOENV = GOOS=linux GOARCH=amd64
MOD= -mod=vendor
SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

default: bin/weaver bin/print-stack bin/tester

bin/weaver: cmd/weaver
> mkdir -p ./bin
> $(GOENV) go build $(MOD) -o ./bin/weaver ./cmd/weaver/...
.PHONY: bin/weaver

bin/print-stack: cmd/print-stack
> mkdir -p ./bin
> $(GOENV) go build $(MOD) -o ./bin/print-stack ./cmd/print-stack/...
.PHONY: bin/print-stack

bin/tester: cmd/tester
> mkdir -p ./bin
> $(GOENV) go build $(MOD) -o ./bin/tester ./cmd/tester/...
.PHONY: bin/tester

# unit tests
test:
> go test -v ./cmd/weaver/...

smoke-test: tests/run_smoke_test.sh
> sh -c "tests/run_smoke_test.sh"

clean:
> rm -f ./bin/*
.PHONY: clean

.PHONY: clean
help:
> @echo  "Targets:"
> @echo  "    default - bin/*"
> @echo  "    test - (REQUIRES ROOT) run smoke tests"
> @echo  "    bin/weaver - build weaver cli to ./bin/weaver"
> @echo  "    bin/tester - build tester program used in test target to ./bin/tester"
> @echo  "    bin/print-stack - build print-stack cli which traces a particular function by printing the first 25 bytes the stack on function enter to ./bin/print-stack"
> @echo  "    clean - clear out bin"
