# ADR 0031 — Vue Router com guarda por perfil

**Status:** aceita

## Contexto

O sistema tem telas diferentes por perfil (usuário, suporte, administrador), telas públicas (login, cadastro, recuperação de senha) e precisa mandar cada um pro lugar certo depois do login.

## Decisão

- **Vue Router** com histórico HTML5 (`createWebHistory`). O backend devolve o `index.html` pras rotas da SPA (ADR 0013).
- Rotas marcadas com `meta.role` ou `meta.requiresAuth`.
- Uma guarda de navegação (`router/navigationGuard.ts`, função pura `resolveNavigation`) manda quem não está logado pro login, prende na troca de senha quem entrou com senha temporária, tira quem já está logado da tela de login e devolve pra home do perfil quem tenta abrir a área de outro perfil.
- O menu lateral é montado a partir do perfil (`router/menu.ts`).

## Alternativas consideradas

- **Roteamento por arquivo** (Nuxt, unplugin-vue-router): menos configuração, mas traz convenções e dependências a mais.

## Consequências

- Cada perfil só navega pelo que é dele. A proteção real continua no backend (ADR 0010).
- A guarda é uma função pura, testada sem navegador.
