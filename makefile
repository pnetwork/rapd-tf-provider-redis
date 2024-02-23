LOCALBIN ?= $(shell pwd)/bin

GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= 1.54.2

GOLANGCI_LINT_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh"

.PHONY: make-bin
make-bin:
	mkdir -p $(LOCALBIN)

.PHONY: golangci-lint
golangci-lint: make-bin $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(LOCALBIN)
	@if test -x $(LOCALBIN)/golangci-lint && ! $(LOCALBIN)/golangci-lint version | grep -q $(GOLANGCI_LINT_VERSION); then \
		echo "$(LOCALBIN)/golangci-lint version is not expected $(GOLANGCI_LINT_VERSION). Removing it before installing."; \
		rm -rf $(LOCALBIN)/golangci-lint; \
	fi
	test -s $(LOCALBIN)/golangci-lint || curl -sSfL $(GOLANGCI_LINT_INSTALL_SCRIPT) | sh -s -- -b $(LOCALBIN) v$(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint: golangci-lint ## Run aggregate linter against code.
	GOFLAGS=-buildvcs=false $(GOLANGCI_LINT) --config .golangci.yml run ./...

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
