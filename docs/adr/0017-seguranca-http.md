# ADR 0017 — Cabeçalhos de segurança e limites nas requisições HTTP

**Status:** aceita

## Contexto

O app serve tela e API na mesma origem (ADR 0013) e guarda a sessão em cookie. XSS, clickjacking e corpos gigantes são riscos comuns nesse cenário.

## Decisão

Um middleware da plataforma aplica em toda resposta:

- **Content-Security-Policy** restrita a `'self'` (scripts, conexões, fontes), `object-src 'none'`, `frame-ancestors 'none'`.
- `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`, `Referrer-Policy`, `Permissions-Policy` (sem câmera, microfone e localização), `Cross-Origin-Opener-Policy`.
- **HSTS** só quando a requisição chegou por HTTPS (direto ou via `X-Forwarded-Proto`).
- **Limite de 1 MB no corpo**: a API só recebe JSON pequeno.
- **gzip** nas respostas (klauspost/compress `gzhttp`).

## Alternativas consideradas

- **Deixar pro proxy reverso** (Nginx, Ingress): funciona, mas o app ficaria desprotegido rodando sem proxy, como no compose local.
- **Pacote de middleware pronto** (ex.: `secure`): menos controle sobre a CSP exata.

## Consequências

- Proteção igual em qualquer ambiente, com ou sem proxy na frente.
- Os cabeçalhos são verificados por teste (`server_test.go`) e pelo smoke test.
- Qualquer script ou recurso externo novo exige ajustar a CSP.
