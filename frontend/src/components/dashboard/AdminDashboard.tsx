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
} from '@mui/icons-material';
import { useTranslation } from 'react-i18next';
import {
  adminApi,
  AdminStats as AdminStatsType,
  MCPStats as MCPStatsType,
  SystemLog as SystemLogType,
} from '../../services/adminApi';

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
  const [systemLogs, setSystemLogs] = useState<SystemLogType[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;
    const load = async () => {
      try {
        setLoading(true);
        setError(null);
        const [s, u, k, m, logs] = await Promise.all([
          adminApi.getStats(),
          adminApi.getUsers(),
          adminApi.getApiKeys(),
          adminApi.getMCPStats(),
          adminApi.getSystemLogs(),
        ]);
        if (!mounted) return;
        setStats(s);
        setUsers(
          u.map((x) => ({
            id: x.id,
            username: x.username,
            email: x.email,
            role: x.role,
            status: x.isActive ? 'active' : 'inactive',
            lastLogin: x.lastLogin,
          }))
        );
        setApiKeys(k as unknown as ApiKey[]);
        setMcpStats(m);
        setSystemLogs(logs);
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
  }, []);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    console.log('Tab changed to:', newValue);
    setTabValue(newValue);
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" sx={{ mb: 3 }}>
        {t('dashboard.title')}
      </Typography>

      {loading && (
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          Loading...
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

      {/* タブナビゲーション */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
        <Tabs value={tabValue} onChange={handleTabChange} aria-label="admin tabs">
          <Tab 
            icon={<PeopleIcon />} 
            label={t('dashboard.users')} 
            iconPosition="start"
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
        </Tabs>
      </Box>

      {/* ユーザー管理タブ */}
      <TabPanel value={tabValue} index={0}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h5">
            {t('dashboard.userManagement')}
          </Typography>
          <Button variant="contained" startIcon={<AddIcon />}>
            {t('dashboard.addUser')}
          </Button>
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
                <TableCell>{t('dashboard.actions')}</TableCell>
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
                  <TableCell>
                    <IconButton size="small" color="primary">
                      <EditIcon />
                    </IconButton>
                    <IconButton size="small" color="error">
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </TabPanel>

      {/* APIキー管理タブ */}
      <TabPanel value={tabValue} index={1}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h5">
            {t('dashboard.apiKeyManagement')}
          </Typography>
          <Button variant="contained" startIcon={<AddIcon />}>
            {t('dashboard.generateApiKey')}
          </Button>
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
                <TableCell>{t('dashboard.actions')}</TableCell>
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
                  <TableCell>
                    <IconButton size="small" color="primary">
                      <EditIcon />
                    </IconButton>
                    <IconButton size="small" color="error">
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </TabPanel>

      {/* 設定タブ */}
      <TabPanel value={tabValue} index={2}>
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
                  <Select defaultValue="public">
                    <MenuItem value="public">Public</MenuItem>
                    <MenuItem value="user">User</MenuItem>
                    <MenuItem value="admin">Admin</MenuItem>
                  </Select>
                </FormControl>
                <Button variant="contained">
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
                  defaultValue={100}
                  sx={{ mb: 2 }}
                />
                <Button variant="contained">
                  {t('dashboard.updateLimits')}
                </Button>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </TabPanel>

      {/* アナリティクスタブ */}
      <TabPanel value={tabValue} index={3}>
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
      <TabPanel value={tabValue} index={4}>
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
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                  {t('dashboard.averageResponseTime')}: 45ms
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                  {t('dashboard.successRate')}: 99.8%
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {t('dashboard.errorRate')}: 0.2%
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </TabPanel>

      {/* システムログタブ */}
      <TabPanel value={tabValue} index={5}>
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
                  <Typography variant="body2" fontFamily="monospace" fontSize="0.8rem">
                    [2024-01-15 14:30:15] INFO: User 'admin' logged in successfully
                  </Typography>
                  <Typography variant="body2" fontFamily="monospace" fontSize="0.8rem">
                    [2024-01-15 14:29:45] WARN: API key 'user_key_456' expired
                  </Typography>
                  <Typography variant="body2" fontFamily="monospace" fontSize="0.8rem">
                    [2024-01-15 14:28:30] INFO: MCP request 'getRules' processed in 23ms
                  </Typography>
                  <Typography variant="body2" fontFamily="monospace" fontSize="0.8rem">
                    [2024-01-15 14:27:15] ERROR: Database connection timeout
                  </Typography>
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
                    <Chip label="1,234" size="small" color="info" />
                  </Box>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Typography variant="body2">WARN</Typography>
                    <Chip label="45" size="small" color="warning" />
                  </Box>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Typography variant="body2">ERROR</Typography>
                    <Chip label="12" size="small" color="error" />
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </TabPanel>
    </Box>
  );
};

export default AdminDashboard;
