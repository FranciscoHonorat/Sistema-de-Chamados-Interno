# ADR 0019 — Estado e dependências do frontend com provide/inject e composables

**Status:** aceita

## Contexto

O frontend precisa compartilhar entre as telas a sessão do usuário e os clientes da API, e os testes precisam trocar a API real por uma falsa sem esforço.

## Decisão

- `main.ts` cria a sessão e os clientes (`createTicketsApi`, `createEmployeesApi`, `createHttpAuthService`) e os entrega com `provide`, cada um com uma `InjectionKey` tipada.
- Componentes pegam o que precisam com `inject`.
- Lógica reaproveitável fica em composables (`useTicketDetail`, `useLoad`, `usePolling`, `useTheme`, `useToast`).
- O cliente HTTP (`api/apiClient.ts`) renova a sessão uma vez quando recebe `401` e repete a chamada.
- Nos testes, `fakeTicketsApi()` e `fakeAuthService()` são injetados no lugar dos reais.

## Alternativas consideradas

- **Pinia**: store global com devtools, mas o estado compartilhado aqui é pequeno (a sessão), e o resto é dado de cada tela.
- **Importar os clientes direto**: mais curto, mas os testes teriam que fazer mock de módulo.

## Consequências

- Nenhuma dependência extra além de Vue e Vue Router.
- Cada teste escolhe exatamente a API que quer simular.
- Se o estado compartilhado crescer muito, Pinia volta a fazer sentido.
