import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { api } from '../services/api';

interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  role: 'user' | 'admin';
  is_active: boolean;
}

interface Permissions {
  manageUsers: boolean;
  manageRules: boolean;
  manageRoles: boolean;
}

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  register: (userData: RegisterData) => Promise<void>;
  refreshUser: () => Promise<void>;
  permissions: Permissions;
}

interface RegisterData {
  username: string;
  email: string;
  password: string;
  full_name: string;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const isAuthenticated = !!user;
  const permissions: Permissions = {
    manageUsers: user?.role === 'admin',
    manageRules: user?.role === 'admin' || user?.role === 'user',
    manageRoles: user?.role === 'admin',
  };

  useEffect(() => {
    // 初期化時にユーザー情報を取得
    checkAuthStatus();
  }, []);

  const checkAuthStatus = async () => {
    try {
      const token = localStorage.getItem('auth_token');
      const userData = localStorage.getItem('user');
      
      if (token && userData) {
        // トークンとユーザー情報が存在する場合、復元
        const user = JSON.parse(userData);
        setUser(user);
        api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      }
    } catch (error) {
      // データが無効な場合、削除
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user');
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (username: string, password: string) => {
    try {
      const response = await api.post('/auth/login', { username, password });
      const { token, user: userData } = response.data;
      
      // トークンとユーザー情報を保存
      localStorage.setItem('auth_token', token);
      localStorage.setItem('user', JSON.stringify(userData));
      
      // ユーザー情報を設定
      setUser(userData);
      
      // APIクライアントにトークンを設定
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } catch (error: any) {
      const errorData = error.response?.data;
      const message = errorData?.message || errorData?.error || 'ログインに失敗しました';
      throw new Error(message);
    }
  };

  const logout = async () => {
    try {
      // ログアウトAPIを呼び出し
      await api.post('/auth/logout');
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      // ローカル状態をクリア
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user');
      setUser(null);
      delete api.defaults.headers.common['Authorization'];
    }
  };

  const register = async (userData: RegisterData) => {
    try {
      const response = await api.post('/auth/register', userData);
      const { message, status } = response.data;
      
      // 承認待ち状態の場合は、トークンは保存せずにメッセージを返す
      if (status === 'pending_approval') {
        throw new Error(message);
      }
      
      // 通常の登録完了の場合（現在は使用されない）
      const { token, user: newUser } = response.data;
      localStorage.setItem('auth_token', token);
      setUser(newUser);
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } catch (error: any) {
      const errorData = error.response?.data;
      const message = errorData?.message || errorData?.error || 'アカウント作成に失敗しました';
      throw new Error(message);
    }
  };

  const refreshUser = async () => {
    try {
      const response = await api.get('/auth/me');
      setUser(response.data);
    } catch (error) {
      // ユーザー情報の取得に失敗した場合、ログアウト
      localStorage.removeItem('auth_token');
      setUser(null);
      delete api.defaults.headers.common['Authorization'];
      throw error;
    }
  };

  const value: AuthContextType = {
    user,
    isAuthenticated,
    isLoading,
    login,
    logout,
    register,
    refreshUser,
    permissions,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};
