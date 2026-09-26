# ADR 0041 — CodeQL pra análise estática de segurança

**Status:** aceita

## Contexto

Alguns erros de segurança (injeção de SQL, XSS, uso inseguro de criptografia) passam por testes e revisão porque o código "funciona".

## Decisão

- Workflow `codeql.yml` com **CodeQL** do GitHub pra **Go** e **JavaScript/TypeScript**.
- Roda em push na `main`, em pull requests e toda semana (segunda-feira), pra pegar regras novas mesmo sem mudança no código.
- Os achados aparecem na aba Security do GitHub.

## Alternativas consideradas

- **Semgrep**: regras fáceis de escrever, mas exige configuração e conta pra ter o painel.
- **SonarQube / SonarCloud**: mede também qualidade, mas é mais um serviço pra configurar.

## Consequências

- Análise de segurança sem custo pra repositórios públicos.
- Complementa o `gosec` do golangci-lint (ADR 0044), que roda localmente.
