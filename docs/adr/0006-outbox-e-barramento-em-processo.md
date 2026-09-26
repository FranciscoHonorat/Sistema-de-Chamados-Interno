# ADR 0006 — Transactional outbox e barramento de eventos em processo

**Status:** aceita

## Contexto

O módulo `tickets` precisa saber quando um funcionário do suporte é cadastrado (pra mantê-lo na lista de atendentes) e quando alguém pede conta ou troca de senha (pra notificar os administradores). Chamar o `tickets` direto do `employees` acoplaria os módulos. Publicar o evento depois de gravar no banco tem o problema clássico: se o processo cair entre as duas coisas, o evento se perde.

## Decisão

- **Outbox transacional**: o `employees` grava o funcionário e o evento na tabela `outbox_events` **na mesma transação**.
- **Relay** (`adapters/in/outbox/relay.go`): a cada 1 s, pega os eventos pendentes e publica no barramento.
- **Barramento em processo** (`internal/platform/eventbus`): mensagens com tipo + payload JSON. Um assinante que falha ou entra em pânico vira erro, sem derrubar os outros.
- **Entrega at-least-once** com **assinantes idempotentes**: cada mensagem leva o ID do evento do outbox, e reentregas não duplicam nada.
- **Várias réplicas**: o relay reserva lotes com `FOR UPDATE SKIP LOCKED` e um *lease* de 15 s (`locked_until`). Duas réplicas nunca entregam o mesmo evento ao mesmo tempo.
- Os eventos publicados ficam em `modules/employees/contracts`, que só depende da biblioteca padrão.

## Alternativas consideradas

- **Chamada direta entre módulos**: simples, mas acopla e perde o evento se o `tickets` falhar.
- **Broker externo** (Kafka, RabbitMQ, NATS): entrega entre processos, mas é mais um serviço pra operar sem necessidade num processo só (ADR 0001).
- **LISTEN/NOTIFY do Postgres**: não guarda a mensagem se ninguém estiver ouvindo.

## Consequências

- Nenhum evento se perde, mesmo com queda no meio da operação.
- A consistência entre módulos é eventual (até ~1 s).
- Um evento que falha sempre volta até ser entregue. Uma fila de mensagens mortas depois de N tentativas é o próximo passo.
- Trocar o barramento por um broker, se um módulo virar serviço, não muda o contrato dos eventos.
