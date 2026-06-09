import { useTranslation } from 'react-i18next';

export type Tab = 'chat' | 'vocab' | 'progress' | 'subscription' | 'errors' | 'ielts';

interface NavBarProps {
  active: Tab;
  onTabChange: (tab: Tab) => void;
  isPremium?: boolean;
}

const tabs: { key: Tab; icon: string; labelKey: string; premium?: boolean }[] = [
  { key: 'chat', icon: '💬', labelKey: 'nav.chat' },
  { key: 'vocab', icon: '📚', labelKey: 'nav.vocab' },
  { key: 'progress', icon: '📊', labelKey: 'nav.progress' },
  { key: 'errors', icon: '🔍', labelKey: 'nav.errors', premium: true },
  { key: 'ielts', icon: '🎯', labelKey: 'nav.ielts', premium: true },
  { key: 'subscription', icon: '⭐', labelKey: 'nav.subscription' },
];

export function NavBar({ active, onTabChange, isPremium }: NavBarProps) {
  const { t } = useTranslation();

  return (
    <nav style={{
      display: 'flex',
      position: 'fixed',
      bottom: 0,
      left: 0,
      right: 0,
      background: 'var(--tg-bg)',
      borderTop: '1px solid var(--tg-border)',
      paddingBottom: 'var(--safe-bottom)',
      zIndex: 100,
    }}>
      {tabs.map(tab => {
        if (tab.premium && !isPremium) return null;

        return (
          <button
            key={tab.key}
            onClick={() => onTabChange(tab.key)}
            style={{
              flex: 1,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              gap: 2,
              padding: '8px 4px',
              border: 'none',
              background: 'transparent',
              color: active === tab.key ? 'var(--tg-button)' : 'var(--tg-hint)',
              fontSize: 10,
              cursor: 'pointer',
              transition: 'color 0.2s',
            }}
          >
            <span style={{ fontSize: 22 }}>{tab.icon}</span>
            <span>{t(tab.labelKey)}</span>
          </button>
        );
      })}
    </nav>
  );
}
