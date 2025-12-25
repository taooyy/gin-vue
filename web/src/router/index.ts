import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from '@/stores/user';
import { usePermissionStore } from '@/stores/permission';
import Login from '@/views/auth/Login.vue';
import { ROLES } from '@/types/config';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: Login },
    { path: '/', redirect: '/platform/dashboard' },
    // 平台布局
    {
      path: '/platform',
      component: () => import('@/layouts/PlatformLayout.vue'),
      meta: { roles: [ROLES.PLATFORM, ROLES.ROOT], title: '平台' }, // 平台管理员和Root可见
      children: [
        {
          path: 'dashboard',
          component: () => import('@/views/platform/dashboard.vue'),
          meta: { title: '平台总览', icon: 'DataLine' },
        },
        {
          path: 'site-management',
          redirect: '/platform/site-management/list', // 默认重定向到站点展示
          meta: { title: '站点管理', icon: 'School' },
          children: [
            {
              path: 'list', // 对应 /platform/site-management/list
              name: 'SiteList',
              component: () => import('@/views/platform/SiteList.vue'),
              meta: { title: '站点展示', icon: 'List' },
            },
            {
              path: 'manage', // 对应 /platform/site-management/manage
              name: 'SiteManage',
              component: () => import('@/views/platform/SiteManage.vue'),
              meta: { title: '站点管理', icon: 'Setting' },
            },
          ],
        },
        {
          path: 'permission',
          redirect: '/platform/permission/site-permission',
          meta: { title: '权限管理', icon: 'Unlock' },
          children: [
            {
              path: 'site-permission',
              name: 'SitePermission',
              component: () => import('@/views/platform/SitePermission.vue'),
              meta: { title: '站点权限管理' },
            },
            {
              path: 'role-management',
              name: 'SiteRoleManagement',
              component: () => import('@/views/platform/RoleManagement.vue'),
              meta: { title: '站点角色管理' },
            },
            {
              path: 'feature-display',
              name: 'SiteFeatureDisplay',
              component: () => import('@/views/platform/FeatureDisplay.vue'),
              meta: { title: '站点功能展示' },
            },
            {
              path: 'role-assignment',
              name: 'RoleFeatureAssignment',
              component: () => import('@/views/platform/RolePermission.vue'),
              meta: { title: '角色功能划分' },
            },
          ],
        },
        {
          path: 'console',
          redirect: '/platform/console/log', // 默认重定向到日志管理
          meta: { title: '控制台', icon: 'Monitor' },
          children: [
            {
              path: 'log', // 对应 /platform/console/log
              component: () => import('@/views/platform/LogManagement.vue'),
              meta: { title: '日志管理', icon: 'Document' },
            },
          ],
        },
        {
          path: 'account',
          component: () => import('@/views/platform/AccountManagement.vue'),
          meta: { title: '账号管理', icon: 'User', roles: [ROLES.PLATFORM, ROLES.ROOT] },
        },
      ],
    },
    // 租户布局
    {
      path: '/workspace',
      component: () => import('@/layouts/TenantLayout.vue'),
      meta: { roles: [ROLES.SCHOOL, ROLES.SUPPLIER, ROLES.CANTEEN, ROLES.ROOT], title: '工作区' },
      children: [
        // --- Shared ---
        {
          path: 'dashboard',
          component: () => import('@/views/workspace/dashboard.vue'),
          meta: {
            title: '工作台',
            icon: 'Odometer',
            roles: [ROLES.SCHOOL, ROLES.SUPPLIER, ROLES.CANTEEN, ROLES.ROOT],
          },
        },
        // --- School Menu ---
        {
          path: 'scm',
          redirect: '/workspace/scm/supplier',
          meta: { title: '供应链管理', icon: 'Box', roles: [ROLES.SCHOOL, ROLES.ROOT] },
          children: [
            {
              path: 'supplier',
              component: () => import('@/views/workspace/scm/SupplierManagement.vue'),
              meta: { title: '供应商管理', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
            {
              path: 'staff',
              component: () => import('@/views/workspace/scm/SupplierStaff.vue'),
              meta: { title: '供应商员工', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
            {
              path: 'order',
              component: () => import('@/views/workspace/scm/SupplierOrder.vue'),
              meta: { title: '供应商订单', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
          ],
        },
        {
          path: 'canteen',
          redirect: '/workspace/canteen/overview',
          meta: { title: '食堂管理', icon: 'OfficeBuilding', roles: [ROLES.SCHOOL, ROLES.ROOT] },
          children: [
            {
              path: 'overview',
              component: () => import('@/views/workspace/canteen/CanteenOverview.vue'),
              meta: { title: '食堂总览', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
            {
              path: 'merchant',
              component: () => import('@/views/workspace/canteen/MerchantList.vue'),
              meta: { title: '商户列表', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
          ],
        },
        {
          path: 'order',
          redirect: '/workspace/order/summary',
          meta: { title: '订单管理', icon: 'Tickets', roles: [ROLES.SCHOOL, ROLES.ROOT] },
          children: [
            {
              path: 'summary',
              component: () => import('@/views/workspace/order/OrderSummary.vue'),
              meta: { title: '订单汇总', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
            {
              path: 'list',
              component: () => import('@/views/workspace/order/OrderList.vue'),
              meta: { title: '商户订单', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
            {
              path: 'aftersales',
              component: () => import('@/views/workspace/order/OrderAftersales.vue'),
              meta: { title: '订单售后', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
          ],
        },
        {
          path: 'product',
          redirect: '/workspace/product/audit',
          meta: { title: '商品管理', icon: 'Goods', roles: [ROLES.SCHOOL, ROLES.ROOT] },
          children: [
            {
              path: 'audit',
              component: () => import('@/views/workspace/product/ProductAudit.vue'),
              meta: { title: '商品审核', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
            {
              path: 'library',
              component: () => import('@/views/workspace/product/ProductLibrary.vue'),
              meta: { title: '商品库', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
            {
              path: 'listing',
              component: () => import('@/views/workspace/product/ProductListing.vue'),
              meta: { title: '上架管理', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
            {
              path: 'pricing',
              component: () => import('@/views/workspace/product/ProductPricing.vue'),
              meta: { title: '商品价格', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
          ],
        },
        {
          path: 'traceability',
          component: () => import('@/views/workspace/traceability/TraceabilityManagement.vue'),
          meta: { title: '溯源管理', icon: 'Guide', roles: [ROLES.SCHOOL, ROLES.ROOT] },
        },
        {
          path: 'school-settlement',
          redirect: '/workspace/school-settlement/statements',
          meta: { title: '结算管理', icon: 'Coin', roles: [ROLES.SCHOOL, ROLES.ROOT] },
          children: [
            {
              path: 'statements',
              component: () => import('@/views/workspace/settlement/SchoolStatements.vue'),
              meta: { title: '对账单', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
            {
              path: 'canteen',
              component: () => import('@/views/workspace/settlement/CanteenSettlement.vue'),
              meta: { title: '食堂结算', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
            {
              path: 'supplier',
              component: () => import('@/views/workspace/settlement/SupplierSettlement.vue'),
              meta: { title: '供应商结算', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
            {
              path: 'merchant',
              component: () => import('@/views/workspace/settlement/MerchantSettlement.vue'),
              meta: { title: '商户结算', roles: [ROLES.SCHOOL, ROLES.ROOT] },
            },
          ],
        },
        {
          path: 'school-account',
          component: () => import('@/views/workspace/school/AccountManagement.vue'),
          meta: { title: '账号管理', icon: 'User', roles: [ROLES.SCHOOL, ROLES.ROOT] },
        },
        // --- Supplier Menu ---
        {
          path: 'supplier-product',
          redirect: '/workspace/supplier-product/upload',
          meta: { title: '商品管理', icon: 'Goods', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
          children: [
            {
              path: 'upload',
              component: () => import('@/views/workspace/supplier/product/ProductUpload.vue'),
              meta: { title: '商品上传', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
            },
            {
              path: 'quotation',
              component: () => import('@/views/workspace/supplier/product/ProductQuotation.vue'),
              meta: { title: '商品报价', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
            },
            {
              path: 'modification',
              component: () => import('@/views/workspace/supplier/product/ProductModification.vue'),
              meta: { title: '商品修改', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
            },
          ],
        },
        {
          path: 'supplier-delivery',
          redirect: '/workspace/supplier-delivery/picking',
          meta: { title: '配送管理', icon: 'Van', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
          children: [
            {
              path: 'picking',
              component: () => import('@/views/workspace/supplier/delivery/OrderPicking.vue'),
              meta: { title: '订单分拣', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
            },
            {
              path: 'delivery',
              component: () => import('@/views/workspace/supplier/delivery/OrderDelivery.vue'),
              meta: { title: '订单配送', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
            },
          ],
        },
        {
          path: 'supplier-order',
          redirect: '/workspace/supplier-order/list',
          meta: { title: '订单管理', icon: 'Tickets', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
          children: [
            {
              path: 'list',
              component: () => import('@/views/workspace/supplier/order/OrderList.vue'),
              meta: { title: '订单', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
            },
            {
              path: 'aftersales',
              component: () => import('@/views/workspace/supplier/order/OrderAftersales.vue'),
              meta: { title: '订单售后', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
            },
          ],
        },
        {
          path: 'supplier-account',
          component: () => import('@/views/workspace/supplier/AccountManagement.vue'),
          meta: { title: '账号管理', icon: 'User', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
        },
        {
          path: 'supplier-settlement',
          redirect: '/workspace/supplier-settlement/statements',
          meta: { title: '结算管理', icon: 'Coin', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
          children: [
            {
              path: 'statements',
              component: () => import('@/views/workspace/supplier/settlement/Statements.vue'),
              meta: { title: '对账单', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
            },
            {
              path: 'statistics',
              component: () =>
                import('@/views/workspace/supplier/settlement/SettlementStatistics.vue'),
              meta: { title: '结算统计', roles: [ROLES.SUPPLIER, ROLES.ROOT] },
            },
          ],
        },
        // --- Canteen Menu ---
        {
          path: 'canteen-merchant',
          redirect: '/workspace/canteen-merchant/list',
          meta: { title: '商户管理', icon: 'Shop', roles: [ROLES.CANTEEN, ROLES.ROOT] },
          children: [
            {
              path: 'list',
              component: () => import('@/views/workspace/canteen/MerchantList.vue'),
              meta: { title: '商户列表', roles: [ROLES.CANTEEN, ROLES.ROOT] },
            },
            {
              path: 'manage',
              component: () => import('@/views/workspace/canteen/MerchantManage.vue'),
              meta: { title: '商户管理', roles: [ROLES.CANTEEN, ROLES.ROOT] },
            },
          ],
        },
        {
          path: 'canteen-order',
          redirect: '/workspace/canteen-order/summary',
          meta: { title: '订单管理', icon: 'Tickets', roles: [ROLES.CANTEEN, ROLES.ROOT] },
          children: [
            {
              path: 'summary',
              component: () => import('@/views/workspace/order/OrderSummary.vue'),
              meta: { title: '订单汇总', roles: [ROLES.CANTEEN, ROLES.ROOT] },
            },
            {
              path: 'list',
              component: () => import('@/views/workspace/order/OrderList.vue'),
              meta: { title: '商户订单', roles: [ROLES.CANTEEN, ROLES.ROOT] },
            },
          ],
        },
        {
          path: 'canteen-account',
          component: () => import('@/views/workspace/canteen/AccountManagement.vue'),
          meta: { title: '账号管理', icon: 'User', roles: [ROLES.CANTEEN, ROLES.ROOT] },
        },
      ],
    },
  ],
});

router.beforeEach(async (to, _from, next) => {
  const userStore = useUserStore();
  const permissionStore = usePermissionStore();
  const hasToken = userStore.token;

  if (hasToken) {
    // 用户已登录
    if (to.path === '/login') {
      // 如果已登录，则重定向到首页
      next({ path: '/' });
    } else {
      // 检查菜单是否已生成
      if (userStore.role && permissionStore.menus.length === 0) {
        // 如果菜单为空（通常是刷新页面导致），则重新生成
        try {
          const layoutType = to.path.startsWith('/platform') ? 'platform' : 'tenant';
          permissionStore.generateMenus(userStore.role, layoutType);
          // 动态添加路由后，需要使用 next({ ...to, replace: true }) 来确保路由完全加载
          next({ ...to, replace: true });
        } catch (error) {
          // 如果生成菜单出错（例如，角色信息有问题），则登出并跳转到登录页
          console.error('生成菜单时出错:', error);
          userStore.logout();
          next('/login');
        }
      } else {
        // 菜单已存在，直接放行
        next();
      }
    }
  } else {
    // 用户未登录
    if (to.path !== '/login') {
      // 如果访问的是非登录页，则重定向到登录页
      next('/login');
    } else {
      // 如果访问的是登录页，直接放行
      next();
    }
  }
});

export default router;
