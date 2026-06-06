import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { UsageIndicator } from '../UsageIndicator';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, vars?: Record<string, unknown>) => {
      if (key === 'chat.daily_limit') {
        return `Сегодня: ${vars?.used}/${vars?.total}`;
      }
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
});
