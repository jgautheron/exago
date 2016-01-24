# Path to go
GO ?= go

# The release number & build date are stamped into the binary.
.PHONY: build
build: LDFLAGS += -X "main.buildTag=$(shell git describe --tags)"
build: LDFLAGS += -X "main.buildDate=$(shell date -u '+%Y/%m/%d %H:%M:%S')"
build:
	$(GO) build -ldflags '$(LDFLAGS)' -v -i -o exago

# Compile the binary for all available OSes and ARCHes.
.PHONY: buildall
buildall: LDFLAGS += -X "main.buildTag=$(shell git describe --tags)"
buildall: LDFLAGS += -X "main.buildDate=$(shell date -u '+%Y/%m/%d %H:%M:%S')"
buildall:
	gox -ldflags '$(LDFLAGS)' -output "dist/exago_{{.OS}}_{{.Arch}}"

# Run several automated source checks to get rid of the most simple issues.
# This helps keeping code review focused on application logic.
.PHONY: check
check:
	@echo "gometalinter"
	@ gometalinter --deadline 10s ./...

# "go test -i" builds dependencies and installs them into GOPATH/pkg,
# but does not run the tests.
.PHONY: test
test:
	$(GO) test

default: build