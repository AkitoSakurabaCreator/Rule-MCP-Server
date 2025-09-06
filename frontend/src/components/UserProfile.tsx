import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Button,
  TextField,
  Alert,
  Divider,
  Chip,
  Grid,
} from '@mui/material';
import {
  Person as PersonIcon,
  Email as EmailIcon,
  Security as SecurityIcon,
  Edit as EditIcon,
  Save as SaveIcon,
  Cancel as CancelIcon,
} from '@mui/icons-material';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../contexts/AuthContext';
import { api } from '../services/api';

interface UserProfileData {
  username: string;
  email: string;
  role: string;
  created_at: string;
  last_login: string;
}

const UserProfile: React.FC = () => {
  const { t } = useTranslation();
  const { user, logout } = useAuth();
  const [profile, setProfile] = useState<UserProfileData | null>(null);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [formData, setFormData] = useState({
    username: '',
    email: '',
  });

  useEffect(() => {
    loadProfile();
  }, []);

  const loadProfile = async () => {
    try {
      setLoading(true);
      const response = await api.get('/auth/me');
      setProfile(response.data);
      setFormData({
        username: response.data.username,
        email: response.data.email,
      });
    } catch (error: any) {
      console.error('Profile load error:', error);
      setError('プロフィールの読み込みに失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    try {
      setError('');
      setSuccess('');
      
      await api.put('/auth/profile', formData);
      setSuccess('プロフィールが更新されました');
      setEditing(false);
      await loadProfile();
    } catch (error: any) {
      console.error('Profile update error:', error);
      setError(error.response?.data?.error || 'プロフィールの更新に失敗しました');
    }
  };

  const handleCancel = () => {
    setEditing(false);
    setFormData({
      username: profile?.username || '',
      email: profile?.email || '',
    });
    setError('');
    setSuccess('');
  };

  const getRoleColor = (role: string) => {
    switch (role) {
      case 'admin':
        return 'error';
      case 'user':
        return 'primary';
      default:
        return 'default';
    }
  };

  const getRoleLabel = (role: string) => {
    switch (role) {
      case 'admin':
        return '管理者';
      case 'user':
        return 'ユーザー';
      default:
        return role;
    }
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '50vh' }}>
        <Typography>読み込み中...</Typography>
      </Box>
    );
  }

  if (!profile) {
    return (
      <Alert severity="error">
        プロフィール情報を取得できませんでした
      </Alert>
    );
  }

  return (
    <Box sx={{ maxWidth: 800, mx: 'auto' }}>
      <Typography variant="h4" component="h1" sx={{ mb: 3, display: 'flex', alignItems: 'center', gap: 1 }}>
        <PersonIcon />
        ユーザープロフィール
      </Typography>

      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
      {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}

      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ mb: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
                <PersonIcon />
                基本情報
              </Typography>
              
              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle2" color="text.secondary">
                  ユーザー名
                </Typography>
                {editing ? (
                  <TextField
                    fullWidth
                    value={formData.username}
                    onChange={(e) => setFormData(prev => ({ ...prev, username: e.target.value }))}
                    size="small"
                    sx={{ mt: 1 }}
                  />
                ) : (
                  <Typography variant="body1">{profile.username}</Typography>
                )}
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle2" color="text.secondary">
                  メールアドレス
                </Typography>
                {editing ? (
                  <TextField
                    fullWidth
                    type="email"
                    value={formData.email}
                    onChange={(e) => setFormData(prev => ({ ...prev, email: e.target.value }))}
                    size="small"
                    sx={{ mt: 1 }}
                  />
                ) : (
                  <Typography variant="body1">{profile.email}</Typography>
                )}
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle2" color="text.secondary">
                  ロール
                </Typography>
                <Chip
                  label={getRoleLabel(profile.role)}
                  color={getRoleColor(profile.role) as any}
                  size="small"
                  sx={{ mt: 1 }}
                />
              </Box>

              <Box sx={{ display: 'flex', gap: 1, mt: 2 }}>
                {editing ? (
                  <>
                    <Button
                      variant="contained"
                      startIcon={<SaveIcon />}
                      onClick={handleSave}
                      size="small"
                    >
                      保存
                    </Button>
                    <Button
                      variant="outlined"
                      startIcon={<CancelIcon />}
                      onClick={handleCancel}
                      size="small"
                    >
                      キャンセル
                    </Button>
                  </>
                ) : (
                  <Button
                    variant="outlined"
                    startIcon={<EditIcon />}
                    onClick={() => setEditing(true)}
                    size="small"
                  >
                    編集
                  </Button>
                )}
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" sx={{ mb: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
                <SecurityIcon />
                アカウント情報
              </Typography>
              
              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle2" color="text.secondary">
                  アカウント作成日
                </Typography>
                <Typography variant="body1">
                  {new Date(profile.created_at).toLocaleDateString('ja-JP')}
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle2" color="text.secondary">
                  最終ログイン
                </Typography>
                <Typography variant="body1">
                  {profile.last_login 
                    ? new Date(profile.last_login).toLocaleString('ja-JP')
                    : '未ログイン'
                  }
                </Typography>
              </Box>

              <Divider sx={{ my: 2 }} />

              <Button
                variant="outlined"
                color="error"
                onClick={logout}
                fullWidth
              >
                ログアウト
              </Button>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default UserProfile;
