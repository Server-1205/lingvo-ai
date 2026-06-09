import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { getIeltsScores } from '../api/client';
import type { IeltsScoreStats } from '../api/client';
import { IeltsWriting } from './IeltsWriting';
import { IeltsSpeaking } from './IeltsSpeaking';
import { IeltsReading } from './IeltsReading';

type Module = 'dashboard' | 'writing' | 'speaking' | 'reading';

export function IeltsDashboard() {
  const { t } = useTranslation();
  const [module, setModule] = useState<Module>('dashboard');

  const scoresQuery = useQuery({
    queryKey: ['ielts-scores'],
    queryFn: () => getIeltsScores(),
    retry: 1,
  });

  if (module === 'writing') {
    return <IeltsWriting onBack={() => setModule('dashboard')} />;
  }
  if (module === 'speaking') {
    return <IeltsSpeaking onBack={() => setModule('dashboard')} />;
  }
  if (module === 'reading') {
    return <IeltsReading onBack={() => setModule('dashboard')} />;
  }

  const is403 = scoresQuery.error && (scoresQuery.error as { status?: number })?.status === 403;

  if (is403) {
    return (
      <div className="scroll-area" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '60vh' }}>
        <div className="placeholder">
          <div className="placeholder-icon">⭐</div>
          <div>{t('ielts.premium_only')}</div>
        </div>
      </div>
    );
  }

  const stats: IeltsScoreStats | undefined = scoresQuery.data?.stats;

  const modules: { key: Module; icon: string; labelKey: string; avg: number | undefined }[] = [
    { key: 'writing', icon: '📝', labelKey: 'ielts.writing', avg: undefined },
    { key: 'speaking', icon: '🎤', labelKey: 'ielts.speaking', avg: undefined },
    { key: 'reading', icon: '📖', labelKey: 'ielts.reading', avg: undefined },
  ];

  if (stats) {
    modules[0].avg = Math.max(stats.writing_task1_avg, stats.writing_task2_avg);
    modules[1].avg = stats.speaking_avg;
    modules[2].avg = stats.reading_avg;
  }

  return (
    <div className="scroll-area">
      <div style={{ padding: 16 }}>
        <div style={{ fontSize: 18, fontWeight: 700, marginBottom: 4 }}>🎯 IELTS</div>
        <div style={{ fontSize: 13, color: 'var(--tg-hint)', marginBottom: 16 }}>
          {scoresQuery.isLoading ? t('common.loading') : `${t('ielts.practices')}: ${stats?.total_practices || 0}`}
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
          {modules.map(m => (
            <button
              key={m.key}
              className="btn card"
              style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                padding: 16,
                cursor: 'pointer',
                textAlign: 'left',
              }}
              onClick={() => setModule(m.key)}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                <span style={{ fontSize: 24 }}>{m.icon}</span>
                <div>
                  <div style={{ fontSize: 15, fontWeight: 600 }}>{t(m.labelKey)}</div>
                  {m.avg !== undefined && m.avg > 0 && (
                    <div style={{ fontSize: 12, color: 'var(--tg-hint)', marginTop: 2 }}>
                      {t('ielts.avg_score')}: {m.avg.toFixed(1)}
                    </div>
                  )}
                </div>
              </div>
              <span style={{ fontSize: 18, color: 'var(--tg-hint)' }}>→</span>
            </button>
          ))}
        </div>

        {scoresQuery.data && scoresQuery.data.entries.length > 0 && (
          <div style={{ marginTop: 24 }}>
            <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 12 }}>{t('ielts.history')}</div>
            {scoresQuery.data.entries.slice(0, 10).map(entry => (
              <div key={entry.id} className="card" style={{
                padding: 12, marginBottom: 8, fontSize: 13,
                display: 'flex', justifyContent: 'space-between', alignItems: 'center',
              }}>
                <div>
                  <div style={{ fontWeight: 500 }}>{entry.module.replace('_', ' ')}</div>
                  <div style={{ fontSize: 11, color: 'var(--tg-hint)' }}>{new Date(entry.created_at).toLocaleDateString()}</div>
                </div>
                <div style={{ fontSize: 18, fontWeight: 700, color: 'var(--c-primary)' }}>
                  {entry.band_score.toFixed(1)}
                </div>
              </div>
            ))}
          </div>
        )}

        {scoresQuery.isLoading && (
          <div className="placeholder" style={{ marginTop: 24 }}>
            <div className="placeholder-icon">⏳</div>
            <div>{t('common.loading')}</div>
          </div>
        )}
      </div>
    </div>
  );
}
