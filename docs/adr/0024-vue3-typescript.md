# ADR 0024 — Vue 3 com TypeScript no frontend

**Status:** aceita

## Contexto

A tela precisa de listagem com filtros e busca, detalhe com conversa, formulários, modais e comportamento diferente por perfil. O dia a dia da Codificar usa Vue.js.

## Decisão

- **Vue 3** com **Composition API** e `<script setup>`.
- **TypeScript** em todo o frontend, checado por `vue-tsc` no CI.
- Componentes de interface próprios e pequenos (`components/ui/`: `AppButton`, `StatusBadge`, `PriorityBadge`, `EmptyState`…) reaproveitados nas telas.
- Tipos da API (`Ticket`, `TicketDetail`, `NewTicket`…) declarados junto do cliente HTTP.

## Alternativas consideradas

- **React**: ecossistema maior, porém o projeto não precisa de SPA.
- **Vue 3 com JavaScript**: menos configuração, mas erros de contrato com a API só apareceriam rodando.
- **Biblioteca de componentes pronta** (Vuetify, PrimeVue): acelera, mas traz um visual e um peso que o projeto não precisa (ADR 0028).

## Consequências

- Mesmo framework da equipe que vai avaliar e manter.
- Mudanças no contrato da API quebram o build em vez de quebrar em produção.
