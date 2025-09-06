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
} from '@mui/material';
import {
  Person as PersonIcon,
  Security as SecurityIcon,
  Edit as EditIcon,
  Save as SaveIcon,
  Cancel as CancelIcon,
} from '@mui/icons-material';
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
  const { logout } = useAuth();
  const [profile, setProfile] = useState<UserProfileData | null>(null);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const [passwordEditing, setPasswordEditing] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [formData, setFormData] = useState({
    username: '',
    email: '',
  });
  const [passwordData, setPasswordData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
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

  const handlePasswordChange = async () => {
    try {
      setError('');
      setSuccess('');
      
      if (passwordData.newPassword !== passwordData.confirmPassword) {
        setError('新しいパスワードが一致しません');
        return;
      }
      
      if (passwordData.newPassword.length < 6) {
        setError('新しいパスワードは6文字以上で入力してください');
        return;
      }
      
      await api.put('/auth/change-password', {
        current_password: passwordData.currentPassword,
        new_password: passwordData.newPassword,
      });
      
      setSuccess('パスワードが変更されました');
      setPasswordEditing(false);
      setPasswordData({
        currentPassword: '',
        newPassword: '',
        confirmPassword: '',
      });
    } catch (error: any) {
      console.error('Password change error:', error);
      setError(error.response?.data?.error || 'パスワードの変更に失敗しました');
    }
  };

  const handlePasswordCancel = () => {
    setPasswordEditing(false);
    setPasswordData({
      currentPassword: '',
      newPassword: '',
      confirmPassword: '',
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

      <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, gap: 3 }}>
        <Box sx={{ flex: 1 }}>
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
        </Box>

        <Box sx={{ flex: 1 }}>
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
                  {profile.created_at && profile.created_at !== '0001-01-01T00:00:00Z'
                    ? new Date(profile.created_at).toLocaleDateString('ja-JP')
                    : '不明'
                  }
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="subtitle2" color="text.secondary">
                  最終ログイン
                </Typography>
                <Typography variant="body1">
                  {profile.last_login && profile.last_login !== '0001-01-01T00:00:00Z'
                    ? new Date(profile.last_login).toLocaleString('ja-JP')
                    : '未ログイン'
                  }
                </Typography>
              </Box>

              <Divider sx={{ my: 2 }} />

              <Typography variant="h6" sx={{ mb: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
                <SecurityIcon />
                パスワード変更
              </Typography>

              {passwordEditing ? (
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                  <TextField
                    fullWidth
                    type="password"
                    label="現在のパスワード"
                    value={passwordData.currentPassword}
                    onChange={(e) => setPasswordData(prev => ({ ...prev, currentPassword: e.target.value }))}
                    size="small"
                  />
                  <TextField
                    fullWidth
                    type="password"
                    label="新しいパスワード"
                    value={passwordData.newPassword}
                    onChange={(e) => setPasswordData(prev => ({ ...prev, newPassword: e.target.value }))}
                    size="small"
                    helperText="6文字以上で入力してください"
                  />
                  <TextField
                    fullWidth
                    type="password"
                    label="新しいパスワード（確認）"
                    value={passwordData.confirmPassword}
                    onChange={(e) => setPasswordData(prev => ({ ...prev, confirmPassword: e.target.value }))}
                    size="small"
                  />
                  <Box sx={{ display: 'flex', gap: 1 }}>
                    <Button
                      variant="contained"
                      onClick={handlePasswordChange}
                      size="small"
                    >
                      パスワード変更
                    </Button>
                    <Button
                      variant="outlined"
                      onClick={handlePasswordCancel}
                      size="small"
                    >
                      キャンセル
                    </Button>
                  </Box>
                </Box>
              ) : (
                <Button
                  variant="outlined"
                  onClick={() => setPasswordEditing(true)}
                  size="small"
                >
                  パスワードを変更
                </Button>
              )}

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
        </Box>
      </Box>
    </Box>
  );
};

export default UserProfile;
