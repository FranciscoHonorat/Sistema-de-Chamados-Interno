# ADR 0012 — Cache de chamados em memória por instância

**Status:** aceita

## Contexto

Com Event Sourcing (ADR 0004), cada leitura reconstrói o chamado a partir dos seus eventos. A lista, a distribuição automática e a carga dos atendentes precisam de todos os chamados.

## Decisão

- `InMemoryTicketCache` (`adapters/out/cache/memory.go`): um mapa protegido por `sync.RWMutex`, atualizado a cada gravação.
- `LoadAllTickets` usa o cache e busca os que faltam numa **única consulta** (`LoadMany`, `aggregate_id = ANY(...)`).
- Ligado por padrão (`TICKETS_CACHE_ENABLED=true`). Com mais de uma réplica, cada instância teria um cache diferente, então o chart Helm troca pelo `NoopTicketCache` automaticamente.

## Alternativas consideradas

- **Redis**: cache compartilhado entre réplicas, mas mais um serviço pra operar.
- **Projeção de leitura no Postgres**: resolve o custo e leva os filtros pro banco. É a evolução documentada nos trade-offs do README.
- **Sem cache**: correto, mas reconstrói tudo a cada lista.

## Consequências

- Leituras rápidas com uma réplica, que é o cenário do compose.
- Com várias réplicas, o custo de reconstruir volta, mas continua em duas idas ao banco.
- O cache ocupa memória proporcional ao número de chamados.
