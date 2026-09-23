import { createRouter, type RouterHistory } from 'vue-router'

import type { Role } from '../auth/authService'
import type { Session } from '../auth/session'
import ChangePasswordView from '../views/ChangePasswordView.vue'
import HomeView from '../views/HomeView.vue'
import LoginPage from '../views/LoginPage.vue'
import NotFoundView from '../views/NotFoundView.vue'
import RecoverPasswordView from '../views/RecoverPasswordView.vue'
import SignUpView from '../views/SignUpView.vue'
import SupportsView from '../views/SupportsView.vue'
import TicketDetailView from '../views/TicketDetailView.vue'
import TicketsView from '../views/TicketsView.vue'
import UsersView from '../views/UsersView.vue'
import { resolveNavigation } from './navigationGuard'
import { paths } from './paths'

declare module 'vue-router' {
  interface RouteMeta {
    role?: Role
    requiresAuth?: boolean
    title?: string
  }
}

export function createAppRouter(session: Session, history: RouterHistory) {
  const router = createRouter({
    history,
    routes: [
      { path: paths.login, component: LoginPage, meta: { title: 'Entrar' } },
      { path: paths.register, component: SignUpView, meta: { title: 'Criar conta' } },
      { path: paths.recoverPassword, component: RecoverPasswordView, meta: { title: 'Recuperar senha' } },
      { path: paths.changePassword, component: ChangePasswordView, meta: { requiresAuth: true, title: 'Trocar senha' } },
      { path: paths.userHome, component: HomeView, meta: { role: 'user', title: 'Início' } },
      { path: paths.supportHome, component: HomeView, meta: { role: 'support', title: 'Início' } },
      { path: paths.adminHome, component: HomeView, meta: { role: 'admin', title: 'Início' } },
      { path: paths.tickets, component: TicketsView, meta: { requiresAuth: true, title: 'Chamados' } },
      { path: paths.ticketDetail, component: TicketDetailView, props: true, meta: { requiresAuth: true, title: 'Chamado' } },
      { path: paths.users, component: UsersView, meta: { role: 'admin', title: 'Usuários' } },
      { path: paths.supports, component: SupportsView, meta: { role: 'admin', title: 'Suportes' } },
      { path: '/:pathMatch(.*)*', component: NotFoundView, meta: { title: 'Página não encontrada' } },
    ],
    scrollBehavior(_to, _from, saved) {
      return saved ?? { top: 0 }
    },
  })

  router.beforeEach((to) => resolveNavigation(session.user.value, to))
  router.afterEach((to) => {
    document.title = to.meta.title ? `${to.meta.title} · CodeTicket` : 'CodeTicket'
  })

  return router
}
