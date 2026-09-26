# ADR 0016 — Logs estruturados, request id e probes de saúde

**Status:** aceita

## Contexto

Quando algo dá errado em produção, é preciso achar a requisição, saber se o app está vivo e se está pronto pra receber tráfego. O orquestrador (Docker ou Kubernetes) também precisa dessas respostas pra reiniciar ou tirar uma instância do balanceamento.

## Decisão

- **Logs com `log/slog`** (biblioteca padrão): texto em desenvolvimento, JSON em produção.
- **Request id** em cada requisição (gerado ou herdado de `X-Request-ID`), presente em todas as linhas de log e devolvido no cabeçalho.
- **Access log** com rota, status e duração; erros inesperados ficam no log e a resposta ao cliente é só `internal error`.
- **Recuperação de pânico**: um pânico num handler vira `500` e é logado, sem derrubar o processo.
- **`/healthz`** (liveness: o processo responde) e **`/readyz`** (readiness: pinga o banco).
- **Desligamento gracioso**: ao receber SIGTERM, para de aceitar conexões e espera as requisições em andamento (`SHUTDOWN_TIMEOUT`).
- O binário tem o comando `healthcheck`, usado pelo `HEALTHCHECK` da imagem distroless, que não tem shell nem curl.
- Métricas ficam no ADR 0043.

## Alternativas consideradas

- **zap / zerolog**: um pouco mais rápidos, mas o `slog` já é estruturado e não adiciona dependência.
- **Um endpoint de saúde só**: mistura "está vivo" com "o banco está fora", e o orquestrador reiniciaria o app à toa quando o banco cai.

## Consequências

- Uma requisição é rastreável de ponta a ponta pelo request id.
- Um banco fora do ar tira a instância do balanceamento sem reiniciá-la.
- Deploys não derrubam requisições em andamento.
