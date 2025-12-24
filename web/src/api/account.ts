import apiClient from './index';

// 定义创建账号时发送到后端的数据类型
// 注意：这应该与后端 model.CreateAccountRequest 保持一致
export interface CreateAccountPayload {
  username: string;
  password: string;
  realName: string;
  mobile?: string;
}

// 定义获取账号列表时的查询参数
export interface ListAccountsParams {
  page?: number;
  pageSize?: number;
}

// 定义获取账号列表接口的响应体结构
export interface ListAccountsResponse {
  list: any[]; // TODO: 应该定义更精确的 User 类型
  total: number;
  page: number;
  pageSize: number;
}

/**
 * @description 调用后端接口创建子账号
 * @param data 包含新账号信息的载荷
 * @returns Promise
 */
export const createAccountApi = (data: CreateAccountPayload) => {
  return apiClient.post('/accounts', data);
};

/**
 * @description 调用后端接口获取子账号列表
 * @param params 包含分页信息的查询参数
 * @returns Promise<ListAccountsResponse>
 */
export const listAccountsApi = (params: ListAccountsParams): Promise<ListAccountsResponse> => {
  return apiClient.get('/accounts', { params });
};

// 定义更新账号状态时发送到后端的数据类型
export interface UpdateAccountStatusPayload {
  status: 1 | 2; // 1=正常, 2=锁定
}

/**
 * @description 调用后端接口更新子账号状态
 * @param id 账号ID
 * @param payload 包含新状态的载荷
 * @returns Promise
 */
export const updateAccountStatusApi = (id: number, payload: UpdateAccountStatusPayload) => {
  return apiClient.put(`/accounts/${id}/status`, payload);
};

/**
 * @description 调用后端接口删除子账号
 * @param id 账号ID
 * @returns Promise
 */
export const deleteAccountApi = (id: number) => {
  return apiClient.delete(`/accounts/${id}`);
};

// 定义更新账号基本信息时发送到后端的数据类型
export interface UpdateAccountPayload {
  realName: string;
  mobile?: string;
}

/**
 * @description 调用后端接口更新子账号基本信息
 * @param id 账号ID
 * @param payload 包含要更新的信息的载荷
 * @returns Promise
 */
export const updateAccountApi = (id: number, payload: UpdateAccountPayload) => {
  return apiClient.put(`/accounts/${id}`, payload);
};

// 定义重置密码时发送到后端的数据类型
export interface ResetPasswordPayload {
  password: string;
}

/**
 * @description 调用后端接口重置子账号密码
 * @param id 账号ID
 * @param payload 包含新密码的载荷
 * @returns Promise
 */
export const resetPasswordApi = (id: number, payload: ResetPasswordPayload) => {
  return apiClient.put(`/accounts/${id}/password`, payload);
};
