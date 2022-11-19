#!/usr/bin/make
# Makefile readme (en): <https://www.gnu.org/software/make/manual/html_node/index.html#SEC_Contents>

SHELL = /bin/sh
LDFLAGS = "-s -w -X github.com/tarampampam/poke/internal/version.version=$(shell git rev-parse HEAD)"

DC_RUN_ARGS = --rm --user "$(shell id -u):$(shell id -g)"
APP_NAME = $(notdir $(CURDIR))

.DEFAULT_GOAL : help

# This will output the help for each task. thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Show this help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[32m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

image: ## Build docker image with the app
	docker build -f ./Dockerfile -t $(APP_NAME):local .
	docker images $(APP_NAME):local # print the info
	@printf "\n   \e[30;42m %s \033[0m\n\n" 'Now you can use image like `docker run --rm $(APP_NAME):local ...`'

generate: ## Generate assets
	docker-compose run $(DC_RUN_ARGS) --no-deps go go generate ./...

build: generate ## Build app binary file
	docker-compose run $(DC_RUN_ARGS) -e "CGO_ENABLED=0" --no-deps go go build -trimpath -ldflags $(LDFLAGS) ./cmd/poke/

test: ## Run tests
	docker-compose run $(DC_RUN_ARGS) --no-deps go go test -v -race -timeout 10s ./...

lint: ## Run code linter
	docker-compose run --rm golint golangci-lint run

fmt: ## Run source code formatting tools
	docker-compose run $(DC_RUN_ARGS) --no-deps go gofmt -s -w -d .
	docker-compose run $(DC_RUN_ARGS) --no-deps go goimports -d -w .
	docker-compose run $(DC_RUN_ARGS) --no-deps go go mod tidy

shell: ## Start shell inside go environment
	docker-compose run $(DC_RUN_ARGS) go sh

# Overall stuff

clean: ## Make clean
	docker-compose down -v -t 1
	-docker rmi $(APP_NAME):local -f
	-rm -R ./poke
