# ADR 0028 — Tailwind CSS 4 com tokens de design próprios

**Status:** aceita

## Contexto

precisava de um framework CSS moderno, que permita montar a interface rápido sem reinventar a roda, e sem ser um custo adicional para o projeto. Não era necessário um design profissional: basta uma interface organizada e funcional. A tela também precisa funcionar no celular e no tema escuro.

## Decisão

- **Tailwind CSS 4**, integrado pelo plugin `@tailwindcss/vite`.
- **Tokens semânticos** em `frontend/src/style.css` (superfície, texto, bordas, marca `--color-brand-*`), redefinidos no **tema escuro**. Os componentes usam `bg-surface`, `text-ink`, `border-line`… e não cores fixas.
- Algumas classes utilitárias compostas (`card`, `field-input`, `field-label`) pra padrões repetidos.
- Layout responsivo: barra lateral no desktop, gaveta no celular; tabelas viram cartões em telas pequenas.

## Alternativas consideradas

- **Bootstrap**: componentes prontos, mas visual reconhecível e difícil de adaptar sem sobrescrever CSS.
- **Biblioteca de componentes Vue** (Vuetify, PrimeVue): acelera formulários e tabelas, mas pesa no bundle e impõe o próprio design.
- **CSS puro / módulos CSS**: controle total, mas mais lento de escrever.

## Consequências

- Interface consistente com pouco CSS escrito à mão.
- Tema escuro sai de graça: basta trocar os tokens.
- As classes utilitárias deixam o template mais longo.
