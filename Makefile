default: oster testers

ENV = GOOS=linux GOARCH=amd64

.PHONY: oster
oster:
	mkdir -p ./bin
	$(ENV) go build -o ./bin/oster ./cmd/oster/...

.PHONY: printstack
printstack:
	mkdir -p ./bin
	$(ENV) go build -o ./bin/print-stack ./cmd/print_stack/...

.PHONY: testers
testers:
	mkdir -p ./bin
	$(ENV) go build -o ./bin/test-prog ./cmd/test-prog/main.go



clean:
	rm ./bin/*

help:
	@echo  "Targets:"
	@echo  "    oster (default) - build oster cli to ./bin/oster"
	@echo  "    testers - build test programs to run oster on"
	@echo  "    clean - clear out bin"
