import { useTranslation } from 'react-i18next';
import type { Correction } from '../api/client';

interface GrammarBlockProps {
  corrections: Correction[];
}

const typeLabels: Record<string, string> = {
  grammar: '📖',
  vocabulary: '📝',
  spelling: '✏️',
};

export function GrammarBlock({ corrections }: GrammarBlockProps) {
  const { t, i18n } = useTranslation();

  if (!corrections.length) {
    return (
      <div className="card" style={{ textAlign: 'center', color: 'var(--tg-hint)' }}>
        {t('chat.no_corrections')}
      </div>
    );
  }

  return (
    <div style={{ margin: '8px 16px' }}>
      <div style={{ fontSize: 15, fontWeight: 600, marginBottom: 8, color: 'var(--tg-text)' }}>
        {t('chat.corrections_title')}
      </div>
      {corrections.map((c, i) => (
        <div
          key={i}
          className="card"
          style={{ margin: '8px 0', padding: 12 }}
        >
          <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginBottom: 6 }}>
            <span>{typeLabels[c.type] || '📌'}</span>
            <span style={{
              fontSize: 11,
              background: 'var(--tg-secondary-bg)',
              padding: '2px 6px',
              borderRadius: 4,
              color: 'var(--tg-hint)',
            }}>
              {c.type}
            </span>
          </div>
          <div style={{ fontSize: 13, marginBottom: 4 }}>
            <span style={{ color: 'var(--tg-destructive)', textDecoration: 'line-through' }}>
              {c.original}
            </span>
            <span style={{ margin: '0 6px', color: 'var(--tg-hint)' }}>→</span>
            <span style={{ color: '#4caf50', fontWeight: 500 }}>{c.corrected}</span>
          </div>
          <div style={{ fontSize: 13, color: 'var(--tg-hint)', marginTop: 4 }}>
            {i18n.language === 'uz' ? c.explanation_uz : c.explanation_ru}
          </div>
        </div>
      ))}
    </div>
  );
}
