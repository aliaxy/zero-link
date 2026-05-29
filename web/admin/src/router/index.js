import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { public: true },
    },
    {
      path: '/',
      redirect: '/links',
      component: () => import('@/components/AppLayout.vue'),
      children: [
        {
          path: 'links',
          name: 'links',
          component: () => import('@/views/LinksView.vue'),
        },
        {
          path: 'links/:id',
          name: 'link-detail',
          component: () => import('@/views/LinkDetailView.vue'),
        },
        {
          path: 'links/:id/stats',
          name: 'link-stats',
          component: () => import('@/views/LinkStatsView.vue'),
        },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.token) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.name === 'login' && auth.token) {
    return { name: 'links' }
  }
})

export default router
