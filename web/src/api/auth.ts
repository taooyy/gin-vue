// web/src/api/auth.ts
import apiClient from './index';
import type { LoginRequest, LoginResponse } from '@/types/api/auth';

/**
 * @description 调用后端登录接口
 * @param data 包含用户名和密码的登录信息
 * @returns Promise<LoginResponse>
 */
export const loginApi = (data: LoginRequest): Promise<LoginResponse> => {
  return apiClient.post('/system/login', data);
};
