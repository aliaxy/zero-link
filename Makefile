GO ?= go
DOCKER_COMPOSE ?= docker compose
INFRA_COMPOSE := deploy/docker-compose.infra.yml

.PHONY: infra-up infra-down migrate-up migrate-down run-api run-rpc test test-integration fmt install-hooks

infra-up:
	$(DOCKER_COMPOSE) -f $(INFRA_COMPOSE) up -d

infra-down:
	$(DOCKER_COMPOSE) -f $(INFRA_COMPOSE) down

migrate-up:
	@echo "No migrations exist in the skeleton stage."

migrate-down:
	@echo "No migrations exist in the skeleton stage."

run-api:
	@echo "link-api has not been generated yet. Generate it with goctl in a later stage."

run-rpc:
	@echo "link-rpc has not been generated yet. Generate it with goctl in a later stage."

test:
	@if find . -name '*.go' -not -path './.direnv/*' -not -path './.cache/*' | grep -q .; then \
		$(GO) test ./...; \
	else \
		echo "No Go packages yet; skeleton foundation is ready."; \
	fi

test-integration:
	@if find . -name '*.go' -not -path './.direnv/*' -not -path './.cache/*' | grep -q .; then \
		$(GO) test -tags=integration ./...; \
	else \
		echo "No Go integration packages yet; skeleton foundation is ready."; \
	fi

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
