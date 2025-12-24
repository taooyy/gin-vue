<template>
  <div style="display: flex; height: 100vh">
    <div class="sidebar-container" :style="{ width: sidebarWidth }">
      <div class="logo-container">
        <el-icon :size="30" color="#2ecc71"><School /></el-icon>
        <transition name="fade">
          <span v-if="!appStore.isSidebarCollapsed" class="logo-title">智慧校园</span>
        </transition>
      </div>
      <Sidebar />
    </div>
    <div style="flex: 1; display: flex; flex-direction: column">
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
  .sidebar-container {
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
