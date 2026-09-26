# ADR 0026 — Vitest e Testing Library nos testes do frontend

**Status:** aceita

## Contexto

Os componentes têm lógica que precisa de teste: validação de formulário, ações visíveis por perfil, filtros da lista, renovação de sessão. Testes presos à estrutura interna dos componentes quebram a cada refatoração.

## Decisão

- **Vitest** com ambiente `jsdom` e cobertura por `@vitest/coverage-v8`.
- **Testing Library** (`@testing-library/vue` + `user-event` + `jest-dom`): os testes procuram elementos **por papel e rótulo** (`getByRole('button', { name: 'Abrir chamado' })`, `getByLabelText('Prioridade')`), como um usuário ou um leitor de tela.
- APIs falsas injetadas por `provide` (ADR 0019).

## Alternativas consideradas

- **Jest**: exigiria configurar a transformação de Vue e TypeScript à parte.
- **Vue Test Utils puro**: acessa a instância do componente, o que leva a testar detalhes internos.

## Consequências

- Testes que sobrevivem a refatorações de marcação.
- A acessibilidade é testada de brinde: um botão sem rótulo não é encontrado.
- Cobertura de ~94% das linhas do frontend.
