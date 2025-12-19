<template>
  <div style="display: flex; height: 100vh;">
    <div style="width: 210px; background: #304156;">
      <div style="color:white; height:50px; line-height:50px; text-align:center; font-weight:bold;">🏫 智慧校园</div>
      <Sidebar />
    </div>
    <div style="flex: 1; display: flex; flex-direction: column;">
      <div style="height: 50px; border-bottom: 1px solid #eee; display:flex; align-items:center; justify-content:flex-end; padding:0 20px;">
        <span style="margin-right:15px">当前角色: {{ userStore.role }}</span>
        <el-button link type="danger" @click="logout">退出</el-button>
      </div>
      <div style="padding: 20px; background: #f0f2f5; flex: 1;">
        <router-view />
      </div>
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