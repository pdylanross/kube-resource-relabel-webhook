
help: # Show this help.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: default
default: help

.PHONY: build
build: # Build a container image
	./ci/build.sh

.PHONY: publish
publish: build # publish the current image to the container registry
	./ci/publish.sh

.PHONY: fix-deps
fix-deps: # Run dependency maintenance commands
	go mod tidy
	go vet ./...

.PHONY: fix-check
fix-check: fix-deps # Run linters
	golangci-lint run --fix

.PHONY: check
check:
	golangci-lint run

.PHONY: unit-test
unit-test:
	go test ./...
