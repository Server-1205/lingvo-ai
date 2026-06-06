import { useState, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation } from '@tanstack/react-query';
import { chatSend } from '../api/client';
import type { ChatResponse, Correction } from '../api/client';
import { GrammarBlock } from './GrammarBlock';
import { UsageIndicator } from './UsageIndicator';

interface Message {
  role: 'user' | 'ai';
  text: string;
  corrections?: Correction[];
}

export function Chat() {
  const { t } = useTranslation();
  const [input, setInput] = useState('');
  const [messages, setMessages] = useState<Message[]>([]);
  const [usage, setUsage] = useState({ daily_used: 0, daily_limit: 10, is_premium: false });
  const listRef = useRef<HTMLDivElement>(null);

  const mutation = useMutation({
    mutationFn: (text: string) => chatSend({ text }),
    onSuccess: (data: ChatResponse) => {
      console.debug('[chat] response received', data);
      setMessages(prev => [...prev, {
        role: 'ai',
        text: data.reply,
        corrections: data.corrections,
      }]);
      setUsage(data.usage);
    },
    onError: (err: Error) => {
      console.debug('[chat] error', err.message);
    },
  });

  useEffect(() => {
    if (listRef.current) {
      listRef.current.scrollTop = listRef.current.scrollHeight;
    }
  }, [messages]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const text = input.trim();
    if (!text || mutation.isPending) return;

    console.debug('[chat] sending...', text);
    setMessages(prev => [...prev, { role: 'user', text }]);
    setInput('');
    mutation.mutate(text);
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%', flex: 1 }}>
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
          <div key={i} style={{ margin: '8px 16px' }}>
            <div style={{
              display: 'flex',
              justifyContent: msg.role === 'user' ? 'flex-end' : 'flex-start',
            }}>
              <div style={{
                maxWidth: '80%',
                padding: '10px 14px',
                borderRadius: msg.role === 'user' ? '16px 16px 4px 16px' : '16px 16px 16px 4px',
                background: msg.role === 'user' ? 'var(--tg-button)' : 'var(--tg-secondary-bg)',
                color: msg.role === 'user' ? 'var(--tg-button-text)' : 'var(--tg-text)',
                fontSize: 15,
                lineHeight: 1.4,
                whiteSpace: 'pre-wrap',
              }}>
                {msg.text}
              </div>
            </div>
            {msg.role === 'ai' && msg.corrections && msg.corrections.length > 0 && (
              <GrammarBlock corrections={msg.corrections} />
            )}
          </div>
        ))}
        {mutation.isPending && (
          <div style={{ margin: '8px 16px', color: 'var(--tg-hint)', fontSize: 14 }}>
            {t('common.loading')}
          </div>
        )}
      </div>

      <UsageIndicator usage={usage} />

      <form
        onSubmit={handleSubmit}
        style={{
          display: 'flex',
          gap: 8,
          padding: '8px 16px',
          borderTop: '1px solid var(--tg-border)',
          background: 'var(--tg-bg)',
        }}
      >
        <input
          value={input}
          onChange={e => setInput(e.target.value)}
          placeholder={t('chat.input_placeholder')}
          style={{
            flex: 1,
            padding: '10px 14px',
            border: '1px solid var(--tg-border)',
            borderRadius: 10,
            background: 'var(--tg-secondary-bg)',
            color: 'var(--tg-text)',
            fontSize: 15,
            outline: 'none',
          }}
        />
        <button
          type="submit"
          disabled={!input.trim() || mutation.isPending}
          className="btn btn-primary"
          style={{ padding: '10px 16px' }}
        >
          {t('chat.send')}
        </button>
      </form>
    </div>
  );
}
