# ADR 0011 — Distribuição automática pelo atendente com menos chamados em aberto

**Status:** aceita

## Contexto

A cliente contou que o suporte reclama da distribuição: uns ficam atolados, outros sem nada. Ela também quer que cada chamado tenha um responsável claro. O enunciado pede a opção de atribuir automaticamente ao responsável com menos chamados em aberto, e deixa pra nós definir o que é "em aberto".

## Decisão

- **"Em aberto" = qualquer status diferente de fechado.** Aberto e em andamento contam como carga, porque os dois são trabalho pendente de quem atende.
- **Escolha**: conta os chamados em aberto por atendente e escolhe o de menor contagem. Empate: o atendente cadastrado há mais tempo.
- **Na abertura**: o campo "Responsável" vem em **Automático** por padrão (`"auto_assign": true`), então o chamado já nasce com dono. Dá pra escolher "Definir depois" e, pro administrador, um atendente específico.
- **Depois de aberto**: o botão "Distribuir automaticamente" (`POST /tickets/:id/assign/auto`) refaz a escolha; "Atribuir a" troca à mão.
- **Usuários comuns** escolhem só entre automático e definir depois. Se pudessem escolher o atendente, todo mundo pediria sempre o mesmo e a distribuição perderia o sentido.
- A mesma função (`pickLeastBusyResponsible`) serve à abertura e ao botão.

## Alternativas consideradas

- **Round-robin**: fácil, mas ignora que alguns chamados demoram muito mais que outros.
- **Peso por prioridade** (alta vale mais): mais justo em tese, mas mais difícil de explicar pra cliente. Fica como evolução.
- **Contar só "aberto"**: quem está atendendo muitos chamados pareceria livre.

## Consequências

- A carga fica equilibrada sem ninguém precisar decidir.
- A regra é simples de explicar e de conferir na tela "Suportes" (carga por atendente).
- Duas aberturas ao mesmo tempo podem escolher o mesmo atendente; com o volume de uma empresa, a diferença é de um chamado e se corrige na próxima distribuição.
