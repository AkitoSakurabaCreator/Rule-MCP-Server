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
  Stepper,
  Step,
  StepLabel,
  CircularProgress,
  Chip,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../../contexts/AuthContext';

const steps = ['account', 'personal', 'confirm'];

const RegisterForm: React.FC = () => {
  const [activeStep, setActiveStep] = useState(0);
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    full_name: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [pendingApproval, setPendingApproval] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<{[key: string]: string}>({});
  
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { register } = useAuth();

  const handleNext = () => {
    if (activeStep === steps.length - 1) {
      handleSubmit();
    } else {
      setActiveStep((prevActiveStep) => prevActiveStep + 1);
    }
  };

  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };

  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    // Clear field error when user starts typing
    if (fieldErrors[field]) {
      setFieldErrors(prev => ({ ...prev, [field]: '' }));
    }
  };

  const validateStep = (step: number): boolean => {
    switch (step) {
      case 0: // account
        return formData.username.length >= 3 && 
               formData.email.includes('@') && 
               formData.password.length >= 8;
      case 1: // personal
        return formData.full_name.length >= 2;
      case 2: // confirm
        return formData.password === formData.confirmPassword;
      default:
        return false;
    }
  };

  const parseErrorResponse = (error: any): { message: string; field?: string; code?: string } => {
    if (error?.response?.data) {
      const data = error.response.data;
      return {
        message: data.message || t('auth.registerError'),
        field: data.field,
        code: data.code
      };
    }
    return {
      message: error instanceof Error ? error.message : t('auth.registerError')
    };
  };

  const getErrorMessage = (errorInfo: { message: string; field?: string; code?: string }): string => {
    const { message, code } = errorInfo;
    
    // Handle specific error codes
    switch (code) {
      case 'validation_error':
        if (message.includes('Username already exists') || message.includes('ユーザー名は既に存在します')) {
          return t('auth.usernameExists');
        }
        if (message.includes('Email already exists') || message.includes('メールアドレスは既に存在します')) {
          return t('auth.emailExists');
        }
        if (message.includes('Invalid email') || message.includes('無効なメールアドレス')) {
          return t('auth.invalidEmail');
        }
        if (message.includes('Weak password') || message.includes('弱いパスワード')) {
          return t('auth.weakPassword');
        }
        break;
      case 'username_exists':
        return t('auth.usernameExists');
      case 'email_exists':
        return t('auth.emailExists');
      case 'invalid_email':
        return t('auth.invalidEmail');
      case 'weak_password':
        return t('auth.weakPassword');
    }
    
    return message;
  };

  const handleSubmit = async () => {
    if (formData.password !== formData.confirmPassword) {
      setError(t('auth.passwordsDoNotMatch'));
      return;
    }

    setLoading(true);
    setError(null);
    setSuccess(null);
    setPendingApproval(false);
    setFieldErrors({});

    try {
      await register({
        username: formData.username,
        email: formData.email,
        password: formData.password,
        full_name: formData.full_name,
      });
      navigate('/');
    } catch (err) {
      const errorInfo = parseErrorResponse(err);
      const errorMessage = getErrorMessage(errorInfo);
      
      // 承認待ちメッセージの場合は成功として表示
      if (errorMessage.includes('管理者の承認をお待ちください') || 
          errorMessage.includes('administrator approval')) {
        setSuccess(t('auth.accountCreated'));
        setPendingApproval(true);
        // 3秒後にログインページに遷移
        setTimeout(() => {
          navigate('/auth/login');
        }, 3000);
      } else {
        // フィールド固有のエラーの場合
        if (errorInfo.field) {
          setFieldErrors(prev => ({ ...prev, [errorInfo.field!]: errorMessage }));
        } else {
          setError(errorMessage);
        }
      }
    } finally {
      setLoading(false);
    }
  };

  const getStepContent = (step: number) => {
    switch (step) {
      case 0:
        return (
          <Box>
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
              error={!!fieldErrors.username}
              helperText={fieldErrors.username || t('auth.usernameHelp')}
            />
            <TextField
              fullWidth
              label={t('auth.email')}
              type="email"
              value={formData.email}
              onChange={(e) => handleInputChange('email', e.target.value)}
              required
              margin="normal"
              variant="outlined"
              autoComplete="email"
              error={!!fieldErrors.email}
              helperText={fieldErrors.email || t('auth.emailHelp')}
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
              autoComplete="new-password"
              error={!!fieldErrors.password}
              helperText={fieldErrors.password || t('auth.passwordHelp')}
            />
          </Box>
        );
      case 1:
        return (
          <Box>
            <TextField
              fullWidth
              label={t('auth.fullName')}
              value={formData.full_name}
              onChange={(e) => handleInputChange('full_name', e.target.value)}
              required
              margin="normal"
              variant="outlined"
              autoComplete="name"
              error={!!fieldErrors.full_name}
              helperText={fieldErrors.full_name || t('auth.fullNameHelp')}
            />
          </Box>
        );
      case 2:
        return (
          <Box>
            <TextField
              fullWidth
              label={t('auth.confirmPassword')}
              type="password"
              value={formData.confirmPassword}
              onChange={(e) => handleInputChange('confirmPassword', e.target.value)}
              required
              margin="normal"
              variant="outlined"
              autoComplete="new-password"
              error={!!fieldErrors.confirmPassword}
              helperText={fieldErrors.confirmPassword || t('auth.confirmPasswordHelp')}
            />
            <Box sx={{ mt: 2, p: 2, bgcolor: 'background.paper', borderRadius: 1 }}>
              <Typography variant="subtitle2" sx={{ mb: 1 }}>
                {t('auth.accountSummary')}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                <strong>{t('auth.username')}:</strong> {formData.username}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                <strong>{t('auth.email')}:</strong> {formData.email}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                <strong>{t('auth.fullName')}:</strong> {formData.full_name}
              </Typography>
            </Box>
          </Box>
        );
      default:
        return 'Unknown step';
    }
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
          maxWidth: 500,
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
            {t('auth.createAccount')}
          </Typography>
          <Typography variant="body1" color="white" sx={{ mt: 1, opacity: 0.9 }}>
            {t('auth.joinOurCommunity')}
          </Typography>
        </Box>

        <CardContent sx={{ p: 4 }}>
          <Stepper activeStep={activeStep} sx={{ mb: 4 }}>
            {steps.map((label) => (
              <Step key={label}>
                <StepLabel>{t(`auth.steps.${label}`)}</StepLabel>
              </Step>
            ))}
          </Stepper>

          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          )}

          {success && (
            <Alert 
              severity="success" 
              sx={{ mb: 3 }}
              icon={pendingApproval ? <CircularProgress size={20} /> : undefined}
            >
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                {success}
                {pendingApproval && (
                  <Chip 
                    label={t('auth.pendingApproval')} 
                    size="small" 
                    color="warning" 
                    variant="outlined"
                  />
                )}
              </Box>
              <Typography variant="body2" sx={{ mt: 1 }}>
                {t('auth.redirectingToLogin')}
              </Typography>
            </Alert>
          )}

          <Box sx={{ opacity: pendingApproval ? 0.6 : 1, pointerEvents: pendingApproval ? 'none' : 'auto' }}>
            {getStepContent(activeStep)}
          </Box>

          <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 4 }}>
            <Button
              disabled={activeStep === 0 || pendingApproval}
              onClick={handleBack}
              sx={{ mr: 1 }}
            >
              {t('common.back')}
            </Button>
            <Box>
              <Button
                variant="contained"
                onClick={handleNext}
                disabled={!validateStep(activeStep) || loading || pendingApproval}
                sx={{
                  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                  '&:hover': {
                    background: 'linear-gradient(135deg, #5a6fd8 0%, #6a4190 100%)',
                  },
                }}
              >
                {loading 
                  ? t('auth.creatingAccount') 
                  : activeStep === steps.length - 1 
                    ? t('auth.createAccount') 
                    : t('common.next')
                }
              </Button>
            </Box>
          </Box>

          <Box sx={{ textAlign: 'center', mt: 3 }}>
            <Typography variant="body2" color="text.secondary">
              {t('auth.alreadyHaveAccount')}{' '}
              <Link
                component="button"
                variant="body2"
                onClick={() => navigate('/auth/login')}
                sx={{ textDecoration: 'none' }}
              >
                {t('auth.signIn')}
              </Link>
            </Typography>
          </Box>
        </CardContent>
      </Paper>
    </Box>
  );
};

export default RegisterForm;
