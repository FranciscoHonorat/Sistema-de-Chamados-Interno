# ADR 0027 — Playwright nos testes ponta a ponta

**Status:** aceita

## Contexto

Os testes unitários não provam que o fluxo inteiro funciona: login real, cookie de sessão, backend, banco e tela juntos, com três perfis interagindo no mesmo chamado.

## Decisão

- **Playwright** (`frontend/e2e`), contra o app de verdade rodando (compose ou binário local), com os usuários de desenvolvimento.
- Dois projetos: **desktop** (Chrome 1440×900) e **celular** (Pixel 7).
- Cenários: login e erros, home por perfil, bloqueio de área de outro perfil, logout, sessão após recarregar, aprovação de conta nova, ciclo completo de um chamado entre os três perfis, distribuição automática, validação, busca, tema e menu no celular.
- Cada teste cria os próprios dados com nomes únicos: rodam em paralelo e podem repetir no mesmo banco.
- No CI, trace e screenshot só quando falha; o relatório é guardado como artefato.

## Alternativas consideradas

- **Cypress**: popular, mas mais lento em paralelo e com suporte limitado a várias abas/contextos, que o teste entre três perfis usa.
- **Selenium**: mais configuração e testes mais instáveis.

## Consequências

- O fluxo principal da cliente é verificado no navegador a cada push.
- Os testes E2E precisam do app no ar, por isso têm um job próprio no CI.
