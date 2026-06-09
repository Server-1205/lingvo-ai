import { useTranslation } from 'react-i18next';
import type { Correction, PremiumAnalysis } from '../api/client';

interface PremiumCorrectionBlockProps {
  corrections: Correction[];
  analysis: PremiumAnalysis;
}

const severityLabels: Record<string, string> = {
  critical: '🔴',
  major: '🟡',
  minor: '🔵',
};

const typeLabels: Record<string, string> = {
  grammar: '📖',
  vocabulary: '📝',
  spelling: '✏️',
  word_order: '🔄',
  punctuation: '🔣',
};

export function PremiumCorrectionBlock({ corrections, analysis }: PremiumCorrectionBlockProps) {
  const { t, i18n } = useTranslation();

  if (!corrections.length) {
    return (
      <div className="card" style={{ textAlign: 'center', color: 'var(--tg-hint)' }}>
        {t('chat.no_corrections')}
      </div>
    );
  }

  const gradeColor = analysis.overall_grade === 'A' ? '#4caf50'
    : analysis.overall_grade === 'B' ? '#8bc34a'
    : analysis.overall_grade === 'C' ? '#ff9800'
    : '#f44336';

  return (
    <div style={{ margin: '8px 16px' }}>
      {/* Overall grade pill */}
      <div style={{
        display: 'inline-block',
        background: gradeColor,
        color: '#fff',
        padding: '4px 12px',
        borderRadius: 12,
        fontSize: 14,
        fontWeight: 600,
        marginBottom: 8,
      }}>
        {t('premium.overall_grade')}: {analysis.overall_grade}
      </div>

      {/* Correction cards with premium details */}
      {corrections.map((c, i) => (
        <div key={i} className="card" style={{ margin: '8px 0', padding: 12 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginBottom: 6 }}>
            <span>{typeLabels[c.type] || '📌'}</span>
            {c.severity && (
              <span title={c.severity}>
                {severityLabels[c.severity] || ''}
              </span>
            )}
            <span style={{
              fontSize: 11,
              background: 'var(--tg-secondary-bg)',
              padding: '2px 6px',
              borderRadius: 4,
              color: 'var(--tg-hint)',
            }}>
              {c.type}
            </span>
            {c.severity && (
              <span style={{
                fontSize: 11,
                background: 'var(--tg-secondary-bg)',
                padding: '2px 6px',
                borderRadius: 4,
                color: 'var(--tg-hint)',
              }}>
                {t(`premium.severity_${c.severity}`, c.severity)}
              </span>
            )}
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

          {c.rule_violated && (
            <div style={{
              fontSize: 12,
              color: 'var(--tg-hint)',
              marginTop: 6,
              padding: '4px 8px',
              background: 'var(--tg-secondary-bg)',
              borderRadius: 6,
              fontStyle: 'italic',
            }}>
              📏 {c.rule_violated}
            </div>
          )}

          {c.learning_tip && (
            <div style={{
              fontSize: 12,
              marginTop: 6,
              padding: '4px 8px',
              background: 'rgba(76, 175, 80, 0.1)',
              borderRadius: 6,
              color: '#4caf50',
            }}>
              💡 {t('premium.learning_tip')}: {c.learning_tip}
            </div>
          )}
        </div>
      ))}

      {/* Strengths */}
      {analysis.strengths.length > 0 && (
        <div className="card" style={{ margin: '8px 0', padding: 12, background: 'rgba(76, 175, 80, 0.05)' }}>
          <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 6, color: '#4caf50' }}>
            ✅ {t('premium.strengths')}
          </div>
          <ul style={{ margin: 0, paddingLeft: 20, fontSize: 13, color: 'var(--tg-text)' }}>
            {analysis.strengths.map((s, i) => <li key={i}>{s}</li>)}
          </ul>
        </div>
      )}

      {/* Areas for improvement */}
      {analysis.areas_for_improvement.length > 0 && (
        <div className="card" style={{ margin: '8px 0', padding: 12, background: 'rgba(255, 152, 0, 0.05)' }}>
          <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 6, color: '#ff9800' }}>
            📈 {t('premium.areas_to_improve')}
          </div>
          <ul style={{ margin: 0, paddingLeft: 20, fontSize: 13, color: 'var(--tg-text)' }}>
            {analysis.areas_for_improvement.map((a, i) => <li key={i}>{a}</li>)}
          </ul>
        </div>
      )}

      {/* Suggested topic */}
      {analysis.suggested_topic && (
        <div className="card" style={{
          margin: '8px 0',
          padding: 12,
          textAlign: 'center',
          background: 'rgba(33, 150, 243, 0.05)',
          border: '1px solid rgba(33, 150, 243, 0.2)',
        }}>
          <div style={{ fontSize: 13, color: 'var(--tg-hint)', marginBottom: 4 }}>
            {t('premium.suggested_topic')}
          </div>
          <div style={{ fontSize: 14, fontWeight: 600, color: '#2196f3' }}>
            🎯 {analysis.suggested_topic}
          </div>
        </div>
      )}
    </div>
  );
}
