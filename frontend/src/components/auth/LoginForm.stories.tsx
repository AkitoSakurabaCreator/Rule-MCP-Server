import type { Meta, StoryObj } from '@storybook/react';
import { within, userEvent } from '@storybook/testing-library';
import LoginForm from './LoginForm';

const meta: Meta<typeof LoginForm> = {
  title: 'Auth/LoginForm',
  component: LoginForm,
  parameters: {
    layout: 'centered',
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

// エラー状態
export const WithError: Story = {
  args: {},
  play: async ({ canvasElement }) => {
    // エラー状態をシミュレート
    const canvas = within(canvasElement);
    const usernameInput = canvas.getByLabelText('ユーザー名');
    const passwordInput = canvas.getByLabelText('パスワード');
    const submitButton = canvas.getByRole('button', { name: 'サインイン' });
    
    await userEvent.type(usernameInput, 'invalid');
    await userEvent.type(passwordInput, 'wrong');
    await userEvent.click(submitButton);
  },
};

// ローディング状態
export const Loading: Story = {
  args: {},
  play: async ({ canvasElement }) => {
    // ローディング状態をシミュレート
    const canvas = within(canvasElement);
    const usernameInput = canvas.getByLabelText('ユーザー名');
    const passwordInput = canvas.getByLabelText('パスワード');
    const submitButton = canvas.getByRole('button', { name: 'サインイン' });
    
    await userEvent.type(usernameInput, 'admin');
    await userEvent.type(passwordInput, 'admin123');
    await userEvent.click(submitButton);
  },
};
