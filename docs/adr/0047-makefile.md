# ADR 0047 — Makefile como ponto de entrada dos comandos

**Status:** aceita

## Contexto

O projeto tem backend, frontend, Docker, Kubernetes e testes em vários níveis, cada um com seus comandos. Quem chega ao projeto não deveria precisar decorar todos.

## Decisão

- Um `makefile` na raiz com os comandos do dia a dia, cada um com uma descrição que o `make help` lista:
  - `make env`, `make up`, `make down`, `make ps`, `make logs`, `make smoke`
  - `make setup`, `make test`, `make test-integration`, `make test-e2e`, `make lint`, `make check`
  - `make dev-backend`, `make dev-frontend`, `make db-test`
  - `make image`, `make helm-lint`, `make k8s-up`, `make k8s-test`, `make k8s-deploy`, `make k8s-down`
- O README usa os mesmos comandos, e `make check` roda localmente o que o CI roda antes da integração.

## Alternativas consideradas

- **Scripts npm**: ficam presos ao frontend.
- **Task / just**: sintaxe mais amigável, mas exigem instalar mais uma ferramenta; o `make` já vem na maioria dos sistemas.

## Consequências

- `make help` é a documentação viva dos comandos.
- Sintaxe do make (tabs, variáveis) é pouco amigável pra quem não conhece.
