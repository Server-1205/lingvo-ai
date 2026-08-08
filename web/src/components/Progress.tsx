import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { LoadingDots } from './LoadingDots';
import { useQuery } from '@tanstack/react-query';
import { getProgress, getProgressHistory } from '../api/client';
import { ProgressChart } from './ProgressChart';
import { debug } from '../lib/debug';

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

  debug('[progress] chart mounted, period=' + period);
  if (history) {
    debug('[progress] history loaded', history.entries.length, 'entries');
  }

  return (
    <div className="scroll-area">
      <div className="section-title">{t('progress.title')}</div>

      {isLoading && (
        <div className="placeholder">
          <LoadingDots size={14} />
          <div>{t('common.loading')}</div>
        </div>
      )}

      {progress && (
        <div style={{
          margin: '0 16px 16px',
          padding: '24px 12px',
          borderRadius: 'var(--round-lg)',
          border: '1px solid var(--c-outline-variant)',
          background: 'var(--c-surface-container-lowest)',
          display: 'flex',
          justifyContent: 'space-around',
        }}>
          {[
            { value: progress.level?.toUpperCase() || '?', label: t('progress.level') },
            { value: progress.messages_sent, label: t('progress.messages') },
            { value: progress.words_learned, label: t('progress.words') },
            { value: progress.streak_days, label: t('progress.streak') },
          ].map((stat, i) => (
            <div key={i} style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 32, fontWeight: 700, color: 'var(--c-primary)' }}>
                {stat.value}
              </div>
              <div style={{ fontSize: 13, color: 'var(--c-on-surface-variant)', marginTop: 4, fontWeight: 500 }}>
                {stat.label}
              </div>
            </div>
          ))}
        </div>
      )}

      {history && (
        <ProgressChart
          data={history.entries}
          onPeriodChange={(d) => {
            debug('[progress] chart updated to period=' + d);
            setPeriod(d);
          }}
        />
      )}

      {!isLoading && (
        <div style={{ padding: '16px' }}>
          <button
            className="btn btn-primary"
            style={{ width: '100%', padding: 14, fontSize: 16, fontWeight: 700 }}
            onClick={onStartLevelTest}
          >
            {t('level.take_test')}
          </button>
        </div>
      )}
    </div>
  );
}
