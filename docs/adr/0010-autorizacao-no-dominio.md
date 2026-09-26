# ADR 0010 — Autorização como regra de domínio

**Status:** aceita

## Contexto

Quem pode ver e mudar um chamado depende do perfil e do próprio chamado: o usuário só vê os que abriu e só edita enquanto ninguém começou o atendimento; o suporte vê os abertos e os atribuídos a ele; só o responsável inicia e fecha. Checar isso nos handlers HTTP espalha a regra e deixa buracos.

## Decisão

- As regras ficam no agregado: `ticket.IsVisibleTo(actor)` (`visibility.go`) e as permissões `CanEdit`, `CanRespond`, `CanManage`, `CanWork` (`permissions.go`).
- Todo comando passa por `EventSourcedUseCase.UpdateTicket`, que recebe a permissão exigida e checa visibilidade e permissão antes de mudar o chamado.
- Um chamado que o ator não pode ver responde **404, igual a um inexistente**: a API não revela que ele existe. Se pode ver mas não pode agir, **403**.
- O autor de uma resposta é sempre quem está autenticado, nunca um campo do corpo da requisição.
- O frontend espelha as regras (`tickets/permissions.ts`) só pra esconder botões; quem decide é o backend.

| Perfil | Vê | Pode |
|---|---|---|
| Usuário | os que abriu | abrir, responder, editar os seus enquanto abertos |
| Suporte | abertos e atribuídos a ele | atribuir, distribuir, mudar prioridade, responder, iniciar e fechar os seus |
| Administrador | todos | editar, atribuir, distribuir, mudar prioridade, responder |

## Alternativas consideradas

- **Middleware por rota**: resolve o perfil, mas não regras que dependem do chamado (dono, responsável, status).
- **Biblioteca de políticas** (Casbin, OPA): poderosa, mas exagerada pra três perfis.

## Consequências

- As regras são testadas no domínio, sem HTTP.
- Um caso de uso novo não esquece a checagem, porque o caminho padrão já a exige.
- A regra existe duas vezes (backend e frontend), mas só a do backend protege.
