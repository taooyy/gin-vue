<template>
  <div class="login-container">
    <div class="login-left">
      <div class="brand-wrapper">
        <div class="brand-logo">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="logo-svg">
            <path d="M12 2L1 9l4 2.5V17h14v-5.5L23 9zM12 4.5l7 4.37V15H5V8.87zM7 18v-5h10v5z"/>
          </svg>
        </div>
        <h1 class="brand-title">SaaS Supply Chain</h1>
        <p class="brand-subtitle">赋能智慧校园，链接高效供应链</p>
      </div>
    </div>
    <div class="login-right">
      <el-card class="login-box">
        <h2 class="login-title">欢迎回来</h2>
        <p class="login-subtitle">请选择您的角色以继续</p>
        <el-form label-position="top" class="login-form">
          <el-form-item label="用户名">
             <el-input placeholder="（模拟登录无需输入）" prefix-icon="User" disabled />
          </el-form-item>
           <el-form-item label="密码">
             <el-input placeholder="（模拟登录无需输入）" type="password" prefix-icon="Lock" disabled />
          </el-form-item>
          <el-form-item label="模拟角色切换">
            <el-select v-model="selectedRole" placeholder="选择模拟登录的角色" style="width:100%">
              <el-option label="平台管理员" value="platform_admin" />
              <el-option label="学校管理员" value="school_admin" />
              <el-option label="供应商" value="supplier" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" class="login-button" @click="handleLogin">登 录</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '@/composables/useAuth'

const { selectedRole, handleLogin } = useAuth()
</script>

<style scoped>
.login-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}

.login-left {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #409EFF, #67C23A);
  color: white;
  padding: 40px;
  text-align: center;
}

.brand-wrapper {
  max-width: 450px;
}

.brand-logo {
  margin-bottom: 20px;
}

.logo-svg {
  width: 80px;
  height: 80px;
}

.brand-title {
  font-size: 2.5rem;
  font-weight: bold;
  margin: 0 0 10px;
}

.brand-subtitle {
  font-size: 1.2rem;
  opacity: 0.9;
}

.login-right {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
}

.login-box {
  width: 400px;
  padding: 30px;
  border: none;
  box-shadow: 0 10px 30px rgba(0,0,0,0.1);
  border-radius: 10px;
}

.login-title {
  font-size: 1.8rem;
  font-weight: 600;
  color: #333;
  text-align: center;
  margin-bottom: 10px;
}

.login-subtitle {
  text-align: center;
  color: #999;
  margin-bottom: 30px;
}

.login-form {
  margin-top: 20px;
}

.login-button {
  width: 100%;
  height: 40px;
  font-size: 1rem;
}

:deep(.el-form-item__label) {
  color: #666;
  font-weight: 500;
}
</style>