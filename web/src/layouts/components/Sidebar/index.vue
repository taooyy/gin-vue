<template>
  <el-menu
    :default-active="activeMenu"
    background-color="#304156"
    text-color="#bfcbd9"
    active-text-color="#409EFF"
    router
    style="height:100vh; border-right:none;"
  >
    <sidebar-item v-for="menu in permissionStore.menus" :key="menu.path" :item="menu" />
  </el-menu>
</template>
<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { usePermissionStore } from '@/stores/permission'
import SidebarItem from './SidebarItem.vue'
const route = useRoute()
const permissionStore = usePermissionStore()
const activeMenu = computed(() => route.path)
</script>