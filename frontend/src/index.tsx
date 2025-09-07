import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

// アプリのパフォーマンス測定を開始したい場合は、結果をログに記録する関数を渡してください
// （例: reportWebVitals(console.log)）
// または分析エンドポイントに送信してください。詳細: https://bit.ly/CRA-vitals
reportWebVitals();
