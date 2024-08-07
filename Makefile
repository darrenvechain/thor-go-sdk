PACKAGE = github.com/vechain/thor-go-sdk
MAJOR = $(shell go version | cut -d' ' -f3 | cut -b 3- | cut -d. -f1)
MINOR = $(shell go version | cut -d' ' -f3 | cut -b 3- | cut -d. -f2)
PACKAGES = `go list ./... | grep -v '/vendor/'`

help:
	@egrep -h '\s#@\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?#@ "}; {printf "\033[36m  %-30s\033[0m %s\n", $$1, $$2}'

go_version_check:
	@if test $(MAJOR) -lt 1; then \
		echo "Go 1.22 or higher required"; \
		exit 1; \
	else \
		if test $(MAJOR) -eq 1 -a $(MINOR) -lt 19; then \
			echo "Go 1.22 or higher required"; \
			exit 1; \
		fi \
	fi

thor_solo:| #@ Start a thor solo node if it's not running
	@docker compose up -d --wait

test:| go_version_check thor_solo #@ Run the tests
	@go test -cover $(PACKAGES)

test-coverage:| go_version_check thor_solo #@ Run the tests with coverage
	@go test -race -coverprofile=coverage.out -covermode=atomic $(PACKAGES)
	@go tool cover -html=coverage.out

lint_command_check:
	@command -v golangci-lint || (echo "golangci-lint not found, please install it from https://golangci-lint.run/usage/install/" && exit 1)

lint: | go_version_check lint_command_check #@ Run 'golangci-lint'
	@golangci-lint run --config .golangci.yml

lint-fix: | go_version_check lint_command_check #@ Run 'golangci-lint' with fix
	@golangci-lint run --config .golangci.yml --fix
