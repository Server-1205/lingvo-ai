import { useState, useDebugValue } from 'react';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts';
import type { DailyProgressEntry } from '../api/client';

interface ProgressChartProps {
  data: DailyProgressEntry[];
  onPeriodChange?: (days: number) => void;
}

function formatDate(dateStr: string): string {
  const d = new Date(dateStr + 'T00:00:00');
  return `${String(d.getDate()).padStart(2, '0')}.${String(d.getMonth() + 1).padStart(2, '0')}`;
}

interface CustomTooltipProps {
  active?: boolean;
  payload?: Array<{ payload: DailyProgressEntry }>;
}

function CustomTooltip({ active, payload }: CustomTooltipProps) {
  if (!active || !payload || payload.length === 0) return null;
  const entry = payload[0].payload;
  const d = new Date(entry.date + 'T00:00:00');
  const dateStr = `${String(d.getDate()).padStart(2, '0')}.${String(d.getMonth() + 1).padStart(2, '0')}.${d.getFullYear()}`;
  return (
    <div style={{
      background: 'var(--tg-bg)',
      border: '1px solid var(--tg-border)',
      borderRadius: 8,
      padding: '8px 12px',
      fontSize: 13,
      boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
    }}>
      <div style={{ fontWeight: 600, marginBottom: 4, color: 'var(--tg-text)' }}>{dateStr}</div>
      <div style={{ color: 'var(--tg-hint)' }}>Сообщения: {entry.messages_sent}</div>
      <div style={{ color: 'var(--tg-hint)' }}>Слова: {entry.words_learned}</div>
      <div style={{ color: 'var(--tg-hint)' }}>Тесты: {entry.quizzes_taken}</div>
    </div>
  );
}

export function ProgressChart({ data, onPeriodChange }: ProgressChartProps) {
  const [period, setPeriod] = useState(7);

  useDebugValue(data.length > 0
    ? `chart: ${data.length} points, ${data[0].date} - ${data[data.length-1].date}`
    : 'chart: empty');

  console.debug('[progress-chart] rendering with', data.length, 'data points');
  if (data.length > 0) {
    console.debug('[progress-chart] date range:', data[0].date, '-', data[data.length - 1].date);
  }

  const handlePeriodChange = (days: number) => {
    console.debug('[progress-chart] period change:', period, '->', days);
    setPeriod(days);
    onPeriodChange?.(days);
  };

  if (data.length === 0) {
    return (
      <div className="card" style={{ padding: 20, textAlign: 'center' }}>
        <div style={{ fontSize: 13, color: 'var(--tg-hint)' }}>
          Нет данных об активности за выбранный период
        </div>
      </div>
    );
  }

  return (
    <div className="card" style={{ padding: '16px 8px' }}>
      <div style={{ display: 'flex', gap: 8, marginBottom: 12, justifyContent: 'center' }}>
        <button
          className={period === 7 ? 'btn btn-primary' : 'btn btn-secondary'}
          style={{ fontSize: 12, padding: '4px 12px' }}
          onClick={() => handlePeriodChange(7)}
        >
          7 дней
        </button>
        <button
          className={period === 30 ? 'btn btn-primary' : 'btn btn-secondary'}
          style={{ fontSize: 12, padding: '4px 12px' }}
          onClick={() => handlePeriodChange(30)}
        >
          30 дней
        </button>
      </div>

      <ResponsiveContainer width="100%" height={200}>
        <BarChart data={data} margin={{ top: 4, right: 8, bottom: 4, left: -16 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="var(--tg-border)" />
          <XAxis
            dataKey="date"
            tickFormatter={formatDate}
            tick={{ fontSize: 11, fill: 'var(--tg-hint)' }}
            interval="preserveStartEnd"
          />
          <YAxis
            allowDecimals={false}
            tick={{ fontSize: 11, fill: 'var(--tg-hint)' }}
          />
          <Tooltip content={<CustomTooltip />} />
          <Bar
            dataKey="messages_sent"
            fill="var(--tg-button)"
            radius={[3, 3, 0, 0]}
            maxBarSize={24}
          />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
