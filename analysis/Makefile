container = api

env_file ?= ./dev.env

serve:
	sh -ac ". $(env_file) && go run ./"
.PHONY: serve

build-container:
	cd ../../ && docker build -t $(container) --file "./src/$(container)/Dockerfile" .
.PHONY: build-container

run-container:
	docker run --rm -it -p 8080:8080 --env-file $(env_file) $(container)
.PHONY: run-container