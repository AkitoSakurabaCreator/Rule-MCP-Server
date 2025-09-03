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

    interpolation: {
      escapeValue: false,
    },

    detection: {
      order: ['localStorage', 'navigator'],
      caches: ['localStorage'],
    },
  });

export default i18n;
