import { useState, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { chatStream } from '../api/client';
import type { Correction, Usage } from '../api/client';
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
  const [usage, setUsage] = useState<Usage>({ daily_used: 0, daily_limit: 10, is_premium: false });
  const [isStreaming, setIsStreaming] = useState(false);
  const listRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (listRef.current) {
      listRef.current.scrollTop = listRef.current.scrollHeight;
    }
  }, [messages]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const text = input.trim();
    if (!text || isStreaming) return;

    console.debug('[chat] sending...', text);
    setInput('');
    setMessages(prev => [...prev, { role: 'user', text }]);

    const aiIndex = messages.length + 1;
    setMessages(prev => [...prev, { role: 'ai', text: '' }]);
    setIsStreaming(true);

    try {
      await chatStream(
        { text },
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
        },
        () => {
          setIsStreaming(false);
        },
      );
    } catch (err) {
      console.debug('[chat] error', (err as Error).message);
      setIsStreaming(false);
    }
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
        {isStreaming && (
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
          disabled={!input.trim() || isStreaming}
          className="btn btn-primary"
          style={{ padding: '10px 16px' }}
        >
          {t('chat.send')}
        </button>
      </form>
    </div>
  );
}
