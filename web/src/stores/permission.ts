import { defineStore } from 'pinia'
import type { MenuItem } from '@/types/config'
import router from '@/router'
import type { RouteRecordRaw } from 'vue-router'
import path from 'path-browserify'

/**
 * 检查角色是否有权限访问路由
 * @param route 路由
 * @param role 角色
 */
function hasPermission(route: RouteRecordRaw, role: string) {
  if (route.meta && route.meta.roles) {
    return (route.meta.roles as string[]).includes(role)
  }
  // 如果没有定义 roles，则默认所有角色都可以访问
  return true
}

/**
 * 递归地从路由构建菜单
 * @param routes 路由配置
 * @param role 角色
 * @param basePath 基础路径
 */
function buildMenusFromRoutes(routes: readonly RouteRecordRaw[], role: string, basePath = '/'): MenuItem[] {
  const res: MenuItem[] = []

  for (const route of routes) {
    // 过滤掉没有 meta 或者 title 的路由，这些通常不是菜单项
    if (!route.meta || !route.meta.title) {
      continue
    }
    
    // 检查权限
    if (hasPermission(route, role)) {
      const fullPath = path.resolve(basePath, route.path)
      const menuItem: MenuItem = {
        path: fullPath,
        title: route.meta.title as string,
        icon: route.meta.icon as string || undefined,
        children: []
      }

      const childrenRoutes = route.children
      if (childrenRoutes && childrenRoutes.length > 0) {
        menuItem.children = buildMenusFromRoutes(childrenRoutes, role, fullPath)
      }
      
      // 如果一个菜单项没有任何有权限的子菜单，那它自身也不应该显示（除非它本身就是个页面）
      if (menuItem.children && menuItem.children.length === 0 && !route.component) {
        continue
      }
      // 如果只有一个子菜单，有些UI设计会选择直接提升子菜单，这里我们保持结构
      
      res.push(menuItem)
    }
  }
  return res
}


export const usePermissionStore = defineStore('permission', {
  state: () => ({ 
    menus: [] as MenuItem[] 
  }),
  actions: {
    generateMenus(role: string, layoutType: 'platform' | 'tenant') {
      const layoutRouteName = layoutType === 'platform' ? '/platform' : '/workspace';
      const layoutRoute = router.options.routes.find(r => r.path === layoutRouteName);
      
      if (layoutRoute && layoutRoute.children) {
        this.menus = buildMenusFromRoutes(layoutRoute.children, role, layoutRoute.path)
      } else {
        this.menus = []
      }
    }
  }
})