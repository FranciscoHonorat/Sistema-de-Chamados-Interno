# ADR 0014 — Migrações SQL embutidas no binário, por schema

**Status:** aceita

## Contexto

O banco precisa estar no formato certo antes do app atender. Em produção podem subir várias réplicas ao mesmo tempo, e cada módulo é dono do seu schema.

## Decisão

- Cada módulo embute seus arquivos `.sql` com `embed.FS` (`adapters/out/postgres/migrations/`) e os expõe por `Migrations()`.
- `internal/platform/migrate` aplica os arquivos em ordem, por schema, e registra cada um em `<schema>.schema_migrations`.
- Um **advisory lock do Postgres** serializa execuções concorrentes: várias réplicas subindo juntas aplicam cada arquivo uma vez só.
- O binário tem o comando `sys-called migrate`. O compose roda um serviço `migrate` antes do app, e o chart usa um *init container*. Fora disso, `MIGRATE_ON_START=true` migra na subida.
- As migrações são idempotentes (`CREATE TABLE IF NOT EXISTS`, `ADD COLUMN IF NOT EXISTS`).

## Alternativas consideradas

- **golang-migrate, goose, Atlas**: completos, mas mais uma dependência e mais um binário na imagem pra fazer o que cabe em um pacote pequeno.
- **ORM com auto-migrate**: esconde o SQL e não lida bem com mudanças destrutivas.

## Consequências

- A imagem carrega o próprio schema: não há arquivo externo pra esquecer.
- Sem migração de volta (*down*): corrigir é escrever uma migração nova.
