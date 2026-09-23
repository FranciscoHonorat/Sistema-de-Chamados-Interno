import type { RouteLocationRaw } from 'vue-router'

import type { Role } from '../auth/authService'
import type { IconName } from '../components/ui/AppIcon.vue'
import { TicketStatus } from '../tickets/status'
import { paths } from './paths'

export interface MenuItem {
  label: string
  to: RouteLocationRaw
  icon: IconName
  description: string
}

function ticketsWithStatus(status: TicketStatus): RouteLocationRaw {
  return { path: paths.tickets, query: { status } }
}

const menus: Record<Role, MenuItem[]> = {
  user: [
    { label: 'Meus chamados', to: paths.tickets, icon: 'inbox', description: 'Acompanhe os chamados que você abriu' },
  ],
  support: [
    { label: 'Chamados abertos', to: ticketsWithStatus(TicketStatus.Open), icon: 'inbox', description: 'Chamados esperando atendimento' },
    { label: 'Em atendimento', to: ticketsWithStatus(TicketStatus.InProgress), icon: 'clock', description: 'O que está com você agora' },
    { label: 'Fechados por mim', to: ticketsWithStatus(TicketStatus.Closed), icon: 'archive', description: 'Histórico do que você resolveu' },
  ],
  admin: [
    { label: 'Todos os chamados', to: paths.tickets, icon: 'inbox', description: 'Todos os chamados da empresa' },
    { label: 'Usuários', to: paths.users, icon: 'users', description: 'Aprove contas e gere senhas temporárias' },
    { label: 'Suportes', to: paths.supports, icon: 'chart', description: 'Carga de trabalho de cada atendente' },
  ],
}

export function menuFor(role: Role): MenuItem[] {
  return menus[role]
}
