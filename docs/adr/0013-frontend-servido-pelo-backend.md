# ADR 0013 — Frontend servido pelo próprio backend, na mesma origem

**Status:** aceita

## Contexto

numa equipe full stack pequena, reduzir o atrito entre frontend e backend economiza muito tempo. Dois deploys, CORS, cookies entre domínios e versões desencontradas entre API e tela são atrito.

## Decisão

- O build do Vite (`frontend/dist`) entra na mesma imagem do backend e é servido pelo binário Go (`internal/platform/httpserver/spa.go`).
- API em `/api/*`, frontend no resto, com fallback pro `index.html` (rotas da SPA). Assets com hash recebem cache imutável; o `index.html`, `no-cache`.
- Em desenvolvimento, o Vite faz proxy de `/api` pro backend, então o código do frontend é o mesmo nos dois modos.

## Alternativas consideradas

- **Laravel + Inertia.js** (stack da Codificar): elimina a API separada, mas prende o frontend ao ciclo de requisição do servidor. Aqui a API REST também serve integrações e testes de ponta a ponta.
- **Frontend em CDN separada**: escala melhor os estáticos, mas exige CORS e cookie entre domínios.
- **Nginx servindo os estáticos**: mais um container e mais uma configuração.

## Consequências

- Um deploy, uma imagem, uma versão: tela e API nunca ficam desencontradas.
- Sem CORS, e o cookie `SameSite=Strict` funciona sem ajuste.
- A CSP pode ser restrita a `'self'`.
- Um bug no frontend exige publicar a imagem inteira.
