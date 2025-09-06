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
  MenuItem,
  FormControl,
  InputLabel,
  Select,
  FormControlLabel,
  Checkbox,
  Alert,
} from '@mui/material';
import { 
  Edit as EditIcon, 
  Delete as DeleteIcon, 
  Add as AddIcon,
  FileDownload as ExportIcon,
  FileUpload as ImportIcon
} from '@mui/icons-material';
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
  const [exportDialogOpen, setExportDialogOpen] = useState(false);
  const [importDialogOpen, setImportDialogOpen] = useState(false);
  const [exportFormat, setExportFormat] = useState('json');
  const [selectedRules, setSelectedRules] = useState<string[]>([]);
  const [importFile, setImportFile] = useState<File | null>(null);
  const [importOverwrite, setImportOverwrite] = useState(false);
  const navigate = useNavigate();
  const { projectId } = useParams<{ projectId: string }>();
  const { t } = useTranslation();
  const { permissions } = useAuth();
  const canManageRules = permissions.manageRules;

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

  const handleExport = async () => {
    try {
      const response = await api.post('/rules/export', {
        project_id: projectId,
        format: exportFormat,
        rule_ids: selectedRules.length > 0 ? selectedRules : undefined,
      });
      
      // ファイルダウンロード
      const blob = new Blob([JSON.stringify(response.data, null, 2)], { type: 'application/json' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `rules-export-${projectId}-${new Date().toISOString().split('T')[0]}.${exportFormat}`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
      
      setExportDialogOpen(false);
      setSelectedRules([]);
    } catch (error: any) {
      console.error('Export failed:', error);
    }
  };

  const handleImport = async () => {
    if (!importFile) return;
    
    try {
      const text = await importFile.text();
      const data = JSON.parse(text);
      
      await api.post('/rules/import', {
        project_id: projectId,
        data,
        overwrite: importOverwrite,
      });
      
      setImportDialogOpen(false);
      setImportFile(null);
      setImportOverwrite(false);
      
      // ルール一覧を再読み込み
      loadRules();
    } catch (error: any) {
      console.error('Import failed:', error);
    }
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
        {canManageRules && (
          <Box sx={{ display: 'flex', gap: 1 }}>
            <Button
              variant="outlined"
              startIcon={<ExportIcon />}
              onClick={() => setExportDialogOpen(true)}
            >
              エクスポート
            </Button>
            <Button
              variant="outlined"
              startIcon={<ImportIcon />}
              onClick={() => setImportDialogOpen(true)}
            >
              インポート
            </Button>
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={() => navigate(`/projects/${projectId}/rules/new`)}
            >
              {t('rules.newRule')}
            </Button>
          </Box>
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
                {canManageRules && (
                  <TableCell padding="checkbox">
                    <Checkbox
                      indeterminate={selectedRules.length > 0 && selectedRules.length < rules.length}
                      checked={rules.length > 0 && selectedRules.length === rules.length}
                      onChange={(e) => {
                        if (e.target.checked) {
                          setSelectedRules(rules.map(rule => rule.rule_id));
                        } else {
                          setSelectedRules([]);
                        }
                      }}
                    />
                  </TableCell>
                )}
                <TableCell>{t('rules.name')}</TableCell>
                <TableCell>{t('rules.description')}</TableCell>
                <TableCell>{t('rules.type')}</TableCell>
                <TableCell>{t('rules.severity')}</TableCell>
                <TableCell>{t('rules.pattern')}</TableCell>
                <TableCell>{t('rules.status')}</TableCell>
                {canManageRules && <TableCell>{t('common.actions')}</TableCell>}
              </TableRow>
            </TableHead>
            <TableBody>
              {rules.map((rule) => (
                <TableRow key={rule.rule_id}>
                  {canManageRules && (
                    <TableCell padding="checkbox">
                      <Checkbox
                        checked={selectedRules.includes(rule.rule_id)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setSelectedRules([...selectedRules, rule.rule_id]);
                          } else {
                            setSelectedRules(selectedRules.filter(id => id !== rule.rule_id));
                          }
                        }}
                      />
                    </TableCell>
                  )}
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
                  {canManageRules && (
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

      {/* エクスポートダイアログ */}
      <Dialog open={exportDialogOpen} onClose={() => setExportDialogOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>ルールエクスポート</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
            <FormControl fullWidth>
              <InputLabel>フォーマット</InputLabel>
              <Select
                value={exportFormat}
                onChange={(e) => setExportFormat(e.target.value)}
                label="フォーマット"
              >
                <MenuItem value="json">JSON</MenuItem>
                <MenuItem value="yaml">YAML</MenuItem>
                <MenuItem value="csv">CSV</MenuItem>
              </Select>
            </FormControl>
            {selectedRules.length > 0 && (
              <Alert severity="info">
                選択されたルール ({selectedRules.length}件) をエクスポートします
              </Alert>
            )}
            {selectedRules.length === 0 && (
              <Alert severity="info">
                プロジェクトの全ルールをエクスポートします
              </Alert>
            )}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setExportDialogOpen(false)}>キャンセル</Button>
          <Button variant="contained" onClick={handleExport}>
            エクスポート
          </Button>
        </DialogActions>
      </Dialog>

      {/* インポートダイアログ */}
      <Dialog open={importDialogOpen} onClose={() => setImportDialogOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>ルールインポート</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
            <Button
              variant="outlined"
              component="label"
              startIcon={<ImportIcon />}
            >
              ファイルを選択
              <input
                type="file"
                hidden
                accept=".json,.yaml,.yml"
                onChange={(e) => setImportFile(e.target.files?.[0] || null)}
              />
            </Button>
            {importFile && (
              <Typography variant="body2" color="text.secondary">
                選択されたファイル: {importFile.name}
              </Typography>
            )}
            <FormControlLabel
              control={
                <Checkbox
                  checked={importOverwrite}
                  onChange={(e) => setImportOverwrite(e.target.checked)}
                />
              }
              label="既存のルールを上書きする"
            />
            <Alert severity="warning">
              インポート前に現在のルールをエクスポートしてバックアップを取ることをお勧めします。
            </Alert>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setImportDialogOpen(false)}>キャンセル</Button>
          <Button 
            variant="contained" 
            onClick={handleImport} 
            disabled={!importFile}
          >
            インポート
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default RuleList;
