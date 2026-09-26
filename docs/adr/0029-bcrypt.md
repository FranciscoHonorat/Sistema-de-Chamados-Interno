# ADR 0029 — bcrypt pra guardar senhas

**Status:** aceita

## Contexto

O sistema tem login com usuário e senha. Senhas nunca podem ser guardadas em texto nem com hash rápido (MD5, SHA), que se quebra por força bruta.

## Decisão

- **bcrypt** de `golang.org/x/crypto`, com o custo padrão (`bcrypt.DefaultCost`).
- Fica atrás da porta `PasswordHasher`, no adaptador `adapters/out/security/bcrypt_hasher.go`.
- Login com usuário inexistente compara a senha contra um hash fictício, pra levar o mesmo tempo que uma senha errada.
- Senhas precisam ter pelo menos 8 caracteres.

## Alternativas consideradas

- **Argon2id**: recomendado hoje pelo OWASP e mais resistente a GPU, mas com mais parâmetros pra ajustar. bcrypt continua aceito e é o padrão da maioria dos frameworks (inclusive Laravel).
- **scrypt**: equivalente, menos comum.

## Consequências

- Senhas protegidas mesmo se o banco vazar.
- Trocar pra Argon2id no futuro é implementar outro adaptador da mesma porta.
