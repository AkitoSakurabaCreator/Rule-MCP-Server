import React, { useState, useEffect, useCallback } from 'react';
import {
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  FormControlLabel,
  Checkbox,
  Box,
  MenuItem,
  Alert,
} from '@mui/material';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Project } from '../types';
import { api } from '../services/api';

interface Language {
  code: string;
  name: string;
  description: string;
  icon: string;
  color: string;
  is_active: boolean;
}

const ProjectForm: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [languages, setLanguages] = useState<Language[]>([]);

  const [formData, setFormData] = useState({
    project_id: '',
    name: '',
    description: '',
    language: 'general',
    apply_global_rules: true,
  });

  const isEditMode = Boolean(projectId);

  const loadLanguages = useCallback(async () => {
    try {
      const response = await api.get('/languages');
      setLanguages(response.data.languages.filter((lang: Language) => lang.is_active));
    } catch (error) {
      console.error('Failed to load languages:', error);
    }
  }, []);

  const loadProject = useCallback(async (id: string) => {
    try {
      const response = await api.get(`/projects/${id}`);
      const project: Project = response.data;
      setFormData({
        project_id: project.project_id,
        name: project.name,
        description: project.description,
        language: project.language,
        apply_global_rules: project.apply_global_rules,
      });
    } catch (error) {
      setError(t('common.error'));
      console.error('Failed to load project:', error);
    }
  }, [t]);

  useEffect(() => {
    loadLanguages();
    if (isEditMode && projectId) {
      loadProject(projectId);
    }
  }, [isEditMode, projectId, loadProject, loadLanguages]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      if (isEditMode) {
        await api.put(`/projects/${projectId}`, formData);
        setSuccess(t('projects.editSuccess'));
      } else {
        await api.post('/projects', formData);
        setSuccess(t('projects.createSuccess'));
      }
      setTimeout(() => {
        navigate('/');
      }, 1500);
    } catch (error: any) {
      setError(error.response?.data?.error || t('common.error'));
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field: string, value: string | boolean) => {
    setFormData(prev => ({
      ...prev,
      [field]: value,
    }));
  };

  return (
    <Card>
      <CardContent>
        <Typography variant="h5" component="h2" sx={{ mb: 3 }}>
          {isEditMode ? t('projects.editProject') : t('projects.newProject')}
        </Typography>

        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
        {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}

        <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <TextField
            label={t('projects.projectId')}
            value={formData.project_id}
            onChange={(e) => handleChange('project_id', e.target.value)}
            required
            disabled={isEditMode}
            helperText={t('projects.projectIdHelp')}
          />

          <TextField
            label={t('projects.name')}
            value={formData.name}
            onChange={(e) => handleChange('name', e.target.value)}
            required
            helperText={t('projects.nameHelp')}
          />

          <TextField
            label={t('projects.description')}
            value={formData.description}
            onChange={(e) => handleChange('description', e.target.value)}
            multiline
            rows={3}
            helperText={t('projects.descriptionHelp')}
          />

          <TextField
            select
            label={t('projects.language')}
            value={formData.language}
            onChange={(e) => handleChange('language', e.target.value)}
            required
            helperText={t('projects.languageHelp')}
          >
            {languages.map((language) => (
              <MenuItem key={language.code} value={language.code}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <Box
                    sx={{
                      width: 16,
                      height: 16,
                      backgroundColor: language.color,
                      borderRadius: 0.5,
                    }}
                  />
                  {language.name}
                </Box>
              </MenuItem>
            ))}
          </TextField>

          <FormControlLabel
            control={
              <Checkbox
                checked={formData.apply_global_rules}
                onChange={(e) => handleChange('apply_global_rules', e.target.checked)}
              />
            }
            label={t('projects.applyGlobalRules')}
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
              {loading ? t('common.loading') : (isEditMode ? t('common.save') : t('common.add'))}
            </Button>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default ProjectForm;
