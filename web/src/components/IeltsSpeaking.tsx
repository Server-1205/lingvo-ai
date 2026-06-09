import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQuery } from '@tanstack/react-query';
import { getIeltsSpeakingQuestions, submitIeltsSpeaking } from '../api/client';
import type { IeltsSpeakingResponse } from '../api/client';

interface IeltsSpeakingProps {
  onBack: () => void;
}

export function IeltsSpeaking({ onBack }: IeltsSpeakingProps) {
  const { t } = useTranslation();
  const [part, setPart] = useState(1);
  const [currentQ, setCurrentQ] = useState(0);
  const [userResponse, setUserResponse] = useState('');
  const [result, setResult] = useState<IeltsSpeakingResponse | null>(null);

  const questionsQuery = useQuery({
    queryKey: ['ielts-speaking', part],
    queryFn: () => getIeltsSpeakingQuestions(part),
    retry: 1,
  });

  const evaluateMutation = useMutation({
    mutationFn: (q: string) => submitIeltsSpeaking({ part, question: q, user_response: userResponse }),
    onSuccess: (data) => setResult(data),
  });

  const handlePartChange = useCallback((newPart: number) => {
    setPart(newPart);
    setCurrentQ(0);
    setUserResponse('');
    setResult(null);
  }, []);

  const handleNext = useCallback(() => {
    const questions = questionsQuery.data?.questions || [];
    if (currentQ + 1 < questions.length) {
      setCurrentQ(prev => prev + 1);
      setUserResponse('');
      setResult(null);
    }
  }, [currentQ, questionsQuery.data]);

  const question = (questionsQuery.data?.questions || [])[currentQ];
  const isLast = currentQ >= (questionsQuery.data?.questions?.length || 1) - 1;

  if (questionsQuery.isLoading) {
    return (
      <div className="scroll-area" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '60vh' }}>
        <div className="placeholder">
          <div className="placeholder-icon">⏳</div>
          <div>{t('common.loading')}</div>
        </div>
      </div>
    );
  }

  if (questionsQuery.error) {
    const is403 = (questionsQuery.error as { status?: number })?.status === 403;
    return (
      <div className="scroll-area" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '60vh' }}>
        <div className="placeholder">
          <div className="placeholder-icon">❌</div>
          <div>{is403 ? t('ielts.premium_only') : t('common.error')}</div>
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

        <div style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
          {[1, 2, 3].map(p => (
            <button
              key={p}
              className={`btn ${part === p ? 'btn-primary' : ''}`}
              style={{ flex: 1, padding: '10px', fontSize: 12 }}
              onClick={() => handlePartChange(p)}
            >
              {t(`ielts.part${p}`)}
            </button>
          ))}
        </div>

        {questionsQuery.data?.cue_card && (
          <div className="card" style={{ padding: 12, marginBottom: 12, fontSize: 13, lineHeight: 1.5, background: 'var(--tg-secondary-bg)' }}>
            <div style={{ fontWeight: 600, marginBottom: 4 }}>🎴 {t('ielts.part2')}</div>
            <div style={{ whiteSpace: 'pre-wrap' }}>{questionsQuery.data.cue_card}</div>
          </div>
        )}

        {question && (
          <div className="card" style={{ padding: 16, marginBottom: 12 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
              <span style={{ fontSize: 12, color: 'var(--tg-hint)' }}>
                {t('level.question_of', { current: currentQ + 1, total: questionsQuery.data?.questions?.length || 1 })}
              </span>
              <span style={{ fontSize: 11, padding: '2px 8px', borderRadius: 4, background: 'var(--tg-secondary-bg)', color: 'var(--tg-hint)' }}>
                IELTS Speaking Part {part}
              </span>
            </div>
            <div style={{ fontSize: 15, fontWeight: 500, marginBottom: 12, lineHeight: 1.5 }}>{question}</div>

            <textarea
              className="input"
              style={{ width: '100%', minHeight: 120, marginBottom: 12, padding: 12, fontSize: 14, lineHeight: 1.5, resize: 'vertical' }}
              placeholder={t('ielts.your_response')}
              value={userResponse}
              onChange={(e) => setUserResponse(e.target.value)}
            />

            <button
              className="btn btn-primary"
              style={{ width: '100%', padding: 14, marginBottom: 8 }}
              onClick={() => evaluateMutation.mutate(question)}
              disabled={evaluateMutation.isPending || !userResponse.trim()}
            >
              {evaluateMutation.isPending ? t('common.loading') : t('ielts.submit')}
            </button>
          </div>
        )}

        {evaluateMutation.error && (
          <div style={{ marginBottom: 12, padding: 12, borderRadius: 10, background: '#5e1b2033', color: 'var(--tg-destructive)', fontSize: 13 }}>
            {t('common.error')}
          </div>
        )}

        {result && (
          <>
            <div className="card" style={{ padding: 16, marginBottom: 12, textAlign: 'center' }}>
              <div style={{ fontSize: 11, color: 'var(--tg-hint)', marginBottom: 4 }}>{t('ielts.band_score')}</div>
              <div style={{ fontSize: 36, fontWeight: 700, color: 'var(--c-primary)' }}>{result.band_score.toFixed(1)}</div>
            </div>

            <div className="card" style={{ padding: 16, marginBottom: 12 }}>
              <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 8 }}>{t('ielts.criteria')}</div>
              {Object.entries(result.criteria).map(([key, val]) => (
                <div key={key} style={{ display: 'flex', justifyContent: 'space-between', fontSize: 13, padding: '4px 0' }}>
                  <span style={{ color: 'var(--tg-hint)' }}>{t(`ielts.${key}`)}</span>
                  <span style={{ fontWeight: 600 }}>{(val as number).toFixed(1)}</span>
                </div>
              ))}
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

            {!isLast && (
              <button className="btn" style={{ width: '100%', padding: 14 }} onClick={handleNext}>
                {t('ielts.new_question')}
              </button>
            )}
          </>
        )}
      </div>
    </div>
  );
}
