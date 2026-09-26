# ADR 0033 — UUIDs como identificadores

**Status:** aceita

## Contexto

Chamados, respostas, eventos, notificações e contas precisam de identificadores. Com Event Sourcing, o ID do chamado precisa existir antes de gravar o primeiro evento.

## Decisão

- **`github.com/google/uuid`** (UUID v4 aleatório).
- O ID é gerado na aplicação, não no banco: o value object `ID` cria um novo quando recebe `uuid.Nil`.
- IDs de ticket inválidos na URL respondem erro antes de consultar o banco.

## Alternativas consideradas

- **ID sequencial do banco** (`SERIAL`): legível, mas só existe depois do `INSERT` e deixa adivinhar os IDs de outros chamados.
- **UUID v7 / ULID**: ordenáveis por tempo, melhores pra índice. Ficam como evolução; o volume atual não justifica.

## Consequências

- O agregado nasce com ID antes de tocar no banco.
- IDs não revelam quantos chamados existem nem permitem adivinhar os dos outros.
- IDs longos na URL.
