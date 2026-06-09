import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation } from '@tanstack/react-query';
import { getIeltsReading, submitIeltsReading } from '../api/client';
import type { IeltsReadingPassage, IeltsReadingResult, IeltsReadingQuestion } from '../api/client';

interface IeltsReadingProps {
  onBack: () => void;
}

export function IeltsReading({ onBack }: IeltsReadingProps) {
  const { t, i18n } = useTranslation();
  const [passage, setPassage] = useState<IeltsReadingPassage | null>(null);
  const [userAnswers, setUserAnswers] = useState<number[]>([]);
  const [result, setResult] = useState<IeltsReadingResult | null>(null);

  const generateMutation = useMutation({
    mutationFn: getIeltsReading,
    onSuccess: (data) => {
      setPassage(data);
      setUserAnswers(new Array(data.questions.length).fill(-1));
      setResult(null);
    },
  });

  const submitMutation = useMutation({
    mutationFn: () => {
      if (!passage) throw new Error('no passage');
      return submitIeltsReading({
        passage: passage.passage,
        questions: passage.questions,
        user_answers: userAnswers,
      });
    },
    onSuccess: (data) => setResult(data),
  });

  const handleAnswer = useCallback((qIndex: number, optionIndex: number) => {
    setUserAnswers(prev => {
      const next = [...prev];
      next[qIndex] = optionIndex;
      return next;
    });
  }, []);

  const allAnswered = userAnswers.every(a => a >= 0);

  const renderQuestion = (q: IeltsReadingQuestion, i: number) => {
    const isTF = q.options?.length === 3 && q.options.some(o => o === 'True' || o === 'False');

    return (
      <div key={i} className="card" style={{ padding: 12, marginBottom: 10, fontSize: 13 }}>
        <div style={{ marginBottom: 8, fontWeight: 500, lineHeight: 1.5 }}>
          <span style={{ fontSize: 11, padding: '2px 6px', borderRadius: 4, background: 'var(--tg-secondary-bg)', color: 'var(--tg-hint)', marginRight: 6 }}>
            Q{i + 1}
          </span>
          {q.question}
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
          {q.options.map((opt, oi) => {
            const isSelected = userAnswers[i] === oi;
            let btnStyle: React.CSSProperties = {
              textAlign: 'left' as const,
              padding: '8px 12px',
              borderRadius: 8,
              border: isSelected ? '2px solid var(--tg-button)' : '1px solid var(--tg-border)',
              background: isSelected ? 'var(--tg-secondary-bg)' : 'var(--tg-bg)',
              color: 'var(--tg-text)',
              fontSize: 13,
              cursor: 'pointer',
            };

            return (
              <button
                key={oi}
                className="btn"
                style={btnStyle}
                onClick={() => handleAnswer(i, oi)}
              >
                {isTF ? opt : `${String.fromCharCode(65 + oi)}. ${opt}`}
              </button>
            );
          })}
        </div>
      </div>
    );
  };

  if (!passage) {
    return (
      <div className="scroll-area">
        <div style={{ padding: 16 }}>
          <button className="btn" onClick={onBack} style={{ marginBottom: 12, fontSize: 13 }}>
            ← {t('common.back')}
          </button>

          <div className="placeholder" style={{ marginTop: 48 }}>
            <div className="placeholder-icon">📖</div>
            <div style={{ marginBottom: 16 }}>{t('ielts.start_reading')}</div>
            <button
              className="btn btn-primary"
              style={{ padding: '12px 24px' }}
              onClick={() => generateMutation.mutate()}
              disabled={generateMutation.isPending}
            >
              {generateMutation.isPending ? t('common.loading') : t('ielts.new_reading')}
            </button>
          </div>

          {generateMutation.error && (
            <div style={{ marginTop: 12, padding: 12, borderRadius: 10, background: '#5e1b2033', color: 'var(--tg-destructive)', fontSize: 13 }}>
              {(generateMutation.error as { status?: number })?.status === 403 ? t('ielts.premium_only') : t('common.error')}
            </div>
          )}
        </div>
      </div>
    );
  }

  return (
    <div className="scroll-area">
      <div style={{ padding: 16 }}>
        <button className="btn" onClick={onBack} style={{ marginBottom: 12, fontSize: 13 }}>
          ← {t('common.back')}
        </button>

        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
          <span style={{ fontSize: 16, fontWeight: 600 }}>{passage.title}</span>
          <span style={{ fontSize: 11, color: 'var(--tg-hint)' }}>{passage.word_count} words</span>
        </div>

        <div className="card" style={{ padding: 16, marginBottom: 16, fontSize: 14, lineHeight: 1.7, maxHeight: 300, overflowY: 'auto' }}>
          {passage.passage}
        </div>

        {passage.questions.map((q, i) => renderQuestion(q, i))}

        <button
          className="btn btn-primary"
          style={{ width: '100%', padding: 14, marginTop: 8 }}
          onClick={() => submitMutation.mutate()}
          disabled={submitMutation.isPending || !allAnswered}
        >
          {submitMutation.isPending ? t('common.loading') : t('ielts.submit')}
        </button>

        {result && (
          <div style={{ marginTop: 16 }}>
            <div className="card" style={{ padding: 16, marginBottom: 12, textAlign: 'center' }}>
              <div style={{ fontSize: 11, color: 'var(--tg-hint)', marginBottom: 4 }}>{t('ielts.band_score')}</div>
              <div style={{ fontSize: 36, fontWeight: 700, color: 'var(--c-primary)' }}>{result.band_score.toFixed(1)}</div>
              <div style={{ fontSize: 13, color: 'var(--tg-hint)', marginTop: 4 }}>
                {t('ielts.correct_answers')}: {result.correct_answers}/{result.total_questions}
              </div>
            </div>

            <div className="card" style={{ padding: 16, marginBottom: 12 }}>
              <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 8 }}>{t('ielts.feedback')}</div>
              <div style={{ fontSize: 13, lineHeight: 1.6 }}>{result.feedback}</div>
            </div>

            {result.results.map((r, i) => (
              <div key={i} className="card" style={{
                padding: 12, marginBottom: 8, fontSize: 13,
                borderLeft: `4px solid ${r.is_correct ? '#4caf50' : 'var(--tg-destructive)'}`,
              }}>
                <div style={{ fontWeight: 500, marginBottom: 4 }}>Q{r.question_index + 1}</div>
                <div style={{ color: r.is_correct ? '#4caf50' : 'var(--tg-destructive)' }}>
                  {r.is_correct ? '✓' : '✗'} {t('ielts.correct_answers')}: {String.fromCharCode(65 + r.correct_answer)}
                </div>
                {!r.is_correct && (
                  <div style={{ color: 'var(--tg-hint)', marginTop: 4, fontSize: 12 }}>
                    {i18n.language === 'uz' ? r.explanation_uz : r.explanation_ru}
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
