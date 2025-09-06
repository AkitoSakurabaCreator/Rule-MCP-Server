import React, { useState, useEffect } from 'react';
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
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { api } from '../services/api';
import { adminApi, RuleOption } from '../services/adminApi';
import { useAuth } from '../contexts/AuthContext';

const RuleForm: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { permissions } = useAuth();
  const canManageRules = permissions.manageRules;
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    project_id: projectId || '',
    rule_id: '',
    name: '',
    description: '',
    type: 'style',
    severity: 'warning',
    pattern: '',
    message: '',
  });

  const [typeOptions, setTypeOptions] = useState<string[]>([]);
  const [severityOptions, setSeverityOptions] = useState<string[]>([]);
  const [newType, setNewType] = useState('');
  const [newSeverity, setNewSeverity] = useState('');

  useEffect(() => {
    if (projectId) {
      setFormData(prev => ({ ...prev, project_id: projectId }));
    }
    loadOptions();
  }, [projectId]);

  const loadOptions = async () => {
    try {
      const [types, severities] = await Promise.all([
        adminApi.getRuleOptions('type'),
        adminApi.getRuleOptions('severity'),
      ]);
      setTypeOptions(types.map((o: RuleOption) => o.value));
      setSeverityOptions(severities.map((o: RuleOption) => o.value));
    } catch (e) {
      // ignore
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      await api.post('/rules', formData);
      setSuccess(t('rules.createSuccess'));
      setTimeout(() => {
        navigate(`/projects/${projectId}/rules`);
      }, 1500);
    } catch (error: any) {
      console.error('Rule creation error:', error);
      const errorData = error.response?.data;
      
      // 認証エラーの場合はログイン画面にリダイレクト
      if (error.response?.status === 401 || error.response?.status === 403) {
        setError('認証エラーが発生しました。再度ログインしてください。');
        setTimeout(() => {
          window.location.href = '/login';
        }, 2000);
        return;
      }
      
      if (errorData?.details && errorData?.suggestion) {
        setError(`${errorData.error}\n\n詳細: ${errorData.details}\n\n提案: ${errorData.suggestion}`);
      } else {
        setError(errorData?.error || t('common.error'));
      }
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

  const addOption = async (kind: 'type' | 'severity', value: string) => {
    const v = (value || '').trim();
    if (!v) return;
    try {
      await adminApi.addRuleOption(kind, v);
      await loadOptions();
      if (kind === 'type') setNewType(''); else setNewSeverity('');
    } catch (e) {
      // ignore
    }
  };

  return (
    <Card>
      <CardContent>
        <Typography variant="h5" component="h2" sx={{ mb: 3 }}>
          {t('rules.newRule')}
        </Typography>

        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
        {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}

        <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <TextField
            label={t('rules.ruleId')}
            value={formData.rule_id}
            onChange={(e) => handleChange('rule_id', e.target.value)}
            required
            helperText={t('rules.ruleIdHelp')}
          />

          <TextField
            label={t('rules.name')}
            value={formData.name}
            onChange={(e) => handleChange('name', e.target.value)}
            required
            helperText={t('rules.nameHelp')}
          />

          <TextField
            label={t('rules.description')}
            value={formData.description}
            onChange={(e) => handleChange('description', e.target.value)}
            multiline
            rows={3}
            helperText={t('rules.descriptionHelp')}
          />

          <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
            <TextField
              select
              label={t('rules.type')}
              value={formData.type}
              onChange={(e) => handleChange('type', e.target.value)}
              required
              helperText={t('rules.typeHelp')}
              sx={{ flex: 1 }}
            >
              {typeOptions.map((v) => (
                <MenuItem key={v} value={v}>{v}</MenuItem>
              ))}
            </TextField>
            {canManageRules && (
              <>
                <TextField size="small" value={newType} onChange={(e) => setNewType(e.target.value)} placeholder={t('rules.addCustomType')} />
                <Button size="small" onClick={() => addOption('type', newType)}>{t('common.add')}</Button>
              </>
            )}
          </Box>

          <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
            <TextField
              select
              label={t('rules.severity')}
              value={formData.severity}
              onChange={(e) => handleChange('severity', e.target.value)}
              required
              helperText={t('rules.severityHelp')}
              sx={{ flex: 1 }}
            >
              {severityOptions.map((v) => (
                <MenuItem key={v} value={v}>{v}</MenuItem>
              ))}
            </TextField>
            {canManageRules && (
              <>
                <TextField size="small" value={newSeverity} onChange={(e) => setNewSeverity(e.target.value)} placeholder={t('rules.addCustomSeverity')} />
                <Button size="small" onClick={() => addOption('severity', newSeverity)}>{t('common.add')}</Button>
              </>
            )}
          </Box>

          <TextField
            label={t('rules.pattern')}
            value={formData.pattern}
            onChange={(e) => handleChange('pattern', e.target.value)}
            helperText={t('rules.patternHelp')}
          />

          <TextField
            label={t('rules.message')}
            value={formData.message}
            onChange={(e) => handleChange('message', e.target.value)}
            multiline
            rows={2}
            helperText={t('rules.messageHelp')}
          />

          <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end', alignItems: 'center' }}>
            <Button
              variant="outlined"
              onClick={() => navigate(`/projects/${projectId}/rules`)}
              disabled={loading}
            >
              {t('common.cancel')}
            </Button>
            {canManageRules ? (
              <Button
                type="submit"
                variant="contained"
                disabled={loading}
              >
                {loading ? t('common.loading') : t('common.add')}
              </Button>
            ) : (
              <Typography variant="body2" color="text.secondary">
                {t('common.permissionDenied')}
              </Typography>
            )}
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default RuleForm;
