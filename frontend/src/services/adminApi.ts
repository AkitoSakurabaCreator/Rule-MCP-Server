import { api } from './api';

export interface AdminStats {
  totalUsers: number;
  totalProjects: number;
  totalRules: number;
  activeApiKeys: number;
  mcpRequests: number;
  activeSessions: number;
  systemLoad: string;
}

export interface AdminUser {
  id: number;
  username: string;
  email: string;
  fullName: string;
  role: string;
  isActive: boolean;
  lastLogin: string;
}

export interface AdminApiKey {
  id: number;
  name: string;
  key: string;
  accessLevel: string;
  status: string;
  createdAt: string;
  lastUsed: string;
}

export interface MCPStats {
  method: string;
  count: number;
  lastUsed: string;
  status: string;
}

export interface SystemLog {
  timestamp: string;
  level: string;
  message: string;
}

// 管理者用APIサービス
export const adminApi = {
  // 統計データ取得
  getStats: async (): Promise<AdminStats> => {
    const response = await api.get('/admin/stats');
    return response.data;
  },

  // ユーザー一覧取得
  getUsers: async (): Promise<AdminUser[]> => {
    const response = await api.get('/admin/users');
    return response.data;
  },

  // APIキー一覧取得
  getApiKeys: async (): Promise<AdminApiKey[]> => {
    const response = await api.get('/admin/api-keys');
    return response.data;
  },

  // MCP統計取得
  getMCPStats: async (): Promise<MCPStats[]> => {
    const response = await api.get('/admin/mcp-stats');
    return response.data;
  },

  // システムログ取得
  getSystemLogs: async (): Promise<SystemLog[]> => {
    const response = await api.get('/admin/system-logs');
    return response.data;
  },

  // ユーザー作成
  createUser: async (userData: Partial<AdminUser>): Promise<AdminUser> => {
    const response = await api.post('/admin/users', userData);
    return response.data;
  },

  // ユーザー更新
  updateUser: async (id: number, userData: Partial<AdminUser>): Promise<AdminUser> => {
    const response = await api.put(`/admin/users/${id}`, userData);
    return response.data;
  },

  // ユーザー削除
  deleteUser: async (id: number): Promise<void> => {
    await api.delete(`/admin/users/${id}`);
  },

  // APIキー生成
  generateApiKey: async (keyData: { name: string; accessLevel: string }): Promise<AdminApiKey> => {
    const response = await api.post('/admin/api-keys', keyData);
    return response.data;
  },

  // APIキー削除
  deleteApiKey: async (id: number): Promise<void> => {
    await api.delete(`/admin/api-keys/${id}`);
  },
};
