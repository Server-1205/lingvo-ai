import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { UsageIndicator } from '../UsageIndicator';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, vars?: Record<string, unknown>) => {
      if (key === 'chat.daily_limit') {
        return `Сегодня: ${vars?.used}/${vars?.total}`;
      }
      if (key === 'chat.limit_exhausted') return 'Дневной лимит исчерпан';
      if (key === 'chat.get_unlimited') return 'Получить безлимит';
      return key;
    },
  }),
}));

describe('UsageIndicator', () => {
  it('shows used/limit for free users', () => {
    render(<UsageIndicator usage={{ daily_used: 4, daily_limit: 10, is_premium: false }} />);
    expect(screen.getByText('Сегодня: 4/10')).toBeInTheDocument();
  });

  it('shows Unlimited for premium users', () => {
    render(<UsageIndicator usage={{ daily_used: 0, daily_limit: 0, is_premium: true }} />);
    expect(screen.getByText('Unlimited')).toBeInTheDocument();
  });

  it('shows exhausted message and upgrade button when limit reached', () => {
    render(<UsageIndicator usage={{ daily_used: 10, daily_limit: 10, is_premium: false }} />);
    expect(screen.getByText('Дневной лимит исчерпан')).toBeInTheDocument();
    expect(screen.getByText('Получить безлимит')).toBeInTheDocument();
  });

  it('calls onUpgrade when upgrade button is clicked', () => {
    const onUpgrade = vi.fn();
    render(<UsageIndicator usage={{ daily_used: 10, daily_limit: 10, is_premium: false }} onUpgrade={onUpgrade} />);
    screen.getByText('Получить безлимит').click();
    expect(onUpgrade).toHaveBeenCalledTimes(1);
  });
});
