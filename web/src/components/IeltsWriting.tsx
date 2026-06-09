import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation } from '@tanstack/react-query';
import { getIeltsWriting } from '../api/client';
import type { IeltsWritingResponse } from '../api/client';

interface IeltsWritingProps {
  onBack: () => void;
}

const taskSamples: Record<string, string> = {
  task1: 'The chart below shows the number of international students in universities in Canada, Australia, and the UK from 2010 to 2020.\n\nSummarise the information by selecting and reporting the main features, and make comparisons where relevant.',
  task2: 'Some people believe that unpaid community service should be a compulsory part of high school programmes.\n\nTo what extent do you agree or disagree?',
};

export function IeltsWriting({ onBack }: IeltsWritingProps) {
  const { t, i18n } = useTranslation();
  const [taskType, setTaskType] = useState<'task1' | 'task2'>('task1');
  const [userText, setUserText] = useState('');
  const [taskDescription, setTaskDescription] = useState(taskSamples.task1);
  const [result, setResult] = useState<IeltsWritingResponse | null>(null);

  const mutation = useMutation({
    mutationFn: () => getIeltsWriting({ type: taskType, user_text: userText, task_description: taskDescription }),
    onSuccess: (data) => setResult(data),
  });

  const handleTaskTypeChange = useCallback((type: 'task1' | 'task2') => {
    setTaskType(type);
    setTaskDescription(taskSamples[type]);
    setResult(null);
    setUserText('');
  }, []);

  const renderCriteria = (criteria: Record<string, number>) => {
    return Object.entries(criteria).map(([key, val]) => (
      <div key={key} style={{ display: 'flex', justifyContent: 'space-between', fontSize: 13, padding: '4px 0' }}>
        <span style={{ color: 'var(--tg-hint)' }}>{t(`ielts.${key}`)}</span>
        <span style={{ fontWeight: 600 }}>{val.toFixed(1)}</span>
      </div>
    ));
  };

  return (
    <div className="scroll-area">
      <div style={{ padding: 16 }}>
        <button className="btn" onClick={onBack} style={{ marginBottom: 12, fontSize: 13 }}>
          ← {t('common.back')}
        </button>

        <div style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
          <button
            className={`btn ${taskType === 'task1' ? 'btn-primary' : ''}`}
            style={{ flex: 1, padding: '10px', fontSize: 13 }}
            onClick={() => handleTaskTypeChange('task1')}
          >
            {t('ielts.task1')}
          </button>
          <button
            className={`btn ${taskType === 'task2' ? 'btn-primary' : ''}`}
            style={{ flex: 1, padding: '10px', fontSize: 13 }}
            onClick={() => handleTaskTypeChange('task2')}
          >
            {t('ielts.task2')}
          </button>
        </div>

        <div className="card" style={{ padding: 12, marginBottom: 12, fontSize: 13, lineHeight: 1.5, color: 'var(--tg-hint)' }}>
          <strong>{t('ielts.task')}:</strong>
          <div style={{ marginTop: 4, whiteSpace: 'pre-wrap' }}>{taskDescription}</div>
        </div>

        <textarea
          className="input"
          style={{ width: '100%', minHeight: 160, marginBottom: 12, padding: 12, fontSize: 14, lineHeight: 1.5, resize: 'vertical' }}
          placeholder={t('ielts.your_response')}
          value={userText}
          onChange={(e) => setUserText(e.target.value)}
        />

        <button
          className="btn btn-primary"
          style={{ width: '100%', padding: 14 }}
          onClick={() => mutation.mutate()}
          disabled={mutation.isPending || !userText.trim()}
        >
          {mutation.isPending ? t('common.loading') : t('ielts.submit')}
        </button>

        {mutation.error && (
          <div style={{ marginTop: 12, padding: 12, borderRadius: 10, background: '#5e1b2033', color: 'var(--tg-destructive)', fontSize: 13 }}>
            {(mutation.error as { status?: number })?.status === 403 ? t('ielts.premium_only') : t('common.error')}
          </div>
        )}

        {result && (
          <div style={{ marginTop: 16 }}>
            <div className="card" style={{ padding: 16, marginBottom: 12, textAlign: 'center' }}>
              <div style={{ fontSize: 11, color: 'var(--tg-hint)', marginBottom: 4 }}>{t('ielts.band_score')}</div>
              <div style={{ fontSize: 36, fontWeight: 700, color: 'var(--c-primary)' }}>{result.band_score.toFixed(1)}</div>
            </div>

            <div className="card" style={{ padding: 16, marginBottom: 12 }}>
              <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 8 }}>{t('ielts.criteria')}</div>
              {renderCriteria(result.criteria as unknown as Record<string, number>)}
            </div>

            <div className="card" style={{ padding: 16, marginBottom: 12 }}>
              <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 8 }}>{t('ielts.feedback')}</div>
              <div style={{ fontSize: 13, lineHeight: 1.6, whiteSpace: 'pre-wrap' }}>{result.feedback}</div>
            </div>

            {result.improvement_tips.length > 0 && (
              <div className="card" style={{ padding: 16, marginBottom: 12 }}>
                <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 8 }}>{t('ielts.improvement_tips')}</div>
                {result.improvement_tips.map((tip, i) => (
                  <div key={i} style={{ fontSize: 13, padding: '4px 0', display: 'flex', gap: 6 }}>
                    <span>💡</span>
                    <span>{tip}</span>
                  </div>
                ))}
              </div>
            )}

            {result.corrections.length > 0 && (
              <div className="card" style={{ padding: 16, marginBottom: 12 }}>
                <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 8 }}>{t('chat.corrections')}</div>
                {result.corrections.map((c, i) => (
                  <div key={i} style={{ fontSize: 13, padding: '8px 0', borderBottom: '1px solid var(--tg-border)' }}>
                    <div><span style={{ color: 'var(--tg-destructive)', textDecoration: 'line-through' }}>{c.original}</span></div>
                    <div><span style={{ color: '#4caf50' }}>{c.corrected}</span></div>
                    <div style={{ color: 'var(--tg-hint)', marginTop: 2, fontSize: 12 }}>
                      {i18n.language === 'uz' ? c.explanation_uz : c.explanation_ru}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
