import { defineStore } from 'pinia'
import { tenantMenus, platformMenus, type MenuItem } from '@/config/menuMap'

function filterMenus(menus: MenuItem[], role: string): MenuItem[] {
  const res: MenuItem[] = []
  menus.forEach(menu => {
    const temp = { ...menu }
    if (!temp.roles || temp.roles.includes(role)) {
      if (temp.children) temp.children = filterMenus(temp.children, role)
      res.push(temp)
    }
  })
  return res
}

export const usePermissionStore = defineStore('permission', {
  state: () => ({ menus: [] as MenuItem[] }),
  actions: {
    generateMenus(role: string, layoutType: 'platform' | 'tenant') {
      this.menus = layoutType === 'platform' ? platformMenus : filterMenus(tenantMenus, role)
    }
  }
})