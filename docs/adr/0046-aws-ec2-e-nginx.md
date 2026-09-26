# ADR 0046 — Demonstração publicada numa EC2 com Nginx na frente

**Status:** aceita

## Contexto

Ter a aplicação acessível por uma URL pública facilita a avaliação e o portfólio. O custo tem que ser zero ou quase.

## Decisão

- Uma instância **Amazon EC2 Free Tier** com Ubuntu, rodando a mesma stack do compose local (`make env` + `docker compose up`).
- **Nginx** na própria máquina como proxy reverso na porta 80, encaminhando pra `127.0.0.1:8000`.
- O app e o Postgres ficam presos ao loopback: nada além do Nginx é acessível de fora.

## Alternativas consideradas

- **Kubernetes gerenciado** (EKS): o chart está pronto (ADR 0036), mas o custo não cabe numa demonstração.
- **PaaS** (Render, Railway, Fly.io): menos trabalho, mas planos gratuitos com limites e hibernação.

## Consequências

- Demonstração pública com a mesma stack que roda localmente.
- Não é alta disponibilidade: uma instância, sem backup externo nem failover.
- Sem HTTPS na configuração descrita; em produção, TLS no Nginx (ex.: Let's Encrypt) ou no Ingress.
