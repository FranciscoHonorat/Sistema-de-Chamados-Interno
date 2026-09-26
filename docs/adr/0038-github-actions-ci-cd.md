# ADR 0038 — GitHub Actions pra integração e entrega contínuas

**Status:** aceita

## Contexto

O código está no GitHub. Cada mudança precisa ser verificada (lint, tipos, testes, imagem) antes de entrar, e a publicação precisa ser repetível, sem passos manuais.

## Decisão

**CI** (`.github/workflows/ci.yml`), em todo push e pull request, com jobs em paralelo:

- **Backend**: `go mod tidy` sem diferença, golangci-lint (ADR 0044), testes com `-race` e cobertura.
- **Backend · integração**: testes contra um Postgres de serviço, com `-p 1`.
- **Backend · vulnerabilidades**: govulncheck (ADR 0040).
- **Frontend**: ESLint, `vue-tsc`, Vitest com cobertura, build e `npm audit` (ADR 0040 e 0045).
- **E2E**: sobe Postgres e o binário servindo o `dist/` e roda o Playwright (ADR 0027).
- **Imagem**: build, teste de que o binário sobe e varredura com Trivy.
- Pull requests cancelam execuções antigas; `main` e tags sempre terminam.

**CD** (`release.yml` + `deploy.yml`), chamado pelo CI quando tudo passa:

- `main` → imagem `edge` e `sha-<commit>` no GHCR → deploy em **staging**.
- Tag `vX.Y.Z` → imagem versionada, chart Helm publicado como OCI, release no GitHub → deploy em **produção** após aprovação manual no environment.
- O deploy é **por digest**, verifica a assinatura antes (ADR 0039), usa `helm upgrade --atomic`, roda `helm test` e o smoke test, e faz **rollback automático** se algo falhar.
- Os deploys só rodam com `DEPLOY_ENABLED=true`; sem isso, o pipeline publica os artefatos e pula os deploys.

## Alternativas consideradas

- **GitLab CI, CircleCI, Jenkins**: equivalentes, mas exigiriam integração extra com o GitHub ou um servidor próprio.
- **Deploy manual**: sujeito a esquecimento de passos.

## Consequências

- Nada chega em `main` sem passar pelos mesmos testes que rodam localmente.
- Publicar uma versão é criar uma tag.
- Os ambientes `staging` e `production` precisam ser configurados no GitHub pra ligar os deploys.
