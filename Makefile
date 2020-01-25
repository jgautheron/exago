container ?= api

env_file ?= ./dev.env

serve:
	sh -ac ". $(env_file) && go run ./cmd/${container}"
.PHONY: serve

build-container:
	docker build -t exago-$(container) --file "./Dockerfile.$(container)" .
.PHONY: build-container

run-container:
	docker run --rm -it -p 8080:8080 --env-file $(env_file) exago-$(container)
.PHONY: run-container