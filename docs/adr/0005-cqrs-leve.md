# ADR 0005 — CQRS leve: comandos e consultas separados

**Status:** aceita

## Contexto

As operações que mudam um chamado (abrir, atribuir, fechar) e as que só leem (listar, detalhar, carga dos atendentes) têm necessidades diferentes. As escritas passam pelo agregado e pelo controle de versão; as leituras filtram por visibilidade e montam respostas pra tela.

## Decisão

Separar os casos de uso do módulo `tickets` em dois pacotes:

- `application/command/` — um caso de uso por intenção: `OpenTicket`, `EditTicket`, `AssignTicket`, `AutoAssignTicket`, `ChangeTicketPriority`, `MoveTicketToInProgress`, `CloseTicket`, `AddTicketResponse`…
- `application/query/` — `ListTickets`, `GetTicket`, `ListResponsibles`, `SupportWorkload`, `ListNotifications`.

Os dois lados usam o mesmo armazenamento de eventos. É CQRS "leve": não há banco de leitura separado.

## Alternativas consideradas

- **Um service por entidade** (`TicketService` com todos os métodos): cresce sem limite e mistura leitura com escrita.
- **CQRS completo com projeções**: mais rápido pra ler, mas exige manter as projeções em dia. Fica como evolução natural (ver trade-offs no README).

## Consequências

- Cada caso de uso é pequeno, com um teste próprio.
- A separação já está pronta pra receber uma projeção de leitura sem tocar nos comandos.
- Mais arquivos do que um service único.
