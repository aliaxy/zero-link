GO ?= go
GO_ENV := GOCACHE=$(CURDIR)/.cache/go-build
DOCKER_COMPOSE ?= docker compose
INFRA_COMPOSE := deploy/docker-compose.infra.yml

.PHONY: infra-up infra-down migrate-up migrate-down run-api run-rpc test test-integration lint fmt install-hooks

infra-up:
	$(DOCKER_COMPOSE) -f $(INFRA_COMPOSE) up -d

infra-down:
	$(DOCKER_COMPOSE) -f $(INFRA_COMPOSE) down

migrate-up:
	@echo "No migrations exist in the skeleton stage."

migrate-down:
	@echo "No migrations exist in the skeleton stage."

run-api:
	$(GO_ENV) $(GO) run ./services/link-api -f etc/link-api-local.yaml

run-rpc:
	$(GO_ENV) $(GO) run ./services/link-rpc -f etc/link-rpc-local.yaml

test:
	@if find . -name '*.go' -not -path './.direnv/*' -not -path './.cache/*' | grep -q .; then \
		$(GO_ENV) $(GO) test ./...; \
	else \
		echo "No Go packages yet; skeleton foundation is ready."; \
	fi

test-integration:
	@if find . -name '*.go' -not -path './.direnv/*' -not -path './.cache/*' | grep -q .; then \
		$(GO_ENV) $(GO) test -tags=integration ./...; \
	else \
		echo "No Go integration packages yet; skeleton foundation is ready."; \
	fi

lint:
	GOCACHE=$(CURDIR)/.cache/go-build GOLANGCI_LINT_CACHE=$(CURDIR)/.cache/golangci-lint golangci-lint run ./...

fmt:
	@if find . -name '*.go' -not -path './.direnv/*' -not -path './.cache/*' | grep -q .; then \
		gofmt -w $$(find . -name '*.go' -not -path './.direnv/*' -not -path './.cache/*'); \
	else \
		echo "No Go files to format."; \
	fi

install-hooks:
	chmod +x .githooks/commit-msg
	@if git config core.hooksPath .githooks; then \
		echo "Git hooks installed from .githooks/"; \
	else \
		echo "Could not update Git config. Run manually: git config core.hooksPath .githooks"; \
		exit 1; \
	fi
