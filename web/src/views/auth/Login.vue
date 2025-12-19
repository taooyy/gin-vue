<template>
  <div class="login-container">
    <el-card class="login-box">
      <h2>SaaS 供应链登录</h2>
      <el-form label-position="top">
        <el-form-item label="模拟角色切换">
          <el-select v-model="selectedRole" placeholder="选择模拟登录的角色" style="width:100%">
            <el-option label="平台管理员" value="platform_admin" />
            <el-option label="学校管理员" value="school_admin" />
            <el-option label="供应商" value="supplier" />
          </el-select>
        </el-form-item>
        <el-button type="primary" style="width:100%" @click="handleLogin">登录</el-button>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { usePermissionStore } from '@/stores/permission'

const router = useRouter()
const userStore = useUserStore()
const permissionStore = usePermissionStore()
const selectedRole = ref('school_admin')

const handleLogin = () => {
  // 模拟后端返回
  const mockOrgType = selectedRole.value === 'platform_admin' ? 1 : 2
  const mockToken = 'token-' + Date.now()

  // 1. 存状态
  userStore.setLoginState(mockToken, selectedRole.value, mockOrgType)

  // 2. 核心逻辑：分流 + 生成菜单
  if (mockOrgType === 1) {
    permissionStore.generateMenus(selectedRole.value, 'platform')
    router.push('/platform/dashboard')
  } else {
    permissionStore.generateMenus(selectedRole.value, 'tenant')
    router.push('/workspace/dashboard')
  }
}
</script>
<style scoped>
.login-container { height: 100vh; display: flex; justify-content: center; align-items: center; background: #eee; }
.login-box { width: 400px; }
</style>