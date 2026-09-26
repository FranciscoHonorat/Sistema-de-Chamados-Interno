# ADR 0045 — ESLint e vue-tsc no frontend

**Status:** aceita

## Contexto

O frontend em TypeScript e Vue precisa do mesmo cuidado que o backend: erros de tipo e padrões ruins apanhados antes de rodar.

## Decisão

- **ESLint** (flat config, `eslint.config.js`) com `@eslint/js`, `typescript-eslint` e `eslint-plugin-vue`.
- **`vue-tsc -b`**: checagem de tipos completa, incluindo os templates `.vue`, que o `tsc` sozinho não entende.
- Os dois rodam no CI e em `make lint`; o build (`npm run build`) também roda o `vue-tsc` antes do Vite.

## Alternativas consideradas

- **Biome**: lint e formatação muito mais rápidos, mas o suporte a Vue ainda é parcial.
- **Só o TypeScript do editor**: não protege o que entra no repositório.

## Consequências

- Um template que usa uma propriedade inexistente quebra o build.
- Mais tempo de CI, compensado por menos erros em produção.
