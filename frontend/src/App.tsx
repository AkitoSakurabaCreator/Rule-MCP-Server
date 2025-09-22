import React, { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider as MuiThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { Container, Box, CircularProgress } from '@mui/material';
import { useTranslation } from 'react-i18next';

import './i18n';
import { ThemeProvider, useTheme } from './contexts/ThemeContext';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import Header from './components/Header';
import ProjectList from './components/ProjectList';
import ProjectForm from './components/ProjectForm';
import RuleList from './components/RuleList';
import RuleEdit from './components/RuleEdit';
import GlobalRuleForm from './components/GlobalRuleForm';
import CodeValidator from './components/CodeValidator';
import LoginForm from './components/auth/LoginForm';
import RegisterForm from './components/auth/RegisterForm';
import AdminDashboard from './components/dashboard/AdminDashboard';
import UserProfile from './components/UserProfile';
import LanguageManagement from './components/LanguageManagement';

const AppContent: React.FC = () => {
  const { i18n } = useTranslation();
  const { themeMode } = useTheme();
  const { isAuthenticated, user, isLoading } = useAuth();
  
  // RTL言語のサポート
  const isRTL = ['ar', 'he', 'fa'].includes(i18n.language);
  
  // 言語変更時にHTMLタグの言語属性を更新
  useEffect(() => {
    const updateHtmlLang = () => {
      const htmlElement = document.documentElement;
      htmlElement.setAttribute('lang', i18n.language);
      htmlElement.setAttribute('dir', isRTL ? 'rtl' : 'ltr');
    };
    
    // 初期設定
    updateHtmlLang();
    
    // 言語変更時のイベントリスナー
    i18n.on('languageChanged', updateHtmlLang);
    
    // クリーンアップ
    return () => {
      i18n.off('languageChanged', updateHtmlLang);
    };
  }, [i18n, isRTL]);
  
  const theme = createTheme({
    direction: isRTL ? 'rtl' : 'ltr',
    palette: {
      mode: themeMode,
      primary: {
        main: themeMode === 'dark' ? '#90caf9' : '#1976d2',
      },
      secondary: {
        main: themeMode === 'dark' ? '#f48fb1' : '#dc004e',
      },
      background: {
        default: themeMode === 'dark' ? '#121212' : '#f5f5f5',
        paper: themeMode === 'dark' ? '#1e1e1e' : '#ffffff',
      },
      text: {
        primary: themeMode === 'dark' ? '#ffffff' : '#000000',
        secondary: themeMode === 'dark' ? '#b0b0b0' : '#666666',
      },
    },
    typography: {
      fontFamily: isRTL 
        ? '"Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif'
        : '"Roboto", "Helvetica", "Arial", sans-serif',
    },
    components: {
      MuiAppBar: {
        styleOverrides: {
          root: {
            backgroundColor: themeMode === 'dark' ? '#1e1e1e' : '#1976d2',
          },
        },
      },
      MuiCard: {
        styleOverrides: {
          root: {
            backgroundColor: themeMode === 'dark' ? '#1e1e1e' : '#ffffff',
            border: themeMode === 'dark' ? '1px solid #333' : '1px solid #e0e0e0',
          },
        },
      },
    },
  });

  // ローディング中は何も表示しない
  if (isLoading) {
    return (
      <MuiThemeProvider theme={theme}>
        <CssBaseline />
        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
          <CircularProgress />
        </Box>
      </MuiThemeProvider>
    );
  }

  return (
    <MuiThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <Routes>
          {/* 認証ルート - Headerなし */}
          <Route path="/auth/login" element={<LoginForm />} />
          <Route path="/auth/register" element={<RegisterForm />} />
          
          {/* 管理者ルート */}
          {user?.role === 'admin' && (
            <Route path="/admin/dashboard" element={
              <>
                <Header />
                <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                  <AdminDashboard />
                </Container>
              </>
            } />
          )}
          
          {/* 保護されたルート */}
          {isAuthenticated ? (
            <>
              <Route path="/" element={
                <>
                  <Header />
                  <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                    <ProjectList />
                  </Container>
                </>
              } />
              <Route path="/projects/new" element={
                <>
                  <Header />
                  <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                    <ProjectForm />
                  </Container>
                </>
              } />
              <Route path="/projects/:projectId/edit" element={
                <>
                  <Header />
                  <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                    <ProjectForm />
                  </Container>
                </>
              } />
              <Route path="/projects/:projectId/rules" element={
                <>
                  <Header />
                  <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                    <RuleList />
                  </Container>
                </>
              } />
              <Route path="/projects/:projectId/rules/new" element={
                <>
                  <Header />
                  <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                    <RuleEdit />
                  </Container>
                </>
              } />
              <Route path="/projects/:projectId/rules/:ruleId/edit" element={
                <>
                  <Header />
                  <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                    <RuleEdit />
                  </Container>
                </>
              } />
              <Route path="/global-rules" element={
                <>
                  <Header />
                  <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                    <GlobalRuleForm />
                  </Container>
                </>
              } />
              <Route path="/validate" element={
                <>
                  <Header />
                  <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                    <CodeValidator />
                  </Container>
                </>
              } />
              <Route path="/profile" element={
                <>
                  <Header />
                  <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                    <UserProfile />
                  </Container>
                </>
              } />
              <Route path="/languages" element={
                <>
                  <Header />
                  <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                    <LanguageManagement />
                  </Container>
                </>
              } />
            </>
          ) : (
            <Route path="*" element={<LoginForm />} />
          )}
        </Routes>
      </Router>
    </MuiThemeProvider>
  );
};

const App: React.FC = () => {
  return (
    <ThemeProvider>
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </ThemeProvider>
  );
};

export default App;
