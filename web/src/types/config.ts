export const ROLES = {
  PLATFORM: 'platform_admin',
  PLATFORM_STAFF: 'platform_staff',
  SCHOOL: 'school_admin',
  SCHOOL_STAFF: 'school_staff',
  SUPPLIER: 'supplier_admin',
  SUPPLIER_STAFF: 'supplier_staff',
  CANTEEN: 'canteen_admin',
  ROOT: 'root',
};

export interface MenuItem {
  path: string;
  title: string;
  icon?: string;
  roles?: string[];
  children?: MenuItem[];
}
