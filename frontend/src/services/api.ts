import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:18080/api/v1';

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for logging
api.interceptors.request.use(
  (config) => {
    console.log('API Request:', config.method?.toUpperCase(), config.url);
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    const data = error.response?.data;
    // Backend unified error shape: { code, message, details?, requestId?, timestamp }
    const unified = {
      code: data?.code ?? 'unknown_error',
      message: data?.message ?? error.message ?? 'エラーが発生しました',
      details: data?.details,
      requestId: data?.requestId,
      status: error.response?.status,
    } as const;

    // Optional: console diagnostics for developers
    if (unified.requestId) {
      console.error(`[API Error] ${unified.code}: ${unified.message} (reqId=${unified.requestId})`, unified.details);
    } else {
      console.error(`[API Error] ${unified.code}: ${unified.message}`, unified.details);
    }

    // Attach normalized error for UI layers
    (error as any).normalized = unified;

    // Example UX hooks (callers can use):
    // if (unified.code === 'unauthorized') redirectToLogin();

    return Promise.reject(error);
  }
);
