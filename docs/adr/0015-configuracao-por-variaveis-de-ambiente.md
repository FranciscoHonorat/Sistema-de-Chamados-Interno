# ADR 0015 — Configuração por variáveis de ambiente, validada na subida

**Status:** aceita

## Contexto

A mesma imagem roda no compose local, no kind, em staging e em produção. Configuração em arquivo dentro da imagem obrigaria uma imagem por ambiente.

## Decisão

- Toda configuração vem do ambiente (`internal/platform/config`), no estilo *twelve-factor*: `DATABASE_URL`, `APP_ENV`, `JWT_PRIVATE_KEY`, `PORT`, `METRICS_PORT`, `STATIC_DIR`, `MIGRATE_ON_START`, `SEED_DEMO_USERS`, `TICKETS_CACHE_ENABLED`, `DB_MAX_CONNS`, `LOG_LEVEL`, `LOG_FORMAT`, `SHUTDOWN_TIMEOUT`.
- A validação junta **todos os problemas de uma vez** e o app sai com código 2, em vez de falhar um por um.
- **Seguro por padrão**: a imagem vem com `APP_ENV=production`, que exige a chave JWT, desliga os usuários de demonstração e usa logs em JSON. O compose muda pra `development` de forma explícita.
- `make env` gera o `.env` a partir do `.env.example`, com uma chave Ed25519 nova.

## Alternativas consideradas

- **Arquivo YAML/TOML**: bom pra muita configuração aninhada, mas precisa ser montado em cada ambiente.
- **Viper**: lê de várias fontes, mas é pesado pra 13 variáveis.

## Consequências

- Uma imagem só pra todos os ambientes.
- Esquecer uma variável obrigatória aparece na subida, com a lista completa.
- Segredos vêm de Secret no Kubernetes e do `.env` (fora do git) no compose.
