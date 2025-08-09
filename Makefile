PKG_PATH = "github.com/containers/podman-tui"
GO := go
FIRST_GOPATH := $(firstword $(subst :, ,$(GOPATH)))
GOPKGDIR := $(FIRST_GOPATH)/src/$(PKG_PATH)
GOPKGBASEDIR ?= $(shell dirname "$(GOPKGDIR)")
GOBIN := $(shell $(GO) env GOBIN)
BUILDFLAGS := -mod=vendor $(BUILDFLAGS)
BUILDTAGS := "exclude_graphdriver_btrfs containers_image_openpgp remote"
COVERAGE_PATH ?= .coverage
TARGET = podman-tui
BIN = ./bin
DESTDIR = /usr/bin
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
SELINUXOPT ?= $(shell test -x /usr/sbin/selinuxenabled && selinuxenabled && echo -Z)
PKG_MANAGER ?= $(shell command -v dnf yum|head -n1)
PRE_COMMIT = $(shell command -v bin/venv/bin/pre-commit ~/.local/bin/pre-commit pre-commit | head -n1)
GINKO_CLI_VERSION = $(shell grep 'ginkgo/v2' go.mod | grep -o ' v.*' | sed 's/ //g')

# Default to the native OS type and architecture unless otherwise specified
NATIVE_GOOS := $(shell env -u GOOS $(GO) env GOOS)
GOOS ?= $(call err_if_empty,NATIVE_GOOS)
# Default to the native architecture type
NATIVE_GOARCH := $(shell env -u GOARCH $(GO) env GOARCH)
GOARCH ?= $(NATIVE_GOARCH)

.PHONY: default
default: all

.PHONY: all
all: binary binary-win binary-darwin

.PHONY: binary
binary: $(TARGET)  ## Build podman-tui binary
	@true

.PHONY: $(TARGET)
$(TARGET): $(SRC)
	@mkdir -p $(BIN)
	@echo "running go build"
	@env CGO_ENABLED=0 $(GO) build $(BUILDFLAGS) -o $(BIN)/$(TARGET) -tags $(BUILDTAGS)

.PHONY: binary-win
binary-win:  ## Build podman-tui.exe windows binary
	@mkdir -p $(BIN)/windows/
	@echo "running go build for windows"
	@env CGO_ENABLED=0 GOOS=windows GOARCH=$(GOARCH) $(GO) build $(BUILDFLAGS) -o $(BIN)/windows/$(TARGET).exe -tags $(BUILDTAGS)

.PHONY: binary-darwin
binary-darwin: ## Build podman-tui for darwin
	@mkdir -p $(BIN)/darwin/
	@echo "running go build for darwin"
	@env CGO_ENABLED=0 GOOS=darwin GOARCH=$(GOARCH) $(GO) build $(BUILDFLAGS) -o $(BIN)/darwin/$(TARGET) -tags $(BUILDTAGS)

.PHONY: clean
clean:
	@rm -rf $(BIN)

.PHONY: install
install:    ## Install podman-tui binary
	@install ${SELINUXOPT} -D -m0755 $(BIN)/$(TARGET) $(DESTDIR)/$(TARGET)

.PHONY: uninstall
uninstall:  ## Uninstall podman-tui binary
	@rm -f $(DESTDIR)/$(TARGET)

#=================================================
# Required tools installation tartgets
#=================================================

.PHONY: install.tools
install.tools: .install.ginkgo .install.bats .install.pre-commit .install.codespell .install.golangci-lint ## Install needed tools

.PHONY: .install.ginkgo
.install.ginkgo:
	if [ ! -x "$(GOBIN)/ginkgo" ]; then \
		$(GO) install -mod=mod github.com/onsi/ginkgo/v2/ginkgo@$(GINKO_CLI_VERSION) ; \
	fi

.PHONY: .install.bats
.install.bats:
	sudo ${PKG_MANAGER} -y install bats

.PHONY: .install.pre-commit
.install.pre-commit:
	if [ -z "$(PRE_COMMIT)" ]; then \
		python3 -m pip install --user pre-commit; \
	fi

.PHONY: .install.golangci-lint
.install.golangci-lint:
	VERSION=2.3.1 ./hack/install_golangci.sh

.PHONY: .install.codespell
.install.codespell:
	sudo ${PKG_MANAGER} -y install codespell

#=================================================
# Testing (units, functionality, ...) targets
#=================================================

.PHONY: test
test: test-unit test-functionality

.PHONY: test-unit
test-unit: ## Run unit tests
	rm -rf ${COVERAGE_PATH} && mkdir -p ${COVERAGE_PATH}
	$(GOBIN)/ginkgo \
		-r \
		--skip-package test/ \
		--cover \
		--covermode atomic \
		--coverprofile coverprofile \
		--output-dir ${COVERAGE_PATH} \
		--succinct
	$(GO) tool cover -html=${COVERAGE_PATH}/coverprofile -o ${COVERAGE_PATH}/coverage.html
	$(GO) tool cover -func=${COVERAGE_PATH}/coverprofile > ${COVERAGE_PATH}/functions
	cat ${COVERAGE_PATH}/functions | sed -n 's/\(total:\).*\([0-9][0-9].[0-9]\)/\1 \2/p'

.PHONY: test-functionality
test-functionality: ## Run functionality tests
	@bats test/

.PHONY: package
package:  ## Build rpm package
	rpkg local

.PHONY: package-install
package-install: package  ## Install rpm package
	sudo ${PKG_MANAGER} -y install ${HOME}/rpmbuild/RPMS/*/*.rpm
	/usr/bin/podman-tui version

#=================================================
# Linting/Formatting/Code Validation targets
#=================================================

.PHONY: validate
validate: gofmt lint pre-commit  ## Validate podman-tui code (fmt, lint, ...)

.PHONY: vendor
vendor: ## Check vendor
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify
	@bash ./hack/tree_status.sh

.PHONY: lint
lint: ## Run golangci-lint
	@echo "running golangci-lint"
	$(BIN)/golangci-lint version
	$(BIN)/golangci-lint run

.PHONY: pre-commit
pre-commit:   ## Run pre-commit
ifeq ($(PRE_COMMIT),)
	@echo "FATAL: pre-commit was not found, make .install.pre-commit to installing it." >&2
	@exit 2
endif
	$(PRE_COMMIT) run -a

.PHONY: gofmt
gofmt:   ## Run gofmt
	@echo -e "gofmt check and fix"
	@gofmt -w $(SRC)

.PHONY: codespell
codespell: ## Run codespell
	@echo "running codespell"
	@codespell -S ./vendor,go.mod,go.sum,./.git,*_test.go

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
