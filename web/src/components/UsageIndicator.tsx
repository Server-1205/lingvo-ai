import { useTranslation } from 'react-i18next';
import type { Usage } from '../api/client';

interface UsageIndicatorProps {
  usage: Usage;
  onUpgrade?: () => void;
}

export function UsageIndicator({ usage, onUpgrade }: UsageIndicatorProps) {
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
        <span style={{ color: 'var(--c-primary)', fontWeight: 600 }}>
          ⭐ {t('chat.unlimited')}
        </span>
      </div>
    );
  }

  const pct = usage.daily_limit > 0 ? (usage.daily_used / usage.daily_limit) * 100 : 0;
  const isExhausted = usage.daily_used >= usage.daily_limit;

  if (isExhausted) {
    return (
      <div style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        gap: 8,
        padding: '8px 16px',
      }}>
        <span style={{ fontSize: 13, color: 'var(--c-error)', fontWeight: 500 }}>
          {t('chat.limit_exhausted')}
        </span>
        <button
          className="btn btn-primary"
          style={{ fontSize: 13, padding: '6px 16px' }}
          onClick={onUpgrade}
        >
          {t('chat.get_unlimited')}
        </button>
      </div>
    );
  }

  return (
    <div style={{
      display: 'flex',
      alignItems: 'center',
      gap: 8,
      padding: '8px 16px',
      fontSize: 13,
      color: 'var(--c-on-surface-variant)',
    }}>
      <span style={{ fontWeight: 500, whiteSpace: 'nowrap' }}>
        {t('chat.daily_limit', { used: usage.daily_used, total: usage.daily_limit })}
      </span>
      <div style={{
        flex: 1,
        height: 6,
        borderRadius: 'var(--round-full)',
        background: 'var(--c-surface-container-highest)',
        overflow: 'hidden',
        maxWidth: 80,
      }}>
        <div style={{
          width: `${Math.min(pct, 100)}%`,
          height: '100%',
          borderRadius: 'var(--round-full)',
          background: 'var(--c-primary)',
          transition: 'width 0.3s',
        }} />
      </div>
    </div>
  );
}
