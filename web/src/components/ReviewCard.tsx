import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { ReviewQuality } from '../api/client';
import type { VocabWord } from '../api/client';
import { useTTS } from '../hooks/useTTS';

interface ReviewCardProps {
  word: VocabWord;
  onRate: (quality: number) => void;
}

export function ReviewCard({ word, onRate }: ReviewCardProps) {
  const { t, i18n } = useTranslation();
  const [flipped, setFlipped] = useState(false);
  const [loading, setLoading] = useState<number | null>(null);
  const { play: playTTS } = useTTS();
  const lang = i18n.language === 'ru' ? 'ru' : 'uz';
  const translation = lang === 'ru' && word.translation_ru ? word.translation_ru : word.translation;
  const example = lang === 'ru' && word.example_ru ? word.example_ru : word.example;

  const handleRate = async (quality: number) => {
    setLoading(quality);
    onRate(quality);
  };

  if (!flipped) {
    return (
      <div
        className="card review-card"
        onClick={() => setFlipped(true)}
        style={{
          cursor: 'pointer',
          textAlign: 'center',
          padding: 40,
          minHeight: 240,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          border: '1px solid var(--c-outline-variant)',
        }}
      >
        <div style={{ fontSize: 28, fontWeight: 700, marginBottom: 12, color: 'var(--c-on-surface)' }}>
          {word.word}
        </div>
        <button
          className="btn"
          style={{
            fontSize: 24,
            padding: '8px 16px',
            marginBottom: 12,
            borderRadius: 'var(--round-md)',
            background: 'var(--c-surface-container)',
          }}
          onClick={(e) => { e.stopPropagation(); playTTS(word.word); }}
        >
          🔊
        </button>
        {word.level && (
          <span style={{
            fontSize: 13,
            fontWeight: 600,
            background: 'var(--c-primary-fixed-dim)',
            color: 'var(--c-on-primary-fixed)',
            padding: '4px 12px',
            borderRadius: 'var(--round-full)',
          }}>
            {word.level.toUpperCase()}
          </span>
        )}
        <div style={{ fontSize: 15, color: 'var(--c-on-surface-variant)', marginTop: 20 }}>
          {t('vocab.review_again')} → {t('vocab.review_easy')}
        </div>
      </div>
    );
  }

  return (
    <div className="card" style={{
      padding: 28,
      minHeight: 240,
      border: '1px solid var(--c-outline-variant)',
    }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 8 }}>
        <div style={{ fontSize: 24, fontWeight: 700, color: 'var(--c-on-surface)' }}>{word.word}</div>
        <button
          className="btn"
          style={{
            fontSize: 18,
            padding: '4px 10px',
            lineHeight: 1,
            borderRadius: 'var(--round-md)',
            background: 'var(--c-surface-container)',
          }}
          onClick={() => playTTS(word.word)}
        >
          🔊
        </button>
      </div>
      <div style={{ fontSize: 17, color: 'var(--c-on-surface-variant)', marginBottom: 14 }}>{translation}</div>
      {example && (
        <div style={{
          fontSize: 15,
          fontStyle: 'italic',
          color: 'var(--c-on-surface-variant)',
          marginBottom: 20,
          padding: '12px 16px',
          background: 'var(--c-surface-container)',
          borderRadius: 'var(--round-md)',
        }}>
          {example}
        </div>
      )}

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 10, marginTop: 20 }}>
        {[
          { label: t('vocab.review_again'), quality: ReviewQuality.Again, color: 'var(--c-error)' },
          { label: t('vocab.review_hard'), quality: ReviewQuality.Hard, color: 'var(--c-secondary)' },
          { label: t('vocab.review_good'), quality: ReviewQuality.Good, color: 'var(--c-primary)' },
          { label: t('vocab.review_easy'), quality: ReviewQuality.Easy, color: '#2979FF' },
        ].map(({ label, quality, color }) => (
          <button
            key={quality}
            className="btn"
            style={{
              background: color,
              color: '#fff',
              padding: '14px 0',
              fontSize: 15,
              fontWeight: 700,
              opacity: loading !== null && loading !== quality ? 0.4 : 1,
            }}
            onClick={() => handleRate(quality)}
            disabled={loading !== null}
          >
            {loading === quality ? '...' : label}
          </button>
        ))}
      </div>
    </div>
  );
}
