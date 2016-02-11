# Path to go
GO ?= go

.PHONY: build
build:
	$(GO) build -v -i -o exago

# Compile the binary for all available OSes and ARCHes.
.PHONY: buildall
buildall:
	gox -output "dist/exago_{{.OS}}_{{.Arch}}"

# Run several automated source checks to get rid of the most simple issues.
# This helps keeping code review focused on application logic.
.PHONY: check
check:
	@echo "gometalinter"
	@! gometalinter --disable gotype,aligncheck,interfacer,structcheck --deadline 10s ./... | \
	  grep -vE 'vendor'

# "go test -i" builds dependencies and installs them into GOPATH/pkg,
# but does not run the tests.
.PHONY: test
test:
	$(GO) test

default: build