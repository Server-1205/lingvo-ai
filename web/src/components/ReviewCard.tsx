import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { ReviewQuality } from '../api/client';
import type { VocabWord } from '../api/client';

interface ReviewCardProps {
  word: VocabWord;
  onRate: (quality: number) => void;
}

export function ReviewCard({ word, onRate }: ReviewCardProps) {
  const { t } = useTranslation();
  const [flipped, setFlipped] = useState(false);
  const [loading, setLoading] = useState<number | null>(null);

  const handleRate = async (quality: number) => {
    setLoading(quality);
    onRate(quality);
  };

  if (!flipped) {
    return (
      <div
        className="card review-card"
        onClick={() => setFlipped(true)}
        style={{ cursor: 'pointer', textAlign: 'center', padding: 32, minHeight: 200, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}
      >
        <div style={{ fontSize: 24, fontWeight: 700, marginBottom: 8 }}>{word.word}</div>
        <span style={{ fontSize: 12, background: 'var(--tg-secondary-bg)', padding: '2px 8px', borderRadius: 4 }}>
          {word.level?.toUpperCase()}
        </span>
        <div style={{ fontSize: 13, color: 'var(--tg-hint)', marginTop: 16 }}>
          {t('vocab.review_again')} → {t('vocab.review_easy')}
        </div>
      </div>
    );
  }

  return (
    <div className="card" style={{ padding: 24, minHeight: 200 }}>
      <div style={{ fontSize: 20, fontWeight: 700, marginBottom: 4 }}>{word.word}</div>
      <div style={{ fontSize: 15, color: 'var(--tg-hint)', marginBottom: 12 }}>{word.translation}</div>
      {word.example && (
        <div style={{ fontSize: 13, fontStyle: 'italic', color: 'var(--tg-hint)', marginBottom: 16, padding: '8px 12px', background: 'var(--tg-secondary-bg)', borderRadius: 8 }}>
          {word.example}
        </div>
      )}

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8, marginTop: 16 }}>
        {[
          { label: t('vocab.review_again'), quality: ReviewQuality.Again, color: '#e53935' },
          { label: t('vocab.review_hard'), quality: ReviewQuality.Hard, color: '#fb8c00' },
          { label: t('vocab.review_good'), quality: ReviewQuality.Good, color: '#43a047' },
          { label: t('vocab.review_easy'), quality: ReviewQuality.Easy, color: '#1e88e5' },
        ].map(({ label, quality, color }) => (
          <button
            key={quality}
            className="btn"
            style={{
              background: color,
              color: '#fff',
              padding: '12px 0',
              fontSize: 14,
              fontWeight: 600,
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
