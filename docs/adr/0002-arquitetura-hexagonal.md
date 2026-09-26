# ADR 0002 — Arquitetura hexagonal dentro de cada módulo

**Status:** aceita

## Contexto

Cada módulo tem regras de negócio (quem pode ver um chamado, como uma senha é validada, quando um chamado pode ser fechado) e várias formas de falar com o mundo: HTTP, Postgres, JWT, barramento de eventos. Se as regras dependerem direto dessas tecnologias, qualquer troca de biblioteca vira uma reescrita, e testar uma regra passa a exigir banco e servidor HTTP.

## Decisão

Organizar cada módulo em portas e adaptadores:

- `internal/domain/` — agregados, value objects e regras. Não conhece nada de fora.
- `internal/application/` — casos de uso. Declaram as portas de saída que precisam em `port/out` (ex.: `EventStore`, `ResponsibleDirectory`, `TokenIssuer`).
- `internal/adapters/in/` — quem chama o módulo: HTTP (Gin), assinantes de eventos, relay do outbox.
- `internal/adapters/out/` — implementações das portas: Postgres, bcrypt, JWT, cache.
- `module.go` — a composição: monta os adaptadores e injeta nos casos de uso, à mão, sem framework de injeção de dependências.

## Alternativas consideradas

- **MVC / camadas por tipo** (controllers, models, services): mais conhecido, mas deixa o modelo acoplado ao ORM e ao HTTP.
- **Framework de DI** (wire, fx): tira código de composição, mas esconde como o módulo é montado. Com dois módulos, montar à mão é curto e explícito.

## Consequências

- Casos de uso são testados com fakes das portas (`port/out/outtest`), sem banco e sem rede.
- Trocar uma tecnologia (ex.: o barramento por um broker) mexe só no adaptador.
- Mais arquivos e interfaces do que numa estrutura MVC. As regras de dependência entre as camadas são garantidas por teste (ADR 0007).
