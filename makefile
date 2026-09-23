BACKEND := backend
FRONTEND := frontend
CHART := infra/helm/sys-called
DOCKER_COMPOSE := docker compose
IMAGE ?= sys-called:local
# Postgres the integration tests may write to (make db-test starts one).
TEST_DATABASE_URL ?= postgres://sys_called:sys_called@localhost:55432/sys_called?sslmode=disable

.DEFAULT_GOAL := help
.PHONY: help setup dev dev-backend dev-frontend build image run test test-backend test-integration \
	test-frontend lint lint-backend lint-frontend fmt tidy check env up down ps logs smoke \
	db-test db-test-down k8s-up k8s-deploy k8s-test k8s-status k8s-down helm-lint

help: ## Lista os comandos
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## --- Desenvolvimento ---------------------------------------------------------

setup: env ## Instala as dependências do frontend e do backend
	cd $(FRONTEND) && npm ci
	cd $(BACKEND) && go mod download

dev-backend: ## Sobe o backend com go run (precisa do Postgres: make db-test)
	cd $(BACKEND) && DATABASE_URL="$(TEST_DATABASE_URL)" go run ./cmd/sys-called

dev-frontend: ## Sobe o Vite (proxy de /api para http://localhost:8080)
	cd $(FRONTEND) && npm run dev

build: ## Compila o frontend e o binário Go em bin/
	cd $(FRONTEND) && npm run build
	cd $(BACKEND) && CGO_ENABLED=0 go build -trimpath -o ../bin/sys-called ./cmd/sys-called

image: ## Build da imagem Docker ($(IMAGE))
	docker build -t $(IMAGE) .

## --- Qualidade ---------------------------------------------------------------

test: test-backend test-frontend ## Testes unitários e de arquitetura (sem dependências externas)

test-backend: ## Testes Go com race detector
	cd $(BACKEND) && go test -race ./...

test-integration: ## Testes Go contra Postgres real (make db-test antes)
	cd $(BACKEND) && SYS_CALLED_TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test -race -count=1 -p 1 ./...

test-frontend: ## Testes do frontend
	cd $(FRONTEND) && npm test

lint: lint-backend lint-frontend ## Todos os linters

lint-backend: ## golangci-lint e go vet
	cd $(BACKEND) && go vet ./... && golangci-lint run ./...

lint-frontend: ## Type-check do frontend
	cd $(FRONTEND) && npm run type-check

fmt: ## Formata o código Go
	cd $(BACKEND) && golangci-lint fmt ./... || gofmt -w .

tidy: ## go mod tidy
	cd $(BACKEND) && go mod tidy

check: lint test helm-lint ## O que o CI roda antes dos testes de integração

## --- Docker Compose ----------------------------------------------------------

env: .env ## Cria o .env com uma chave JWT nova (não sobrescreve)

.env:
	@cp .env.example .env
	@{ printf 'JWT_PRIVATE_KEY="'; openssl genpkey -algorithm ed25519; printf '"\n'; } > .env.jwt
	@grep -v '^JWT_PRIVATE_KEY=' .env > .env.tmp && cat .env.tmp .env.jwt > .env && rm -f .env.tmp .env.jwt
	@echo "created .env"

up: .env ## Sobe a stack (http://localhost:8000)
	$(DOCKER_COMPOSE) up --build -d --wait

down: ## Derruba a stack (os dados ficam no volume)
	$(DOCKER_COMPOSE) down

ps: ## Estado e saúde dos containers
	$(DOCKER_COMPOSE) ps

logs: ## Acompanha os logs do app
	$(DOCKER_COMPOSE) logs -f app

smoke: ## Smoke test contra a stack no ar
	./scripts/smoke-test.sh http://localhost:$${APP_PORT:-8000}

db-test: ## Sobe um Postgres descartável na porta 55432
	docker run -d --rm --name sys-called-testdb -e POSTGRES_USER=sys_called -e POSTGRES_PASSWORD=sys_called \
		-e POSTGRES_DB=sys_called -p 127.0.0.1:55432:5432 postgres:16-alpine

db-test-down: ## Remove o Postgres descartável
	docker rm -f sys-called-testdb

## --- Kubernetes (kind + Helm) ------------------------------------------------

k8s-up: ## Cria o cluster kind, faz o build e instala o chart (http://localhost:30080)
	./infra/scripts/k8s.sh up

k8s-deploy: ## Novo build e helm upgrade
	./infra/scripts/k8s.sh deploy

k8s-test: ## helm test + smoke test no cluster
	./infra/scripts/k8s.sh test

k8s-status: ## Pods, services e volumes
	./infra/scripts/k8s.sh status

k8s-down: ## Apaga o cluster kind
	./infra/scripts/k8s.sh down

helm-lint: ## Lint e render do chart (dev e produção)
	helm lint $(CHART)
	helm lint $(CHART) -f $(CHART)/values-dev.yaml
	helm template sys-called $(CHART) > /dev/null
	helm template sys-called $(CHART) -f $(CHART)/values-dev.yaml > /dev/null
