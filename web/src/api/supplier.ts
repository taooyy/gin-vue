import apiClient from './index';

// 对应后端的 SysOrganization 结构
export interface Supplier {
  ID: number;
  Name: string;
  OrgType: number;
  ParentID: number;
  AdminUserID: number;
  ContactName: string;
  ContactPhone: string;
  Address: string;
  IsEnabled: boolean;
  CreatedAt: string;
  UpdatedAt: string;
}

// 获取供应商详情时，可能包含关联的用户信息
export interface SupplierDetails extends Supplier {
  adminUser?: {
    Username: string;
    RealName: string;
  };
}

// 对应后端 service 层的 CreateSupplierRequest 结构
export interface CreateSupplierPayload {
  name: string;
  contactName: string;
  contactPhone: string;
  address: string;
  username: string;
  password: string;
  realName: string;
}

export interface ListSuppliersResponse {
  items: Supplier[];
  total: number;
}

/**
 * 获取供应商列表
 * @param params 分页参数
 */
export function listSuppliersApi(params: {
  page: number;
  pageSize: number;
}): Promise<ListSuppliersResponse> {
  return apiClient.get('/suppliers', { params });
}

/**
 * 创建一个新的供应商及其管理员
 * @param data 供应商和管理员信息
 */
export function createSupplierApi(data: CreateSupplierPayload): Promise<any> {
  return apiClient.post('/suppliers', data);
}

/**
 * 根据ID获取供应商详细信息
 * @param id 供应商ID
 */
export function getSupplierApi(id: number): Promise<SupplierDetails> {
  return apiClient.get(`/suppliers/${id}`);
}

/**
 * 更新供应商信息
 * @param id 供应商ID
 * @param data 要更新的数据
 */
export function updateSupplierApi(
  id: number,
  data: Partial<CreateSupplierPayload>
): Promise<any> {
  return apiClient.put(`/suppliers/${id}`, data);
}

/**
 * 更新供应商状态
 * @param id 供应商ID
 * @param isEnabled 启用状态
 */
export function updateSupplierStatusApi(id: number, isEnabled: boolean): Promise<any> {
  return apiClient.put(`/suppliers/${id}/status`, { isEnabled });
}
