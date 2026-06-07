import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { getProgress, getProgressHistory } from '../api/client';
import { ProgressChart } from './ProgressChart';

interface ProgressViewProps {
  onStartLevelTest?: () => void;
}

export function ProgressView({ onStartLevelTest }: ProgressViewProps) {
  const { t } = useTranslation();

  const { data: progress, isLoading } = useQuery({
    queryKey: ['progress'],
    queryFn: getProgress,
    refetchInterval: 30_000,
  });

  const [period, setPeriod] = useState(7);

  const { data: history } = useQuery({
    queryKey: ['progress-history', period],
    queryFn: () => getProgressHistory(period),
  });

  console.debug('[progress] chart mounted, period=' + period);
  if (history) {
    console.debug('[progress] history loaded', history.entries.length, 'entries');
  }

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

      {history && (
        <ProgressChart
          data={history.entries}
          onPeriodChange={(d) => {
            console.debug('[progress] chart updated to period=' + d);
            setPeriod(d);
          }}
        />
      )}

      {!isLoading && (
        <div style={{ padding: '16px' }}>
          <button
            className="btn btn-primary"
            style={{ width: '100%', padding: 12 }}
            onClick={onStartLevelTest}
          >
            {t('level.take_test')}
          </button>
        </div>
      )}
    </div>
  );
}
