import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { NavBar } from '../NavBar';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      const map: Record<string, string> = {
        'nav.chat': 'Чат',
        'nav.vocab': 'Слова',
        'nav.progress': 'Прогресс',
        'nav.subscription': 'Подписка',
      };
      return map[key] || key;
    },
  }),
}));

describe('NavBar', () => {
  it('renders 4 tabs', () => {
    render(<NavBar active="chat" onTabChange={() => {}} />);
    expect(screen.getByText('Чат')).toBeInTheDocument();
    expect(screen.getByText('Слова')).toBeInTheDocument();
    expect(screen.getByText('Прогресс')).toBeInTheDocument();
    expect(screen.getByText('Подписка')).toBeInTheDocument();
  });
});
