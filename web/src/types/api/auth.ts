// web/src/types/api/auth.ts

// 对应后端的 model.LoginRequest
export interface LoginRequest {
  username: string;
  password: string;
  role: string;
}

// 对应后端的 model.UserInfo
export interface UserInfo {
  id: number;
  username: string;
  real_name: string;
  role: string;
  org_id: number;
}

// 对应后端的 model.LoginResponse
export interface LoginResponse {
  token: string;
  user_info: UserInfo;
}
