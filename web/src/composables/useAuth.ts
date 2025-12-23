// web/src/composables/useAuth.ts
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';
import { usePermissionStore } from '@/stores/permission';
import { loginApi } from '@/api/auth';
import { ElMessage } from 'element-plus';

/**
 * @description 处理用户认证（登录、登出）相关逻辑的 Composable
 */
export function useAuth() {
  const router = useRouter();
  const userStore = useUserStore();
  const permissionStore = usePermissionStore();

  const username = ref('admin'); // 默认用户名
  const password = ref('password123'); // 默认密码
  const selectedRole = ref('school_admin'); // 登录表单中选择的角色

  /**
   * @description 处理登录操作，调用真实后端 API
   */
  const handleLogin = async () => {
    try {
      // 实际项目中，用户名和密码应该由用户输入
      // 这里我们为了方便，暂时使用角色名作为用户名
      const loginData = {
        username: selectedRole.value, 
        password: password.value, // 使用一个固定的密码
        role: selectedRole.value,
      };

      const response = await loginApi(loginData);
      const { token, user_info } = response;

      // 1. 调用 user store 保存登录状态
      // 后端返回的用户信息中的 role 可能与登录时选择的 role 不完全一致，以后端为准
      const orgType = user_info.role === 'platform_admin' ? 1 : 2;
      userStore.setLoginState(token, user_info.role, orgType, user_info);

      // 2. 根据组织类型进行分流，并生成对应的菜单
      if (orgType === 1) { // 平台侧
        permissionStore.generateMenus(user_info.role, 'platform');
        router.push('/platform/dashboard');
      } else { // 租户侧
        permissionStore.generateMenus(user_info.role, 'tenant');
        router.push('/workspace/dashboard');
      }
       ElMessage.success('登录成功！');

    } catch (error) {
      console.error('登录失败:', error);
      // API 拦截器中已经处理了错误消息提示，这里可以不用重复提示
    }
  };

  // 将组件模板需要用到的状态和方法暴露出去
  return {
    username,
    password,
    selectedRole,
    handleLogin,
  };
}
