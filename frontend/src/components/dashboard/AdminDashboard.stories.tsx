import type { Meta, StoryObj } from '@storybook/react';
import AdminDashboard from './AdminDashboard';

const meta: Meta<typeof AdminDashboard> = {
  title: 'Dashboard/AdminDashboard',
  component: AdminDashboard,
  parameters: {
    layout: 'fullscreen',
  },
  tags: ['autodocs'],
  argTypes: {
    // 必要に応じてコントロールを追加
  },
};

export default meta;
type Story = StoryObj<typeof meta>;

// 基本的なストーリー
export const Default: Story = {
  args: {},
};

// ダークテーマでの表示
export const DarkTheme: Story = {
  args: {},
  parameters: {
    backgrounds: {
      default: 'dark',
    },
  },
};

// モバイル表示
export const Mobile: Story = {
  args: {},
  parameters: {
    viewport: {
      defaultViewport: 'mobile1',
    },
  },
};

// タブレット表示
export const Tablet: Story = {
  args: {},
  parameters: {
    viewport: {
      defaultViewport: 'tablet',
    },
  },
};
