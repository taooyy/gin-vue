// web/src/api/log.ts
import apiClient from './index';

export interface OpLog {
  ID: number;
  UserID: number;
  Username: string;
  OrgID: number;
  Module: string;
  Action: string;
  Params: string;
  CreatedAt: string;
}

export interface ListLogsParams {
  page: number;
  pageSize: number;
  orgId?: number;
}

export interface ListLogsResponse {
  list: OpLog[];
  total: number;
}

/**
 * 获取操作日志列表
 */
export function listLogsApi(params: ListLogsParams): Promise<ListLogsResponse> {
  return apiClient.get('/logs', { params });
}
