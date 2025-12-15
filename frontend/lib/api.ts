import axios from 'axios';

// API 基础 URL（从环境变量读取，默认为本地后端地址）
const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';

// 创建 axios 实例
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器：自动添加 JWT Token
apiClient.interceptors.request.use(
  (config) => {
    // 从 localStorage 获取 token（仅在浏览器环境）
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器：统一处理错误
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // 如果是 401 错误，清除 token 并重定向到登录页
    if (error.response?.status === 401 && typeof window !== 'undefined') {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// API 响应类型定义
export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  error?: string;
}

// 用户类型定义
export interface User {
  id: number;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

// 登录响应类型
export interface LoginResponse {
  token: string;
  user: User;
}

// 注册请求
export const register = async (email: string, password: string, name?: string) => {
  const response = await apiClient.post<ApiResponse<User>>('/api/auth/register', {
    email,
    password,
    name,
  });
  return response.data;
};

// 登录请求
export const login = async (email: string, password: string) => {
  const response = await apiClient.post<ApiResponse<LoginResponse>>('/api/auth/login', {
    email,
    password,
  });
  return response.data;
};

// 获取当前用户信息
export const getCurrentUser = async () => {
  const response = await apiClient.get<ApiResponse<User>>('/api/auth/me');
  return response.data;
};

// 健康检查
export const healthCheck = async () => {
  const response = await apiClient.get<ApiResponse>('/health');
  return response.data;
};

export default apiClient;

