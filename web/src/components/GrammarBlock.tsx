import { useTranslation } from 'react-i18next';
import type { Correction } from '../api/client';

interface GrammarBlockProps {
  corrections: Correction[];
}

const typeColors: Record<string, string> = {
  grammar: 'var(--c-primary)',
  vocabulary: 'var(--c-secondary)',
  spelling: 'var(--c-error)',
};

const typeLabels: Record<string, string> = {
  grammar: 'Grammar',
  vocabulary: 'Vocabulary',
  spelling: 'Spelling',
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
    <div style={{ marginTop: 8 }}>
      <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 6, color: 'var(--c-on-surface-variant)' }}>
        {t('chat.corrections_title')}
      </div>
      {corrections.map((c, i) => {
        const color = typeColors[c.type] || 'var(--c-primary)';
        return (
          <div
            key={i}
            style={{
              padding: 12,
              marginBottom: 6,
              borderRadius: 'var(--round-md)',
              background: 'var(--c-surface-container-lowest)',
              border: '1px solid var(--c-outline-variant)',
            }}
          >
            <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginBottom: 6 }}>
              <span style={{
                fontSize: 11,
                fontWeight: 600,
                color: color,
                background: `${color}1A`,
                padding: '2px 8px',
                borderRadius: 'var(--round-full)',
                letterSpacing: '0.03em',
                textTransform: 'uppercase',
              }}>
                {typeLabels[c.type] || c.type}
              </span>
            </div>
            <div style={{ fontSize: 14, marginBottom: 4, lineHeight: 1.5 }}>
              <span style={{ color: 'var(--c-error)', textDecoration: 'line-through' }}>
                {c.original}
              </span>
              <span style={{ margin: '0 6px', color: 'var(--c-outline)' }}>→</span>
              <span style={{ color: 'var(--c-tertiary)', fontWeight: 600 }}>{c.corrected}</span>
            </div>
            <div style={{ fontSize: 13, color: 'var(--c-on-surface-variant)', marginTop: 4, lineHeight: 1.4 }}>
              {i18n.language === 'uz' ? c.explanation_uz : c.explanation_ru}
            </div>
          </div>
        );
      })}
    </div>
  );
}
