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

  const username = ref(''); // 用户名
  const password = ref(''); // 密码

  /**
   * @description 处理登录操作，调用真实后端 API
   */
  const handleLogin = async () => {
    try {
      const loginData = {
        username: username.value,
        password: password.value,
      };

      // 前端简单校验
      if (!loginData.username || !loginData.password) {
        ElMessage.warning('请输入用户名和密码');
        return;
      }

      const response = await loginApi(loginData);
      const { token, user_info } = response;

      // 1. 调用 user store 保存登录状态
      // 后端返回的用户信息中的 role 是权威信息
      const orgType = user_info.role === 'platform_admin' ? 1 : 2; // 简单判断布局类型
      userStore.setLoginState(token, user_info.role, orgType, user_info);

      // 2. 根据角色信息生成菜单并跳转
      // 注意：这里仍然可以根据角色分流，但更健壮的逻辑已在 router.beforeEach 中处理
      if (user_info.role.startsWith('platform')) {
        // 平台侧
        permissionStore.generateMenus(user_info.role, 'platform');
        router.push('/platform/dashboard');
      } else {
        // 租户侧
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
    handleLogin,
  };
}
