import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { initData } from '@telegram-apps/sdk';
import { NavBar } from './components/NavBar';
import type { Tab } from './components/NavBar';
import { Onboarding } from './components/Onboarding';
import { Chat } from './components/Chat';
import { Vocabulary } from './components/Vocabulary';
import { ProgressView } from './components/Progress';
import { SubscriptionPlans } from './components/Subscription';
import { LevelTest } from './components/LevelTest';
import { DailyLesson } from './components/DailyLesson';
import { useTelegram } from './hooks/useTelegram';

function App() {
  const [activeTab, setActiveTab] = useState<Tab>('chat');
  const [showLevelTest, setShowLevelTest] = useState(false);
  const [showDailyLesson, setShowDailyLesson] = useState(false);
  const [showOnboarding, setShowOnboarding] = useState(() => !localStorage.getItem('lingvo_onboarding_done'));
  const [vocabInitialTab, setVocabInitialTab] = useState<'my' | 'lookup' | 'review' | undefined>(undefined);
  const { t, i18n } = useTranslation();
  const { user, theme } = useTelegram();

  useEffect(() => {
    const tgLang = initData.user()?.language_code;
    if (tgLang === 'ru' || tgLang === 'uz') {
      i18n.changeLanguage(tgLang);
    }
  }, [i18n]);

  useEffect(() => {
    try {
      const tg = (window as any).Telegram?.WebApp;
      const tgStartParam = tg?.initDataUnsafe?.start_param;
      const urlParams = new URLSearchParams(window.location.search);
      const urlStartParam = urlParams.get('startapp') || tgStartParam;
      if (urlStartParam === 'review') {
        console.debug('[app] deep-link start_param=review → navigating to vocab/review');
        setActiveTab('vocab');
        setVocabInitialTab('review');
      } else if (urlStartParam === 'daily') {
        console.debug('[app] deep-link start_param=daily → showing daily lesson');
        setShowDailyLesson(true);
      } else if (urlStartParam) {
        console.debug('[app] deep-link start_param=' + urlStartParam + ' (unhandled)');
      }
    } catch {
      // initData not available outside Telegram
    }
  }, []);

  const textColor = theme?.text_color || 'var(--tg-text)';

  if (showOnboarding) {
    return <Onboarding onDone={() => setShowOnboarding(false)} />;
  }

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
        padding: '10px 16px',
        paddingTop: 'calc(var(--safe-top) + 6px)',
        background: 'linear-gradient(135deg, var(--c-primary-dark), var(--c-primary))',
        borderBottom: 'none',
        minHeight: 'var(--header-height)',
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
          <span style={{ fontSize: 24 }}>✨</span>
          <span style={{ fontSize: 19, fontWeight: 700, color: '#fff' }}>{t('app.title')}</span>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          {user && (
            <span style={{ fontSize: 13, color: 'rgba(255,255,255,0.85)' }}>
              {user.first_name}
            </span>
          )}
        </div>
      </header>

      <main style={{
        flex: 1,
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden',
      }}>
        {showDailyLesson ? (
          <DailyLesson onBack={() => setShowDailyLesson(false)} />
        ) : showLevelTest ? (
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
