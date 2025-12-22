import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Login from '@/views/auth/Login.vue'
import { ROLES } from '@/config/types'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: Login },
    { path: '/', redirect: '/login' },
    // 平台布局
    {
      path: '/platform',
      component: () => import('@/layouts/PlatformLayout.vue'),
      meta: { roles: [ROLES.PLATFORM, ROLES.ROOT], title: '平台' }, // 平台管理员和Root可见
      children: [
        { 
          path: 'dashboard', 
          component: () => import('@/views/platform/dashboard.vue'),
          meta: { title: '平台总览', icon: 'DataLine' }
        },
        { 
          path: 'site-management', 
          redirect: '/platform/site-management/list', // 默认重定向到站点展示
          meta: { title: '站点管理', icon: 'School' },
          children: [
            { 
              path: 'list', // 对应 /platform/site-management/list
              component: () => import('@/views/platform/SiteList.vue'),
              meta: { title: '站点展示', icon: 'List' }
            },
            { 
              path: 'manage', // 对应 /platform/site-management/manage
              component: () => import('@/views/platform/SiteManage.vue'),
              meta: { title: '站点管理', icon: 'Setting' }
            }
          ]
        },
        {
          path: 'permission',
          redirect: '/platform/permission/user', // 默认重定向到用户权限管理
          meta: { title: '权限管理', icon: 'Unlock' },
          children: [
            {
              path: 'user', // 对应 /platform/permission/user
              component: () => import('@/views/platform/UserPermission.vue'),
              meta: { title: '用户权限管理', icon: 'UserFilled' }
            },
            {
              path: 'role', // 对应 /platform/permission/role
              component: () => import('@/views/platform/RoleManagement.vue'),
              meta: { title: '角色管理', icon: 'Avatar' }
            }
          ]
        },
        {
          path: 'console',
          redirect: '/platform/console/log', // 默认重定向到日志管理
          meta: { title: '控制台', icon: 'Monitor' },
          children: [
            {
              path: 'log', // 对应 /platform/console/log
              component: () => import('@/views/platform/LogManagement.vue'),
              meta: { title: '日志管理', icon: 'Document' }
            }
          ]
        }
      ]
    },
    // 租户布局
    {
      path: '/workspace',
      component: () => import('@/layouts/TenantLayout.vue'),
      meta: { roles: [ROLES.SCHOOL, ROLES.SUPPLIER, ROLES.CANTEEN, ROLES.ROOT], title: '工作区' },
      children: [
        { 
          path: 'dashboard', 
          component: () => import('@/views/workspace/dashboard.vue'),
          meta: { title: '工作台', icon: 'Odometer', roles: [ROLES.SCHOOL, ROLES.SUPPLIER, ROLES.CANTEEN, ROLES.ROOT] }
        },
        { 
          path: 'scm', 
          // 该路由只用于创建菜单父级，无实际页面
          redirect: '/workspace/scm/audit',
          meta: { title: '供应链管理', icon: 'Box' },
          children: [
            {
              path: 'audit',
              component: () => import('@/views/workspace/scm/ProductAudit.vue'),
              meta: { title: '商品审核', roles: [ROLES.SCHOOL, ROLES.ROOT] }
            },
            {
              path: 'entry',
              component: () => import('@/views/workspace/scm/ProductEntry.vue'),
              meta: { title: '商品录入', roles: [ROLES.SUPPLIER] }
            }
          ]
        },
        {
          path: 'order',
          redirect: '/workspace/order/list',
          meta: { title: '订单中心', icon: 'Tickets' },
          children: [
            {
              path: 'list',
              component: () => import('@/views/workspace/order/OrderList.vue'),
              meta: { title: '我的订单', roles: [ROLES.SCHOOL, ROLES.SUPPLIER, ROLES.CANTEEN, ROLES.ROOT] }
            }
          ]
        }
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