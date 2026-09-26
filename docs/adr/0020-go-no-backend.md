# ADR 0020 — Go no backend

**Status:** aceita

## Contexto

O backend precisa ser simples de rodar pra qualquer pessoa do time avaliador, fácil de manter por uma equipe pequena e leve em produção. A escolha de linguagem é livre; o dia a dia da Codificar usa Laravel, e Go aparece como diferencial.

## Decisão

Usar **Go 1.26** no backend.

- **Binário estático único** (`CGO_ENABLED=0`): a imagem final não precisa de runtime, só do binário (ADR 0034).
- **`internal/` garantido pelo compilador**: é a base das fronteiras entre módulos (ADR 0007).
- **Biblioteca padrão forte**: `log/slog`, `embed`, `net/http`, `crypto/ed25519`, `context` e `testing` cobrem boa parte do que o projeto precisa.
- **Concorrência nativa**: relay do outbox, servidor HTTP e servidor de métricas rodam juntos no mesmo processo.
- **Testes com `-race`**: o detector de corrida roda no CI.
- **Cross-compile**: a imagem sai pra amd64 e arm64 sem emulação.

## Alternativas consideradas

- **PHP + Laravel**: a stack da Codificar, com muito pronto (ORM, filas, auth). Em troca, precisa de runtime (PHP-FPM + servidor web) e as fronteiras entre módulos dependem de convenção.
- **Node.js (NestJS)**: mesma linguagem no front e no back, mas precisa de runtime na imagem e tem tipagem só em tempo de compilação do TypeScript.
- **Java / Kotlin (Spring)**: robusto, mas pesado pra um sistema desse tamanho.

## Consequências

- Imagem final de ~45 MB, que sobe em milissegundos.
- Código explícito e às vezes mais verboso (tratamento de erro a cada chamada).
- Quem avalia precisa conhecer Go pra ler o backend. O README e os ADRs explicam a estrutura.
