// web/src/api/index.ts
import axios from 'axios';
import { useUserStore } from '@/stores/user';
import { ElMessage } from 'element-plus';

// 创建 axios 实例
const apiClient = axios.create({
  // baseURL 将从 .env 文件中根据当前环境动态加载
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000, // 请求超时时间
});

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    const userStore = useUserStore();
    // 如果存在 token，则在请求头中添加 Authorization
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`;
    }
    return config;
  },
  (error) => {
    // 对请求错误做些什么
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  // 对于 2xx 范围内的状态码，直接返回响应数据
  (response) => {
    return response.data;
  },
  // 对于超出 2xx 范围的状态码
  (error) => {
    // 处理网络错误
    if (!error.response) {
      ElMessage.error('网络错误，请检查您的连接');
      return Promise.reject(error);
    }
    
    // 处理后端返回的错误
    const { status, data } = error.response;
    let message = `请求失败，状态码：${status}`;
    if (data && data.error) {
      message = data.error;
    }

    // 特殊处理 401 Unauthorized
    if (status === 401) {
        message = '认证失败或 Token 已过期，请重新登录';
        // 在这里可以触发登出逻辑
        const userStore = useUserStore();
        userStore.logout();
    }
    
    ElMessage.error(message);
    
    return Promise.reject(error);
  }
);

export default apiClient;
