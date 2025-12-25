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
  .platform-layout {
    height: 100vh;
    width: 100%;
    overflow-x: hidden;
  }

  .content-wrapper {
    display: flex;
    flex-direction: column;
    height: 100%;
    transition: margin-left 0.3s ease;
  }

  .sidebar-container {
    position: fixed;
    top: 0;
    left: 0;
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: #2a384a;
    box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);
    z-index: 10;
    transition: width 0.3s ease;
    overflow-x: hidden;
  }
  .logo-container {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
    height: 60px;
    flex-shrink: 0;
  }
  .logo-title {
    color: white;
    font-size: 1.2rem;
    font-weight: bold;
    white-space: nowrap;
  }
  .main-container {
    padding: 20px;
    background: #f0f2f5;
    flex: 1;
    overflow-y: auto;
  }

  .fade-enter-active,
  .fade-leave-active {
    transition: opacity 0.2s ease;
  }
  .fade-enter-from,
  .fade-leave-to {
    opacity: 0;
  }
</style>
