# ADR 0018 — Notificações gravadas no banco e consultadas a cada 30 segundos

**Status:** aceita

## Contexto

Quem está envolvido num chamado precisa saber quando algo muda: o suporte e o administrador quando um chamado é aberto, o responsável quando recebe um chamado, o solicitante quando o atendimento começa ou termina, e os administradores quando alguém pede conta ou recuperação de senha.

## Decisão

- As notificações são geradas a partir dos eventos do chamado (`notification.ForTicketEvent`) e dos eventos do `employees` recebidos pelo barramento, e gravadas na tabela de notificações.
- Quem executou a ação não é notificado sobre ela mesma.
- O sino no cabeçalho (`NotificationBell.vue`) consulta `GET /notifications` a cada 30 s e mostra o número de não lidas; `POST /notifications/read` marca como lidas.

## Alternativas consideradas

- **Server-Sent Events / WebSocket**: tempo real, mas exige conexões longas, cuidado com proxies e com várias réplicas. É o próximo passo documentado.
- **E-mail**: exigiria um servidor de e-mail só pra isso.

## Consequências

- Simples, funciona atrás de qualquer proxy e com várias réplicas.
- Atraso de até 30 s pra uma notificação aparecer.
- Uma consulta por usuário logado a cada 30 s, desprezível no volume de uma empresa.
