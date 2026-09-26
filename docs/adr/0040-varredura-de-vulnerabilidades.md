# ADR 0040 — Varredura de vulnerabilidades: Trivy, govulncheck e npm audit

**Status:** aceita

## Contexto

O código depende de bibliotecas Go, pacotes npm e uma imagem base. Qualquer um deles pode ter uma vulnerabilidade publicada depois que o código foi escrito.

## Decisão

Três verificações no CI, cada uma na camada que conhece melhor:

- **govulncheck**: analisa o código Go e só aponta vulnerabilidades em funções que o projeto realmente chama.
- **npm audit** (`--omit=dev --audit-level=high`): dependências de produção do frontend.
- **Trivy** na imagem final: pacotes do sistema e binários. Gera relatório SARIF (visível na aba Security do GitHub) e **falha o CI** com vulnerabilidade alta ou crítica que já tenha correção.

## Alternativas consideradas

- **Snyk**: mais completo, mas serviço pago pra repositórios privados.
- **Só Dependabot** (ADR 0042): avisa de versões novas, mas não bloqueia o CI.

## Consequências

- Uma vulnerabilidade conhecida com correção não chega em produção.
- Vulnerabilidades sem correção ficam no relatório mas não bloqueiam (`ignore-unfixed`), pra não travar o time por algo que ninguém pode resolver ainda.
