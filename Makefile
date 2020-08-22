PROJECTNAME ?= $(shell basename "$(PWD)")
BASE         = $(shell pwd)
BUILD_DIR   ?= $(BASE)/build
VETARGS     ?= -all
GOFMT_FILES ?= $$(find . -name '*.go' | grep -v vendor)
TAG         ?= latest

# Ensure GOPATH is set
GOPATH            ?= $(shell go env GOPATH)

GOBIN             ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO111MODULE       ?= on
export GO111MODULE
GOPROXY           ?= https://proxy.golang.org
export GOPROXY

.PHONY: help
help: ## Display usage and help message.
help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-z0-9A-Z_-]+:.*?##/ { printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: run
run: ## Runs the program.
run:
	@echo ">> starting $(PROJECTNAME)"
	@go run .


.PHONY: build
build: ## Builds the binary.
build:
	@echo ">> building $(PROJECTNAME) binary"
	@go build -o $(BUILD_DIR)/$(PROJECTNAME) .


.PHONY: docker
docker: ## Builds the docker image.
docker:
	@echo ">> building $(PROJECTNAME) docker image"
	@docker build -t $(PROJECTNAME):$(TAG) .


.PHONY: docker-push
docker-push: ## Builds the docker image and publish it.
docker-push:
	@echo ">> building $(PROJECTNAME) docker image"
	@docker build -t prmsrswt/$(PROJECTNAME):$(TAG) .
	@docker push prmsrswt/$(PROJECTNAME):$(TAG)


.PHONY: vet
vet: ## Runs go vet against all packages.
vet:
	@echo ">> running go vet on packages"
	@go vet $(VETARGS) ./pkg/... . ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi


.PHONY: fmt
fmt: ## Format all go files using go fmt.
fmt:
	@echo ">> running go fmt on all go files"
	@gofmt -w $(GOFMT_FILES)


.PHONY: test
test: ## Run all unit tests
test:
	@echo ">> running unit tests"
	@go test ./...


.PHONY: coverage
coverage: ## Generate and open a HTML test coverage report
coverage:
	@echo ">> generating coverage report for all tests"
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@xdg-open coverage.html


.PHONY: clean
clean: ## Deletes temporary files created by this Makefile's targets
clean:
	@echo ">> deleting files made by target 'coverage'"
	@rm coverage.out || true
	@rm coverage.html || true
	@echo ">> deleting binaries"
	@rm build/* || true
