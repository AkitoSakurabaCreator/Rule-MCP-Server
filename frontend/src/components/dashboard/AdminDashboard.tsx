import React, { useEffect, useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  IconButton,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Tabs,
  Tab,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Snackbar,
  Tooltip,
  FormControlLabel,
  Switch,
  Alert,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  People as PeopleIcon,
  Security as SecurityIcon,
  Settings as SettingsIcon,
  Analytics as AnalyticsIcon,
  Code as CodeIcon,
  FileDownload as ExportIcon,
  FileUpload as ImportIcon,
} from '@mui/icons-material';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import { useTranslation } from 'react-i18next';
import { adminApi, AdminStats as AdminStatsType, MCPStats as MCPStatsType, SystemLog as SystemLogType, Role as RoleType } from '../../services/adminApi';
import { useAuth } from '../../contexts/AuthContext';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`admin-tabpanel-${index}`}
      aria-labelledby={`admin-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

interface User {
  id: number;
  username: string;
  email: string;
  fullName?: string;
  role: string;
  status?: string;
  lastLogin: string;
}

interface ApiKey {
  id: number;
  name: string;
  key: string;
  accessLevel: string;
  status: string;
  createdAt: string;
  lastUsed: string;
}

const AdminDashboard: React.FC = () => {
  const { t } = useTranslation();
  const { permissions } = useAuth();
  const canManageUsers = permissions.manageUsers;
  const canManageRoles = permissions.manageRoles;
  const [tabValue, setTabValue] = useState(0);
  const [users, setUsers] = useState<User[]>([]);
  const [apiKeys, setApiKeys] = useState<ApiKey[]>([]);
  const [stats, setStats] = useState<AdminStatsType>({
    totalUsers: 0,
    totalProjects: 0,
    totalRules: 0,
    activeApiKeys: 0,
    mcpRequests: 0,
    activeSessions: 0,
    systemLoad: '0%',
  });
  const [mcpStats, setMcpStats] = useState<MCPStatsType[]>([]);
  const [mcpPerf, setMcpPerf] = useState<{ avgMs: number; successRate: number; errorRate: number; p95Ms: number }>({ avgMs: 0, successRate: 0, errorRate: 0, p95Ms: 0 });
  const [languages, setLanguages] = useState<{ code: string; name: string; description?: string; icon?: string; color?: string; isActive?: boolean }[]>([]);
  const [openAddLanguage, setOpenAddLanguage] = useState(false);
  const [editLanguageCode, setEditLanguageCode] = useState<string | null>(null);
  const [languageForm, setLanguageForm] = useState<{ code: string; name: string; description: string; icon: string; color: string; isActive: boolean }>({ code: '', name: '', description: '', icon: '', color: '#007acc', isActive: true });
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [systemLogs, setSystemLogs] = useState<SystemLogType[]>([]);
  const [logCounts, setLogCounts] = useState<{ INFO: number; WARN: number; ERROR: number }>({ INFO: 0, WARN: 0, ERROR: 0 });
  const [openAddUser, setOpenAddUser] = useState(false);
  const [openEditUser, setOpenEditUser] = useState<null | User>(null);
  const [userForm, setUserForm] = useState({
    username: '',
    email: '',
    fullName: '',
    role: 'user',
    isActive: true,
  });
  const [snackbar, setSnackbar] = useState<{ open: boolean; message: string; severity: 'success' | 'error' }>({
    open: false,
    message: '',
    severity: 'success',
  });
  const [bulkExportDialogOpen, setBulkExportDialogOpen] = useState(false);
  const [bulkImportDialogOpen, setBulkImportDialogOpen] = useState(false);
  const [exportFormat, setExportFormat] = useState('json');
  const [exportScope, setExportScope] = useState('all');
  const [importFile, setImportFile] = useState<File | null>(null);
  const [importOverwrite, setImportOverwrite] = useState(false);
  const [openAddApiKey, setOpenAddApiKey] = useState(false);
  const [apiKeyForm, setApiKeyForm] = useState({ name: '', accessLevel: 'user' });
  const [openDeleteApiKey, setOpenDeleteApiKey] = useState<null | ApiKey>(null);
  const [openEditApiKey, setOpenEditApiKey] = useState<null | ApiKey>(null);
  const [apiKeyEditForm, setApiKeyEditForm] = useState<{ name: string; description: string; isActive: boolean }>({ name: '', description: '', isActive: true });
  const [roles, setRoles] = useState<RoleType[]>([]);
  const [openAddRole, setOpenAddRole] = useState(false);
  const [openEditRole, setOpenEditRole] = useState<null | RoleType>(null);
  const [roleForm, setRoleForm] = useState<RoleType>({ name: '', description: '', permissions: { manage_users: false, manage_rules: true, manage_roles: false }, is_active: true });
  const [pendingUsers, setPendingUsers] = useState<User[]>([]);
  const applyPreset = (preset: 'readonly' | 'editor' | 'admin') => {
    if (preset === 'readonly') setRoleForm((r) => ({ ...r, permissions: { manage_users: false, manage_rules: false, manage_roles: false } }));
    if (preset === 'editor') setRoleForm((r) => ({ ...r, permissions: { manage_users: false, manage_rules: true, manage_roles: false } }));
    if (preset === 'admin') setRoleForm((r) => ({ ...r, permissions: { manage_users: true, manage_rules: true, manage_roles: true } }));
  };

  // 設定状態
  const [settings, setSettings] = useState<{ defaultAccessLevel: string; requestsPerMinute: number }>({ defaultAccessLevel: 'public', requestsPerMinute: 100 });

  useEffect(() => {
    let mounted = true;
    const load = async () => {
      try {
        setLoading(true);
        setError(null);
        const [s, u, k, m, logs, perf, conf, langs, pending] = await Promise.all([
          adminApi.getStats(),
          canManageUsers ? adminApi.getUsers() : Promise.resolve([]),
          adminApi.getApiKeys(),
          adminApi.getMCPStats(),
          adminApi.getSystemLogs(),
          adminApi.getMCPPerformance(),
          adminApi.getSettings(),
          adminApi.getLanguages(),
          canManageUsers ? adminApi.getPendingUsers() : Promise.resolve([]),
        ]);
        const rls = canManageRoles ? await adminApi.getRoles() : [];
        if (!mounted) return;
        setStats(s);
        setUsers(
          (u as any[]).map((x: any) => ({
            id: x.id,
            username: x.username,
            email: x.email,
            fullName: x.fullName,
            role: x.role,
            status: x.isActive ? 'active' : 'inactive',
            lastLogin: x.lastLogin,
          }))
        );
        setApiKeys(k as unknown as ApiKey[]);
        setMcpStats(m);
        setSystemLogs(Array.isArray(logs) ? logs : []);
        // 集計（INFO/WARN/ERROR）
        const counts = (Array.isArray(logs) ? logs : []).reduce((acc: any, l: any) => {
          const level = (l.level || '').toUpperCase();
          if (level === 'WARN') acc.WARN += 1;
          else if (level === 'ERROR') acc.ERROR += 1;
          else acc.INFO += 1;
          return acc;
        }, { INFO: 0, WARN: 0, ERROR: 0 });
        setLogCounts(counts);
        setRoles(rls);
        // MCPパフォーマンス
        setMcpPerf(perf as any);
        // 設定
        const da = (conf as any).defaultAccessLevel || 'public';
        const rpm = parseInt((conf as any).requestsPerMinute || '100', 10) || 100;
        setSettings({ defaultAccessLevel: da, requestsPerMinute: rpm });
        setLanguages(Array.isArray(langs) ? langs : []);
        setPendingUsers(
          (pending as any[]).map((x: any) => ({
            id: x.id,
            username: x.username,
            email: x.email,
            fullName: x.fullName,
            role: x.role,
            status: 'pending',
            lastLogin: x.createdAt,
          }))
        );
      } catch (e: any) {
        setError(e?.message || 'Failed to load admin data');
      } finally {
        if (mounted) setLoading(false);
      }
    };
    load();
    return () => {
      mounted = false;
    };
  }, [canManageUsers, canManageRoles]);

  // CPU率（systemLoad）とMCPパフォーマンスを定期更新（5秒間隔）
  useEffect(() => {
    let timer: any;
    const tick = async () => {
      try {
        const s = await adminApi.getStats();
        setStats(s);
        const perf = await adminApi.getMCPPerformance();
        setMcpPerf(perf as any);
      } catch (e: any) {
        // 一時的なエラーを無視
      }
    };
    // 初回即時 + 5秒間隔
    tick();
    timer = setInterval(tick, 5000);
    return () => {
      if (timer) clearInterval(timer);
    };
  }, []);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const handleBulkExport = async () => {
    try {
      setLoading(true);
      const response = await adminApi.bulkExportRules({
        format: exportFormat,
        scope: exportScope,
      });
      
      // ファイルダウンロード
      const blob = new Blob([JSON.stringify(response, null, 2)], { type: 'application/json' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `rules-export-${exportScope}-${new Date().toISOString().split('T')[0]}.${exportFormat}`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
      
      setSnackbar({ open: true, message: 'ルールのエクスポートが完了しました', severity: 'success' });
      setBulkExportDialogOpen(false);
    } catch (error: any) {
      setSnackbar({ open: true, message: error.message || 'エクスポートに失敗しました', severity: 'error' });
    } finally {
      setLoading(false);
    }
  };

  const handleBulkImport = async () => {
    if (!importFile) return;
    
    try {
      setLoading(true);
      const text = await importFile.text();
      const data = JSON.parse(text);
      
      await adminApi.bulkImportRules({
        data,
        overwrite: importOverwrite,
      });
      
      setSnackbar({ open: true, message: 'ルールのインポートが完了しました', severity: 'success' });
      setBulkImportDialogOpen(false);
      setImportFile(null);
      setImportOverwrite(false);
      
      // データを再読み込み
      window.location.reload();
    } catch (error: any) {
      setSnackbar({ open: true, message: error.message || 'インポートに失敗しました', severity: 'error' });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" sx={{ mb: 3 }}>
        {t('dashboard.title')}
      </Typography>

      {loading && (
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          {t('common.loading')}
        </Typography>
      )}
      {error && (
        <Typography variant="body2" color="error" sx={{ mb: 2 }}>
          {error}
        </Typography>
      )}

      {/* 統計カード */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid sx={{ width: { xs: '100%', sm: '50%', md: '25%' } }}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                {t('dashboard.totalUsers')}
              </Typography>
              <Typography variant="h4" component="div">
                {stats.totalUsers || 0}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid sx={{ width: { xs: '100%', sm: '50%', md: '25%' } }}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                {t('dashboard.totalProjects')}
              </Typography>
              <Typography variant="h4" component="div">
                {stats.totalProjects || 0}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid sx={{ width: { xs: '100%', sm: '50%', md: '25%' } }}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                {t('dashboard.totalRules')}
              </Typography>
              <Typography variant="h4" component="div">
                {stats.totalRules || 0}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid sx={{ width: { xs: '100%', sm: '50%', md: '25%' } }}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                {t('dashboard.activeApiKeys')}
              </Typography>
              <Typography variant="h4" component="div">
                {stats.activeApiKeys || 0}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid sx={{ width: { xs: '100%', sm: '50%', md: '25%' } }}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                {t('dashboard.mcpRequests')}
              </Typography>
              <Typography variant="h4" component="div">
                {stats.mcpRequests || 0}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                {t('dashboard.last24Hours')}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid sx={{ width: { xs: '100%', sm: '50%', md: '25%' } }}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                {t('dashboard.activeSessions')}
              </Typography>
              <Typography variant="h4" component="div">
                {stats.activeSessions || 0}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                {t('dashboard.currentUsers')}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid sx={{ width: { xs: '100%', sm: '50%', md: '25%' } }}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                {t('dashboard.systemLoad')}
              </Typography>
              <Typography variant="h4" component="div">
                {stats.systemLoad || '0%'}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                {t('dashboard.cpuUsage')}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* 一括操作ボタン */}
      <Box sx={{ mb: 3, display: 'flex', gap: 2 }}>
        <Button
          variant="outlined"
          startIcon={<ExportIcon />}
          onClick={() => setBulkExportDialogOpen(true)}
          disabled={loading}
        >
          全ルールエクスポート
        </Button>
        <Button
          variant="outlined"
          startIcon={<ImportIcon />}
          onClick={() => setBulkImportDialogOpen(true)}
          disabled={loading}
        >
          全ルールインポート
        </Button>
      </Box>

      {/* タブナビゲーション */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
        <Tabs value={tabValue} onChange={handleTabChange} aria-label="admin tabs">
          <Tab 
            icon={<PeopleIcon />} 
            label={t('dashboard.users')} 
            iconPosition="start"
          />
          <Tab 
            icon={<PeopleIcon />} 
            label="承認待ちユーザー" 
            iconPosition="start"
            disabled={!canManageUsers}
          />
          <Tab 
            icon={<SecurityIcon />} 
            label={t('dashboard.apiKeys')} 
            iconPosition="start"
          />
          <Tab 
            icon={<SettingsIcon />} 
            label={t('dashboard.settings')} 
            iconPosition="start"
          />
          <Tab 
            icon={<AnalyticsIcon />} 
            label={t('dashboard.analytics')} 
            iconPosition="start"
          />
          <Tab 
            icon={<CodeIcon />} 
            label={t('dashboard.mcpMonitoring')} 
            iconPosition="start"
          />
          <Tab 
            icon={<SecurityIcon />} 
            label={t('dashboard.systemLogs')} 
            iconPosition="start"
          />
          <Tab 
            icon={<SecurityIcon />} 
            label={t('dashboard.roles') || 'Roles'} 
            disabled={!canManageRoles}
            iconPosition="start"
          />
          <Tab 
            icon={<CodeIcon />} 
            label={t('dashboard.languages') || 'Languages'} 
            iconPosition="start"
          />
        </Tabs>
      </Box>

      {/* ユーザー管理タブ */}
      <TabPanel value={tabValue} index={0}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h5">
            {t('dashboard.userManagement')}
          </Typography>
          {canManageUsers && (
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => {
              setUserForm({ username: '', email: '', fullName: '', role: 'user', isActive: true });
              setOpenAddUser(true);
            }}
          >
            {t('dashboard.addUser')}
          </Button>
          )}
        </Box>
        
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>{t('dashboard.username')}</TableCell>
                <TableCell>{t('dashboard.email')}</TableCell>
                <TableCell>{t('dashboard.role')}</TableCell>
                <TableCell>{t('dashboard.status')}</TableCell>
                <TableCell>{t('dashboard.lastLogin')}</TableCell>
                {canManageUsers && <TableCell>{t('dashboard.actions')}</TableCell>}
              </TableRow>
            </TableHead>
            <TableBody>
              {users.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>{user.username}</TableCell>
                  <TableCell>{user.email}</TableCell>
                  <TableCell>
                    <Chip 
                      label={user.role} 
                      color={user.role === 'admin' ? 'error' : 'default'}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>
                    <Chip label={user.status || 'inactive'} color={user.status === 'active' ? 'success' : 'default'} size="small" />
                  </TableCell>
                  <TableCell>{user.lastLogin}</TableCell>
                  {canManageUsers && (
                  <TableCell>
                    <IconButton
                      size="small"
                      color="primary"
                      onClick={() => {
                        setUserForm({
                          username: user.username,
                          email: user.email,
                          fullName: user.fullName || '',
                          role: user.role,
                          isActive: user.status === 'active',
                        });
                        setOpenEditUser(user);
                      }}
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      size="small"
                      color="error"
                      onClick={async () => {
                        if (!window.confirm(t('dashboard.deleteUserConfirm') || 'Delete this user?')) return;
                        try {
                          await adminApi.deleteUser(user.id);
                          setUsers(prev => prev.filter(u => u.id !== user.id));
                        } catch (e: any) {
                          setError(e?.message || 'Failed to delete user');
                        }
                      }}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                  )}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </TabPanel>

      {/* 承認待ちユーザー管理タブ */}
      <TabPanel value={tabValue} index={1}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h5">
            承認待ちユーザー管理
          </Typography>
        </Box>
        
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ユーザー名</TableCell>
                <TableCell>メールアドレス</TableCell>
                <TableCell>フルネーム</TableCell>
                <TableCell>登録日時</TableCell>
                <TableCell>操作</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {pendingUsers.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>{user.username}</TableCell>
                  <TableCell>{user.email}</TableCell>
                  <TableCell>{user.fullName || '-'}</TableCell>
                  <TableCell>{new Date(user.lastLogin).toLocaleString()}</TableCell>
                  <TableCell>
                    <Button
                      size="small"
                      variant="contained"
                      color="success"
                      sx={{ mr: 1 }}
                      onClick={async () => {
                        try {
                          await adminApi.approveUser(user.id, true);
                          setPendingUsers(prev => prev.filter(u => u.id !== user.id));
                          setSnackbar({ open: true, message: 'ユーザーを承認しました', severity: 'success' });
                        } catch (e: any) {
                          setSnackbar({ open: true, message: e?.message || '承認に失敗しました', severity: 'error' });
                        }
                      }}
                    >
                      承認
                    </Button>
                    <Button
                      size="small"
                      variant="contained"
                      color="error"
                      onClick={async () => {
                        if (!window.confirm('このユーザーを拒否しますか？')) return;
                        try {
                          await adminApi.approveUser(user.id, false);
                          setPendingUsers(prev => prev.filter(u => u.id !== user.id));
                          setSnackbar({ open: true, message: 'ユーザーを拒否しました', severity: 'success' });
                        } catch (e: any) {
                          setSnackbar({ open: true, message: e?.message || '拒否に失敗しました', severity: 'error' });
                        }
                      }}
                    >
                      拒否
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
        
        {pendingUsers.length === 0 && (
          <Box sx={{ textAlign: 'center', py: 4 }}>
            <Typography variant="body1" color="text.secondary">
              承認待ちのユーザーはありません
            </Typography>
          </Box>
        )}
      </TabPanel>

      {/* 言語管理タブ */}
      <TabPanel value={tabValue} index={8}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h5">{t('dashboard.languages')}</Typography>
          {canManageRoles && (
            <Button variant="contained" startIcon={<AddIcon />} onClick={() => { setLanguageForm({ code: '', name: '', description: '', icon: '', color: '#007acc', isActive: true }); setOpenAddLanguage(true); }}> 
              {t('dashboard.addLanguage')}
            </Button>
          )}
        </Box>
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>{t('dashboard.code')}</TableCell>
                <TableCell>{t('dashboard.name')}</TableCell>
                <TableCell>{t('dashboard.description')}</TableCell>
                <TableCell>{t('dashboard.icon')}</TableCell>
                <TableCell>{t('dashboard.color')}</TableCell>
                <TableCell>{t('dashboard.status')}</TableCell>
                {canManageRoles && <TableCell>{t('dashboard.actions')}</TableCell>}
              </TableRow>
            </TableHead>
            <TableBody>
              {(languages || []).map((lang) => (
                <TableRow key={lang.code}>
                  <TableCell>{lang.code}</TableCell>
                  <TableCell>{lang.name}</TableCell>
                  <TableCell>{lang.description || ''}</TableCell>
                  <TableCell>{lang.icon || ''}</TableCell>
                  <TableCell>
                    <Chip label={lang.color || ''} size="small" />
                  </TableCell>
                  <TableCell>
                    <Chip label={lang.isActive ? 'active' : 'inactive'} color={lang.isActive ? 'success' : 'default'} size="small" />
                  </TableCell>
                  {canManageRoles && (
                    <TableCell>
                      <IconButton size="small" color="primary" onClick={() => { setLanguageForm({ code: lang.code, name: lang.name, description: lang.description || '', icon: lang.icon || '', color: lang.color || '#007acc', isActive: !!lang.isActive }); setEditLanguageCode(lang.code); setOpenAddLanguage(true); }}>
                        <EditIcon />
                      </IconButton>
                      <IconButton size="small" color="error" onClick={async () => { await adminApi.deleteLanguage(lang.code); setLanguages((prev) => prev.filter((l) => l.code !== lang.code)); }}>
                        <DeleteIcon />
                      </IconButton>
                    </TableCell>
                  )}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>

        <Dialog open={openAddLanguage} onClose={() => { setOpenAddLanguage(false); setEditLanguageCode(null); }} fullWidth maxWidth="sm">
          <DialogTitle>{editLanguageCode ? t('dashboard.editLanguage') : t('dashboard.addLanguage')}</DialogTitle>
          <DialogContent>
            <Grid container spacing={2} sx={{ mt: 1 }}>
              {!editLanguageCode && (
                <Grid sx={{ width: '100%' }}>
                  <TextField label={t('dashboard.code')} fullWidth value={languageForm.code} onChange={(e) => setLanguageForm({ ...languageForm, code: e.target.value })} />
                </Grid>
              )}
              <Grid sx={{ width: '100%' }}>
                <TextField label={t('dashboard.name')} fullWidth value={languageForm.name} onChange={(e) => setLanguageForm({ ...languageForm, name: e.target.value })} />
              </Grid>
              <Grid sx={{ width: '100%' }}>
                <TextField label={t('dashboard.description')} fullWidth value={languageForm.description} onChange={(e) => setLanguageForm({ ...languageForm, description: e.target.value })} />
              </Grid>
              <Grid sx={{ width: '100%' }}>
                <TextField label={t('dashboard.icon')} fullWidth value={languageForm.icon} onChange={(e) => setLanguageForm({ ...languageForm, icon: e.target.value })} />
              </Grid>
              <Grid sx={{ width: '100%' }}>
                <TextField label={t('dashboard.color')} fullWidth value={languageForm.color} onChange={(e) => setLanguageForm({ ...languageForm, color: e.target.value })} />
              </Grid>
              <Grid sx={{ width: '100%' }}>
                <FormControlLabel control={<Switch checked={languageForm.isActive} onChange={(e) => setLanguageForm({ ...languageForm, isActive: e.target.checked })} />} label={t('dashboard.active')} />
              </Grid>
            </Grid>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => { setOpenAddLanguage(false); setEditLanguageCode(null); }}>{t('common.cancel')}</Button>
            <Button variant="contained" onClick={async () => {
              if (editLanguageCode) {
                await adminApi.updateLanguage(editLanguageCode, { name: languageForm.name, description: languageForm.description, icon: languageForm.icon, color: languageForm.color, isActive: languageForm.isActive });
                const updated = await adminApi.getLanguages();
                setLanguages(updated);
              } else {
                await adminApi.createLanguage({ code: languageForm.code, name: languageForm.name, description: languageForm.description, icon: languageForm.icon, color: languageForm.color, isActive: languageForm.isActive });
                const updated = await adminApi.getLanguages();
                setLanguages(updated);
              }
              setOpenAddLanguage(false); setEditLanguageCode(null);
            }}>{t('common.save')}</Button>
          </DialogActions>
        </Dialog>
      </TabPanel>

      {/* APIキー管理タブ */}
      <TabPanel value={tabValue} index={2}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h5">
            {t('dashboard.apiKeyManagement')}
          </Typography>
          {canManageUsers && (
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => {
              setApiKeyForm({ name: '', accessLevel: 'user' });
              setOpenAddApiKey(true);
            }}
          >
            {t('dashboard.generateApiKey')}
          </Button>
          )}
        </Box>
        
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>{t('dashboard.name')}</TableCell>
                <TableCell>{t('dashboard.key')}</TableCell>
                <TableCell>{t('dashboard.accessLevel')}</TableCell>
                <TableCell>{t('dashboard.status')}</TableCell>
                <TableCell>{t('dashboard.createdAt')}</TableCell>
                <TableCell>{t('dashboard.lastUsed')}</TableCell>
                {canManageUsers && <TableCell>{t('dashboard.actions')}</TableCell>}
              </TableRow>
            </TableHead>
            <TableBody>
              {apiKeys.map((apiKey) => (
                <TableRow key={apiKey.id}>
                  <TableCell>{apiKey.name}</TableCell>
                  <TableCell>
                    <Typography variant="body2" fontFamily="monospace" fontSize="0.8rem">
                      {apiKey.key}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Chip 
                      label={apiKey.accessLevel} 
                      color={apiKey.accessLevel === 'admin' ? 'error' : 'default'}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>
                    <Chip 
                      label={apiKey.status} 
                      color={apiKey.status === 'active' ? 'success' : 'warning'}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>{apiKey.createdAt}</TableCell>
                  <TableCell>{apiKey.lastUsed}</TableCell>
                  {canManageUsers && (
                  <TableCell>
                    <IconButton size="small" color="primary" onClick={() => { setOpenEditApiKey(apiKey); setApiKeyEditForm({ name: apiKey.name, description: '', isActive: apiKey.status === 'active' }); }}>
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      size="small"
                      color="error"
                      onClick={() => setOpenDeleteApiKey(apiKey)}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                  )}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </TabPanel>
      {/* Add API Key Dialog */}
      <Dialog open={openAddApiKey} onClose={() => setOpenAddApiKey(false)} fullWidth maxWidth="sm">
        <DialogTitle>{t('dashboard.generateApiKey')}</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid sx={{ width: '100%' }}>
              <TextField label="Name" fullWidth value={apiKeyForm.name} onChange={(e) => setApiKeyForm({ ...apiKeyForm, name: e.target.value })} />
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <FormControl fullWidth>
                <InputLabel>Access Level</InputLabel>
                <Select value={apiKeyForm.accessLevel} label="Access Level" onChange={(e) => setApiKeyForm({ ...apiKeyForm, accessLevel: e.target.value as string })}>
                  <MenuItem value="public">Public</MenuItem>
                  <MenuItem value="user">User</MenuItem>
                  <MenuItem value="admin">Admin</MenuItem>
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenAddApiKey(false)}>{t('common.cancel')}</Button>
          <Button
            variant="contained"
            onClick={async () => {
              try {
                if (!apiKeyForm.name) {
                  setSnackbar({ open: true, message: t('dashboard.nameRequired') || 'Name is required', severity: 'error' });
                  return;
                }
                const created = await adminApi.generateApiKey({ name: apiKeyForm.name, accessLevel: apiKeyForm.accessLevel });
                setApiKeys(prev => [...prev, created as unknown as ApiKey]);
                setOpenAddApiKey(false);
                setSnackbar({ open: true, message: 'APIキーを生成しました', severity: 'success' });
              } catch (e: any) {
                setSnackbar({ open: true, message: e?.normalized?.message || e?.message || (t('dashboard.apiKeyGenerateFailed') as string) || 'Failed to generate', severity: 'error' });
              }
            }}
          >{t('common.create')}</Button>
        </DialogActions>
      </Dialog>

      {/* Edit API Key Dialog */}
      <Dialog open={!!openEditApiKey} onClose={() => setOpenEditApiKey(null)} fullWidth maxWidth="sm">
        <DialogTitle>{t('dashboard.apiKeyManagement')}</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid sx={{ width: '100%' }}>
              <TextField label={t('dashboard.name')} fullWidth value={apiKeyEditForm.name} onChange={(e) => setApiKeyEditForm({ ...apiKeyEditForm, name: e.target.value })} />
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <TextField label={t('dashboard.description') || 'Description'} fullWidth value={apiKeyEditForm.description} onChange={(e) => setApiKeyEditForm({ ...apiKeyEditForm, description: e.target.value })} />
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <FormControlLabel control={<Switch checked={apiKeyEditForm.isActive} onChange={(e) => setApiKeyEditForm({ ...apiKeyEditForm, isActive: e.target.checked })} />} label={apiKeyEditForm.isActive ? (t('dashboard.active') as string) : (t('dashboard.inactive') as string)} />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenEditApiKey(null)}>{t('common.cancel')}</Button>
          <Button
            variant="contained"
            onClick={async () => {
              if (!openEditApiKey) return;
              try {
                await adminApi.updateApiKey(openEditApiKey.id, { name: apiKeyEditForm.name, description: apiKeyEditForm.description, isActive: apiKeyEditForm.isActive });
                setApiKeys(prev => prev.map(k => k.id === openEditApiKey.id ? { ...k, name: apiKeyEditForm.name, status: apiKeyEditForm.isActive ? 'active' : 'inactive' } : k));
                setOpenEditApiKey(null);
                setSnackbar({ open: true, message: t('dashboard.userUpdateSuccess') as string, severity: 'success' });
              } catch (e: any) {
                setSnackbar({ open: true, message: t('dashboard.apiKeyDeleteFailed') as string, severity: 'error' });
              }
            }}
          >{t('common.save')}</Button>
        </DialogActions>
      </Dialog>

      {/* Delete API Key Dialog */}
      <Dialog open={!!openDeleteApiKey} onClose={() => setOpenDeleteApiKey(null)}>
        <DialogTitle>{t('common.delete') || 'Delete'}</DialogTitle>
        <DialogContent>
          <Typography>{t('dashboard.deleteApiKeyConfirm') || 'Delete this API key?'}</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDeleteApiKey(null)}>{t('common.cancel')}</Button>
          <Button
            color="error"
            variant="contained"
            onClick={async () => {
              if (!openDeleteApiKey) return;
              try {
                await adminApi.deleteApiKey(openDeleteApiKey.id);
                setApiKeys(prev => prev.filter(k => k.id !== openDeleteApiKey.id));
                setOpenDeleteApiKey(null);
                setSnackbar({ open: true, message: 'APIキーを削除しました', severity: 'success' });
              } catch (e: any) {
                setSnackbar({ open: true, message: e?.normalized?.message || e?.message || (t('dashboard.apiKeyDeleteFailed') as string) || 'Failed to delete', severity: 'error' });
              }
            }}
          >{t('common.delete')}</Button>
        </DialogActions>
      </Dialog>

      {/* 設定タブ */}
      <TabPanel value={tabValue} index={3}>
        <Typography variant="h5" sx={{ mb: 3 }}>
          {t('dashboard.systemSettings')}
        </Typography>
        
        <Grid container spacing={3}>
          <Grid sx={{ width: { xs: '100%', md: '50%' } }}>
            <Card>
              <CardContent>
                <Typography variant="h6" sx={{ mb: 2 }}>
                  {t('dashboard.securitySettings')}
                </Typography>
                <FormControl fullWidth sx={{ mb: 2 }}>
                  <InputLabel>{t('dashboard.defaultAccessLevel')}</InputLabel>
                  <Select value={settings.defaultAccessLevel} label={t('dashboard.defaultAccessLevel')} onChange={(e) => setSettings(s => ({ ...s, defaultAccessLevel: e.target.value as string }))}>
                    <MenuItem value="public">Public</MenuItem>
                    <MenuItem value="user">User</MenuItem>
                    <MenuItem value="admin">Admin</MenuItem>
                  </Select>
                </FormControl>
                <Button variant="contained" onClick={async () => { try { await adminApi.updateSettings(settings as any); setSnackbar({ open: true, message: t('dashboard.saveSettings') as string, severity: 'success' }); } catch { setSnackbar({ open: true, message: t('dashboard.loadError') as string, severity: 'error' }); } }}>
                  {t('dashboard.saveSettings')}
                </Button>
              </CardContent>
            </Card>
          </Grid>
          
          <Grid sx={{ width: { xs: '100%', md: '50%' } }}>
            <Card>
              <CardContent>
                <Typography variant="h6" sx={{ mb: 2 }}>
                  {t('dashboard.rateLimiting')}
                </Typography>
                <TextField
                  fullWidth
                  label={t('dashboard.requestsPerMinute')}
                  type="number"
                  value={settings.requestsPerMinute}
                  onChange={(e) => setSettings(s => ({ ...s, requestsPerMinute: parseInt(e.target.value || '0', 10) || 0 }))}
                  sx={{ mb: 2 }}
                />
                <Button variant="contained" onClick={async () => { try { await adminApi.updateSettings(settings as any); setSnackbar({ open: true, message: t('dashboard.updateLimits') as string, severity: 'success' }); } catch { setSnackbar({ open: true, message: t('dashboard.loadError') as string, severity: 'error' }); } }}>
                  {t('dashboard.updateLimits')}
                </Button>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </TabPanel>

      {/* アナリティクスタブ */}
      <TabPanel value={tabValue} index={4}>
        <Typography variant="h5" sx={{ mb: 3 }}>
          {t('dashboard.systemAnalytics')}
        </Typography>
        
        <Grid container spacing={3}>
          <Grid sx={{ width: { xs: '100%', md: '50%' } }}>
            <Card>
              <CardContent>
                <Typography variant="h6" sx={{ mb: 2 }}>
                  {t('dashboard.recentActivity')}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {t('dashboard.noRecentActivity')}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          
          <Grid sx={{ width: { xs: '100%', md: '50%' } }}>
            <Card>
              <CardContent>
                <Typography variant="h6" sx={{ mb: 2 }}>
                  {t('dashboard.systemHealth')}
                </Typography>
                <Chip 
                  label={t('dashboard.healthy')} 
                  color="success" 
                  sx={{ mr: 1 }}
                />
                <Chip 
                  label={t('dashboard.uptime')} 
                  variant="outlined"
                />
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </TabPanel>

      {/* MCP監視タブ */}
      <TabPanel value={tabValue} index={5}>
        <Typography variant="h5" sx={{ mb: 3 }}>
          {t('dashboard.mcpMonitoring')}
        </Typography>
        
        <Grid container spacing={3}>
          <Grid sx={{ width: { xs: '100%', md: '50%' } }}>
            <Card>
              <CardContent>
                <Typography variant="h6" sx={{ mb: 2 }}>
                  {t('dashboard.mcpRequests')}
                </Typography>
                <TableContainer component={Paper}>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>{t('dashboard.method')}</TableCell>
                        <TableCell>{t('dashboard.count')}</TableCell>
                        <TableCell>{t('dashboard.lastUsed')}</TableCell>
                        <TableCell>{t('dashboard.status')}</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {mcpStats.map((row, idx) => (
                        <TableRow key={idx}>
                          <TableCell>{row.method}</TableCell>
                          <TableCell>{row.count}</TableCell>
                          <TableCell>{row.lastUsed}</TableCell>
                          <TableCell>
                            <Chip label={row.status} color={row.status === '正常' || row.status === 'ok' ? 'success' : 'warning'} size="small" />
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              </CardContent>
            </Card>
          </Grid>
          
          <Grid sx={{ width: { xs: '100%', md: '50%' } }}>
            <Card>
              <CardContent>
                <Typography variant="h6" sx={{ mb: 2 }}>
                  {t('dashboard.mcpPerformance')}
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                  {t('dashboard.averageResponseTime')}: {mcpPerf.avgMs}ms (p95: {mcpPerf.p95Ms}ms)
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                  {t('dashboard.successRate')}: {mcpPerf.successRate.toFixed(1)}%
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {t('dashboard.errorRate')}: {mcpPerf.errorRate.toFixed(1)}%
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </TabPanel>

      {/* システムログタブ */}
      <TabPanel value={tabValue} index={6}>
        <Typography variant="h5" sx={{ mb: 3 }}>
          {t('dashboard.systemLogs')}
        </Typography>
        
        <Grid container spacing={3}>
          <Grid sx={{ width: { xs: '100%', md: '50%' } }}>
            <Card>
              <CardContent>
                <Typography variant="h6" sx={{ mb: 2 }}>
                  {t('dashboard.recentLogs')}
                </Typography>
                <Box sx={{ maxHeight: 400, overflow: 'auto' }}>
                  {(systemLogs?.length ?? 0) === 0 ? (
                    <Typography variant="body2" color="text.secondary">
                      {t('dashboard.noRecentActivity')}
                    </Typography>
                  ) : (
                    (systemLogs || []).map((log, idx) => (
                      <Typography key={idx} variant="body2" fontFamily="monospace" fontSize="0.8rem">
                        [{new Date(log.timestamp).toLocaleString()}] {log.level}: {log.message}
                      </Typography>
                    ))
                  )}
                </Box>
              </CardContent>
            </Card>
          </Grid>
          
          <Grid sx={{ width: { xs: '100%', md: '50%' } }}>
            <Card>
              <CardContent>
                <Typography variant="h6" sx={{ mb: 2 }}>
                  {t('dashboard.logLevels')}
                </Typography>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Typography variant="body2">INFO</Typography>
                    <Chip label={String(logCounts.INFO)} size="small" color="info" />
                  </Box>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Typography variant="body2">WARN</Typography>
                    <Chip label={String(logCounts.WARN)} size="small" color="warning" />
                  </Box>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Typography variant="body2">ERROR</Typography>
                    <Chip label={String(logCounts.ERROR)} size="small" color="error" />
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </TabPanel>
      {/* Add User Dialog */}
      <Dialog open={openAddUser} onClose={() => setOpenAddUser(false)} fullWidth maxWidth="sm">
        <DialogTitle>{t('dashboard.addUser')}</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid sx={{ width: '100%' }}>
              <TextField label={t('dashboard.username')} fullWidth value={userForm.username} onChange={(e) => setUserForm({ ...userForm, username: e.target.value })} />
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <TextField label={t('dashboard.email')} fullWidth value={userForm.email} onChange={(e) => setUserForm({ ...userForm, email: e.target.value })} />
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <TextField label={t('dashboard.fullName')} fullWidth value={userForm.fullName} onChange={(e) => setUserForm({ ...userForm, fullName: e.target.value })} />
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <FormControl fullWidth>
                <InputLabel>Role</InputLabel>
                <Select value={userForm.role} label="Role" onChange={(e) => setUserForm({ ...userForm, role: e.target.value as string })}>
                  <MenuItem value="user">User</MenuItem>
                  <MenuItem value="admin">Admin</MenuItem>
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenAddUser(false)}>{t('common.cancel')}</Button>
          <Button
            variant="contained"
            onClick={async () => {
              try {
                if (!userForm.username || !userForm.email) {
                  setSnackbar({ open: true, message: 'usernameとemailは必須です', severity: 'error' });
                  return;
                }
                const created = await adminApi.createUser({
                  username: userForm.username,
                  email: userForm.email,
                  fullName: userForm.fullName,
                  role: userForm.role,
                  isActive: userForm.isActive,
                });
                setUsers(prev => ([...prev, {
                  id: created.id,
                  username: created.username,
                  email: created.email,
                  fullName: created.fullName,
                  role: created.role,
                  status: created.isActive ? 'active' : 'inactive',
                  lastLogin: created.lastLogin,
                }]));
                setOpenAddUser(false);
                setSnackbar({ open: true, message: 'ユーザーを作成しました', severity: 'success' });
              } catch (e: any) {
                setSnackbar({ open: true, message: e?.normalized?.message || e?.message || (t('dashboard.userCreateFailed') as string) || 'Failed to create', severity: 'error' });
              }
            }}
          >{t('common.create')}</Button>
        </DialogActions>
      </Dialog>

      {/* Edit User Dialog */}
      <Dialog open={!!openEditUser} onClose={() => setOpenEditUser(null)} fullWidth maxWidth="sm">
        <DialogTitle>{t('dashboard.editUser') || 'Edit User'}</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid sx={{ width: '100%' }}>
              <TextField label={t('dashboard.username')} fullWidth value={userForm.username} onChange={(e) => setUserForm({ ...userForm, username: e.target.value })} />
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <TextField label={t('dashboard.email')} fullWidth value={userForm.email} onChange={(e) => setUserForm({ ...userForm, email: e.target.value })} />
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <TextField label={t('dashboard.fullName')} fullWidth value={userForm.fullName} onChange={(e) => setUserForm({ ...userForm, fullName: e.target.value })} />
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <FormControl fullWidth>
                <InputLabel>Role</InputLabel>
                <Select value={userForm.role} label="Role" onChange={(e) => setUserForm({ ...userForm, role: e.target.value as string })}>
                  <MenuItem value="user">User</MenuItem>
                  <MenuItem value="admin">Admin</MenuItem>
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenEditUser(null)}>{t('common.cancel')}</Button>
          <Button
            variant="contained"
            onClick={async () => {
              if (!openEditUser) return;
              try {
                const updated = await adminApi.updateUser(openEditUser.id, {
                  username: userForm.username,
                  email: userForm.email,
                  fullName: userForm.fullName,
                  role: userForm.role,
                  isActive: userForm.isActive,
                });
                setUsers(prev => prev.map(u => u.id === openEditUser.id ? {
                  id: updated.id,
                  username: updated.username,
                  email: updated.email,
                  fullName: updated.fullName,
                  role: updated.role,
                  status: updated.isActive ? 'active' : 'inactive',
                  lastLogin: updated.lastLogin,
                } : u));
                setOpenEditUser(null);
                setSnackbar({ open: true, message: 'ユーザーを更新しました', severity: 'success' });
              } catch (e: any) {
                setSnackbar({ open: true, message: e?.normalized?.message || e?.message || (t('dashboard.userUpdateFailed') as string) || 'Failed to update', severity: 'error' });
              }
            }}
          >{t('common.save')}</Button>
        </DialogActions>
      </Dialog>

      {/* Roles 管理タブ */}
      <TabPanel value={tabValue} index={7}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h5">
            {t('dashboard.roles') || 'Roles'}
          </Typography>
          {canManageRoles && (
          <Button variant="contained" startIcon={<AddIcon />} onClick={() => { setRoleForm({ name: '', description: '', permissions: { manage_users: false, manage_rules: true, manage_roles: false }, is_active: true }); setOpenAddRole(true); }}>
            {t('common.add')}
          </Button>
          )}
        </Box>
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>{t('dashboard.name')}</TableCell>
                <TableCell>{t('dashboard.description') || 'Description'}</TableCell>
                <TableCell>{t('dashboard.permissions') || 'Permissions'}</TableCell>
                <TableCell>{t('dashboard.status')}</TableCell>
                {canManageRoles && <TableCell>{t('dashboard.actions')}</TableCell>}
              </TableRow>
            </TableHead>
            <TableBody>
              {roles.map((r) => (
                <TableRow key={r.name}>
                  <TableCell>{r.name}</TableCell>
                  <TableCell>{r.description}</TableCell>
                  <TableCell>
                    <Chip label={`manage_users: ${r.permissions?.manage_users ? 'on' : 'off'}`} size="small" sx={{ mr: 1 }} />
                    <Chip label={`manage_rules: ${r.permissions?.manage_rules ? 'on' : 'off'}`} size="small" sx={{ mr: 1 }} />
                    <Chip label={`manage_roles: ${r.permissions?.manage_roles ? 'on' : 'off'}`} size="small" />
                  </TableCell>
                  <TableCell>
                    <Chip label={r.is_active ? (t('dashboard.active') as string) : (t('dashboard.inactive') as string)} color={r.is_active ? 'success' : 'default'} size="small" />
                  </TableCell>
                  {canManageRoles && (
                  <TableCell>
                    <IconButton size="small" color="primary" onClick={() => { setRoleForm({ name: r.name, description: r.description, permissions: r.permissions, is_active: r.is_active }); setOpenEditRole(r); }}>
                      <EditIcon />
                    </IconButton>
                    <IconButton size="small" onClick={() => { setRoleForm({ name: `${r.name}-copy`, description: r.description, permissions: r.permissions, is_active: r.is_active }); setOpenAddRole(true); }}>
                      <ContentCopyIcon />
                    </IconButton>
                    <IconButton size="small" color="error" onClick={async () => { if (!window.confirm(t('dashboard.deleteRoleConfirm') || 'Delete this role?')) return; await adminApi.deleteRole(r.name); setRoles(prev => prev.filter(x => x.name !== r.name)); }}>
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                  )}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </TabPanel>

      {/* Add Role Dialog */}
      <Dialog open={openAddRole} onClose={() => setOpenAddRole(false)} fullWidth maxWidth="sm">
        <DialogTitle>{t('dashboard.addRole') || 'Add Role'}</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid sx={{ width: '100%' }}>
              <TextField label={t('dashboard.name')} fullWidth value={roleForm.name} onChange={(e) => setRoleForm({ ...roleForm, name: e.target.value })} />
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <TextField label={t('dashboard.description') || 'Description'} fullWidth value={roleForm.description} onChange={(e) => setRoleForm({ ...roleForm, description: e.target.value })} />
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <Typography variant="subtitle2" sx={{ mb: 1 }}>{t('dashboard.rolePresets') || 'Role Presets'}</Typography>
              <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', mb: 2 }}>
                <Button size="small" variant="outlined" onClick={() => applyPreset('readonly')}>{t('dashboard.presetReadOnly') || 'ReadOnly'}</Button>
                <Button size="small" variant="outlined" onClick={() => applyPreset('editor')}>{t('dashboard.presetEditor') || 'Editor'}</Button>
                <Button size="small" variant="outlined" onClick={() => applyPreset('admin')}>{t('dashboard.presetAdmin') || 'Admin'}</Button>
              </Box>
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <Typography variant="subtitle2" sx={{ mb: 1 }}>{t('dashboard.permissions') || 'Permissions'}</Typography>
              <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                <Tooltip title={t('permissions.manage_users.desc') || ''} arrow>
                  <Button variant={roleForm.permissions?.manage_users ? 'contained' : 'outlined'} onClick={() => setRoleForm({ ...roleForm, permissions: { ...roleForm.permissions, manage_users: !roleForm.permissions?.manage_users } })}>
                    {t('permissions.manage_users.name') || 'manage_users'}
                  </Button>
                </Tooltip>
                <Tooltip title={t('permissions.manage_rules.desc') || ''} arrow>
                  <Button variant={roleForm.permissions?.manage_rules ? 'contained' : 'outlined'} onClick={() => setRoleForm({ ...roleForm, permissions: { ...roleForm.permissions, manage_rules: !roleForm.permissions?.manage_rules } })}>
                    {t('permissions.manage_rules.name') || 'manage_rules'}
                  </Button>
                </Tooltip>
                <Tooltip title={t('permissions.manage_roles.desc') || ''} arrow>
                  <Button variant={roleForm.permissions?.manage_roles ? 'contained' : 'outlined'} onClick={() => setRoleForm({ ...roleForm, permissions: { ...roleForm.permissions, manage_roles: !roleForm.permissions?.manage_roles } })}>
                    {t('permissions.manage_roles.name') || 'manage_roles'}
                  </Button>
                </Tooltip>
              </Box>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenAddRole(false)}>{t('common.cancel')}</Button>
          <Button variant="contained" onClick={async () => { await adminApi.createRole({ name: roleForm.name, description: roleForm.description, permissions: roleForm.permissions, is_active: roleForm.is_active }); setOpenAddRole(false); const list = await adminApi.getRoles(); setRoles(list); }}>{t('common.create')}</Button>
        </DialogActions>
      </Dialog>

      {/* Edit Role Dialog */}
      <Dialog open={!!openEditRole} onClose={() => setOpenEditRole(null)} fullWidth maxWidth="sm">
        <DialogTitle>{t('dashboard.editRole') || 'Edit Role'}</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid sx={{ width: '100%' }}>
              <TextField label={t('dashboard.name')} fullWidth value={roleForm.name} disabled />
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <TextField label={t('dashboard.description') || 'Description'} fullWidth value={roleForm.description} onChange={(e) => setRoleForm({ ...roleForm, description: e.target.value })} />
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <Typography variant="subtitle2" sx={{ mb: 1 }}>{t('dashboard.rolePresets') || 'Role Presets'}</Typography>
              <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', mb: 2 }}>
                <Button size="small" variant="outlined" onClick={() => applyPreset('readonly')}>{t('dashboard.presetReadOnly') || 'ReadOnly'}</Button>
                <Button size="small" variant="outlined" onClick={() => applyPreset('editor')}>{t('dashboard.presetEditor') || 'Editor'}</Button>
                <Button size="small" variant="outlined" onClick={() => applyPreset('admin')}>{t('dashboard.presetAdmin') || 'Admin'}</Button>
              </Box>
            </Grid>
            <Grid sx={{ width: '100%' }}>
              <Typography variant="subtitle2" sx={{ mb: 1 }}>{t('dashboard.permissions') || 'Permissions'}</Typography>
              <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                <Tooltip title={t('permissions.manage_users.desc') || ''} arrow>
                  <Button variant={roleForm.permissions?.manage_users ? 'contained' : 'outlined'} onClick={() => setRoleForm({ ...roleForm, permissions: { ...roleForm.permissions, manage_users: !roleForm.permissions?.manage_users } })}>
                    {t('permissions.manage_users.name') || 'manage_users'}
                  </Button>
                </Tooltip>
                <Tooltip title={t('permissions.manage_rules.desc') || ''} arrow>
                  <Button variant={roleForm.permissions?.manage_rules ? 'contained' : 'outlined'} onClick={() => setRoleForm({ ...roleForm, permissions: { ...roleForm.permissions, manage_rules: !roleForm.permissions?.manage_rules } })}>
                    {t('permissions.manage_rules.name') || 'manage_rules'}
                  </Button>
                </Tooltip>
                <Tooltip title={t('permissions.manage_roles.desc') || ''} arrow>
                  <Button variant={roleForm.permissions?.manage_roles ? 'contained' : 'outlined'} onClick={() => setRoleForm({ ...roleForm, permissions: { ...roleForm.permissions, manage_roles: !roleForm.permissions?.manage_roles } })}>
                    {t('permissions.manage_roles.name') || 'manage_roles'}
                  </Button>
                </Tooltip>
              </Box>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenEditRole(null)}>{t('common.cancel')}</Button>
          <Button variant="contained" onClick={async () => { if (!openEditRole) return; await adminApi.updateRole(openEditRole.name, { description: roleForm.description, permissions: roleForm.permissions, is_active: roleForm.is_active }); setOpenEditRole(null); const list = await adminApi.getRoles(); setRoles(list); }}>{t('common.save')}</Button>
        </DialogActions>
      </Dialog>

      {/* 一括エクスポートダイアログ */}
      <Dialog open={bulkExportDialogOpen} onClose={() => setBulkExportDialogOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>全ルールエクスポート</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
            <FormControl fullWidth>
              <InputLabel>エクスポート範囲</InputLabel>
              <Select
                value={exportScope}
                onChange={(e) => setExportScope(e.target.value)}
                label="エクスポート範囲"
              >
                <MenuItem value="all">全プロジェクト + グローバルルール</MenuItem>
                <MenuItem value="projects">プロジェクトルールのみ</MenuItem>
                <MenuItem value="global">グローバルルールのみ</MenuItem>
              </Select>
            </FormControl>
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
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setBulkExportDialogOpen(false)}>キャンセル</Button>
          <Button variant="contained" onClick={handleBulkExport} disabled={loading}>
            {loading ? 'エクスポート中...' : 'エクスポート'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* 一括インポートダイアログ */}
      <Dialog open={bulkImportDialogOpen} onClose={() => setBulkImportDialogOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>全ルールインポート</DialogTitle>
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
                <Switch
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
          <Button onClick={() => setBulkImportDialogOpen(false)}>キャンセル</Button>
          <Button 
            variant="contained" 
            onClick={handleBulkImport} 
            disabled={loading || !importFile}
          >
            {loading ? 'インポート中...' : 'インポート'}
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={snackbar.open}
        onClose={() => setSnackbar(prev => ({ ...prev, open: false }))}
        message={snackbar.message}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
        autoHideDuration={3000}
      />
    </Box>
  );
};

export default AdminDashboard;
