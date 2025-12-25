<template>
  <div class="platform-layout">
    <div class="sidebar-container" :style="{ width: sidebarWidth }">
      <div class="logo-container">
        <el-icon :size="32" color="#2ecc71"><ElementPlus /></el-icon>
        <transition name="fade">
          <span v-if="!appStore.isSidebarCollapsed" class="logo-title">SaaS 总控台</span>
        </transition>
      </div>
      <Sidebar />
    </div>
    <div class="content-wrapper" :style="{ marginLeft: sidebarWidth }">
      <Header />
      <div class="main-container">
        <router-view />
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
  import { computed } from 'vue';
  import Sidebar from './components/Sidebar/index.vue';
  import Header from './components/Header/index.vue';
  import { useAppStore } from '@/stores/app';

  const appStore = useAppStore();
  const sidebarWidth = computed(() => (appStore.isSidebarCollapsed ? '65px' : '210px'));
</script>
<style scoped>
  /* 整个平台布局的基础容器 */
  .platform-layout {
    height: 100vh;
    width: 100%;
    /* overflow-x: hidden; 防止内容过宽时出现水平滚动条 */
  }

  /* 右侧内容包裹器 */
  .content-wrapper {
    display: flex;
    flex-direction: column;
    height: 100%;
    /* 通过动态的 margin-left 来为左侧的固定侧边栏留出空间 */
    transition: margin-left 0.3s ease;
  }

  /* 左侧侧边栏容器 */
  .sidebar-container {
    /* 使用固定定位，使其脱离文档流，固定在屏幕左侧 */
    position: fixed;
    top: 0;
    left: 0;
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: #2a384a;
    box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);
    z-index: 10; /* 确保在最上层 */
    /* 宽度变化的过渡动画 */
    transition: width 0.3s ease;
    overflow-x: hidden; /* 收起时隐藏内部多余内容 */
  }
  .logo-container {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
    height: 60px;
    flex-shrink: 0; /* 防止 logo 容器被压缩 */
  }
  .logo-title {
    color: white;
    font-size: 1.2rem;
    font-weight: bold;
    white-space: nowrap; /* 防止标题换行 */
  }

  /* 主内容区域，包含 router-view */
  .main-container {
    padding: 20px;
    background: #f0f2f5;
    flex: 1; /* 占据 content-wrapper 中除了 Header 外的所有剩余空间 */
    overflow-y: auto; /* 当内容超长时，仅在此区域出现垂直滚动条 */
  }

  /* logo标题的淡入淡出动画 */
  .fade-enter-active,
  .fade-leave-active {
    transition: opacity 0.2s ease;
  }
  .fade-enter-from,
  .fade-leave-to {
    opacity: 0;
  }
</style>
