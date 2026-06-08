import { useTranslation } from 'react-i18next';

export type Tab = 'chat' | 'vocab' | 'progress' | 'subscription';

interface NavBarProps {
  active: Tab;
  onTabChange: (tab: Tab) => void;
}

const tabs: { key: Tab; icon: string; activeIcon: string; labelKey: string }[] = [
  {
    key: 'chat',
    icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>',
    activeIcon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>',
    labelKey: 'nav.chat',
  },
  {
    key: 'vocab',
    icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/></svg>',
    activeIcon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/></svg>',
    labelKey: 'nav.vocab',
  },
  {
    key: 'progress',
    icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>',
    activeIcon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>',
    labelKey: 'nav.progress',
  },
  {
    key: 'subscription',
    icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>',
    activeIcon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>',
    labelKey: 'nav.subscription',
  },
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
      background: 'rgba(255,255,255,0.92)',
      backdropFilter: 'blur(10px)',
      WebkitBackdropFilter: 'blur(10px)',
      borderTop: 'none',
      paddingBottom: 'var(--safe-bottom)',
      height: 'calc(64px + var(--safe-bottom))',
      zIndex: 100,
      boxShadow: '0 -2px 8px rgba(0,0,0,0.06)',
    }}>
      {tabs.map(tab => {
        const isActive = active === tab.key;
        return (
          <button
            key={tab.key}
            onClick={() => onTabChange(tab.key)}
            style={{
              flex: 1,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              gap: 2,
              padding: '8px 4px',
              border: 'none',
              background: 'transparent',
              color: isActive ? 'var(--c-primary)' : '#9E9E9E',
              fontSize: 12,
              fontWeight: isActive ? 600 : 500,
              cursor: 'pointer',
              position: 'relative',
              transition: 'color 0.2s',
              letterSpacing: '0.05em',
            }}
          >
            <span
              style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', width: 24, height: 24 }}
              dangerouslySetInnerHTML={{ __html: isActive ? tab.activeIcon : tab.icon }}
            />
            <span>{t(tab.labelKey)}</span>
            {isActive && (
              <span style={{
                position: 'absolute',
                bottom: 2,
                width: 4,
                height: 4,
                borderRadius: '50%',
                background: 'var(--c-primary)',
              }} />
            )}
          </button>
        );
      })}
    </nav>
  );
}
