import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider as MuiThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { Container } from '@mui/material';
import { useTranslation } from 'react-i18next';

import './i18n';
import { ThemeProvider, useTheme } from './contexts/ThemeContext';
import Header from './components/Header';
import ProjectList from './components/ProjectList';
import ProjectForm from './components/ProjectForm';
import RuleList from './components/RuleList';
import RuleEdit from './components/RuleEdit';
import GlobalRuleForm from './components/GlobalRuleForm';
import CodeValidator from './components/CodeValidator';

const AppContent: React.FC = () => {
  const { i18n } = useTranslation();
  const { themeMode } = useTheme();
  
  // RTL言語のサポート
  const isRTL = ['ar', 'he', 'fa'].includes(i18n.language);
  
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

  return (
    <MuiThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <Header />
        <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
          <Routes>
            <Route path="/" element={<ProjectList />} />
            <Route path="/projects/new" element={<ProjectForm />} />
            <Route path="/projects/:projectId/edit" element={<ProjectForm />} />
            <Route path="/projects/:projectId/rules" element={<RuleList />} />
            <Route path="/projects/:projectId/rules/new" element={<RuleEdit />} />
            <Route path="/projects/:projectId/rules/:ruleId/edit" element={<RuleEdit />} />
            <Route path="/global-rules" element={<GlobalRuleForm />} />
            <Route path="/validate" element={<CodeValidator />} />
          </Routes>
        </Container>
      </Router>
    </MuiThemeProvider>
  );
};

const App: React.FC = () => {
  return (
    <ThemeProvider>
      <AppContent />
    </ThemeProvider>
  );
};

export default App;
