import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { NavBar } from './components/NavBar';
import type { Tab } from './components/NavBar';
import { LanguageSwitcher } from './components/LanguageSwitcher';
import { Chat } from './components/Chat';
import { Vocabulary } from './components/Vocabulary';
import { ProgressView } from './components/Progress';
import { SubscriptionPlans } from './components/Subscription';
import { useTelegram } from './hooks/useTelegram';

function App() {
  const [activeTab, setActiveTab] = useState<Tab>('chat');
  const { t } = useTranslation();
  const { user, theme } = useTelegram();

  const headerBg = theme?.header_bg_color || theme?.bg_color || 'var(--tg-bg)';
  const textColor = theme?.text_color || 'var(--tg-text)';

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      minHeight: '100svh',
      background: 'var(--tg-bg)',
      color: textColor,
    }}>
      <header style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        padding: '8px 16px',
        paddingTop: 'calc(var(--safe-top) + 4px)',
        background: headerBg,
        borderBottom: '1px solid var(--tg-border)',
        minHeight: 'var(--header-height)',
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <span style={{ fontSize: 20 }}>✨</span>
          <span style={{ fontSize: 17, fontWeight: 600 }}>{t('app.title')}</span>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          {user && (
            <span style={{ fontSize: 12, color: 'var(--tg-hint)' }}>
              {user.first_name}
            </span>
          )}
          <LanguageSwitcher />
        </div>
      </header>

      <main style={{
        flex: 1,
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden',
      }}>
        {activeTab === 'chat' && <Chat />}
        {activeTab === 'vocab' && <Vocabulary />}
        {activeTab === 'progress' && <ProgressView />}
        {activeTab === 'subscription' && <SubscriptionPlans />}
      </main>

      <NavBar active={activeTab} onTabChange={setActiveTab} />
    </div>
  );
}

export default App;
