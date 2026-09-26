# ADR 0008 — TDD e estratégia de testes em camadas

**Status:** aceita

## Contexto

O enunciado diz que testes, SOLID e DRY importam mais que features extras, e que menos features com qualidade valem mais que muitas sem teste. O sistema tem regras de negócio (permissões, transições de status, distribuição de carga) que precisam continuar certas a cada mudança.

## Decisão

Escrever os testes junto com o código (TDD) e cobrir cada camada com o tipo de teste mais barato que a exercita:

| Camada | Como |
|---|---|
| Domínio e casos de uso | Testes unitários com fakes das portas (`port/out/outtest`), sem banco |
| Arquitetura | Testes de dependência entre camadas e módulos (ADR 0007) |
| Contratos | `modules/employees/contracts_test.go` |
| Adaptadores | HTTP com `httptest`; Postgres com testes de integração |
| App inteiro | `internal/app/app_test.go`, por HTTP, contra Postgres real |
| Frontend | Vitest + Testing Library, consultando a tela por papel e rótulo (ADR 0026) |
| Navegador | Playwright, desktop e celular (ADR 0027) |
| Implantação | `scripts/smoke-test.sh` contra o compose, o kind e os ambientes do CD |

- Os testes de integração só rodam com `SYS_CALLED_TEST_DATABASE_URL` definida; sem ela, são pulados. Assim `make test` não precisa de nada instalado além de Go e Node.
- O backend roda com `-race` no CI.

## Alternativas consideradas

- **Só testes ponta a ponta**: pegam integração, mas são lentos e apontam mal onde está o erro.
- **Só unitários com mocks**: rápidos, mas não pegam SQL errado nem problema de cookie.

## Consequências

- A maior parte dos testes roda em segundos, sem infraestrutura.
- Integração e E2E confirmam que as peças encaixam de verdade.
- Mais código de teste que de produção em alguns pacotes, o que é esperado.
