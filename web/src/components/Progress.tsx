import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { getProgress } from '../api/client';

export function ProgressView() {
  const { t } = useTranslation();

  const { data: progress, isLoading } = useQuery({
    queryKey: ['progress'],
    queryFn: getProgress,
    refetchInterval: 30_000,
  });

  return (
    <div className="scroll-area">
      <div className="section-title">{t('progress.title')}</div>

      {isLoading && (
        <div className="placeholder">
          <div className="placeholder-icon">⏳</div>
          <div>{t('common.loading')}</div>
        </div>
      )}

      {progress && (
        <div className="card" style={{ display: 'flex', justifyContent: 'space-around', padding: 20 }}>
          <div style={{ textAlign: 'center' }}>
            <div style={{ fontSize: 28, fontWeight: 700, color: 'var(--tg-button)' }}>
              {progress.level?.toUpperCase()}
            </div>
            <div style={{ fontSize: 12, color: 'var(--tg-hint)', marginTop: 4 }}>{t('progress.level')}</div>
          </div>
          <div style={{ textAlign: 'center' }}>
            <div style={{ fontSize: 28, fontWeight: 700, color: 'var(--tg-button)' }}>
              {progress.messages_sent}
            </div>
            <div style={{ fontSize: 12, color: 'var(--tg-hint)', marginTop: 4 }}>{t('progress.messages')}</div>
          </div>
          <div style={{ textAlign: 'center' }}>
            <div style={{ fontSize: 28, fontWeight: 700, color: 'var(--tg-button)' }}>
              {progress.words_learned}
            </div>
            <div style={{ fontSize: 12, color: 'var(--tg-hint)', marginTop: 4 }}>{t('progress.words')}</div>
          </div>
          <div style={{ textAlign: 'center' }}>
            <div style={{ fontSize: 28, fontWeight: 700, color: 'var(--tg-button)' }}>
              {progress.streak_days}
            </div>
            <div style={{ fontSize: 12, color: 'var(--tg-hint)', marginTop: 4 }}>{t('progress.streak')}</div>
          </div>
        </div>
      )}
    </div>
  );
}
