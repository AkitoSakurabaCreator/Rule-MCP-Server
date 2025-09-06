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

  const handleSubmit = async () => {
    if (formData.password !== formData.confirmPassword) {
      setError(t('auth.passwordsDoNotMatch'));
      return;
    }

    setLoading(true);
    setError(null);

    try {
      await register({
        username: formData.username,
        email: formData.email,
        password: formData.password,
        full_name: formData.full_name,
      });
      navigate('/');
    } catch (err) {
      setError(err instanceof Error ? err.message : t('auth.registerError'));
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
              helperText={t('auth.usernameHelp')}
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
              helperText={t('auth.emailHelp')}
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
              helperText={t('auth.passwordHelp')}
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
              helperText={t('auth.fullNameHelp')}
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
              helperText={t('auth.confirmPasswordHelp')}
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

          {getStepContent(activeStep)}

          <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 4 }}>
            <Button
              disabled={activeStep === 0}
              onClick={handleBack}
              sx={{ mr: 1 }}
            >
              {t('common.back')}
            </Button>
            <Box>
              <Button
                variant="contained"
                onClick={handleNext}
                disabled={!validateStep(activeStep) || loading}
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
