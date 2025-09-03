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

const RuleForm: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();
  const { t } = useTranslation();
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

  useEffect(() => {
    if (projectId) {
      setFormData(prev => ({ ...prev, project_id: projectId }));
    }
  }, [projectId]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      await api.post('/rules', formData);
      setSuccess(t('rules.createSuccess'));
      setTimeout(() => {
        navigate('/');
      }, 1500);
    } catch (error: any) {
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

          <TextField
            select
            label={t('rules.type')}
            value={formData.type}
            onChange={(e) => handleChange('type', e.target.value)}
            required
            helperText={t('rules.typeHelp')}
          >
            {Object.entries(t('types', { returnObjects: true })).map(([key, value]) => (
              <MenuItem key={key} value={key}>
                {value as string}
              </MenuItem>
            ))}
          </TextField>

          <TextField
            select
            label={t('rules.severity')}
            value={formData.severity}
            onChange={(e) => handleChange('severity', e.target.value)}
            required
            helperText={t('rules.severityHelp')}
          >
            {Object.entries(t('severity', { returnObjects: true })).map(([key, value]) => (
              <MenuItem key={key} value={key}>
                {value as string}
              </MenuItem>
            ))}
          </TextField>

          <TextField
            label={t('rules.pattern')}
            value={formData.pattern}
            onChange={(e) => handleChange('pattern', e.target.value)}
            required
            helperText={t('rules.patternHelp')}
          />

          <TextField
            label={t('rules.message')}
            value={formData.message}
            onChange={(e) => handleChange('message', e.target.value)}
            required
            multiline
            rows={2}
            helperText={t('rules.messageHelp')}
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

export default RuleForm;
