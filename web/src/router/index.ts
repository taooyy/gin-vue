import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Login from '@/views/auth/Login.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: Login },
    { path: '/', redirect: '/login' },
    // 平台布局
    {
      path: '/platform',
      component: () => import('@/layouts/PlatformLayout.vue'),
      children: [
        { path: 'dashboard', component: () => import('@/views/platform/dashboard.vue') },
        { path: 'sites', component: () => import('@/views/platform/SiteList.vue') }
      ]
    },
    // 租户布局
    {
      path: '/workspace',
      component: () => import('@/layouts/TenantLayout.vue'),
      children: [
        { path: 'dashboard', component: () => import('@/views/workspace/dashboard.vue') },
        { path: 'scm/audit', component: () => import('@/views/workspace/scm/ProductAudit.vue') },
        { path: 'scm/entry', component: () => import('@/views/workspace/scm/ProductEntry.vue') },
        { path: 'order/list', component: () => import('@/views/workspace/order/OrderList.vue') }
      ]
    }
  ]
})

router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()
  if (to.path !== '/login' && !userStore.token) {
    next('/login')
  } else {
    next()
  }
})

export default router