import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:18080/api/v1';

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// ログ用のリクエストインターセプター
api.interceptors.request.use(
  (config) => {
    console.log('API Request:', config.method?.toUpperCase(), config.url);
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// エラーハンドリング用のレスポンスインターセプター
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    const data = error.response?.data;
    // バックエンド統一エラー形式: { code, message, details?, requestId?, timestamp }
    const unified = {
      code: data?.code ?? 'unknown_error',
      message: data?.message ?? error.message ?? 'エラーが発生しました',
      details: data?.details,
      requestId: data?.requestId,
      status: error.response?.status,
    } as const;

    // オプション: 開発者向けコンソール診断
    if (unified.requestId) {
      console.error(`[API Error] ${unified.code}: ${unified.message} (reqId=${unified.requestId})`, unified.details);
    } else {
      console.error(`[API Error] ${unified.code}: ${unified.message}`, unified.details);
    }

    // UI層用に正規化されたエラーを添付
    (error as any).normalized = unified;

    // 認証エラーの場合は自動的にログイン画面にリダイレクト
    if (unified.status === 401 || unified.status === 403) {
      console.warn('Authentication error detected, redirecting to login...');
      // localStorageをクリア
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user');
      // ログイン画面にリダイレクト（現在のページがログイン画面でない場合のみ）
      if (window.location.pathname !== '/login') {
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  }
);
