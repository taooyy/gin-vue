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
