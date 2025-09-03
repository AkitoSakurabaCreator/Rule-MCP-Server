import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Button,
  Grid,
  Chip,
  Box,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material';
import { Edit as EditIcon, Delete as DeleteIcon, Add as AddIcon } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Project } from '../types';
import { api } from '../services/api';

const ProjectList: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [projectToDelete, setProjectToDelete] = useState<string | null>(null);
  const navigate = useNavigate();
  const { t } = useTranslation();

  useEffect(() => {
    loadProjects();
  }, []);

  const loadProjects = async () => {
    try {
      const response = await api.get('/projects');
      setProjects(response.data.projects);
    } catch (error) {
      console.error('Failed to load projects:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!projectToDelete) return;

    try {
      await api.delete(`/projects/${projectToDelete}`);
      await loadProjects();
      setDeleteDialogOpen(false);
      setProjectToDelete(null);
    } catch (error) {
      console.error('Failed to delete project:', error);
    }
  };

  const openDeleteDialog = (projectId: string) => {
    setProjectToDelete(projectId);
    setDeleteDialogOpen(true);
  };

  if (loading) {
    return <Typography>{t('projects.loadingProjects')}</Typography>;
  }

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1">
          {t('projects.title')}
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => navigate('/projects/new')}
        >
          {t('projects.newProject')}
        </Button>
      </Box>

      <Grid container spacing={3}>
        {projects.map((project) => (
          <Grid key={project.project_id} sx={{ width: { xs: '100%', md: '50%', lg: '33.333%' } }}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
                  <Typography variant="h6" component="h2">
                    {project.name}
                  </Typography>
                  <Box>
                    <IconButton
                      size="small"
                      onClick={() => navigate(`/projects/${project.project_id}/edit`)}
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      size="small"
                      color="error"
                      onClick={() => openDeleteDialog(project.project_id)}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Box>
                </Box>

                <Typography color="text.secondary" sx={{ mb: 2 }}>
                  {project.description || t('projects.noDescription')}
                </Typography>

                <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                  <Chip label={project.language} color="primary" size="small" />
                  <Chip
                    label={project.apply_global_rules ? t('projects.globalRules') : t('projects.noGlobalRules')}
                    color={project.apply_global_rules ? 'success' : 'default'}
                    size="small"
                  />
                </Box>

                <Typography variant="body2" color="text.secondary">
                  {t('projects.rules')}: {project.rule_count}
                </Typography>

                <Button
                  variant="outlined"
                  size="small"
                  sx={{ mt: 2 }}
                  onClick={() => navigate(`/projects/${project.project_id}/rules/new`)}
                >
                  {t('projects.addRule')}
                </Button>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>{t('projects.deleteProject')}</DialogTitle>
        <DialogContent>
          <Typography>
            {t('projects.deleteConfirm')}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>{t('common.cancel')}</Button>
          <Button onClick={handleDelete} color="error" variant="contained">
            {t('common.delete')}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default ProjectList;
