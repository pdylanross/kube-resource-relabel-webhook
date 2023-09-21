ifneq (,$(wildcard ./.env))
    include .env
    export
endif


help: # Show this help.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: default
default: help

.PHONY: build
build:
	docker build -t relabel:dev .

.PHONY: fix-deps
fix-deps: # Run dependency maintenance commands
	go mod tidy
	go vet ./...

.PHONY: fix-check
fix-check: fix-deps # Run linters and autofix issues
	golangci-lint run --fix

.PHONY: check
check: # Run linters, no autofix
	golangci-lint run

.PHONY: unit-test
unit-test: # Run unit tests
	go clean -testcache
	go test ./...

.PHONY: integration-test
integration-test: build # Run integration tests
	go clean -testcache
	go test ./... --tags=integration_tests -v -p=1

.PHONY: pre-commit
pre-commit: fix-check unit-test # Run all standard cleanups before a commit

.PHONY: doc-setup
doc-setup: # Fetch required deps for doc
	cd doc && $(MAKE) setup

.PHONY: doc-gen
doc-gen: # Codegen doc
	cd doc && $(MAKE) gen

.PHONY: doc-serve
doc-serve: # Serve doc pages locally
	cd doc && $(MAKE) serve

.PHONY: local-create-cluster
local-create-cluster:
	kind create cluster --config ./integration/kind-config.yaml --name test1

.PHONY: local-load-image build
local-load-image: build
	kind load docker-image --name test1 relabel:dev

.PHONY: local-install-relabel
local-install-relabel: local-install-cert-manager local-load-image
	helm upgrade -i relabel ./chart \
		--set image.tag=dev \
		--set image.repository=relabel \
		--set image.pullPolicy=Never \
		--set fullnameOverride=kube-resource-relabel-webhook \
		--atomic \
		-n default \
		-f ./integration/common-tests/values.yaml \
		-f ./integration/test-cases/cert-manager/relabel-values.yaml

.PHONY: local-install-cert-manager
local-install-cert-manager:
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
	kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=cert-manager -n cert-manager

.PHONY: local-setup
local-setup: local-create-cluster local-install-cert-manager local-install-relabel # setup local dev cluster

.PHONY: local-teardown
local-teardown: # teardown local dev cluster
	kind delete cluster --name test1