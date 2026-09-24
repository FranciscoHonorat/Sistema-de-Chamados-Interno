# ADR 0001 — Monólito modular em vez de microsserviços

**Status:** aceita

## Contexto

O sistema tem dois domínios bem separados: funcionários (contas, sessões e perfis) e chamados. Uma opção seria um serviço por domínio, cada um com seu banco, conversando por um broker de mensagens.

Pro tamanho do domínio (funcionários e chamados de uma empresa), essa separação cobraria caro:

- containers de infraestrutura extras, incluindo um broker de mensagens pra operar e monitorar;
- validação de token por HTTP entre os serviços;
- deploys que precisariam ser coordenados quando um contrato mudasse.

O time que vai manter o produto é pequeno, e o enunciado pede que o projeto rode localmente com facilidade.

## Decisão

Construir um **monólito modular**: um binário e uma imagem, que também servem o frontend. Cada domínio vira um módulo em `backend/modules/<nome>`, com schema próprio no Postgres.

O que uma arquitetura de serviços teria de bom fica garantido dentro do monólito:

- **Fronteiras claras**: um módulo não importa o `internal/` de outro, o que o compilador Go garante. Os testes de arquitetura garantem que a conversa entre módulos passa só pelo pacote `contracts`.
- **Comunicação assíncrona**: o `employees` grava os eventos numa tabela de outbox na mesma transação. Um relay publica esses eventos num barramento em processo, e o `tickets` assina os que precisa.
- **Contrato de eventos próprio**: tipo + payload JSON, independente do meio de transporte.

## Consequências

- Um container do app e um do Postgres bastam pra rodar tudo; não há broker nem rede entre os módulos.
- Um bug que derrube o processo derruba os dois módulos. As réplicas e o PDB do chart cuidam da disponibilidade.
- A consistência entre módulos é eventual (o relay roda a cada 1 s).
- Extrair um módulo pra um serviço no futuro é trocar o barramento em processo por um broker e o `TokenVerifier` em processo por um verificador JWKS. O domínio e os casos de uso não mudam.
