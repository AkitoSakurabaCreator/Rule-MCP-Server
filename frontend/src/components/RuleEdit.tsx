import React, { useState, useEffect } from 'react';
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
  const [customLanguages, setCustomLanguages] = useState<string[]>([]);
  const [customTypes, setCustomTypes] = useState<string[]>([]);
  const [customSeverities, setCustomSeverities] = useState<string[]>([]);
  
  const navigate = useNavigate();
  const { projectId, ruleId } = useParams<{ projectId: string; ruleId: string }>();
  const { t } = useTranslation();

  const isEditMode = !!ruleId;

  useEffect(() => {
    if (isEditMode && projectId && ruleId) {
      loadRule();
    } else {
      setLoading(false);
    }
    loadCustomOptions();
  }, [isEditMode, projectId, ruleId]);

  // projectIdが変更された時にルールのproject_idを更新
  useEffect(() => {
    if (projectId) {
      setRule(prev => ({ ...prev, project_id: projectId }));
    }
  }, [projectId]);

  const loadRule = async () => {
    try {
      const response = await api.get(`/projects/${projectId}/rules/${ruleId}`);
      const loadedRule = response.data;
      // is_activeがundefinedの場合はtrueに設定
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
  };

  const loadCustomOptions = async () => {
    try {
      // カスタムオプションをロード（必要に応じて実装）
      setCustomLanguages([]);
      setCustomTypes([]);
      setCustomSeverities([]);
    } catch (error) {
      console.error('Failed to load custom options:', error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setError(null);
    setSuccess(null);

    try {
      // ルールデータにproject_idを確実に設定
      const ruleData = {
        ...rule,
        project_id: projectId
      };

      if (isEditMode) {
        await api.put(`/projects/${projectId}/rules/${ruleId}`, ruleData);
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

  const addCustomOption = (type: 'language' | 'type' | 'severity', value: string) => {
    if (!value.trim()) return;
    
    switch (type) {
      case 'language':
        if (!customLanguages.includes(value)) {
          setCustomLanguages(prev => [...prev, value]);
        }
        break;
      case 'type':
        if (!customTypes.includes(value)) {
          setCustomTypes(prev => [...prev, value]);
        }
        break;
      case 'severity':
        if (!customSeverities.includes(value)) {
          setCustomSeverities(prev => [...prev, value]);
        }
        break;
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
                required
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
                  {Object.entries(t('types', { returnObjects: true })).map(([key, value]) => (
                    <MenuItem key={key} value={key}>
                      {value}
                    </MenuItem>
                  ))}
                  {customTypes.map((type) => (
                    <MenuItem key={type} value={type}>
                      {type}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
              <Box sx={{ mt: 1 }}>
                <TextField
                  size="small"
                  placeholder={t('rules.addCustomType')}
                  onKeyPress={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault();
                      addCustomOption('type', (e.target as HTMLInputElement).value);
                      (e.target as HTMLInputElement).value = '';
                    }
                  }}
                />
                <Button
                  size="small"
                  onClick={() => {
                    const input = document.querySelector('input[placeholder*="Type"]') as HTMLInputElement;
                    if (input) {
                      addCustomOption('type', input.value);
                      input.value = '';
                    }
                  }}
                >
                  {t('common.add')}
                </Button>
              </Box>
            </Grid>

            <Grid sx={{ width: { xs: '100%', md: '33.333%' } }}>
              <FormControl fullWidth required>
                <InputLabel>{t('rules.severity')}</InputLabel>
                <Select
                  value={rule.severity || ''}
                  onChange={(e) => handleInputChange('severity', e.target.value)}
                  label={t('rules.severity')}
                >
                  {Object.entries(t('severity', { returnObjects: true })).map(([key, value]) => (
                    <MenuItem key={key} value={key}>
                      {value}
                    </MenuItem>
                  ))}
                  {customSeverities.map((severity) => (
                    <MenuItem key={severity} value={severity}>
                      {severity}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
              <Box sx={{ mt: 1 }}>
                <TextField
                  size="small"
                  placeholder={t('rules.addCustomSeverity')}
                  onKeyPress={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault();
                      addCustomOption('severity', (e.target as HTMLInputElement).value);
                      (e.target as HTMLInputElement).value = '';
                    }
                  }}
                />
                <Button
                  size="small"
                  onClick={() => {
                    const input = document.querySelector('input[placeholder*="Severity"]') as HTMLInputElement;
                    if (input) {
                      addCustomOption('severity', input.value);
                      input.value = '';
                    }
                  }}
                >
                  {t('common.add')}
                </Button>
              </Box>
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
                required
                helperText={
                  <Box>
                    <Typography variant="body2" component="span">
                      {t('rules.patternHelp')}
                    </Typography>
                    <Box sx={{ mt: 1 }}>
                      <Typography variant="caption" component="div" color="text.secondary">
                        <strong>例:</strong> console\.log, TODO:, api_key, SELECT \* FROM
                      </Typography>
                      <Typography variant="caption" component="div" color="text.secondary">
                        <strong>ヒント:</strong> 特殊文字は \ でエスケープしてください
                      </Typography>
                    </Box>
                  </Box>
                }
                placeholder="例: console\.log"
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
                required
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

          <Box sx={{ mt: 3, display: 'flex', gap: 2 }}>
            <Button
              type="submit"
              variant="contained"
              disabled={saving}
            >
              {saving ? t('common.saving') : (isEditMode ? t('common.save') : t('common.create'))}
            </Button>
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
