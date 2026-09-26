# ADR 0003 — DDD tático: agregados e value objects

**Status:** aceita

## Contexto

As regras do chamado são o coração do sistema: um chamado fechado não aceita edição nem resposta, só pode ir pra "em andamento" a partir de "aberto", precisa de um relatório pra ser fechado, só aceita prioridades e status conhecidos. Espalhar essas validações pelos handlers HTTP faz a mesma regra aparecer em vários lugares e ser esquecida em algum.

## Decisão

Usar os padrões táticos do Domain-Driven Design:

- **Agregado `Ticket`** (`modules/tickets/internal/domain/ticket`): toda mudança passa por um método dele (`Edit`, `AssignTo`, `ChangePriority`, `MoveToInProgress`, `Close`, `AddResponse`), que valida a transição antes de aplicar.
- **Value objects** (`domain/valueobjects`): `Title`, `Description`, `Priority`, `Status`, `AssigneeID`, `Content`… Um valor inválido não chega a existir: o construtor devolve erro.
- **Erros de domínio** nomeados (`domain/domain-errors`), traduzidos pra status HTTP só no adaptador.
- **Políticas no agregado**: permissões (`permissions.go`) e visibilidade (`visibility.go`) ficam junto do chamado (ADR 0010).
- O módulo `employees` segue a mesma ideia com o agregado `Employee` e o `RefreshToken` de sessão.

## Alternativas consideradas

- **Modelo anêmico** (structs com dados e regras nos services): menos código, mas a regra fica fora de quem tem o estado, e duplicar é fácil.
- **Validação só na borda** (tags de binding do Gin): pega formato, não regra de negócio.

## Consequências

- Uma regra existe num lugar só e é testada sem infraestrutura.
- Os value objects deixam as assinaturas dos métodos autoexplicativas.
- Mais tipos pra escrever, e cada novo campo pede um value object.
