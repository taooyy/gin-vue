import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { UserInfo } from '@/types/api/auth';
import router from '@/router';

export const useUserStore = defineStore('user', () => {
  // 从 sessionStorage 初始化 state
  const token = ref(sessionStorage.getItem('token') || '');
  const role = ref(sessionStorage.getItem('role') || '');
  const orgType = ref(Number(sessionStorage.getItem('orgType')) || 0);
  const userInfo = ref<UserInfo | null>(JSON.parse(sessionStorage.getItem('userInfo') || 'null'));

  /**
   * @description 设置登录状态，并将信息持久化到 sessionStorage
   * @param newToken 新的 token
   * @param newRole 新的角色
   * @param newType 新的组织类型
   * @param newUserInfo 新的用户信息
   */
  function setLoginState(newToken: string, newRole: string, newType: number, newUserInfo: UserInfo) {
    token.value = newToken;
    role.value = newRole;
    orgType.value = newType;
    userInfo.value = newUserInfo;
    
    sessionStorage.setItem('token', newToken);
    sessionStorage.setItem('role', newRole);
    sessionStorage.setItem('orgType', String(newType));
    sessionStorage.setItem('userInfo', JSON.stringify(newUserInfo));
  }

  /**
   * @description 清除登录状态并重定向到登录页
   */
  function logout() {
    token.value = '';
    role.value = '';
    orgType.value = 0;
    userInfo.value = null;
    sessionStorage.clear();
    router.push('/login');
  }

  return { token, role, orgType, userInfo, setLoginState, logout };
});