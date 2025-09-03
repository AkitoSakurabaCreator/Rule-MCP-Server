import type { Preview } from '@storybook/react-webpack5';
import React from 'react';
import { ThemeProvider } from '../src/contexts/ThemeContext';
import { AuthProvider } from '../src/contexts/AuthContext';
import { I18nextProvider } from 'react-i18next';
import i18n from '../src/i18n';

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
       color: /(background|color)$/i,
       date: /Date$/i,
      },
    },
    backgrounds: {
      default: 'light',
      values: [
        { name: 'light', value: '#ffffff' },
        { name: 'dark', value: '#121212' },
      ],
    },
  },
  decorators: [
    (Story) => (
      <I18nextProvider i18n={i18n}>
        <ThemeProvider>
          <AuthProvider>
            <Story />
          </AuthProvider>
        </ThemeProvider>
      </I18nextProvider>
    ),
  ],
};

export default preview;