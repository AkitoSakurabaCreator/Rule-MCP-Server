import React, { useState, useEffect } from 'react';
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
import { Project } from '../types';
import { api } from '../services/api';

const languages = ['general', 'javascript', 'go', 'python', 'java', 'csharp', 'cpp', 'rust'];

const ProjectForm: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    project_id: '',
    name: '',
    description: '',
    language: 'general',
    apply_global_rules: true,
  });

  const isEditMode = Boolean(projectId);

  useEffect(() => {
    if (isEditMode && projectId) {
      loadProject(projectId);
    }
  }, [isEditMode, projectId]);

  const loadProject = async (id: string) => {
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
      setError('Failed to load project');
      console.error('Failed to load project:', error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      if (isEditMode) {
        await api.put(`/projects/${projectId}`, formData);
        setSuccess('Project updated successfully!');
      } else {
        await api.post('/projects', formData);
        setSuccess('Project created successfully!');
      }

      setTimeout(() => {
        navigate('/');
      }, 1500);
    } catch (error: any) {
      setError(error.response?.data?.error || 'An error occurred');
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
          {isEditMode ? 'Edit Project' : 'New Project'}
        </Typography>

        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
        {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}

        <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <TextField
            label="Project ID"
            value={formData.project_id}
            onChange={(e) => handleChange('project_id', e.target.value)}
            required
            disabled={isEditMode}
            helperText="Unique identifier for the project"
          />

          <TextField
            label="Name"
            value={formData.name}
            onChange={(e) => handleChange('name', e.target.value)}
            required
            helperText="Display name for the project"
          />

          <TextField
            label="Description"
            value={formData.description}
            onChange={(e) => handleChange('description', e.target.value)}
            multiline
            rows={3}
            helperText="Optional description of the project"
          />

          <TextField
            select
            label="Language"
            value={formData.language}
            onChange={(e) => handleChange('language', e.target.value)}
            required
            helperText="Programming language for this project"
          >
            {languages.map((lang) => (
              <MenuItem key={lang} value={lang}>
                {lang.charAt(0).toUpperCase() + lang.slice(1)}
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
            label="Apply global rules for this language"
          />

          <Box sx={{ display: 'flex', gap: 2, mt: 2 }}>
            <Button
              type="submit"
              variant="contained"
              disabled={loading}
              sx={{ minWidth: 120 }}
            >
              {loading ? 'Saving...' : (isEditMode ? 'Update' : 'Create')}
            </Button>
            <Button
              variant="outlined"
              onClick={() => navigate('/')}
              disabled={loading}
            >
              Cancel
            </Button>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default ProjectForm;
