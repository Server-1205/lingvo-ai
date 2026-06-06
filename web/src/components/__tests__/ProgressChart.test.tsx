import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { ProgressChart } from '../ProgressChart';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (k: string) => k }),
}));

const mockData = [
  { date: '2026-06-01', messages_sent: 5, words_learned: 3, quizzes_taken: 1 },
  { date: '2026-06-02', messages_sent: 8, words_learned: 2, quizzes_taken: 0 },
  { date: '2026-06-03', messages_sent: 3, words_learned: 5, quizzes_taken: 2 },
];

describe('ProgressChart', () => {
  it('renders chart with data', () => {
    render(<ProgressChart data={mockData} />);
    expect(screen.getByText('7 дней')).toBeInTheDocument();
    expect(screen.getByText('30 дней')).toBeInTheDocument();
  });

  it('shows empty state when no data', () => {
    render(<ProgressChart data={[]} />);
    expect(screen.getByText('Нет данных об активности за выбранный период')).toBeInTheDocument();
  });

  it('calls onPeriodChange when period button clicked', () => {
    const onPeriodChange = vi.fn();
    render(<ProgressChart data={mockData} onPeriodChange={onPeriodChange} />);
    screen.getByText('30 дней').click();
    expect(onPeriodChange).toHaveBeenCalledWith(30);
  });
});
