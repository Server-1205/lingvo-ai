import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { NavBar } from './components/NavBar';
import type { Tab } from './components/NavBar';
import { LanguageSwitcher } from './components/LanguageSwitcher';
import { Chat } from './components/Chat';
import { Vocabulary } from './components/Vocabulary';
import { ProgressView } from './components/Progress';
import { SubscriptionPlans } from './components/Subscription';
import { LevelTest } from './components/LevelTest';
import { useTelegram } from './hooks/useTelegram';

function App() {
  const [activeTab, setActiveTab] = useState<Tab>('chat');
  const [showLevelTest, setShowLevelTest] = useState(false);
  const [vocabInitialTab, setVocabInitialTab] = useState<'my' | 'lookup' | 'review' | undefined>(undefined);
  const { t } = useTranslation();
  const { user, theme } = useTelegram();

  useEffect(() => {
    try {
      const tg = (window as any).Telegram?.WebApp;
      const startParam = tg?.initDataUnsafe?.start_param;
      if (startParam === 'review') {
        console.debug('[app] deep-link start_param=review → navigating to vocab/review');
        setActiveTab('vocab');
        setVocabInitialTab('review');
      } else if (startParam) {
        console.debug('[app] deep-link start_param=' + startParam + ' (unhandled)');
      }
    } catch {
      // initData not available outside Telegram
    }
  }, []);

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
        {showLevelTest ? (
          <LevelTest onDone={() => setShowLevelTest(false)} />
        ) : (
          <>
            {activeTab === 'chat' && <Chat
              onUpgrade={() => setActiveTab('subscription')}
              onStartReview={() => { setActiveTab('vocab'); setVocabInitialTab('review'); }}
              onStartLevelTest={() => setShowLevelTest(true)}
            />}
            {activeTab === 'vocab' && <Vocabulary initialTab={vocabInitialTab} />}
            {activeTab === 'progress' && <ProgressView onStartLevelTest={() => setShowLevelTest(true)} />}
            {activeTab === 'subscription' && <SubscriptionPlans />}
          </>
        )}
      </main>

      {!showLevelTest && <NavBar active={activeTab} onTabChange={setActiveTab} />}
    </div>
  );
}

export default App;
