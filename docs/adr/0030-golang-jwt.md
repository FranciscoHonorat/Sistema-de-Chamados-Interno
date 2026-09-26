# ADR 0030 — golang-jwt pra emitir e validar tokens

**Status:** aceita

## Contexto

A autenticação usa JWT assinado com Ed25519 (ADR 0009). Implementar o formato JWT à mão é fácil de errar em pontos de segurança.

## Decisão

- **`github.com/golang-jwt/jwt/v5`**, a biblioteca JWT mais usada em Go.
- Emissão em `jwt_issuer.go` e validação em `jwt_verifier.go`, ambos em `adapters/out/security`.
- A validação fixa o algoritmo (`EdDSA`), o emissor (`codeticket-employees`), a audiência e a expiração.
- A chave vem de `JWT_PRIVATE_KEY` (PKCS8 PEM). Em desenvolvimento, sem chave, uma efêmera é gerada; em produção a chave é obrigatória.

## Alternativas consideradas

- **lestrrat-go/jwx**: suporta JWK/JWS/JWE completos, mais do que o projeto precisa.
- **PASETO**: evita escolhas perigosas de algoritmo por projeto, mas é menos conhecido e sem JWKS padrão.

## Consequências

- Formato padrão, que gateways e outros serviços entendem.
- Fixar o algoritmo na validação fecha o ataque de confusão de algoritmo.
