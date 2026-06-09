import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { getErrorStats, getErrorHistory } from '../api/client';
import type { ErrorStatsResponse, ErrorHistoryEntry } from '../api/client';

const categoryIcons: Record<string, string> = {
  grammar: '📖',
  vocabulary: '📝',
  spelling: '✏️',
  word_order: '🔄',
  punctuation: '🔣',
};

const categoryColors: Record<string, string> = {
  grammar: '#4a90d9',
  vocabulary: '#50c878',
  spelling: '#ff9800',
  word_order: '#9c27b0',
  punctuation: '#607d8b',
};

const severityColors: Record<string, string> = {
  critical: '#f44336',
  major: '#ff9800',
  minor: '#4caf50',
};

export function ErrorDashboard() {
  const { t } = useTranslation();
  const [isLoading, setIsLoading] = useState(true);
  const [stats, setStats] = useState<ErrorStatsResponse | null>(null);
  const [history, setHistory] = useState<ErrorHistoryEntry[]>([]);
  const [historyTotal, setHistoryTotal] = useState(0);
  const [historyPage, setHistoryPage] = useState(1);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadStats();
    loadHistory(1);
  }, []);

  const loadStats = async () => {
    try {
      const data = await getErrorStats(30);
      setStats(data);
    } catch (err) {
      console.debug('[errors] failed to load stats', (err as Error).message);
      if ((err as any)?.status === 403) {
        setError('premium_only');
      } else {
        setError('load_error');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const loadHistory = async (page: number) => {
    try {
      const data = await getErrorHistory(page, 10);
      setHistory(data.entries);
      setHistoryTotal(data.total);
      setHistoryPage(page);
    } catch {
      console.debug('[errors] failed to load history');
    }
  };

  if (isLoading) {
    return (
      <div style={{ padding: 24, textAlign: 'center', color: 'var(--tg-hint)' }}>
        {t('common.loading')}
      </div>
    );
  }

  if (error === 'premium_only') {
    return (
      <div style={{ padding: 40, textAlign: 'center' }}>
        <div style={{ fontSize: 48, marginBottom: 16 }}>⭐</div>
        <div style={{ fontSize: 16, color: 'var(--tg-hint)', lineHeight: 1.6 }}>
          {t('errors.premium_only')}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ padding: 40, textAlign: 'center', color: 'var(--tg-destructive)' }}>
        {t('common.error')}
      </div>
    );
  }

  if (!stats || stats.total_errors === 0) {
    return (
      <div style={{ padding: 40, textAlign: 'center' }}>
        <div style={{ fontSize: 48, marginBottom: 16 }}>✅</div>
        <div style={{ fontSize: 16, color: 'var(--tg-hint)' }}>
          {t('errors.no_errors')}
        </div>
      </div>
    );
  }

  const maxCategoryCount = Math.max(
    ...Object.values(stats.category_counts).filter(v => v > 0),
    1,
  );

  const maxTrendValue = stats.category_trend
    ? Math.max(
        ...stats.category_trend.map(
          d => d.grammar + d.vocabulary + d.spelling + d.word_order + d.punctuation,
        ),
        1,
      )
    : 1;

  const totalPages = Math.ceil(historyTotal / 10);

  return (
    <div style={{ padding: '12px 16px', paddingBottom: 'calc(var(--nav-height) + var(--safe-bottom) + 80px)' }}>
      {/* Total errors */}
      <div className="card" style={{ textAlign: 'center', padding: 20, marginBottom: 12 }}>
        <div style={{ fontSize: 36, fontWeight: 700, color: 'var(--c-primary)' }}>
          {stats.total_errors}
        </div>
        <div style={{ fontSize: 13, color: 'var(--tg-hint)', marginTop: 4 }}>
          {t('errors.total_errors')}
        </div>
      </div>

      {/* By category */}
      <div className="card" style={{ padding: 16, marginBottom: 12 }}>
        <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 12 }}>
          {t('errors.by_category')}
        </div>
        {Object.entries(stats.category_counts).length === 0 ? (
          <div style={{ fontSize: 13, color: 'var(--tg-hint)' }}>{t('errors.no_errors')}</div>
        ) : (
          Object.entries(stats.category_counts).map(([cat, count]) => (
            <div key={cat} style={{ marginBottom: 10 }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 13, marginBottom: 4 }}>
                <span>
                  {categoryIcons[cat] || '📌'} {t(`errors.${cat}`, cat)}
                </span>
                <span style={{ fontWeight: 600 }}>{count}</span>
              </div>
              <div style={{
                height: 6,
                background: 'var(--c-surface-container)',
                borderRadius: 3,
                overflow: 'hidden',
              }}>
                <div style={{
                  width: `${(count / maxCategoryCount) * 100}%`,
                  height: '100%',
                  background: categoryColors[cat] || 'var(--c-primary)',
                  borderRadius: 3,
                  transition: 'width 0.3s',
                }} />
              </div>
            </div>
          ))
        )}
      </div>

      {/* By severity */}
      <div className="card" style={{ padding: 16, marginBottom: 12 }}>
        <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 12 }}>
          {t('errors.by_severity')}
        </div>
        <div style={{ display: 'flex', gap: 8 }}>
          {Object.entries(stats.severity_counts).length === 0 ? (
            <div style={{ fontSize: 13, color: 'var(--tg-hint)' }}>{t('errors.no_errors')}</div>
          ) : (
            Object.entries(stats.severity_counts).map(([sev, count]) => {
              const total = Object.values(stats.severity_counts).reduce((a, b) => a + b, 0);
              const pct = total > 0 ? Math.round((count / total) * 100) : 0;
              return (
                <div key={sev} style={{
                  flex: 1,
                  textAlign: 'center',
                  padding: '10px 6px',
                  borderRadius: 10,
                  background: `${severityColors[sev] || '#999'}15`,
                }}>
                  <div style={{
                    fontSize: 20,
                    fontWeight: 700,
                    color: severityColors[sev] || '#999',
                  }}>
                    {count}
                  </div>
                  <div style={{ fontSize: 11, color: 'var(--tg-hint)', marginTop: 2 }}>
                    {t(`errors.${sev}`, sev)}
                  </div>
                  <div style={{ fontSize: 10, color: severityColors[sev] || '#999' }}>
                    {pct}%
                  </div>
                </div>
              );
            })
          )}
        </div>
      </div>

      {/* Most frequent rules */}
      {stats.most_frequent_rules.length > 0 && (
        <div className="card" style={{ padding: 16, marginBottom: 12 }}>
          <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 8 }}>
            {t('errors.frequent_rules')}
          </div>
          {stats.most_frequent_rules.map((rule, i) => (
            <div key={i} style={{
              fontSize: 13,
              padding: '6px 0',
              borderBottom: i < stats.most_frequent_rules.length - 1 ? '1px solid var(--c-outline-variant)' : 'none',
              color: 'var(--tg-text)',
            }}>
              {i + 1}. {rule}
            </div>
          ))}
        </div>
      )}

      {/* Category trend chart */}
      {stats.category_trend && stats.category_trend.length > 0 && (
        <div className="card" style={{ padding: 16, marginBottom: 12 }}>
          <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 12 }}>
            {t('errors.trend')}
          </div>
          <div style={{ display: 'flex', alignItems: 'flex-end', gap: 3, height: 120, overflowX: 'auto' }}>
            {stats.category_trend.map((day, i) => {
              const total = day.grammar + day.vocabulary + day.spelling + day.word_order + day.punctuation;
              const height = maxTrendValue > 0 ? (total / maxTrendValue) * 100 : 0;
              return (
                <div key={i} style={{
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  minWidth: 28,
                }}>
                  <div style={{
                    width: 20,
                    height: `${Math.max(height, 2)}%`,
                    background: 'var(--c-primary)',
                    borderRadius: '3px 3px 0 0',
                    opacity: 0.7 + (height / 100) * 0.3,
                    position: 'relative',
                  }} />
                  <div style={{
                    fontSize: 9,
                    color: 'var(--tg-hint)',
                    marginTop: 4,
                    whiteSpace: 'nowrap',
                    transform: 'rotate(-45deg)',
                    transformOrigin: 'left',
                  }}>
                    {day.date.slice(5)}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Error history list */}
      {history.length > 0 && (
        <div className="card" style={{ padding: 16, marginBottom: 12 }}>
          <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 12 }}>
            {t('errors.history', 'Recent errors')}
          </div>
          {history.map((entry) => (
            <div key={entry.id} style={{
              padding: '8px 0',
              borderBottom: '1px solid var(--c-outline-variant)',
              fontSize: 13,
            }}>
              <div style={{ display: 'flex', gap: 6, marginBottom: 4 }}>
                <span>{categoryIcons[entry.category] || '📌'}</span>
                <span style={{
                  background: severityColors[entry.severity] || '#999',
                  color: '#fff',
                  padding: '1px 6px',
                  borderRadius: 4,
                  fontSize: 10,
                }}>
                  {entry.severity}
                </span>
                <span style={{ color: 'var(--tg-hint)', fontSize: 11 }}>
                  {new Date(entry.created_at).toLocaleDateString()}
                </span>
              </div>
              <div style={{ marginLeft: 22 }}>
                <span style={{ color: 'var(--tg-destructive)', textDecoration: 'line-through' }}>
                  {entry.original}
                </span>
                <span style={{ margin: '0 6px', color: 'var(--tg-hint)' }}>→</span>
                <span style={{ color: '#4caf50', fontWeight: 500 }}>{entry.corrected}</span>
              </div>
              {entry.rule_violated && (
                <div style={{ marginLeft: 22, fontSize: 11, color: 'var(--tg-hint)', fontStyle: 'italic' }}>
                  📏 {entry.rule_violated}
                </div>
              )}
            </div>
          ))}

          {/* Pagination */}
          {totalPages > 1 && (
            <div style={{ display: 'flex', justifyContent: 'center', gap: 8, marginTop: 12 }}>
              <button
                className="btn btn-secondary"
                style={{ padding: '6px 14px', fontSize: 13 }}
                disabled={historyPage <= 1}
                onClick={() => loadHistory(historyPage - 1)}
              >
                ←
              </button>
              <div style={{ fontSize: 13, alignSelf: 'center', color: 'var(--tg-hint)' }}>
                {historyPage}/{totalPages}
              </div>
              <button
                className="btn btn-secondary"
                style={{ padding: '6px 14px', fontSize: 13 }}
                disabled={historyPage >= totalPages}
                onClick={() => loadHistory(historyPage + 1)}
              >
                →
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
