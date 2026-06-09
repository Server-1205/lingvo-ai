import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { ErrorDashboard } from '../ErrorDashboard';

const mockApi = {
  getErrorStats: vi.fn(),
  getErrorHistory: vi.fn(),
};

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      const map: Record<string, string> = {
        'common.loading': 'Загрузка...',
        'common.error': 'Произошла ошибка',
        'errors.premium_only': 'Только для Premium',
        'errors.no_errors': 'Ошибок пока нет!',
        'errors.total_errors': 'Всего ошибок',
        'errors.by_category': 'По категориям',
        'errors.by_severity': 'По серьёзности',
        'errors.frequent_rules': 'Частые правила',
        'errors.trend': 'Динамика',
        'errors.critical': 'Критические',
        'errors.major': 'Важные',
        'errors.minor': 'Незначительные',
        'errors.grammar': 'Грамматика',
        'errors.vocabulary': 'Лексика',
        'errors.spelling': 'Орфография',
        'errors.word_order': 'Порядок слов',
        'errors.punctuation': 'Пунктуация',
      };
      return map[key] || key;
    },
  }),
}));

vi.mock('../../api/client', () => ({
  getErrorStats: (...args: unknown[]) => mockApi.getErrorStats(...args),
  getErrorHistory: (...args: unknown[]) => mockApi.getErrorHistory(...args),
}));

const mockStats = {
  total_errors: 15,
  category_counts: { grammar: 8, vocabulary: 4, spelling: 2, word_order: 1, punctuation: 0 },
  severity_counts: { critical: 2, major: 5, minor: 8 },
  most_frequent_rules: ['Subject-Verb Agreement (5x)', 'Past Simple vs Present Perfect (3x)'],
  category_trend: [
    { date: '2026-06-01', grammar: 2, vocabulary: 1, spelling: 0, word_order: 0, punctuation: 0 },
    { date: '2026-06-02', grammar: 1, vocabulary: 0, spelling: 1, word_order: 0, punctuation: 0 },
  ],
};

const mockHistory = {
  entries: [
    {
      id: 1,
      user_id: 1,
      original: 'I goes',
      corrected: 'I go',
      category: 'grammar',
      severity: 'major',
      rule_violated: 'Subject-Verb Agreement',
      learning_tip: '',
      context: '',
      created_at: '2026-06-02T10:00:00Z',
    },
  ],
  total: 1,
};

describe('ErrorDashboard', () => {
  beforeEach(() => {
    mockApi.getErrorStats.mockReset();
    mockApi.getErrorHistory.mockReset();
  });

  it('shows loading state initially', () => {
    mockApi.getErrorStats.mockReturnValue(new Promise(() => {}));
    mockApi.getErrorHistory.mockReturnValue(new Promise(() => {}));

    render(<ErrorDashboard />);
    expect(screen.getByText('Загрузка...')).toBeInTheDocument();
  });

  it('shows premium only message when 403', async () => {
    mockApi.getErrorStats.mockRejectedValue({ status: 403 });
    mockApi.getErrorHistory.mockResolvedValue({ entries: [], total: 0 });

    render(<ErrorDashboard />);
    await waitFor(() => {
      expect(screen.getByText('Только для Premium')).toBeInTheDocument();
    });
  });

  it('shows no errors message when empty', async () => {
    mockApi.getErrorStats.mockResolvedValue({
      total_errors: 0,
      category_counts: {},
      severity_counts: {},
      most_frequent_rules: [],
    });
    mockApi.getErrorHistory.mockResolvedValue({ entries: [], total: 0 });

    render(<ErrorDashboard />);
    await waitFor(() => {
      expect(screen.getByText('Ошибок пока нет!')).toBeInTheDocument();
    });
  });

  it('shows total error count', async () => {
    mockApi.getErrorStats.mockResolvedValue(mockStats);
    mockApi.getErrorHistory.mockResolvedValue(mockHistory);

    render(<ErrorDashboard />);
    await waitFor(() => {
      expect(screen.getByText('15')).toBeInTheDocument();
    });
  });

  it('shows error by category', async () => {
    mockApi.getErrorStats.mockResolvedValue(mockStats);
    mockApi.getErrorHistory.mockResolvedValue(mockHistory);

    render(<ErrorDashboard />);
    await waitFor(() => {
      expect(screen.getByText(/Грамматика/)).toBeInTheDocument();
      expect(screen.getByText(/Лексика/)).toBeInTheDocument();
    });
  });

  it('shows error history entries', async () => {
    mockApi.getErrorStats.mockResolvedValue(mockStats);
    mockApi.getErrorHistory.mockResolvedValue(mockHistory);

    render(<ErrorDashboard />);
    await waitFor(() => {
      expect(screen.getByText('I goes')).toBeInTheDocument();
      expect(screen.getByText('I go')).toBeInTheDocument();
    });
  });
});
