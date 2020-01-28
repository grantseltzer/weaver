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

default: oster dummies printstack

bin/oster: cmd/oster
> mkdir -p ./bin
> $(GOENV) go build -o ./bin/oster ./cmd/oster/...

bin/printstack: cmd/oster
> mkdir -p ./bin
> $(GOENV) go build -o ./bin/print-stack ./cmd/print_stack/...

bin/dummies: cmd/dummies
> mkdir -p ./bin
> $(GOENV) go build -o ./bin/dummies ./cmd/dummies/main.go

.PHONY: clean
clean:
> rm ./bin/*

.PHONY: clean
help:
> @echo  "Targets:"
> @echo  "    oster (default) - build oster cli to ./bin/oster"
> @echo  "    dummies - build dummy programs to run oster on"
> @echo  "    clean - clear out bin"
