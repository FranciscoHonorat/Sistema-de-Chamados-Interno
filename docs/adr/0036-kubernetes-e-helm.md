# ADR 0036 — Kubernetes com um chart Helm próprio

**Status:** aceita

## Contexto

Além de rodar localmente, o produto precisa de um caminho pra produção com várias réplicas, atualização sem indisponibilidade e configuração por ambiente.

## Decisão

Um chart Helm em `infra/helm/sys-called`:

- **Deployment** com init container `migrate`, probes de *startup*, *readiness* e *liveness* separadas, `preStop` pra drenar, rolling update com `maxUnavailable: 0` e `topologySpreadConstraints`.
- **Pods endurecidos**: sem root, sem capabilities, `seccomp: RuntimeDefault`, raiz só leitura, sem token de service account.
- **HPA** opcional e **PodDisruptionBudget** quando pode haver mais de uma réplica.
- **Postgres** em StatefulSet, ou banco gerenciado via `externalDatabase.existingSecret`.
- **Secret** com senha do banco e chave Ed25519 gerados na primeira instalação e reaproveitados nos upgrades (`lookup`).
- **NetworkPolicy**: só o app fala com o Postgres.
- **Ingress** com TLS opcional e **ServiceMonitor** opcional (Prometheus Operator).
- `helm test` checa o `/readyz`.
- Valores por ambiente: `values.yaml` (padrões de produção), `values-dev.yaml`, `values-staging.yaml`, `values-production.yaml`.
- O cache de chamados é desligado automaticamente com mais de uma réplica (ADR 0012).

## Alternativas consideradas

- **Kustomize**: sem templates, mas geração de segredos e lógica condicional ficam mais difíceis.
- **Só Docker Compose em produção**: simples, mas sem réplicas, autoescala nem rolling update. É o que roda na EC2 de demonstração (ADR 0046).
- **PaaS** (Render, Fly.io): menos trabalho, mas prende a um fornecedor.

## Consequências

- Deploy reproduzível em qualquer cluster.
- Kubernetes é mais complexo do que um sistema desse tamanho precisa hoje: é preparação pra crescer, não requisito.
