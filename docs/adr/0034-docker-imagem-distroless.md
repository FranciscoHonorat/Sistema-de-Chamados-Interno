# ADR 0034 — Imagem Docker multi-stage com base distroless

**Status:** aceita

## Contexto

A aplicação tem que rodar igual na máquina de quem avalia, no CI e em produção. A imagem deve ser pequena e ter a menor superfície de ataque possível.

## Decisão

Um `Dockerfile` na raiz, em três estágios:

1. **Node 24 (alpine)**: `npm ci` e build do frontend.
2. **Go 1.26 (alpine)**: build do binário com `CGO_ENABLED=0`, `-trimpath` e versão/commit embutidos. Roda na plataforma da máquina e faz cross-compile, então o build amd64 + arm64 não precisa de emulação.
3. **`gcr.io/distroless/static-debian12:nonroot`**: só o binário e o `dist/` do frontend.

- Roda como usuário sem privilégios.
- `HEALTHCHECK` com o próprio binário (`sys-called healthcheck`), já que a imagem não tem shell nem curl.
- Seguro por padrão: `APP_ENV=production`.
- Cache de dependências do npm e do Go com `--mount=type=cache`.

## Alternativas consideradas

- **Alpine como base final**: tem shell pra depurar, mas também mais pacotes com vulnerabilidades.
- **`scratch`**: ainda menor, mas sem certificados CA, sem fuso horário e sem usuário `nonroot` prontos.
- **Imagens separadas pra frontend e backend**: ver ADR 0013.

## Consequências

- Imagem de ~45 MB, sem shell nem gerenciador de pacotes.
- Depurar dentro do container exige uma imagem de debug ou `kubectl debug`.
