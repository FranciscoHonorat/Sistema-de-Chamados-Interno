# ADR 0035 — Docker Compose pra rodar o projeto localmente

**Status:** aceita

## Contexto

Seria uma boa prática qualquer pessoa do time consiga rodar a aplicação localmente. Pedir Go, Node e Postgres instalados na versão certa é atrito.

## Decisão

- `docker-compose.yml` com três serviços, cada um só depois do anterior ficar saudável: `postgres` → `migrate` (aplica as migrações e termina) → `app`.
- `make env` gera o `.env` e `make up` sobe tudo (`docker compose up --build -d --wait`). O app fica em `http://localhost:8000`, com usuários de desenvolvimento prontos.
- **Endurecido como em produção**: `read_only`, `cap_drop: [ALL]`, `no-new-privileges`, limite de memória, healthcheck e rotação de logs.
- O `make env` detecta hosts onde `no-new-privileges` quebra o `exec` (Docker 29 com kernel 7) e desliga só essa opção no `.env`.
- Perfil opcional `observability` com um Prometheus já coletando as métricas.

## Alternativas consideradas

- **Instalar tudo na máquina**: sem Docker, mas com versões variando de pessoa pra pessoa.
- **Dev containers**: bons pra desenvolver, mas exigem VS Code ou ferramenta compatível só pra rodar.

## Consequências

- Só Docker é pré-requisito pra rodar tudo.
- O compose local se comporta como produção (mesma imagem, mesmas restrições).
