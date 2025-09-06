import React, { useEffect, useState } from 'react';
import {
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Box,
  MenuItem,
  Alert,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { api } from '../services/api';
import { adminApi, RuleOption } from '../services/adminApi';
import { useAuth } from '../contexts/AuthContext';

const GlobalRuleForm: React.FC = () => {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const { permissions } = useAuth();
  const canManageRules = permissions.manageRules;

  const [formData, setFormData] = useState({
    language: 'general',
    rule_id: '',
    name: '',
    description: '',
    type: '',
    severity: '',
    pattern: '',
    message: '',
  });

  const [typeOptions, setTypeOptions] = useState<string[]>([]);
  const [severityOptions, setSeverityOptions] = useState<string[]>([]);
  const [newType, setNewType] = useState('');
  const [newSeverity, setNewSeverity] = useState('');

  useEffect(() => {
    const loadOptions = async () => {
      try {
        const [types, severities] = await Promise.all([
          adminApi.getRuleOptions('type'),
          adminApi.getRuleOptions('severity'),
        ]);
        setTypeOptions(types.map((o: RuleOption) => o.value));
        setSeverityOptions(severities.map((o: RuleOption) => o.value));
      } catch (_) {
      }
    };
    loadOptions();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      await api.post('/global-rules', formData);
      setSuccess(t('globalRules.createSuccess'));
      setTimeout(() => {
        navigate('/');
      }, 1500);
    } catch (error: any) {
      console.error('Global rule creation error:', error);
      
      // 認証エラーの場合はログイン画面にリダイレクト
      if (error.response?.status === 401 || error.response?.status === 403) {
        setError('認証エラーが発生しました。再度ログインしてください。');
        setTimeout(() => {
          window.location.href = '/login';
        }, 2000);
        return;
      }
      
      setError(error.response?.data?.error || t('common.error'));
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field: string, value: string) => {
    setFormData(prev => ({
      ...prev,
      [field]: value,
    }));
  };

  return (
    <Card>
      <CardContent>
        <Typography variant="h5" component="h2" sx={{ mb: 3 }}>
          {t('globalRules.newGlobalRule')}
        </Typography>

        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
        {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}

        <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <TextField
            select
            label={t('globalRules.language')}
            value={formData.language}
            onChange={(e) => handleChange('language', e.target.value)}
            required
            helperText={t('globalRules.languageHelp')}
          >
            {Object.entries(t('languages', { returnObjects: true })).map(([key, value]) => (
              <MenuItem key={key} value={key}>
                {value as string}
              </MenuItem>
            ))}
          </TextField>

          <TextField
            label={t('rules.ruleId')}
            value={formData.rule_id}
            onChange={(e) => handleChange('rule_id', e.target.value)}
            required
            helperText={t('globalRules.ruleIdHelp')}
          />

          <TextField
            label={t('rules.name')}
            value={formData.name}
            onChange={(e) => handleChange('name', e.target.value)}
            required
            helperText={t('globalRules.nameHelp')}
          />

          <TextField
            label={t('rules.description')}
            value={formData.description}
            onChange={(e) => handleChange('description', e.target.value)}
            multiline
            rows={3}
            helperText={t('globalRules.descriptionHelp')}
          />

          <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
            <TextField
              select
              label={t('rules.type')}
              value={formData.type}
              onChange={(e) => handleChange('type', e.target.value)}
              helperText={t('globalRules.typeHelp')}
              sx={{ flex: 1 }}
            >
              {typeOptions.map((v) => (
                <MenuItem key={v} value={v}>{v}</MenuItem>
              ))}
            </TextField>
            {canManageRules && (
              <>
                <TextField size="small" value={newType} onChange={(e) => setNewType(e.target.value)} placeholder={t('rules.addCustomType')} />
                <Button size="small" onClick={async () => { const v = (newType || '').trim(); if (!v) return; try { await adminApi.addRuleOption('type', v); setNewType(''); const types = await adminApi.getRuleOptions('type'); setTypeOptions(types.map((o: RuleOption) => o.value)); } catch (_) {} }}>{t('common.add')}</Button>
              </>
            )}
          </Box>

          <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
            <TextField
              select
              label={t('rules.severity')}
              value={formData.severity}
              onChange={(e) => handleChange('severity', e.target.value)}
              helperText={t('globalRules.severityHelp')}
              sx={{ flex: 1 }}
            >
              {severityOptions.map((v) => (
                <MenuItem key={v} value={v}>{v}</MenuItem>
              ))}
            </TextField>
            {canManageRules && (
              <>
                <TextField size="small" value={newSeverity} onChange={(e) => setNewSeverity(e.target.value)} placeholder={t('rules.addCustomSeverity')} />
                <Button size="small" onClick={async () => { const v = (newSeverity || '').trim(); if (!v) return; try { await adminApi.addRuleOption('severity', v); setNewSeverity(''); const severities = await adminApi.getRuleOptions('severity'); setSeverityOptions(severities.map((o: RuleOption) => o.value)); } catch (_) {} }}>{t('common.add')}</Button>
              </>
            )}
          </Box>

          <TextField
            label={t('rules.pattern')}
            value={formData.pattern}
            onChange={(e) => handleChange('pattern', e.target.value)}
            helperText={t('globalRules.patternHelp')}
          />

          <TextField
            label={t('rules.message')}
            value={formData.message}
            onChange={(e) => handleChange('message', e.target.value)}
            multiline
            rows={2}
            helperText={t('globalRules.messageHelp')}
          />

          <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end' }}>
            <Button
              variant="outlined"
              onClick={() => navigate('/')}
              disabled={loading}
            >
              {t('common.cancel')}
            </Button>
            <Button
              type="submit"
              variant="contained"
              disabled={loading}
            >
              {loading ? t('common.loading') : t('common.add')}
            </Button>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default GlobalRuleForm;
