import { useState, useRef, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { chatStream, getSubscription, ApiError } from '../api/client';
import type { Correction, PremiumAnalysis, Usage } from '../api/client';
import { GrammarBlock } from './GrammarBlock';
import { PremiumCorrectionBlock } from './PremiumCorrectionBlock';
import { UsageIndicator } from './UsageIndicator';
import { debug } from '../lib/debug';

interface Message {
  role: 'user' | 'ai';
  text: string;
  corrections?: Correction[];
  premiumAnalysis?: PremiumAnalysis;
}

interface ChatProps {
  onUpgrade?: () => void;
  onStartReview?: () => void;
  onStartLevelTest?: () => void;
}

export function Chat({ onUpgrade, onStartReview, onStartLevelTest }: ChatProps) {
  const { t, i18n } = useTranslation();
  const [input, setInput] = useState('');
  const [messages, setMessages] = useState<Message[]>([]);
  const { data: sub } = useQuery({
    queryKey: ['subscription'],
    queryFn: getSubscription,
    retry: false,
    staleTime: 30_000,
  });
  const [usage, setUsage] = useState<Usage>({ daily_used: 0, daily_limit: 10, is_premium: false });
  useEffect(() => {
    if (sub?.active) {
      setUsage(prev => ({ ...prev, is_premium: true }));
    }
  }, [sub?.active]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [isListening, setIsListening] = useState(false);
  const listRef = useRef<HTMLDivElement>(null);
  const recognitionRef = useRef<any>(null);

  const SpeechRecognitionAPI = (window as any).SpeechRecognition || (window as any).webkitSpeechRecognition;
  const hasSpeechSupport = !!SpeechRecognitionAPI;

  const handleVoice = useCallback(() => {
    if (isListening) {
      recognitionRef.current?.stop();
      setIsListening(false);
      return;
    }

    if (!SpeechRecognitionAPI) return;

    const recognition = new SpeechRecognitionAPI();
    recognition.lang = i18n.language === 'ru' ? 'ru-RU' : 'uz-UZ';
    recognition.continuous = false;
    recognition.interimResults = false;

    recognition.onresult = (event: any) => {
      const result = event.results[event.results.length - 1];
      if (result.isFinal) {
        setInput(result[0].transcript);
      }
    };

    recognition.onerror = () => {
      setIsListening(false);
    };

    recognition.onend = () => {
      setIsListening(false);
    };

    recognitionRef.current = recognition;
    setInput('');
    recognition.start();
    setIsListening(true);
  }, [isListening, i18n.language, SpeechRecognitionAPI]);

  useEffect(() => {
    if (listRef.current) {
      listRef.current.scrollTop = listRef.current.scrollHeight;
    }
  }, [messages]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const text = input.trim();
    if (!text || isStreaming) return;

    debug('[chat] sending...', text);
    setInput('');
    setMessages(prev => [...prev, { role: 'user', text }]);

    const aiIndex = messages.length + 1;
    setMessages(prev => [...prev, { role: 'ai', text: '' }]);
    setIsStreaming(true);

    try {
      await chatStream(
        { text, lang: i18n.language },
        (token) => {
          setMessages(prev => {
            const next = [...prev];
            if (next[aiIndex]) {
              next[aiIndex] = { ...next[aiIndex], text: next[aiIndex].text + token };
            }
            return next;
          });
        },
        (corrections) => {
          setMessages(prev => {
            const next = [...prev];
            if (next[aiIndex]) {
              next[aiIndex] = { ...next[aiIndex], corrections };
            }
            return next;
          });
        },
        (newUsage) => {
          setUsage(newUsage);
          debug('[chat] premium path:', newUsage.is_premium);
        },
        () => {
          setIsStreaming(false);
        },
        (analysis) => {
          setMessages(prev => {
            const next = [...prev];
            if (next[aiIndex]) {
              next[aiIndex] = { ...next[aiIndex], premiumAnalysis: analysis };
            }
            return next;
          });
          debug('[chat] premium analysis received', analysis);
        },
      );
    } catch (err) {
      debug('[FIX] chat error', (err as Error).message);
      if (err instanceof ApiError && err.status === 429) {
        const errData = err.data as Record<string, unknown>;
        if (errData.daily_used !== undefined && errData.daily_limit !== undefined) {
          setUsage({
            daily_used: errData.daily_used as number,
            daily_limit: errData.daily_limit as number,
            is_premium: errData.is_premium as boolean || false,
          });
        }
      }
      setIsStreaming(false);
    }
  };

  const sendIcon = '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/></svg>';
  const micIcon = '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 1a3 3 0 0 0-3 3v8a3 3 0 0 0 6 0V4a3 3 0 0 0-3-3z"/><path d="M19 10v2a7 7 0 0 1-14 0v-2"/><line x1="12" y1="19" x2="12" y2="23"/></svg>';

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%', flex: 1, paddingBottom: 'calc(var(--nav-height) + var(--safe-bottom))' }}>
      <div
        ref={listRef}
        className="scroll-area"
        style={{ flex: 1, padding: '8px 0' }}
      >
        {messages.length === 0 && (
          <div className="placeholder">
            <div className="placeholder-icon">💬</div>
            <div>{t('chat.input_placeholder')}</div>
          </div>
        )}
        {messages.map((msg, i) => (
          <div key={i} style={{ margin: '10px 16px' }}>
            <div style={{
              display: 'flex',
              justifyContent: msg.role === 'user' ? 'flex-end' : 'flex-start',
            }}>
              <div style={{
                maxWidth: '85%',
                padding: '12px 18px',
                borderRadius: msg.role === 'user'
                  ? 'var(--round-lg) var(--round-lg) var(--round-sm) var(--round-lg)'
                  : 'var(--round-lg) var(--round-lg) var(--round-lg) var(--round-sm)',
                background: msg.role === 'user' ? 'var(--c-primary)' : 'var(--c-surface-container)',
                color: msg.role === 'user' ? '#fff' : 'var(--c-on-surface)',
                fontSize: 16,
                lineHeight: 1.5,
                whiteSpace: 'pre-wrap',
                boxShadow: msg.role === 'ai' ? 'var(--shadow-sm)' : 'var(--shadow-sm)',
              }}>
                {msg.text}
              </div>
            </div>
            {msg.role === 'ai' && msg.corrections && msg.corrections.length > 0 && (
              msg.premiumAnalysis ? (
                <PremiumCorrectionBlock
                  corrections={msg.corrections}
                  analysis={msg.premiumAnalysis}
                />
              ) : (
                <GrammarBlock corrections={msg.corrections} />
              )
            )}
          </div>
        ))}
        {isStreaming && (
          <div style={{ margin: '8px 16px', color: 'var(--tg-hint)', fontSize: 14 }}>
            {t('common.loading')}
          </div>
        )}
      </div>

      <UsageIndicator usage={usage} onUpgrade={onUpgrade} />

      {usage.daily_used >= usage.daily_limit && !usage.is_premium ? (
        <div style={{
          padding: '16px',
          borderTop: '1px solid var(--tg-border)',
          background: 'var(--tg-bg)',
        }}>
          <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 12, color: 'var(--tg-text)' }}>
            {t('chat.limit_alternatives')}
          </div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
            <button
              className="btn btn-secondary"
              style={{ width: '100%', textAlign: 'left', padding: '10px 14px', fontSize: 14 }}
              onClick={onStartReview}
            >
              📚 {t('vocab.review')}
            </button>
            <button
              className="btn btn-secondary"
              style={{ width: '100%', textAlign: 'left', padding: '10px 14px', fontSize: 14 }}
              onClick={onStartLevelTest}
            >
              📊 {t('level.take_test')}
            </button>
            <button
              className="btn btn-primary"
              style={{ width: '100%', padding: '10px 14px', fontSize: 14 }}
              onClick={onUpgrade}
            >
              ⭐ {t('chat.get_unlimited')}
            </button>
          </div>
        </div>
      ) : (
        <form
          onSubmit={handleSubmit}
          style={{
            display: 'flex',
            gap: 8,
            padding: '10px 16px',
            paddingBottom: 'calc(10px + var(--safe-bottom))',
            borderTop: '1px solid var(--tg-border)',
            background: 'var(--tg-bg)',
            alignItems: 'flex-end',
          }}
        >
          <div style={{
            flex: 1,
            display: 'flex',
            alignItems: 'center',
            gap: 8,
            padding: '0 4px 0 16px',
            borderRadius: 'var(--round-full)',
            background: 'var(--c-surface-container)',
            border: '1px solid var(--c-outline-variant)',
          }}>
            <input
              value={input}
              onChange={e => setInput(e.target.value)}
              placeholder={t('chat.input_placeholder')}
              style={{
                flex: 1,
                padding: '12px 0',
                border: 'none',
                background: 'transparent',
                color: 'var(--c-on-surface)',
                fontSize: 16,
                outline: 'none',
              }}
            />
            {hasSpeechSupport && (
              <button
                type="button"
                onClick={handleVoice}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  width: 36,
                  height: 36,
                  border: 'none',
                  borderRadius: '50%',
                  background: isListening ? 'var(--c-error)' : 'transparent',
                  color: isListening ? '#fff' : 'var(--c-outline)',
                  cursor: 'pointer',
                  flexShrink: 0,
                }}
                title={t('chat.voice_input')}
                dangerouslySetInnerHTML={{ __html: micIcon }}
              />
            )}
          </div>
          <button
            type="submit"
            disabled={!input.trim() || isStreaming}
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              width: 44,
              height: 44,
              border: 'none',
              borderRadius: '50%',
              background: input.trim() && !isStreaming ? 'var(--c-primary)' : 'var(--c-surface-container)',
              color: input.trim() && !isStreaming ? '#fff' : 'var(--c-outline)',
              cursor: input.trim() && !isStreaming ? 'pointer' : 'default',
              flexShrink: 0,
              transition: 'background 0.2s, color 0.2s',
            }}
            dangerouslySetInnerHTML={{ __html: sendIcon }}
          />
        </form>
      )}
    </div>
  );
}
