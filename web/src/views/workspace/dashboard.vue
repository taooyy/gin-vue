<template>
  <el-card>
    <h2>欢迎回来！</h2>
    <p>当前角色: <el-tag>{{ userStore.role }}</el-tag></p>
  </el-card>
</template>
<script setup lang="ts">
import { useUserStore } from '@/stores/user'
const userStore = useUserStore()
</script>