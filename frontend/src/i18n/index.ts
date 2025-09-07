import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import en from './locales/en.json';
import ja from './locales/ja.json';
import zhCN from './locales/zh-CN.json';
import hi from './locales/hi.json';
import es from './locales/es.json';
import ar from './locales/ar.json';

const resources = {
  en: {
    translation: en,
  },
  ja: {
    translation: ja,
  },
  'zh-CN': {
    translation: zhCN,
  },
  hi: {
    translation: hi,
  },
  es: {
    translation: es,
  },
  ar: {
    translation: ar,
  },
};

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'en',
    debug: process.env.NODE_ENV === 'development',
    
    // サポートする言語の制限（トップレベルで設定）
    supportedLngs: ['en', 'ja', 'zh-CN', 'hi', 'es', 'ar'],
    
    // 言語コードの正規化（トップレベルで設定）
    load: 'languageOnly',
    
    // 開発環境でのデバッグ情報
    ...(process.env.NODE_ENV === 'development' && {
      saveMissing: true,
      missingKeyHandler: (lngs: readonly string[], ns: string, key: string, fallbackValue: string, updateMissing: boolean, options: any) => {
        console.warn(`Missing translation: ${lngs.join(',')}.${ns}.${key}`);
      },
    }),

    interpolation: {
      escapeValue: false,
    },

    detection: {
      // 検出の優先順位: localStorage > URL > HTMLタグ > navigator > fallback
      order: ['localStorage', 'querystring', 'htmlTag', 'navigator', 'path', 'subdomain'],
      caches: ['localStorage'],
      
      // ブラウザの言語コードとアプリの言語コードのマッピング
      lookupLocalStorage: 'i18nextLng',
      lookupQuerystring: 'lng',
      lookupFromPathIndex: 0,
      lookupFromSubdomainIndex: 0,
      
      // ブラウザの言語コードをアプリの言語コードにマッピング
      convertDetectedLanguage: (lng: string) => {
        // ブラウザの言語コードをアプリでサポートする言語コードに変換
        const languageMap: { [key: string]: string } = {
          'ja': 'ja',
          'ja-JP': 'ja',
          'en': 'en',
          'en-US': 'en',
          'en-GB': 'en',
          'zh': 'zh-CN',
          'zh-CN': 'zh-CN',
          'zh-Hans': 'zh-CN',
          'zh-Hans-CN': 'zh-CN',
          'hi': 'hi',
          'hi-IN': 'hi',
          'es': 'es',
          'es-ES': 'es',
          'es-MX': 'es',
          'ar': 'ar',
          'ar-SA': 'ar',
          'ar-AE': 'ar',
        };
        
        // 完全一致をチェック
        if (languageMap[lng]) {
          return languageMap[lng];
        }
        
        // 言語コードのみをチェック（例: 'ja-JP' -> 'ja'）
        const baseLang = lng.split('-')[0];
        if (languageMap[baseLang]) {
          return languageMap[baseLang];
        }
        
        // デフォルトは英語
        return 'en';
      },
    },
  });

export default i18n;
