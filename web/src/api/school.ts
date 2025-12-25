// web/src/api/school.ts
import apiClient from './index';

export interface School {
  ID: number;
  Name: string;
  OrgType: number;
  ContactName: string;
  ContactPhone: string;
  Address: string;
  IsEnabled: boolean;
  AdminUserID: number;
  AdminUsername: string;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface CreateSchoolPayload {
  name: string;
  contactName: string;
  contactPhone: string;
  address: string;
  adminUsername: string;
  adminPassword: string;
  adminRealName: string;
}

export interface UpdateSchoolPayload {
  Name: string;
  ContactName: string;
  ContactPhone: string;
  Address: string;
  IsEnabled: boolean;
}

export interface ListSchoolsResponse {
  list: School[];
  total: number;
}

/**
 * 获取学校/站点列表
 */
export function listSchoolsApi(params: {
  page: number;
  pageSize: number;
}): Promise<ListSchoolsResponse> {
  return apiClient.get('/schools', { params });
}

/**
 * 创建一个新的学校/站点
 */
export function createSchoolApi(data: CreateSchoolPayload): Promise<any> {
  return apiClient.post('/schools', data);
}

/**
 * 更新一个学校/站点的信息
 */
export function updateSchoolApi(id: number, data: UpdateSchoolPayload): Promise<any> {
  return apiClient.put(`/schools/${id}`, data);
}

/**
 * 删除一个学校/站点
 */
export function deleteSchoolApi(id: number): Promise<any> {
  return apiClient.delete(`/schools/${id}`);
}
