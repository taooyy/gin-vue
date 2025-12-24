<template>
  <div class="navbar">
    <div class="left-menu">
      <!-- 菜单收缩按钮 -->
      <el-icon class="sidebar-toggler" :size="24" @click="appStore.toggleSidebar()">
        <component :is="appStore.isSidebarCollapsed ? Expand : Fold" />
      </el-icon>
      <!-- 面包屑导航 -->
      <el-breadcrumb separator="/">
        <el-breadcrumb-item v-for="item in breadcrumbs" :key="item.path" :to="{ path: item.path }">
          {{ item.title }}
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <div class="right-menu">
      <!-- 功能图标 -->
      <div class="action-icons">
        <el-icon class="action-icon" :size="20"><Search /></el-icon>
        <el-icon class="action-icon" :size="20"><Bell /></el-icon>
        <el-icon class="action-icon" :size="20"><FullScreen /></el-icon>
      </div>

      <!-- 右侧用户菜单 -->
      <el-dropdown class="user-menu" @command="handleCommand">
        <div class="user-avatar-wrapper">
          <el-avatar size="small" :icon="UserFilled" class="user-avatar" />
          <span class="user-name">{{ userStore.role }}</span>
          <el-icon><arrow-down /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile" disabled>
              <el-icon><User /></el-icon>个人中心
            </el-dropdown-item>
            <el-dropdown-item command="logout" divided>
              <el-icon><SwitchButton /></el-icon>退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { useUserStore } from '@/stores/user';
  import { useAppStore } from '@/stores/app';
  import {
    ArrowDown,
    User,
    SwitchButton,
    UserFilled,
    Search,
    Bell,
    FullScreen,
    Fold,
    Expand,
  } from '@element-plus/icons-vue';

  const route = useRoute();
  const router = useRouter();
  const userStore = useUserStore();
  const appStore = useAppStore();

  // 计算面包屑
  const breadcrumbs = computed(() => {
    return route.matched
      .filter(
        (item) =>
          item.meta && item.meta.title && item.path !== '/platform' && item.path !== '/workspace'
      )
      .map((item) => ({
        path: item.path,
        title: item.meta.title as string,
      }));
  });

  // 处理下拉菜单命令
  const handleCommand = (command: string) => {
    if (command === 'logout') {
      userStore.logout();
      router.push('/login');
    }
  };
</script>

<style scoped>
  .navbar {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 16px 0 0; /* Adjusted padding */
    background-color: #2a384a; /* Dark background */
    color: #bdc3c7; /* Light text color */
  }

  .left-menu {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .sidebar-toggler {
    cursor: pointer;
    padding: 10px;
    transition: transform 0.3s;
  }

  .sidebar-toggler:hover {
    transform: scale(1.1);
  }

  .right-menu {
    display: flex;
    align-items: center;
    gap: 20px;
  }

  .action-icons {
    display: flex;
    align-items: center;
    gap: 15px;
  }

  .action-icon {
    cursor: pointer;
    transition:
      color 0.2s ease,
      transform 0.2s ease;
  }

  .action-icon:hover {
    color: #ffffff;
    transform: scale(1.1);
  }

  .user-menu {
    cursor: pointer;
  }

  .user-avatar-wrapper {
    display: flex;
    align-items: center;
  }

  .user-avatar {
    background-color: #34495e;
    color: #ecf0f1;
  }

  .user-name {
    margin: 0 8px;
    font-size: 14px;
    color: #ecf0f1; /* Light user name color */
  }

  /* Breadcrumb styling for dark theme */
  :deep(.el-breadcrumb__inner),
  :deep(.el-breadcrumb__separator) {
    color: #bdc3c7 !important; /* Inactive color */
  }

  :deep(.el-breadcrumb__item:last-child .el-breadcrumb__inner) {
    color: #ffffff !important; /* Active color */
    font-weight: 600;
  }

  :deep(.el-breadcrumb__inner:hover) {
    color: #ffffff !important;
  }
</style>
