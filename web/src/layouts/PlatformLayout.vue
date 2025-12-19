<template>
  <div style="display: flex; height: 100vh;">
    <div style="width: 200px; background: #1f2d3d;">
      <div style="color:white; height:60px; line-height:60px; text-align:center; font-weight:bold;">🔧 平台总控</div>
      <Sidebar />
    </div>
    <div style="flex: 1; padding: 20px;">
      <div style="margin-bottom:20px; text-align:right;">
        <el-button type="danger" @click="logout">退出登录</el-button>
      </div>
      <router-view />
    </div>
  </div>
</template>
<script setup lang="ts">
import Sidebar from './components/Sidebar/index.vue'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
const userStore = useUserStore()
const router = useRouter()
const logout = () => { userStore.logout(); router.push('/login') }
</script>