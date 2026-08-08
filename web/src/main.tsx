import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { init, isTMA, mockTelegramEnv, initData } from '@telegram-apps/sdk';
import { debug } from './lib/debug';
import './i18n';
import './index.css';
import App from './App.tsx';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 30_000,
    },
  },
});

async function initTelegramSDK() {
  try {
    if (await isTMA()) {
      init();
      initData.restore();
      debug('[tg] SDK initialized in Telegram WebView');
    } else {
      const launchParams = 'tgWebAppData=auth_date%3D1729680000%26hash%3D3124d540cea212bcf0245f130c5213e59f21fd4c97a4083ecff68e548a57fbe3%26query_id%3DAAHcFpETAAAAANwWkRMB5HMN%26user%3D%257B%2522id%2522%253A12345%252C%2522first_name%2522%253A%2522Dev%2522%252C%2522username%2522%253A%2522devuser%2522%252C%2522language_code%2522%253A%2522ru%2522%257D&tgWebAppVersion=8.0&tgWebAppPlatform=web&tgWebAppThemeParams=%7B%22bg_color%22%3A%22%23ffffff%22%2C%22text_color%22%3A%22%23000000%22%2C%22hint_color%22%3A%22%23999999%22%2C%22button_color%22%3A%22%232481cc%22%2C%22button_text_color%22%3A%22%23ffffff%22%2C%22secondary_bg_color%22%3A%22%23f4f4f5%22%7D';
      mockTelegramEnv({ launchParams });
      init();
      initData.restore();
      debug('[tg] mock Telegram environment for development');
    }
  } catch {
    debug('[tg] not in Telegram WebView, running in browser mode');
  }
}

initTelegramSDK().then(() => {
  createRoot(document.getElementById('root')!).render(
    <StrictMode>
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
    </StrictMode>,
  );
});
