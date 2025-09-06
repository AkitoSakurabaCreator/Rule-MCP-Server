import React, { useState } from 'react';
import {
  CardContent,
  Typography,
  TextField,
  Button,
  Box,
  Alert,
  Link,
  Paper,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../../contexts/AuthContext';

const LoginForm: React.FC = () => {
  const [formData, setFormData] = useState({
    username: '',
    password: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      await login(formData.username, formData.password);
      // ログイン後はトップページ（ルール一覧）に遷移
      navigate('/');
    } catch (err) {
      setError(err instanceof Error ? err.message : t('auth.loginError'));
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        p: 2,
      }}
    >
      <Paper
        elevation={8}
        sx={{
          borderRadius: 3,
          overflow: 'hidden',
          maxWidth: 400,
          width: '100%',
        }}
      >
        <Box
          sx={{
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            p: 3,
            textAlign: 'center',
          }}
        >
          <Typography variant="h4" component="h1" color="white" fontWeight="bold">
            {t('auth.welcomeBack')}
          </Typography>
          <Typography variant="body1" color="white" sx={{ mt: 1, opacity: 0.9 }}>
            {t('auth.signInToContinue')}
          </Typography>
        </Box>

        <CardContent sx={{ p: 4 }}>
          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          )}

          <Box component="form" onSubmit={handleSubmit}>
            <TextField
              fullWidth
              label={t('auth.username')}
              value={formData.username}
              onChange={(e) => handleInputChange('username', e.target.value)}
              required
              margin="normal"
              variant="outlined"
              autoComplete="username"
              autoFocus
            />

            <TextField
              fullWidth
              label={t('auth.password')}
              type="password"
              value={formData.password}
              onChange={(e) => handleInputChange('password', e.target.value)}
              required
              margin="normal"
              variant="outlined"
              autoComplete="current-password"
            />

            <Button
              type="submit"
              fullWidth
              variant="contained"
              size="large"
              disabled={loading}
              sx={{
                mt: 3,
                mb: 2,
                py: 1.5,
                background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                '&:hover': {
                  background: 'linear-gradient(135deg, #5a6fd8 0%, #6a4190 100%)',
                },
              }}
            >
              {loading ? t('auth.signingIn') : t('auth.signIn')}
            </Button>

            <Box sx={{ textAlign: 'center', mt: 2 }}>
              <Typography variant="body2" color="text.secondary">
                {t('auth.noAccount')}{' '}
                <Link
                  component="button"
                  variant="body2"
                  onClick={() => navigate('/auth/register')}
                  sx={{ textDecoration: 'none' }}
                >
                  {t('auth.createAccount')}
                </Link>
              </Typography>
            </Box>

            <Box sx={{ textAlign: 'center', mt: 1 }}>
              <Link
                component="button"
                variant="body2"
                onClick={() => navigate('/auth/forgot-password')}
                sx={{ textDecoration: 'none' }}
              >
                {t('auth.forgotPassword')}
              </Link>
            </Box>
          </Box>
        </CardContent>
      </Paper>
    </Box>
  );
};

export default LoginForm;
