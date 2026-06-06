import { useTranslation } from 'react-i18next';
import type { Usage } from '../api/client';

interface UsageIndicatorProps {
  usage: Usage;
}

export function UsageIndicator({ usage }: UsageIndicatorProps) {
  const { t } = useTranslation();

  if (usage.is_premium) {
    return (
      <div style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: 6,
        padding: '8px 16px',
        fontSize: 13,
      }}>
        <span style={{ color: 'var(--tg-accent)' }}>⭐</span>
        <span style={{ color: 'var(--tg-hint)' }}>Unlimited</span>
      </div>
    );
  }

  const pct = usage.daily_limit > 0 ? (usage.daily_used / usage.daily_limit) * 100 : 0;
  const isExhausted = usage.daily_used >= usage.daily_limit;

  return (
    <div style={{
      display: 'flex',
      alignItems: 'center',
      gap: 8,
      padding: '8px 16px',
      fontSize: 13,
      color: isExhausted ? 'var(--tg-destructive)' : 'var(--tg-hint)',
    }}>
      <span>{t('chat.daily_limit', { used: usage.daily_used, total: usage.daily_limit })}</span>
      <div style={{
        flex: 1,
        height: 4,
        borderRadius: 2,
        background: 'var(--tg-secondary-bg)',
        overflow: 'hidden',
        maxWidth: 80,
      }}>
        <div style={{
          width: `${Math.min(pct, 100)}%`,
          height: '100%',
          borderRadius: 2,
          background: isExhausted ? 'var(--tg-destructive)' : 'var(--tg-button)',
          transition: 'width 0.3s',
        }} />
      </div>
    </div>
  );
}
