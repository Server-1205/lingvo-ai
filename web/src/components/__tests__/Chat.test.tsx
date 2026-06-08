import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Chat } from '../Chat';

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      const map: Record<string, string> = {
        'chat.input_placeholder': 'Напишите...',
        'chat.send': 'Отправить',
        'common.loading': 'Загрузка...',
      };
      return map[key] || key;
    },
    i18n: {
      language: 'ru',
      changeLanguage: vi.fn(),
    },
  }),
}));

vi.mock('../../hooks/useTelegram', () => ({
  useTelegram: () => ({
    user: null,
    theme: null,
  }),
}));

vi.mock('../../api/client', () => ({
  chatStream: vi.fn(),
  ApiError: class {
    name = 'ApiError';
    status: number;
    data: Record<string, unknown>;
    message: string;
    constructor(message: string, status: number, data: Record<string, unknown>) {
      this.message = message;
      this.status = status;
      this.data = data;
    }
  },
}));

function renderWithProviders(ui: React.ReactElement) {
  return render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
}

describe('Chat', () => {
  it('renders input field and send button on mount', () => {
    renderWithProviders(<Chat />);
    expect(screen.getByPlaceholderText('Напишите...')).toBeInTheDocument();
    expect(screen.getByText('Отправить')).toBeInTheDocument();
  });
});
