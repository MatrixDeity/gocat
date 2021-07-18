.DEFAULT_GOAL = help
SHELL = bash

MAINTAINER_NAME = matrixdeity
PROJECT_NAME = gocat

.PHONY: help
help: $(MAKEFILE_LIST) ## show this help (default goal).
	@echo -e "Awailable \033[33mmake\033[0m commands:\n"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## .*$$/ {printf "  \033[33m%-8s\033[0m - %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
	@echo -e "\nMatrixDety, 2021."

.PHONY: image
image: ## build a docker image of GoCat service.
	@docker build -t $(MAINTAINER_NAME)/$(PROJECT_NAME):latest .

.PHONY: build
build: ## build an executable file of GoCat service.
	@go build -o ./$(PROJECT_NAME) ./cmd

.PHONY: start
start: ## build and run an executable file of GoCat service.
	@test -n "$(GOCAT_TELEGRAM_TOKEN)" || (echo "Set GOCAT_TELEGRAM_TOKEN env" >&2 && exit 1)
	@go run ./cmd -token="$(GOCAT_TELEGRAM_TOKEN)"

.PHONY: clean
clean: ## clean a working directory.
	@rm -rf ./$(PROJECT_NAME)

.PHONY: push
push: ## push builded docker image.
	@docker push $(MAINTAINER_NAME)/$(PROJECT_NAME)
