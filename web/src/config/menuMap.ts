export const ROLES = {
  PLATFORM: 'platform_admin',
  SCHOOL: 'school_admin',
  SUPPLIER: 'supplier',
  CANTEEN: 'canteen_admin',
  ROOT: 'root'
}

export interface MenuItem {
  path: string
  title: string
  icon?: string
  roles?: string[]
  children?: MenuItem[]
}

// 租户端菜单 (学校/供应商/食堂)
export const tenantMenus: MenuItem[] = [
  {
    path: '/workspace/dashboard',
    title: '工作台',
    icon: 'Odometer',
    roles: [ROLES.SCHOOL, ROLES.SUPPLIER, ROLES.CANTEEN, ROLES.ROOT]
  },
  {
    path: '/workspace/scm',
    title: '供应链管理',
    icon: 'Box',
    children: [
      {
        path: '/workspace/scm/audit',
        title: '商品审核',
        roles: [ROLES.SCHOOL, ROLES.ROOT] // 仅学校/Root可见
      },
      {
        path: '/workspace/scm/entry',
        title: '商品录入',
        roles: [ROLES.SUPPLIER] // 仅供应商可见
      }
    ]
  },
  {
    path: '/workspace/order',
    title: '订单中心',
    icon: 'Tickets',
    children: [
      {
        path: '/workspace/order/list',
        title: '我的订单',
        roles: [ROLES.SCHOOL, ROLES.SUPPLIER, ROLES.CANTEEN, ROLES.ROOT]
      }
    ]
  }
]

// 平台端菜单
export const platformMenus: MenuItem[] = [
  { path: '/platform/dashboard', title: '平台总览', icon: 'DataLine' },
  { path: '/platform/sites', title: '站点管理', icon: 'School' }
]