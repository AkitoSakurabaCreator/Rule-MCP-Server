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
} from '@mui/material';
import { useNavigate, useParams } from 'react-router-dom';
import { api } from '../services/api';

const ruleTypes = ['style', 'security', 'performance', 'maintainability', 'accessibility'];
const severities = ['error', 'warning', 'info'];

const RuleForm: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    project_id: projectId || '',
    rule_id: '',
    name: '',
    description: '',
    type: 'style',
    severity: 'warning',
    pattern: '',
    message: '',
  });

  useEffect(() => {
    if (projectId) {
      setFormData(prev => ({ ...prev, project_id: projectId }));
    }
  }, [projectId]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      await api.post('/rules', formData);
      setSuccess('Rule created successfully!');
      setTimeout(() => {
        navigate('/');
      }, 1500);
    } catch (error: any) {
      setError(error.response?.data?.error || 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field: string, value: string) => {
    setFormData(prev => ({
      ...prev,
      [field]: value,
    }));
  };

  return (
    <Card>
      <CardContent>
        <Typography variant="h5" component="h2" sx={{ mb: 3 }}>
          New Rule for Project: {projectId}
        </Typography>

        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
        {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}

        <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <TextField
            label="Rule ID"
            value={formData.rule_id}
            onChange={(e) => handleChange('rule_id', e.target.value)}
            required
            helperText="Unique identifier for the rule"
          />

          <TextField
            label="Name"
            value={formData.name}
            onChange={(e) => handleChange('name', e.target.value)}
            required
            helperText="Display name for the rule"
          />

          <TextField
            label="Description"
            value={formData.description}
            onChange={(e) => handleChange('description', e.target.value)}
            multiline
            rows={3}
            helperText="Description of what this rule checks"
          />

          <TextField
            select
            label="Type"
            value={formData.type}
            onChange={(e) => handleChange('type', e.target.value)}
            required
            helperText="Category of the rule"
          >
            {ruleTypes.map((type) => (
              <MenuItem key={type} value={type}>
                {type.charAt(0).toUpperCase() + type.slice(1)}
              </MenuItem>
            ))}
          </TextField>

          <TextField
            select
            label="Severity"
            value={formData.severity}
            onChange={(e) => handleChange('severity', e.target.value)}
            required
            helperText="How serious is this rule violation"
          >
            {severities.map((severity) => (
              <MenuItem key={severity} value={severity}>
                {severity.charAt(0).toUpperCase() + severity.slice(1)}
              </MenuItem>
            ))}
          </TextField>

          <TextField
            label="Pattern (Regex)"
            value={formData.pattern}
            onChange={(e) => handleChange('pattern', e.target.value)}
            required
            helperText="Regular expression pattern to match violations"
          />

          <TextField
            label="Message"
            value={formData.message}
            onChange={(e) => handleChange('message', e.target.value)}
            required
            helperText="Message to display when rule is violated"
          />

          <Box sx={{ display: 'flex', gap: 2, mt: 2 }}>
            <Button
              type="submit"
              variant="contained"
              disabled={loading}
              sx={{ minWidth: 120 }}
            >
              {loading ? 'Creating...' : 'Create Rule'}
            </Button>
            <Button
              variant="outlined"
              onClick={() => navigate('/')}
              disabled={loading}
            >
              Cancel
            </Button>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default RuleForm;
