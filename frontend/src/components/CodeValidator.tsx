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
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Paper,
} from '@mui/material';
import { CheckCircle as CheckIcon, Error as ErrorIcon, Warning as WarningIcon } from '@mui/icons-material';
import { Project, ValidationResult } from '../types';
import { api } from '../services/api';

const CodeValidator: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProject, setSelectedProject] = useState<string>('');
  const [code, setCode] = useState<string>('');
  const [validationResult, setValidationResult] = useState<ValidationResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadProjects();
  }, []);

  const loadProjects = async () => {
    try {
      const response = await api.get('/projects');
      setProjects(response.data.projects);
      if (response.data.projects.length > 0) {
        setSelectedProject(response.data.projects[0].project_id);
      }
    } catch (error) {
      setError('Failed to load projects');
      console.error('Failed to load projects:', error);
    }
  };

  const handleValidate = async () => {
    if (!selectedProject || !code.trim()) {
      setError('Please select a project and enter code to validate');
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
      setError(error.response?.data?.error || 'An error occurred during validation');
    } finally {
      setLoading(false);
    }
  };

  const getSeverityIcon = (severity: string) => {
    switch (severity) {
      case 'error':
        return <ErrorIcon color="error" />;
      case 'warning':
        return <WarningIcon color="warning" />;
      default:
        return <CheckIcon color="info" />;
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'error':
        return 'error.main';
      case 'warning':
        return 'warning.main';
      default:
        return 'info.main';
    }
  };

  return (
    <Card>
      <CardContent>
        <Typography variant="h5" component="h2" sx={{ mb: 3 }}>
          Code Validator
        </Typography>

        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <TextField
            select
            label="Project"
            value={selectedProject}
            onChange={(e) => setSelectedProject(e.target.value)}
            required
            helperText="Select a project to validate code against its rules"
          >
            {projects.map((project) => (
              <MenuItem key={project.project_id} value={project.project_id}>
                {project.name} ({project.language})
              </MenuItem>
            ))}
          </TextField>

          <TextField
            label="Code to Validate"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            multiline
            rows={10}
            required
            placeholder="Enter your code here..."
            helperText="Paste the code you want to validate"
          />

          <Button
            variant="contained"
            onClick={handleValidate}
            disabled={loading || !selectedProject || !code.trim()}
            sx={{ minWidth: 120 }}
          >
            {loading ? 'Validating...' : 'Validate Code'}
          </Button>
        </Box>

        {validationResult && (
          <Paper sx={{ mt: 3, p: 2 }}>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Validation Results
            </Typography>

            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
              {validationResult.valid ? (
                <CheckIcon color="success" sx={{ mr: 1 }} />
              ) : (
                <ErrorIcon color="error" sx={{ mr: 1 }} />
              )}
              <Typography
                variant="h6"
                color={validationResult.valid ? 'success.main' : 'error.main'}
              >
                {validationResult.valid ? 'Code is Valid' : 'Code has Issues'}
              </Typography>
            </Box>

            {validationResult.errors.length > 0 && (
              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle1" color="error" sx={{ mb: 1 }}>
                  Errors ({validationResult.errors.length}):
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

            {validationResult.warnings.length > 0 && (
              <Box>
                <Typography variant="subtitle1" color="warning.main" sx={{ mb: 1 }}>
                  Warnings ({validationResult.warnings.length}):
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
          </Paper>
        )}
      </CardContent>
    </Card>
  );
};

export default CodeValidator;
