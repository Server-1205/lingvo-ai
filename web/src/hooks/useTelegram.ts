import { useEffect, useState } from 'react';
import { initData, themeParams } from '@telegram-apps/sdk';
import { debug } from '../lib/debug';

export interface TelegramContext {
  user: { id: number; first_name: string } | null;
  initDataRaw: string;
  theme: Record<string, string> | null;
  isReady: boolean;
}

export function useTelegram(): TelegramContext {
  const [theme, setTheme] = useState<Record<string, string> | null>(null);

  useEffect(() => {
    try {
      themeParams.mount().then(() => {
        setTheme(themeParams.state() as Record<string, string>);
      });
    } catch {
      debug('[tg] themeParams.mount() not available');
    }
  }, []);

  const rawUser = initData.user();
  const user = rawUser ? { id: rawUser.id, first_name: rawUser.first_name ?? '' } : null;

  return {
    user,
    initDataRaw: initData.raw() ?? '',
    theme,
    isReady: true,
  };
}
