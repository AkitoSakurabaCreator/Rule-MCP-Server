import React, { useState, useEffect, useCallback } from 'react';
import {
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormControlLabel,
  Switch,
  Alert,
  Grid,
} from '@mui/material';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Rule } from '../types';
import { api } from '../services/api';
import { adminApi, RuleOption } from '../services/adminApi';
import { useAuth } from '../contexts/AuthContext';

const RuleEdit: React.FC = () => {
  const [rule, setRule] = useState<Partial<Rule>>({
    is_active: true, // デフォルトで有効状態
    type: 'style',
    severity: 'warning'
  });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [typeOptions, setTypeOptions] = useState<string[]>([]);
  const [severityOptions, setSeverityOptions] = useState<string[]>([]);
  const [newType, setNewType] = useState('');
  const [newSeverity, setNewSeverity] = useState('');
  
  const navigate = useNavigate();
  const { projectId, ruleId } = useParams<{ projectId: string; ruleId: string }>();
  const { t } = useTranslation();
  const { user } = useAuth();
  const isAdmin = user?.role === 'admin';

  const isEditMode = !!ruleId;

  const loadRule = useCallback(async () => {
    if (!projectId || !ruleId) return;
    try {
      const response = await api.get(`/rules/${projectId}/${ruleId}`);
      const loadedRule = response.data;
      if (loadedRule.is_active === undefined) {
        loadedRule.is_active = true;
      }
      setRule(loadedRule);
    } catch (error) {
      setError(t('rules.loadError'));
      console.error('Failed to load rule:', error);
    } finally {
      setLoading(false);
    }
  }, [projectId, ruleId, t]);

  useEffect(() => {
    if (isEditMode && projectId && ruleId) {
      loadRule();
    } else {
      setLoading(false);
    }
    loadOptions();
  }, [isEditMode, projectId, ruleId, loadRule]);

  // projectIdが変更された時にルールのproject_idを更新
  useEffect(() => {
    if (projectId) {
      setRule(prev => ({ ...prev, project_id: projectId }));
    }
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
      // 失敗してもデフォルトは維持
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setError(null);
    setSuccess(null);

    try {
      const ruleData = {
        ...rule,
        project_id: projectId
      };

      if (isEditMode) {
        await api.put(`/rules/${projectId}/${ruleId}`, ruleData);
        setSuccess(t('rules.editSuccess'));
      } else {
        await api.post('/rules', ruleData);
        setSuccess(t('rules.createSuccess'));
      }
      
      setTimeout(() => {
        navigate(`/projects/${projectId}/rules`);
      }, 1500);
    } catch (error: any) {
      console.error('Save error:', error);
      setError(error.response?.data?.error || t('rules.saveError'));
    } finally {
      setSaving(false);
    }
  };

  const handleInputChange = (field: keyof Rule, value: any) => {
    setRule(prev => ({ ...prev, [field]: value }));
  };

  const addOption = async (kind: 'type' | 'severity', value: string) => {
    const v = (value || '').trim();
    if (!v) return;
    try {
      await adminApi.addRuleOption(kind, v);
      await loadOptions();
      if (kind === 'type') setNewType(''); else setNewSeverity('');
    } catch (e) {
      // 失敗時も特に止めない
    }
  };

  if (loading) {
    return <Typography>{t('rules.loadingRule')}</Typography>;
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h4" component="h1" sx={{ mb: 3 }}>
          {isEditMode ? t('rules.editRule') : t('rules.newRule')}
        </Typography>
        
        <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
          {isEditMode 
            ? t('rules.editRuleHelp') 
            : t('rules.createRuleHelp')
          }
        </Typography>

        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
        {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}

        <Box component="form" onSubmit={handleSubmit}>
          <Grid container spacing={3}>
            <Grid sx={{ width: { xs: '100%', md: '50%' } }}>
              <TextField
                fullWidth
                label={t('rules.ruleId')}
                value={rule.rule_id || ''}
                onChange={(e) => handleInputChange('rule_id', e.target.value)}
                required
                helperText={t('rules.ruleIdHelp')}
                disabled={isEditMode} // 編集時は変更不可
              />
            </Grid>

            <Grid sx={{ width: { xs: '100%', md: '50%' } }}>
              <TextField
                fullWidth
                label={t('rules.name')}
                value={rule.name || ''}
                onChange={(e) => handleInputChange('name', e.target.value)}
                required
                helperText={t('rules.nameHelp')}
              />
            </Grid>

            <Grid sx={{ width: '100%' }}>
              <TextField
                fullWidth
                multiline
                rows={3}
                label={t('rules.description')}
                value={rule.description || ''}
                onChange={(e) => handleInputChange('description', e.target.value)}
                helperText={t('rules.descriptionHelp')}
              />
            </Grid>

            <Grid sx={{ width: { xs: '100%', md: '33.333%' } }}>
              <FormControl fullWidth required>
                <InputLabel>{t('rules.type')}</InputLabel>
                <Select
                  value={rule.type || ''}
                  onChange={(e) => handleInputChange('type', e.target.value)}
                  label={t('rules.type')}
                >
                  {typeOptions.map((v) => (
                    <MenuItem key={v} value={v}>{v}</MenuItem>
                  ))}
                </Select>
              </FormControl>
              {isAdmin && (
                <Box sx={{ mt: 1, display: 'flex', gap: 1 }}>
                  <TextField
                    size="small"
                    value={newType}
                    onChange={(e) => setNewType(e.target.value)}
                    placeholder={t('rules.addCustomType')}
                  />
                  <Button size="small" onClick={() => addOption('type', newType)}>
                    {t('common.add')}
                  </Button>
                </Box>
              )}
            </Grid>

            <Grid sx={{ width: { xs: '100%', md: '33.333%' } }}>
              <FormControl fullWidth required>
                <InputLabel>{t('rules.severity')}</InputLabel>
                <Select
                  value={rule.severity || ''}
                  onChange={(e) => handleInputChange('severity', e.target.value)}
                  label={t('rules.severity')}
                >
                  {severityOptions.map((v) => (
                    <MenuItem key={v} value={v}>{v}</MenuItem>
                  ))}
                </Select>
              </FormControl>
              {isAdmin && (
                <Box sx={{ mt: 1, display: 'flex', gap: 1 }}>
                  <TextField
                    size="small"
                    value={newSeverity}
                    onChange={(e) => setNewSeverity(e.target.value)}
                    placeholder={t('rules.addCustomSeverity')}
                  />
                  <Button size="small" onClick={() => addOption('severity', newSeverity)}>
                    {t('common.add')}
                  </Button>
                </Box>
              )}
            </Grid>

            <Grid sx={{ width: { xs: '100%', md: '33.333%' } }}>
              <Box sx={{ p: 2, border: '1px solid', borderColor: 'divider', borderRadius: 1, bgcolor: 'background.paper' }}>
                <Typography variant="subtitle2" sx={{ mb: 1, fontWeight: 'bold' }}>
                  {t('rules.ruleStatus')}
                </Typography>
                <FormControlLabel
                  control={
                    <Switch
                      checked={rule.is_active ?? true}
                      onChange={(e) => handleInputChange('is_active', e.target.checked)}
                      color="success"
                    />
                  }
                  label={
                    <Box>
                      <Typography variant="body2" sx={{ fontWeight: rule.is_active ? 'bold' : 'normal' }}>
                        {rule.is_active ? t('rules.statusActive') : t('rules.statusInactive')}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        {rule.is_active ? t('rules.statusActiveHelp') : t('rules.statusInactiveHelp')}
                      </Typography>
                    </Box>
                  }
                />
              </Box>
            </Grid>

            <Grid sx={{ width: '100%' }}>
              <TextField
                fullWidth
                label={t('rules.pattern')}
                value={rule.pattern || ''}
                onChange={(e) => handleInputChange('pattern', e.target.value)}
                helperText={
                  <Box>
                    <Typography variant="body2" component="span">
                      {t('rules.patternHelp')}
                    </Typography>
                    <Box sx={{ mt: 1 }}>
                      <Typography variant="caption" component="div" color="text.secondary">
                        <strong>例:</strong> console\.log, TODO:, api_key, SELECT * FROM
                      </Typography>
                      <Typography variant="caption" component="div" color="text.secondary">
                        <strong>ヒント:</strong> 特殊文字は \ でエスケープしてください
                      </Typography>
                    </Box>
                  </Box>
                }
                placeholder="例: console\\.log"
              />
            </Grid>

            <Grid sx={{ width: '100%' }}>
              <TextField
                fullWidth
                multiline
                rows={2}
                label={t('rules.message')}
                value={rule.message || ''}
                onChange={(e) => handleInputChange('message', e.target.value)}
                helperText={
                  <Box>
                    <Typography variant="body2" component="span">
                      {t('rules.messageHelp')}
                    </Typography>
                    <Box sx={{ mt: 1 }}>
                      <Typography variant="caption" component="div" color="text.secondary">
                        <strong>例:</strong> Console.log detected. Use proper logging framework in production.
                      </Typography>
                      <Typography variant="caption" component="div" color="text.secondary">
                        <strong>ヒント:</strong> 開発者にとって分かりやすく、修正方法も含めてください
                      </Typography>
                    </Box>
                  </Box>
                }
                placeholder="例: Console.log detected. Use proper logging framework in production."
              />
            </Grid>
          </Grid>

          <Box sx={{ mt: 3, display: 'flex', gap: 2, alignItems: 'center' }}>
            {isAdmin ? (
              <Button
                type="submit"
                variant="contained"
                disabled={saving}
              >
                {saving ? t('common.saving') : (isEditMode ? t('common.save') : t('common.create'))}
              </Button>
            ) : (
              <Typography variant="body2" color="text.secondary">
                {t('common.permissionDenied') || '権限がありません'}
              </Typography>
            )}
            <Button
              variant="outlined"
              onClick={() => navigate(`/projects/${projectId}/rules`)}
            >
              {t('common.cancel')}
            </Button>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default RuleEdit;
