import { useTranslation } from 'react-i18next';

export type Tab = 'chat' | 'vocab' | 'progress' | 'subscription';

interface NavBarProps {
  active: Tab;
  onTabChange: (tab: Tab) => void;
}

const tabs: { key: Tab; icon: string; labelKey: string }[] = [
  { key: 'chat', icon: '💬', labelKey: 'nav.chat' },
  { key: 'vocab', icon: '📚', labelKey: 'nav.vocab' },
  { key: 'progress', icon: '📊', labelKey: 'nav.progress' },
  { key: 'subscription', icon: '⭐', labelKey: 'nav.subscription' },
];

export function NavBar({ active, onTabChange }: NavBarProps) {
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
      boxShadow: '0 -2px 8px rgba(0,0,0,0.06)',
    }}>
      {tabs.map(tab => (
        <button
          key={tab.key}
          onClick={() => onTabChange(tab.key)}
          style={{
            flex: 1,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            gap: 3,
            padding: '10px 4px',
            border: 'none',
            background: 'transparent',
            color: active === tab.key ? 'var(--c-primary)' : '#9E9E9E',
            fontSize: 12,
            fontWeight: active === tab.key ? 700 : 500,
            cursor: 'pointer',
            transition: 'color 0.2s',
          }}
        >
          <span style={{ fontSize: 24 }}>{tab.icon}</span>
          <span>{t(tab.labelKey)}</span>
        </button>
      ))}
    </nav>
  );
}
