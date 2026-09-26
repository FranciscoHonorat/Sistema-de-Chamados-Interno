# ADR 0022 — PostgreSQL, com um schema por módulo

**Status:** aceita

## Contexto

O sistema precisa guardar funcionários, sessões, eventos de chamados, notificações e o outbox, com transações (outbox na mesma transação do funcionário), JSON (payload dos eventos) e controle de concorrência.

## Decisão

- **PostgreSQL 16**, um banco só.
- **Um schema por módulo** (`employees`, `tickets`). Cada módulo abre seu pool com `search_path` no próprio schema e não enxerga as tabelas do outro.
- Recursos usados: transações, `JSONB` (payload dos eventos), `UNIQUE (aggregate_id, version)` (concorrência otimista), `FOR UPDATE SKIP LOCKED` (outbox com várias réplicas), advisory locks (migrações), `= ANY($1)` (carga em lote).

## Alternativas consideradas

- **MySQL**: comum com Laravel, mas `SKIP LOCKED` e JSON são menos maduros.
- **SQLite**: zero configuração, mas não serve pra várias réplicas.
- **Um banco por módulo**: isolamento maior, mas mais um banco pra operar.
- **EventStoreDB** pros eventos: especializado, mas mais uma tecnologia pra um volume pequeno.

## Consequências

- Um container de banco no compose e um StatefulSet (ou banco gerenciado) no Kubernetes.
- A separação por schema deixa extrair um módulo pra outro banco quase sem mudança.
