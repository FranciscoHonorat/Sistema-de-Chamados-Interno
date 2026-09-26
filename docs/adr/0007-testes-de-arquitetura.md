# ADR 0007 — Fronteiras garantidas por testes de arquitetura

**Status:** aceita

## Contexto

Regras como "o domínio não depende de adaptadores" ou "um módulo só fala com outro pelo `contracts`" se perdem com o tempo se dependerem só de revisão de código. Basta um import apressado.

## Decisão

Garantir as regras de dependência em três níveis:

1. **Compilador Go**: o que está em `modules/<m>/internal/` só pode ser importado de dentro de `modules/<m>/`.
2. **`backend/architecture_test.go`**: um módulo só importa outro pelo pacote `contracts`; o `contracts` só usa a biblioteca padrão; `internal/platform` não conhece nenhum módulo.
3. **`modules/<m>/internal/architecture_test.go`**: o domínio não depende da aplicação nem dos adaptadores; a aplicação não depende de adaptadores, de frameworks (Gin, pgx, JWT, Prometheus), da plataforma nem de outros módulos; um adaptador não importa outro.

Um **teste de contrato** (`modules/employees/contracts_test.go`) garante que os eventos gravados no outbox continuam decodificando nas structs publicadas em `contracts`.

## Alternativas consideradas

- **Só convenção e revisão**: funciona até a primeira pressa.
- **Ferramenta externa** (go-arch-lint, depguard): mais uma configuração pra manter. Os testes leem os imports de cada arquivo com `go/parser`, da biblioteca padrão, e rodam junto com `go test`.

## Consequências

- Violar a arquitetura quebra o `go test` e o CI.
- Os testes documentam as regras em código executável.
- Mudar uma regra de propósito exige mudar o teste, o que é desejável.
