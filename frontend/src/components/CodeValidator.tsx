import React, { useState, useEffect, useCallback } from 'react';
import {
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Box,
  MenuItem,
  Alert,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Paper,
} from '@mui/material';
import { Error as ErrorIcon, Warning as WarningIcon } from '@mui/icons-material';
import { useTranslation } from 'react-i18next';
import { Project, ValidationResult } from '../types';
import { api } from '../services/api';

const CodeValidator: React.FC = () => {
  const { t } = useTranslation();
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProject, setSelectedProject] = useState<string>('');
  const [code, setCode] = useState<string>('');
  const [validationResult, setValidationResult] = useState<ValidationResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadProjects = useCallback(async () => {
    try {
      const response = await api.get('/projects');
      setProjects(response.data.projects);
      if (response.data.projects.length > 0) {
        setSelectedProject(response.data.projects[0].project_id);
      }
    } catch (error) {
      setError(t('codeValidator.loadProjectsError'));
      console.error('Failed to load projects:', error);
    }
  }, [t]);

  useEffect(() => {
    loadProjects();
  }, [loadProjects]);

  const handleValidate = async () => {
    if (!selectedProject || !code.trim()) {
      setError(t('codeValidator.selectProjectAndCode'));
      return;
    }

    setLoading(true);
    setError(null);
    setValidationResult(null);

    try {
      const response = await api.post('/rules/validate', {
        project_id: selectedProject,
        code: code.trim(),
      });
      setValidationResult(response.data);
    } catch (error: any) {
      setError(error.response?.data?.error || t('codeValidator.validationError'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card>
      <CardContent>
        <Typography variant="h5" component="h2" sx={{ mb: 3 }}>
          {t('codeValidator.title')}
        </Typography>

        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <TextField
            select
            label={t('codeValidator.selectProject')}
            value={selectedProject}
            onChange={(e) => setSelectedProject(e.target.value)}
            required
            helperText={t('codeValidator.selectProjectHelp')}
          >
            {projects.map((project) => (
              <MenuItem key={project.project_id} value={project.project_id}>
                {project.name} ({project.language})
              </MenuItem>
            ))}
          </TextField>

          <TextField
            label={t('codeValidator.codeToValidate')}
            value={code}
            onChange={(e) => setCode(e.target.value)}
            multiline
            rows={10}
            required
            placeholder={t('codeValidator.enterCode')}
            helperText={t('codeValidator.pasteCode')}
          />

          <Button
            variant="contained"
            onClick={handleValidate}
            disabled={loading || !selectedProject || !code.trim()}
            sx={{ minWidth: 120 }}
          >
            {loading ? t('codeValidator.validating') : t('codeValidator.validateCode')}
          </Button>
        </Box>

        {validationResult && (
          <Paper sx={{ mt: 3, p: 2 }}>
            <Typography variant="h6" sx={{ mb: 2 }}>
              {t('codeValidator.validationResults')}
            </Typography>
            
            {validationResult.valid ? (
              <Alert severity="success" sx={{ mb: 2 }}>
                {t('codeValidator.codeIsValid')}
              </Alert>
            ) : (
              <Alert severity="warning" sx={{ mb: 2 }}>
                {t('codeValidator.codeHasIssues')}
              </Alert>
            )}

            {validationResult.errors && validationResult.errors.length > 0 && (
              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle1" color="error" sx={{ mb: 1 }}>
                  {t('codeValidator.errors')}:
                </Typography>
                <List dense>
                  {validationResult.errors.map((error, index) => (
                    <ListItem key={index}>
                      <ListItemIcon>
                        <ErrorIcon color="error" />
                      </ListItemIcon>
                      <ListItemText primary={error} />
                    </ListItem>
                  ))}
                </List>
              </Box>
            )}

            {validationResult.warnings && validationResult.warnings.length > 0 && (
              <Box>
                <Typography variant="subtitle1" color="warning.main" sx={{ mb: 1 }}>
                  {t('codeValidator.warnings')}:
                </Typography>
                <List dense>
                  {validationResult.warnings.map((warning, index) => (
                    <ListItem key={index}>
                      <ListItemIcon>
                        <WarningIcon color="warning" />
                      </ListItemIcon>
                      <ListItemText primary={warning} />
                    </ListItem>
                  ))}
                </List>
              </Box>
            )}

            {validationResult.valid && (!validationResult.errors || validationResult.errors.length === 0) && 
             (!validationResult.warnings || validationResult.warnings.length === 0) && (
              <Alert severity="success">
                {t('codeValidator.noIssues')}
              </Alert>
            )}
          </Paper>
        )}
      </CardContent>
    </Card>
  );
};

export default CodeValidator;
