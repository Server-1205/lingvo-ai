import { useState, useCallback, useRef } from 'react';
import { initData } from '@telegram-apps/sdk';

export function useTTS() {
  const [isPlaying, setIsPlaying] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const currentAudio = useRef<HTMLAudioElement | null>(null);

  const play = useCallback(async (text: string, lang: string = 'uz') => {
    if (!text || isPlaying) return;

    setIsPlaying(true);
    setError(null);

    try {
      const raw = initData.raw();
      const headers: Record<string, string> = {};
      if (raw) headers['X-Telegram-Init-Data'] = raw;

      const res = await fetch(`/api/tts?text=${encodeURIComponent(text)}&lang=${lang}`, { headers });

      if (!res.ok) {
        const errBody = await res.json().catch(() => ({ error: 'tts_failed' }));
        throw new Error((errBody as Record<string, unknown>).error as string || 'tts_failed');
      }

      const blob = await res.blob();
      const url = URL.createObjectURL(blob);

      const audio = new Audio(url);
      currentAudio.current = audio;

      audio.onended = () => {
        URL.revokeObjectURL(url);
        setIsPlaying(false);
        currentAudio.current = null;
      };

      audio.onerror = () => {
        URL.revokeObjectURL(url);
        setIsPlaying(false);
        currentAudio.current = null;
        setError('playback_failed');
      };

      await audio.play();
    } catch (err) {
      setIsPlaying(false);
      setError((err as Error).message);
      console.debug('[tts] error:', (err as Error).message);
    }
  }, [isPlaying]);

  const stop = useCallback(() => {
    if (currentAudio.current) {
      currentAudio.current.pause();
      currentAudio.current = null;
    }
    setIsPlaying(false);
  }, []);

  return { play, stop, isPlaying, error };
}
