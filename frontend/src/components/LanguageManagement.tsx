import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  FormControlLabel,
  Checkbox,
  Chip,
  Box,
  Alert,
  Snackbar,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Save as SaveIcon,
  Cancel as CancelIcon,
} from '@mui/icons-material';
import { useTranslation } from 'react-i18next';
import { api } from '../services/api';

interface Language {
  code: string;
  name: string;
  description: string;
  icon: string;
  color: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

const LanguageManagement: React.FC = () => {
  const { t } = useTranslation();
  const [languages, setLanguages] = useState<Language[]>([]);
  const [loading, setLoading] = useState(false);
  const [open, setOpen] = useState(false);
  const [editingLanguage, setEditingLanguage] = useState<Language | null>(null);
  const [snackbar, setSnackbar] = useState<{
    open: boolean;
    message: string;
    severity: 'success' | 'error';
  }>({ open: false, message: '', severity: 'success' });

  const [formData, setFormData] = useState({
    code: '',
    name: '',
    description: '',
    icon: '',
    color: '#666666',
    is_active: true,
  });

  const loadLanguages = async () => {
    try {
      setLoading(true);
      const response = await api.get('/languages');
      setLanguages(response.data.languages);
    } catch (error) {
      console.error('Failed to load languages:', error);
      setSnackbar({
        open: true,
        message: '言語の読み込みに失敗しました',
        severity: 'error',
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadLanguages();
  }, []);

  const handleOpenDialog = (language?: Language) => {
    if (language) {
      setEditingLanguage(language);
      setFormData({
        code: language.code,
        name: language.name,
        description: language.description,
        icon: language.icon,
        color: language.color,
        is_active: language.is_active,
      });
    } else {
      setEditingLanguage(null);
      setFormData({
        code: '',
        name: '',
        description: '',
        icon: '',
        color: '#666666',
        is_active: true,
      });
    }
    setOpen(true);
  };

  const handleCloseDialog = () => {
    setOpen(false);
    setEditingLanguage(null);
    setFormData({
      code: '',
      name: '',
      description: '',
      icon: '',
      color: '#666666',
      is_active: true,
    });
  };

  const handleSubmit = async () => {
    try {
      setLoading(true);
      if (editingLanguage) {
        await api.put(`/languages/${editingLanguage.code}`, formData);
        setSnackbar({
          open: true,
          message: '言語が更新されました',
          severity: 'success',
        });
      } else {
        await api.post('/languages', formData);
        setSnackbar({
          open: true,
          message: '言語が作成されました',
          severity: 'success',
        });
      }
      handleCloseDialog();
      loadLanguages();
    } catch (error: any) {
      console.error('Failed to save language:', error);
      setSnackbar({
        open: true,
        message: error.response?.data?.error || '保存に失敗しました',
        severity: 'error',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (code: string) => {
    if (!window.confirm('この言語を削除してもよろしいですか？')) {
      return;
    }

    try {
      setLoading(true);
      await api.delete(`/languages/${code}`);
      setSnackbar({
        open: true,
        message: '言語が削除されました',
        severity: 'success',
      });
      loadLanguages();
    } catch (error: any) {
      console.error('Failed to delete language:', error);
      setSnackbar({
        open: true,
        message: error.response?.data?.error || '削除に失敗しました',
        severity: 'error',
      });
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
    <Box>
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
            <Typography variant="h5" component="h2">
              言語管理
            </Typography>
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={() => handleOpenDialog()}
              disabled={loading}
            >
              新しい言語を追加
            </Button>
          </Box>

          <TableContainer component={Paper}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>コード</TableCell>
                  <TableCell>名前</TableCell>
                  <TableCell>説明</TableCell>
                  <TableCell>アイコン</TableCell>
                  <TableCell>色</TableCell>
                  <TableCell>ステータス</TableCell>
                  <TableCell>操作</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {languages.map((language) => (
                  <TableRow key={language.code}>
                    <TableCell>{language.code}</TableCell>
                    <TableCell>{language.name}</TableCell>
                    <TableCell>{language.description}</TableCell>
                    <TableCell>{language.icon}</TableCell>
                    <TableCell>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Box
                          sx={{
                            width: 20,
                            height: 20,
                            backgroundColor: language.color,
                            borderRadius: 1,
                          }}
                        />
                        {language.color}
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={language.is_active ? 'アクティブ' : '非アクティブ'}
                        color={language.is_active ? 'success' : 'default'}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      <IconButton
                        onClick={() => handleOpenDialog(language)}
                        disabled={loading}
                        size="small"
                      >
                        <EditIcon />
                      </IconButton>
                      <IconButton
                        onClick={() => handleDelete(language.code)}
                        disabled={loading}
                        size="small"
                        color="error"
                      >
                        <DeleteIcon />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>

      <Dialog open={open} onClose={handleCloseDialog} maxWidth="sm" fullWidth>
        <DialogTitle>
          {editingLanguage ? '言語を編集' : '新しい言語を追加'}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
            <TextField
              label="言語コード"
              value={formData.code}
              onChange={(e) => handleChange('code', e.target.value)}
              required
              disabled={!!editingLanguage}
              helperText="言語の一意識別子（例: dart, scala）"
            />
            <TextField
              label="言語名"
              value={formData.name}
              onChange={(e) => handleChange('name', e.target.value)}
              required
              helperText="言語の表示名（例: Dart, Scala）"
            />
            <TextField
              label="説明"
              value={formData.description}
              onChange={(e) => handleChange('description', e.target.value)}
              multiline
              rows={2}
              helperText="言語の説明（オプション）"
            />
            <TextField
              label="アイコン"
              value={formData.icon}
              onChange={(e) => handleChange('icon', e.target.value)}
              helperText="アイコン名（例: dart, scala）"
            />
            <TextField
              label="色"
              type="color"
              value={formData.color}
              onChange={(e) => handleChange('color', e.target.value)}
              helperText="言語のテーマカラー"
            />
            <FormControlLabel
              control={
                <Checkbox
                  checked={formData.is_active}
                  onChange={(e) => handleChange('is_active', e.target.checked)}
                />
              }
              label="アクティブ"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog} disabled={loading}>
            キャンセル
          </Button>
          <Button
            onClick={handleSubmit}
            variant="contained"
            disabled={loading}
            startIcon={editingLanguage ? <SaveIcon /> : <AddIcon />}
          >
            {editingLanguage ? '更新' : '作成'}
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={() => setSnackbar({ ...snackbar, open: false })}
      >
        <Alert
          onClose={() => setSnackbar({ ...snackbar, open: false })}
          severity={snackbar.severity}
        >
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default LanguageManagement;
