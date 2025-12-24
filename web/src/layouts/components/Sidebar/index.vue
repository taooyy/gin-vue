<template>
  <el-menu
    :default-active="activeMenu"
    :collapse="appStore.isSidebarCollapsed"
    background-color="#2a384a"
    text-color="#bdc3c7"
    active-text-color="#ffffff"
    router
    style="height: calc(100vh - 60px); border-right: none"
  >
    <sidebar-item
      v-for="menu in permissionStore.menus"
      :key="menu.path"
      :item="menu"
      :class="{ 'top-level-active': isTopLevelActive(menu, route.path) }"
    />
  </el-menu>
</template>

<script setup lang="ts">
  import { computed } from 'vue';
  import { useRoute } from 'vue-router';
  import { usePermissionStore } from '@/stores/permission';
  import { useAppStore } from '@/stores/app';
  import SidebarItem from './SidebarItem.vue';
  import type { MenuItem } from '@/types/config';

  const route = useRoute();
  const permissionStore = usePermissionStore();
  const appStore = useAppStore();
  const activeMenu = computed(() => route.path);

  /**
   * Checks if a top-level menu item is the parent of the currently active route.
   */
  function isTopLevelActive(topLevelMenu: MenuItem, activePath: string): boolean {
    if (!topLevelMenu.children || topLevelMenu.children.length === 0) {
      return topLevelMenu.path === activePath;
    }
    // A top-level menu is considered "active" if the current route path starts with its path.
    // e.g., menu path: /workspace/scm, active path: /workspace/scm/audit -> true
    return activePath.startsWith(topLevelMenu.path);
  }
</script>

<style scoped>
  :deep(.el-menu-item),
  :deep(.el-sub-menu__title) {
    margin: 4px 8px;
    border-radius: 6px;
    transition: all 0.2s ease;
    height: 48px;
    line-height: 48px;
  }

  :deep(.el-menu-item:hover),
  :deep(.el-sub-menu__title:hover) {
    background-color: #34495e;
    color: #ffffff;
  }

  :deep(.el-menu-item.is-active) {
    background: linear-gradient(90deg, #2ecc71, #f39c12);
    box-shadow: 0 0 10px rgba(46, 204, 113, 0.6);
    color: #ffffff;
    font-weight: bold;
  }

  /* NEW: Highlight top-level active icon when collapsed */
  :deep(.el-menu--collapse .top-level-active.is-active > .el-sub-menu__title .el-icon) {
    color: #ffd700 !important; /* Brighter gold/yellow */
    filter: drop-shadow(0 0 8px rgba(255, 215, 0, 0.9)); /* Stronger glow */
    text-shadow: 0 0 3px #ffd700; /* Additional text shadow for more pop */
  }

  /* Explicitly hide text when collapsed */
  :deep(.el-menu--collapse .el-menu-item span),
  :deep(.el-menu--collapse .el-sub-menu__title span) {
    display: none !important;
  }
</style>
