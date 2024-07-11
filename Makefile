EXTENSION:=
ifeq ($(OS),Windows_NT)
  EXTENSION:=.exe
endif

BINARY?=docker
BUILDER=buildx-multi-arch

STATIC_FLAGS=CGO_ENABLED=0
LDFLAGS="-s -w"
GO_BUILD=$(STATIC_FLAGS) go build -trimpath -ldflags=$(LDFLAGS)

bin: ## Build the binary for the current plarform
	@echo "Building..."
	$(GO_BUILD) -o "bin/$(BINARY)$(EXTENSION)" .
help: ## Show this help
	@echo Please specify a build target. The choices are:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.PHONY: bin help