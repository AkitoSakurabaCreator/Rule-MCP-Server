import React from 'react';
import { IconButton, Tooltip } from '@mui/material';
import { LightMode as LightModeIcon, DarkMode as DarkModeIcon } from '@mui/icons-material';
import { useTheme } from '../contexts/ThemeContext';
import { useTranslation } from 'react-i18next';

const ThemeToggle: React.FC = () => {
  const { themeMode, toggleTheme } = useTheme();
  const { t } = useTranslation();

  const isDark = themeMode === 'dark';
  const tooltipText = isDark ? t('theme.switchToLight') : t('theme.switchToDark');

  return (
    <Tooltip title={tooltipText}>
      <IconButton
        color="inherit"
        onClick={toggleTheme}
        sx={{ ml: 1 }}
      >
        {isDark ? <LightModeIcon /> : <DarkModeIcon />}
      </IconButton>
    </Tooltip>
  );
};

export default ThemeToggle;
