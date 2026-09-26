# ADR 0004 — Event Sourcing no módulo de chamados

**Status:** aceita

## Contexto

A cliente reclamou que "ninguém sabe o que ficou pra fazer" e que falta clareza sobre quem é o responsável. Num chamado importa não só o estado atual, mas a história: quem abriu, quem atribuiu a quem, quando começou o atendimento, quem respondeu o quê, quem fechou e com que resolução.

## Decisão

Guardar cada chamado como a sequência dos seus eventos, e não como uma linha atualizada no lugar:

- Eventos: `TicketOpened`, `TicketEdited`, `TicketAssigned`, `TicketPriorityChanged`, `TicketMovedToInProgress`, `TicketResponseAdded`, `TicketClosed` (`domain/event`).
- Tabela `ticket_events` com `aggregate_id`, `version`, `event_type`, `payload` (JSONB), `occurred_at` e `actor_id` (quem executou a ação).
- O estado é reconstruído com `ticket.LoadFromHistory`, que aplica os eventos em ordem.
- **Concorrência otimista**: `UNIQUE (aggregate_id, version)`. Duas gravações simultâneas no mesmo chamado: a segunda recebe `ErrConcurrencyConflict` (HTTP 409) em vez de sobrescrever a primeira.
- A data de abertura é a `occurred_at` do `TicketOpened`, e a de fechamento, a do `TicketClosed`.

## Alternativas consideradas

- **CRUD com tabela de auditoria**: mais simples de consultar, mas a auditoria vira um segundo registro que pode divergir do principal.
- **CRUD sem histórico**: perde exatamente a informação que a cliente precisa.

## Consequências

- Trilha de auditoria completa por construção, com o autor de cada ação.
- Novas visões (tempo médio de atendimento, reatribuições) saem dos eventos que já existem.
- Listar exige reconstruir os chamados. A mitigação é o cache (ADR 0012) e a carga de todos os eventos numa consulta só (`LoadMany`). O próximo passo, documentado nos trade-offs do README, é uma projeção de leitura.
- Eventos gravados não mudam: uma mudança de formato pede versionamento do evento.
