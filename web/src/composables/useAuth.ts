// web/src/composables/useAuth.ts
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { usePermissionStore } from '@/stores/permission'

/**
 * @description 处理用户认证（登录、登出）相关逻辑的 Composable
 */
export function useAuth() {
  const router = useRouter()
  const userStore = useUserStore()
  const permissionStore = usePermissionStore()
  
  // 登录表单中选择的角色，默认为 'school_admin'
  const selectedRole = ref('school_admin')

  /**
   * @description 处理登录操作
   * 当前为模拟登录，未来需替换为真实的 API 调用
   */
  const handleLogin = () => {
    // 模拟后端根据角色返回不同的组织类型
    const mockOrgType = selectedRole.value === 'platform_admin' ? 1 : 2
    const mockToken = 'mock-token-' + Date.now()

    // 1. 调用 user store 保存登录状态 (token, role, orgType)
    userStore.setLoginState(mockToken, selectedRole.value, mockOrgType)

    // 2. 核心业务：根据组织类型进行分流，并生成对应的菜单
    if (mockOrgType === 1) { // 平台侧
      permissionStore.generateMenus(selectedRole.value, 'platform')
      router.push('/platform/dashboard')
    } else { // 租户侧
      permissionStore.generateMenus(selectedRole.value, 'tenant')
      router.push('/workspace/dashboard')
    }
  }
  
  // 将组件模板需要用到的状态和方法暴露出去
  return {
    selectedRole,
    handleLogin
  }
}
