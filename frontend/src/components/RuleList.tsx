import React, { useState, useEffect, useCallback } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Button,
  Chip,
  Box,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
} from '@mui/material';
import { Edit as EditIcon, Delete as DeleteIcon, Add as AddIcon } from '@mui/icons-material';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Rule } from '../types';
import { api } from '../services/api';
import { useAuth } from '../contexts/AuthContext';

const RuleList: React.FC = () => {
  const [rules, setRules] = useState<Rule[]>([]);
  const [loading, setLoading] = useState(true);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [ruleToDelete, setRuleToDelete] = useState<string | null>(null);
  const [projectName, setProjectName] = useState<string>('');
  const navigate = useNavigate();
  const { projectId } = useParams<{ projectId: string }>();
  const { t } = useTranslation();
  const { user } = useAuth();
  const isAdmin = user?.role === 'admin';

  const loadRules = useCallback(async () => {
    if (!projectId) return;
    try {
      const response = await api.get(`/rules?project_id=${projectId}`);
      setRules(response.data.rules || []);
    } catch (error) {
      console.error('Failed to load rules:', error);
    } finally {
      setLoading(false);
    }
  }, [projectId]);

  const loadProjectInfo = useCallback(async () => {
    if (!projectId) return;
    try {
      const response = await api.get(`/projects`);
      const project = response.data.projects.find((p: any) => p.project_id === projectId);
      if (project) {
        setProjectName(project.name);
      }
    } catch (error) {
      console.error('Failed to load project info:', error);
    }
  }, [projectId]);

  useEffect(() => {
    if (projectId) {
      loadRules();
      loadProjectInfo();
    }
  }, [projectId, loadRules, loadProjectInfo]);

  const handleDelete = async () => {
    if (!projectId || !ruleToDelete) return;
    try {
      await api.delete(`/rules/${projectId}/${ruleToDelete}`);
      await loadRules();
      setDeleteDialogOpen(false);
      setRuleToDelete(null);
    } catch (error) {
      console.error('Failed to delete rule:', error);
    }
  };

  const openDeleteDialog = (ruleId: string) => {
    setRuleToDelete(ruleId);
    setDeleteDialogOpen(true);
  };

  const getSeverityColor = (severity: string) => {
    switch (severity.toLowerCase()) {
      case 'error':
        return 'error';
      case 'warning':
        return 'warning';
      case 'info':
        return 'info';
      default:
        return 'default';
    }
  };

  const getTypeColor = (type: string) => {
    switch (type.toLowerCase()) {
      case 'security':
        return 'error';
      case 'performance':
        return 'warning';
      case 'naming':
        return 'primary';
      case 'formatting':
        return 'secondary';
      default:
        return 'default';
    }
  };

  if (loading) {
    return <Typography>{t('rules.loadingRules')}</Typography>;
  }

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Box>
          <Typography variant="h4" component="h1">
            {t('rules.title')} - {projectName}
          </Typography>
          <Typography variant="subtitle1" color="text.secondary">
            {t('rules.projectRules')}: {projectId}
          </Typography>
        </Box>
        {isAdmin && (
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => navigate(`/projects/${projectId}/rules/new`)}
          >
            {t('rules.newRule')}
          </Button>
        )}
      </Box>

      {rules.length === 0 ? (
        <Card>
          <CardContent>
            <Typography variant="h6" color="text.secondary" align="center">
              {t('rules.noRules')}
            </Typography>
            <Typography color="text.secondary" align="center">
              {t('rules.createFirstRule')}
            </Typography>
          </CardContent>
        </Card>
      ) : (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>{t('rules.name')}</TableCell>
                <TableCell>{t('rules.description')}</TableCell>
                <TableCell>{t('rules.type')}</TableCell>
                <TableCell>{t('rules.severity')}</TableCell>
                <TableCell>{t('rules.pattern')}</TableCell>
                <TableCell>{t('rules.status')}</TableCell>
                {isAdmin && <TableCell>{t('common.actions')}</TableCell>}
              </TableRow>
            </TableHead>
            <TableBody>
              {rules.map((rule) => (
                <TableRow key={rule.rule_id}>
                  <TableCell>
                    <Typography variant="subtitle2" fontWeight="bold">
                      {rule.name}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {rule.rule_id}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2">
                      {rule.description || t('rules.noDescription')}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Chip 
                      label={rule.type} 
                      color={getTypeColor(rule.type) as any}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>
                    <Chip 
                      label={rule.severity} 
                      color={getSeverityColor(rule.severity) as any}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2" fontFamily="monospace" fontSize="0.8rem">
                      {rule.pattern}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Chip 
                      label={rule.is_active ? t('rules.active') : t('rules.inactive')}
                      color={rule.is_active ? 'success' : 'default'}
                      size="small"
                    />
                  </TableCell>
                  {isAdmin && (
                    <TableCell>
                      <Box sx={{ display: 'flex', gap: 0.5 }}>
                        <IconButton
                          size="small"
                          onClick={() => navigate(`/projects/${projectId}/rules/${rule.rule_id}/edit`)}
                        >
                          <EditIcon />
                        </IconButton>
                        <IconButton
                          size="small"
                          color="error"
                          onClick={() => openDeleteDialog(rule.rule_id)}
                        >
                          <DeleteIcon />
                        </IconButton>
                      </Box>
                    </TableCell>
                  )}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}

      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>{t('rules.deleteRule')}</DialogTitle>
        <DialogContent>
          <Typography>
            {t('rules.deleteConfirm')}
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

export default RuleList;
