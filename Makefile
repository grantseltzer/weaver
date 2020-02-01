GOENV = GOOS=linux GOARCH=amd64
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

default: bin/oster bin/print-stack bin/tester

bin/oster: cmd/oster
> mkdir -p ./bin
> $(GOENV) go build -o ./bin/oster ./cmd/oster/...
.PHONY: bin/oster

bin/print-stack: cmd/print-stack
> mkdir -p ./bin
> $(GOENV) go build -o ./bin/print-stack ./cmd/print-stack/...
.PHONY: bin/print-stack

bin/tester: cmd/tester
> mkdir -p ./bin
> $(GOENV) go build -o ./bin/tester ./cmd/tester/main.go
.PHONY: bin/tester

clean:
> rm ./bin/*
.PHONY: clean

.PHONY: clean
help:
> @echo  "Targets:"
> @echo  "    oster (default) - build oster cli to ./bin/oster"
> @echo  "    tester - build dummy programs to run oster on"
> @echo  "    print-stack - build print-stack cli which traces a particular function by printing the first 25 bytes the stack on function enter"
> @echo  "    clean - clear out bin"
