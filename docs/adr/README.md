# Registros de decisões de arquitetura (ADRs)

Cada ADR registra uma decisão: o contexto, o que foi decidido, as alternativas consideradas e as consequências. Todos seguem o mesmo formato e estão com status **aceita**.

## Técnicas e arquitetura

| ADR | Decisão |
|---|---|
| [0001](0001-monolito-modular.md) | Monólito modular em vez de microsserviços |
| [0002](0002-arquitetura-hexagonal.md) | Arquitetura hexagonal dentro de cada módulo |
| [0003](0003-ddd-tatico.md) | DDD tático: agregados e value objects |
| [0004](0004-event-sourcing-chamados.md) | Event Sourcing no módulo de chamados |
| [0005](0005-cqrs-leve.md) | CQRS leve: comandos e consultas separados |
| [0006](0006-outbox-e-barramento-em-processo.md) | Transactional outbox e barramento de eventos em processo |
| [0007](0007-testes-de-arquitetura.md) | Fronteiras garantidas por testes de arquitetura |
| [0008](0008-tdd-e-estrategia-de-testes.md) | TDD e estratégia de testes em camadas |
| [0009](0009-autenticacao-jwt-e-refresh-token.md) | Autenticação com JWT Ed25519 e refresh token rotativo |
| [0010](0010-autorizacao-no-dominio.md) | Autorização como regra de domínio |
| [0011](0011-distribuicao-automatica.md) | Distribuição automática pelo atendente com menos chamados em aberto |
| [0012](0012-cache-de-chamados-em-memoria.md) | Cache de chamados em memória por instância |
| [0013](0013-frontend-servido-pelo-backend.md) | Frontend servido pelo próprio backend, na mesma origem |
| [0014](0014-migracoes-embutidas-no-binario.md) | Migrações SQL embutidas no binário, por schema |
| [0015](0015-configuracao-por-variaveis-de-ambiente.md) | Configuração por variáveis de ambiente, validada na subida |
| [0016](0016-observabilidade-logs-e-probes.md) | Logs estruturados, request id e probes de saúde |
| [0017](0017-seguranca-http.md) | Cabeçalhos de segurança e limites nas requisições HTTP |
| [0018](0018-notificacoes-por-consulta-periodica.md) | Notificações gravadas no banco e consultadas a cada 30 segundos |
| [0019](0019-estado-no-frontend-com-provide-inject.md) | Estado e dependências do frontend com provide/inject e composables |

## Tecnologias e ferramentas

| ADR | Decisão |
|---|---|
| [0020](0020-go-no-backend.md) | Go no backend |
| [0021](0021-gin.md) | Gin como framework HTTP |
| [0022](0022-postgresql-schema-por-modulo.md) | PostgreSQL, com um schema por módulo |
| [0023](0023-pgx.md) | pgx como driver do Postgres, com SQL escrito à mão |
| [0024](0024-vue3-typescript.md) | Vue 3 com TypeScript no frontend |
| [0025](0025-vite.md) | Vite como ferramenta de build e servidor de desenvolvimento |
| [0026](0026-vitest-e-testing-library.md) | Vitest e Testing Library nos testes do frontend |
| [0027](0027-playwright.md) | Playwright nos testes ponta a ponta |
| [0028](0028-tailwind-css.md) | Tailwind CSS 4 com tokens de design próprios |
| [0029](0029-bcrypt.md) | bcrypt pra guardar senhas |
| [0030](0030-golang-jwt.md) | golang-jwt pra emitir e validar tokens |
| [0031](0031-vue-router.md) | Vue Router com guarda por perfil |
| [0032](0032-testify.md) | testify nas asserções dos testes em Go |
| [0033](0033-google-uuid.md) | UUIDs como identificadores |
| [0034](0034-docker-imagem-distroless.md) | Imagem Docker multi-stage com base distroless |
| [0035](0035-docker-compose.md) | Docker Compose pra rodar o projeto localmente |
| [0036](0036-kubernetes-e-helm.md) | Kubernetes com um chart Helm próprio |
| [0037](0037-kind.md) | kind pra testar o chart num cluster local |
| [0038](0038-github-actions-ci-cd.md) | GitHub Actions pra integração e entrega contínuas |
| [0039](0039-assinatura-e-proveniencia-da-imagem.md) | Imagem assinada com cosign, com SBOM e proveniência |
| [0040](0040-varredura-de-vulnerabilidades.md) | Varredura de vulnerabilidades: Trivy, govulncheck e npm audit |
| [0041](0041-codeql.md) | CodeQL pra análise estática de segurança |
| [0042](0042-dependabot.md) | Dependabot pra manter as dependências em dia |
| [0043](0043-prometheus.md) | Métricas com Prometheus |
| [0044](0044-golangci-lint.md) | golangci-lint no backend |
| [0045](0045-eslint-e-vue-tsc.md) | ESLint e vue-tsc no frontend |
| [0046](0046-aws-ec2-e-nginx.md) | Demonstração publicada numa EC2 com Nginx na frente |
| [0047](0047-makefile.md) | Makefile como ponto de entrada dos comandos |

## Formato

```markdown
# ADR NNNN — Título

**Status:** proposta | aceita | substituída por ADR XXXX

## Contexto
## Decisão
## Alternativas consideradas
## Consequências
```

Uma decisão nova ganha o próximo número. Uma decisão que muda não é editada: um ADR novo a substitui, e o antigo passa a apontar pra ele.
