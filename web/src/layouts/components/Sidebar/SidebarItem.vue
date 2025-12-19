<template>
  <div v-if="!item.children || item.children.length === 0">
    <el-menu-item :index="item.path">
      <el-icon v-if="item.icon"><component :is="item.icon" /></el-icon>
      <span>{{ item.title }}</span>
    </el-menu-item>
  </div>
  <div v-else>
    <el-sub-menu :index="item.path">
      <template #title>
        <el-icon v-if="item.icon"><component :is="item.icon" /></el-icon>
        <span>{{ item.title }}</span>
      </template>
      <sidebar-item v-for="child in item.children" :key="child.path" :item="child" />
    </el-sub-menu>
  </div>
</template>
<script setup lang="ts">
import type { MenuItem } from '@/config/menuMap'
defineProps<{ item: MenuItem }>()
</script>