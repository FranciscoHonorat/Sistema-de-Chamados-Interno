# ADR 0009 — Autenticação com JWT Ed25519 e refresh token rotativo

**Status:** aceita

## Contexto

Cada perfil (usuário, suporte, administrador) vê e faz coisas diferentes, então o sistema precisa saber com segurança quem está chamando. O módulo `tickets` precisa validar essa identidade sem depender dos detalhes internos do `employees`.

## Decisão

- **Senhas com bcrypt** (ADR 0029). Login com usuário inexistente e com senha errada respondem igual (`401`) e levam o mesmo tempo, comparando contra um hash fictício.
- **Access token JWT de 15 minutos**, assinado com **Ed25519 (EdDSA)** (ADR 0030). Carrega `sub`, `name`, `role`, `must_change_password`, `iss`, `aud`, `exp` e o `kid` da chave. O verificador exige `EdDSA` explicitamente, o que fecha o ataque de confusão de algoritmo.
- **Chave pública no JWKS** (`/api/employees/.well-known/jwks.json`): qualquer serviço ou gateway futuro valida tokens sem conhecer a chave privada.
- **Refresh token opaco de 7 dias**, com **rotação a cada uso** e **detecção de reuso**: um token já revogado reapresentado revoga todas as sessões do funcionário. Tolerância de 10 s pra duas abas renovando juntas. No banco fica só o hash SHA-256.
- **No navegador**: o access token fica só em memória; o refresh token, num cookie `HttpOnly; Secure; SameSite=Strict` restrito a `/api/employees/auth`.
- **Validação em processo**: o `tickets` recebe um `contracts.TokenVerifier` do `employees` e o traduz pro seu próprio `Actor` (camada anticorrupção em `adapters/out/employees`).
- **Contas novas aguardam aprovação** do administrador, e a recuperação de senha gera uma senha temporária que obriga a troca no primeiro acesso.

## Alternativas consideradas

- **Sessão no servidor com cookie**: revogação imediata, mas cada requisição consulta o banco e extrair um serviço exigiria sessão compartilhada.
- **JWT HS256**: mais simples, mas quem valida precisa da mesma chave que assina.
- **Token no `localStorage`**: sobrevive ao recarregar, mas um XSS consegue ler.

## Consequências

- O access token vale até expirar (15 min) mesmo depois do logout; o refresh é revogado na hora.
- Recarregar a página recupera a sessão chamando o refresh.
- Sem servidor de e-mail, a recuperação de senha depende do administrador.
