import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import App from '../../App';

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      const map: Record<string, string> = {
        'app.title': 'Lingvo AI',
        'nav.chat': 'Чат',
        'nav.vocab': 'Слова',
        'nav.progress': 'Прогресс',
        'nav.subscription': 'Подписка',
        'common.loading': 'Загрузка...',
        'vocab.review_again': 'Снова',
        'vocab.review_hard': 'Сложно',
        'vocab.review_good': 'Хорошо',
        'vocab.review_easy': 'Легко',
        'chat.input_placeholder': 'Напишите...',
        'chat.send': 'Отправить',
        'onboarding.skip': 'Пропустить',
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
  getVocab: vi.fn().mockResolvedValue([]),
  addVocab: vi.fn(),
  deleteVocab: vi.fn(),
  lookupVocab: vi.fn(),
  getDueWords: vi.fn().mockResolvedValue([]),
  submitReview: vi.fn(),
  getProgress: vi.fn().mockResolvedValue({ messages_sent: 0, words_learned: 0, quizzes_taken: 0, streak_days: 0, level: 'a1' }),
  getSubscription: vi.fn().mockResolvedValue({ active: false }),
  chatSend: vi.fn().mockResolvedValue({ reply: '', corrections: [], usage: { daily_used: 0, daily_limit: 10, is_premium: false } }),
  createInvoice: vi.fn(),
  checkGrammar: vi.fn(),
  getQuiz: vi.fn(),
  getLevelTestQuestions: vi.fn(),
  saveLevel: vi.fn(),
  ReviewQuality: { Again: 0, Hard: 2, Good: 4, Easy: 5 },
}));

function renderWithProviders(ui: React.ReactElement) {
  return render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
}

describe('App', () => {
  it('renders onboarding on mount when not completed', () => {
    renderWithProviders(<App />);
    expect(screen.getByText('Пропустить')).toBeInTheDocument();
  });
});
