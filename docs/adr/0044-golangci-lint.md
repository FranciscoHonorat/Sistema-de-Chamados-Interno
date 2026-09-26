# ADR 0044 — golangci-lint no backend

**Status:** aceita

## Contexto

Formatação e erros comuns (erro não verificado, corpo de resposta não fechado, comparação de erros com `==`) não deveriam depender de revisão humana.

## Decisão

- **golangci-lint v2** (`backend/.golangci.yml`), rodado por `make lint` e no CI.
- Linters: o conjunto padrão (`errcheck`, `govet`, `ineffassign`, `staticcheck`, `unused`) mais `bodyclose`, `errorlint`, `gosec`, `misspell`, `nilerr`, `noctx`, `unconvert` e `usestdlibvars`.
- Formatadores: `gofmt` e `goimports`, com os imports do projeto agrupados à parte.
- Nos testes, algumas regras (`gosec`, `noctx`, `bodyclose`, `errcheck`) são relaxadas.

## Alternativas consideradas

- **Só `go vet` e `gofmt`**: pegam menos problemas.
- **staticcheck sozinho**: bom, mas sem as regras de segurança do `gosec`.

## Consequências

- Estilo uniforme sem discussão em revisão.
- Problemas de segurança e de tratamento de erro aparecem antes do commit.
