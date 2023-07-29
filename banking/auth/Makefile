.PHONY: generate
generate:
	@echo "==> Generate: sqlc queries"
	@go install github.com/kyleconroy/sqlc/cmd/sqlc@latest && sqlc generate -f sqlc.yaml

.PHONY: build
build:
	@DOCKER_BUILDKIT=1 docker build --pull -t lucasgsousa/apps-banking-auth:v1 -f Dockerfile .

.PHONY: push
push:
	@docker push lucasgsousa/apps-banking-auth:v1
