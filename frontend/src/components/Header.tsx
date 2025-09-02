import React from 'react';
import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { Add as AddIcon, Code as CodeIcon, Language as LanguageIcon } from '@mui/icons-material';

const Header: React.FC = () => {
  const navigate = useNavigate();

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          Rule MCP Server
        </Typography>
        <Box sx={{ display: 'flex', gap: 1 }}>
          <Button
            color="inherit"
            startIcon={<AddIcon />}
            onClick={() => navigate('/projects/new')}
          >
            New Project
          </Button>
          <Button
            color="inherit"
            startIcon={<LanguageIcon />}
            onClick={() => navigate('/global-rules/new')}
          >
            Global Rules
          </Button>
          <Button
            color="inherit"
            startIcon={<CodeIcon />}
            onClick={() => navigate('/validate')}
          >
            Validate Code
          </Button>
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
