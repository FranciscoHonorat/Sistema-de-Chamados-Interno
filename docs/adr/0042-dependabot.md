# ADR 0042 — Dependabot pra manter as dependências em dia

**Status:** aceita

## Contexto

Dependências desatualizadas acumulam vulnerabilidades e tornam as atualizações futuras maiores e mais arriscadas.

## Decisão

`.github/dependabot.yml` com atualização **semanal** pra cinco ecossistemas:

- `gomod` (backend), `npm` (frontend), `docker` (imagens do Dockerfile), `docker-compose` (imagens do compose) e `github-actions` (versões das actions).

Cada atualização vira um pull request, que passa pelo CI completo antes de ser aceito.

## Alternativas consideradas

- **Renovate**: mais configurável (agrupamento, automerge), mas é mais uma configuração pra aprender.
- **Atualizar à mão**: fácil de esquecer.

## Consequências

- Atualizações pequenas e frequentes, cada uma testada.
- Alguns pull requests por semana pra revisar.
