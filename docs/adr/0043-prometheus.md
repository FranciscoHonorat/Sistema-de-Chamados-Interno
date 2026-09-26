# ADR 0043 — Métricas com Prometheus

**Status:** aceita

## Contexto

Pra saber se o sistema está saudável é preciso de números: quantas requisições, quantos erros, quanto tempo cada rota leva.

## Decisão

- **`prometheus/client_golang`** na plataforma HTTP: `http_requests_total` e `http_request_duration_seconds`, por rota (o padrão da rota, não a URL com IDs), método e status.
- Métricas do runtime Go (memória, goroutines, GC) incluídas.
- **Porta separada** (`METRICS_PORT`, padrão 9090): `/metrics` não passa pelo Ingress nem fica exposto na internet. No compose, só em `127.0.0.1`.
- No compose, o perfil `observability` sobe um Prometheus já configurado (`infra/observability/prometheus.yml`); no Kubernetes, um `ServiceMonitor` opcional.

## Alternativas consideradas

- **OpenTelemetry**: padrão mais amplo (métricas, traces e logs), mas mais configuração e um coletor a mais.
- **APM comercial** (Datadog, New Relic): pronto, mas pago e com agente.

## Consequências

- Formato padrão, compatível com Grafana e alertas do Prometheus.
- Rotas com IDs não explodem a cardinalidade das métricas.
- Sem tracing distribuído, que num processo só faz pouca falta.
