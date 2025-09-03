import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
} from '@mui/material';
import {
  Home as HomeIcon,
  Language as LanguageIcon,
  Code as CodeIcon,
  Dashboard as DashboardIcon,
  Logout as LogoutIcon,
  Login as LoginIcon,
  Person as PersonIcon,
} from '@mui/icons-material';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../contexts/AuthContext';
import LanguageSwitcher from './LanguageSwitcher';
import ThemeToggle from './ThemeToggle';

const Header: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { t } = useTranslation();
  const { isAuthenticated, user, logout } = useAuth();

  // 現在のページがホームページかどうかを判定
  const isHomePage = location.pathname === '/';

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          Rule MCP Server
        </Typography>

        <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
          {/* ホームページ以外にいる場合のみホームボタンを表示 */}
          {!isHomePage && (
            <Button
              color="inherit"
              startIcon={<HomeIcon />}
              onClick={() => navigate('/')}
            >
              {t('navigation.home')}
            </Button>
          )}

          <Button
            color="inherit"
            startIcon={<LanguageIcon />}
            onClick={() => navigate('/global-rules')}
          >
            {t('navigation.globalRules')}
          </Button>

          <Button
            color="inherit"
            startIcon={<CodeIcon />}
            onClick={() => navigate('/validate')}
          >
            {t('navigation.validateCode')}
          </Button>

          {/* 認証状態に応じたボタン */}
          {isAuthenticated ? (
            <>
              {user?.role === 'admin' && (
                <Button
                  color="inherit"
                  startIcon={<DashboardIcon />}
                  onClick={() => navigate('/admin/dashboard')}
                >
                  {t('navigation.adminDashboard')}
                </Button>
              )}
              
              <Button
                color="inherit"
                startIcon={<PersonIcon />}
                onClick={() => navigate('/profile')}
              >
                {user?.username}
              </Button>
              
              <Button
                color="inherit"
                startIcon={<LogoutIcon />}
                onClick={logout}
              >
                {t('navigation.logout')}
              </Button>
            </>
          ) : (
            <Button
              color="inherit"
              startIcon={<LoginIcon />}
              onClick={() => navigate('/auth/login')}
            >
              {t('navigation.login')}
            </Button>
          )}

          <ThemeToggle />
          <LanguageSwitcher />
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
