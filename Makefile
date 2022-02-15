PKG_PATH = "github.com/containers/podman-tui"
TARGET = podman-tui
BIN = ./bin
DESTDIR = /usr/local/bin
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
SELINUXOPT ?= $(shell test -x /usr/sbin/selinuxenabled && selinuxenabled && echo -Z)

.PHONY: default
default: all

.PHONY: all
all: binary

.PHONY: binary
binary: $(TARGET)  ## Build podman-tui binary
	@true

.PHONY: $(TARGET)
$(TARGET): $(SRC)
	@mkdir -p $(BIN)
	@echo "running go build"
	@go build -o $(BIN)/$(TARGET)

.PHONY: clean
clean:
	@rm -rf $(BIN)

.PHONY: install   
install:    ## Install podman-tui binary
	install ${SELINUXOPT} -D -m0755 $(BIN)/$(TARGET) $(DESTDIR)/$(TARGET)

.PHONY: uninstall 
uninstall:  ## Uninstall podman-tui binary
	rm -f $(DESTDIR)/$(TARGET)

.PHONY: validate  
validate:   ## Validate podman-tui code (fmt, lint, ...)
	@echo "running gofmt"
	@test -z "$(shell gofmt -l $(SRC) | tee /dev/stderr | tr '\n' ' ')" || echo "[WARN] Fix formatting issues with 'make fmt'"

	@echo "running golint"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done

	@echo "running go vet"
	@go vet ../$(TARGET)

.PHONY: test
test: functionality

.PHONY: functionality
functionality:
	bats test/

.PHONY: fmt      
fmt:   ## Run gofmt
	@echo -e "gofmt check and fix"
	
	@gofmt -w $(SRC)

_HLP_TGTS_RX = '^[[:print:]]+:.*?\#\# .*$$'
_HLP_TGTS_CMD = grep -E $(_HLP_TGTS_RX) $(MAKEFILE_LIST)
_HLP_TGTS_LEN = $(shell $(_HLP_TGTS_CMD) | cut -d : -f 1 | wc -L)
_HLPFMT = "%-$(_HLP_TGTS_LEN)s %s\n"
.PHONY: help   
help: ## Print listing of key targets with their descriptions
	@printf $(_HLPFMT) "Target:" "Description:"
	@printf $(_HLPFMT) "--------------" "--------------------"
	@$(_HLP_TGTS_CMD) | sort | \
		awk 'BEGIN {FS = ":(.*)?## "}; \
			{printf $(_HLPFMT), $$1, $$2}'
