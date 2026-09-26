# ADR 0039 — Imagem assinada com cosign, com SBOM e proveniência

**Status:** aceita

## Contexto

Quem implanta precisa ter certeza de que a imagem em produção é exatamente a que o CI construiu e testou, e saber o que tem dentro dela.

## Decisão

No job de publicação (`release.yml`):

- Build multi-arquitetura (amd64 + arm64) com **SBOM** e **proveniência** (`provenance: mode=max`) gerados pelo BuildKit.
- **Atestado de proveniência** do GitHub (`actions/attest-build-provenance`).
- **Assinatura keyless com cosign** (Sigstore), usando a identidade OIDC do próprio workflow, sem chave pra guardar.
- O deploy (`deploy.yml`) roda `cosign verify` antes do `helm upgrade` e implanta **pelo digest**, não pela tag.
- As notas de release trazem o comando pra verificar a assinatura.

## Alternativas consideradas

- **Docker Content Trust (Notary v1)**: exige gerenciar chaves e está em desuso.
- **Não assinar**: comum, mas uma tag pode ser sobrescrita no registry sem ninguém perceber.

## Consequências

- Uma imagem adulterada ou de outra origem é recusada no deploy.
- O SBOM permite responder "estamos afetados?" quando sai uma vulnerabilidade nova.
