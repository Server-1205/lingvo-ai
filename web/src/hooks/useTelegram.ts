import { useEffect, useState } from 'react';
import { initData, themeParams } from '@telegram-apps/sdk';
import type { User, ThemeParams } from '@telegram-apps/types';

export interface TelegramContext {
  user: User | null;
  initDataRaw: string;
  theme: ThemeParams | null;
  isReady: boolean;
}

export function useTelegram(): TelegramContext {
  const [theme, setTheme] = useState<ThemeParams | null>(null);

  useEffect(() => {
    themeParams.mount().then(() => {
      setTheme(themeParams.state());
    });
  }, []);

  return {
    user: initData.user() ?? null,
    initDataRaw: initData.raw() ?? '',
    theme,
    isReady: true,
  };
}
