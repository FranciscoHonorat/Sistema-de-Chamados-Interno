# ADR 0032 — testify nas asserções dos testes em Go

**Status:** aceita

## Contexto

O pacote `testing` do Go não tem asserções: cada verificação vira um `if` com `t.Errorf`. Com centenas de verificações, isso deixa os testes longos e as mensagens de falha irregulares.

## Decisão

- **`github.com/stretchr/testify`**: `assert` pra verificações que podem continuar e `require` pra pré-condições que param o teste.
- Testes organizados com `t.Run` e nomes que descrevem o comportamento ("should refuse automatic assignment when there are no responsibles").
- Fakes escritos à mão (`port/out/outtest`) em vez de mocks gerados.

## Alternativas consideradas

- **`testing` puro**: nenhuma dependência, mas mais verboso.
- **testify/mock ou gomock**: geram mocks, mas amarram o teste à sequência de chamadas. Fakes com estado real testam comportamento.

## Consequências

- Testes curtos e mensagens de falha com o valor esperado e o recebido.
- Uma dependência só de teste, que não vai pro binário.
