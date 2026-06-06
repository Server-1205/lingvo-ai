import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQuery } from '@tanstack/react-query';
import { saveLevel, getLevelTestQuestions } from '../api/client';

interface LevelTestProps {
  onDone: () => void;
}

function determineLevel(correct: number): string {
  if (correct <= 3) return 'a1';
  if (correct <= 5) return 'a2';
  if (correct <= 7) return 'b1';
  if (correct <= 9) return 'b2';
  return 'c1';
}

export function LevelTest({ onDone }: LevelTestProps) {
  const { t, i18n } = useTranslation();
  const [step, setStep] = useState<'test' | 'result'>('test');
  const [currentIdx, setCurrentIdx] = useState(0);
  const [answers, setAnswers] = useState<number[]>([]);
  const [selectedOption, setSelectedOption] = useState<number | null>(null);

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['level-test'],
    queryFn: getLevelTestQuestions,
    retry: 1,
  });

  const saveMutation = useMutation({
    mutationFn: (level: string) => saveLevel({ level }),
    onSuccess: onDone,
  });

  const questions = data?.questions ?? [];

  const handleAnswer = useCallback(() => {
    if (selectedOption === null) return;
    setAnswers(prev => [...prev, selectedOption]);
    setSelectedOption(null);

    if (currentIdx + 1 >= questions.length) {
      setStep('result');
    } else {
      setCurrentIdx(prev => prev + 1);
    }
  }, [selectedOption, currentIdx, questions.length]);

  if (isLoading) {
    return (
      <div className="scroll-area" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '60vh' }}>
        <div className="placeholder">
          <div className="placeholder-icon">⏳</div>
          <div>{t('common.loading')}</div>
        </div>
      </div>
    );
  }

  if (error || questions.length === 0) {
    return (
      <div className="scroll-area" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '60vh', flexDirection: 'column', gap: 12 }}>
        <div className="placeholder">
          <div className="placeholder-icon">❌</div>
          <div>{t('common.error')}</div>
        </div>
        <button className="btn btn-primary" onClick={() => refetch()}>
          {t('common.retry')}
        </button>
      </div>
    );
  }

  if (step === 'result') {
    const correctCount = questions.filter((q, i) => q.correct === answers[i]).length;
    const level = determineLevel(correctCount);
    const total = questions.length;

    const levelLabels: Record<string, string> = {
      a1: 'A1 — Beginner',
      a2: 'A2 — Elementary',
      b1: 'B1 — Intermediate',
      b2: 'B2 — Upper Intermediate',
      c1: 'C1 — Advanced',
    };

    return (
      <div className="scroll-area">
        <div style={{ padding: 24, textAlign: 'center' }}>
          <div style={{ fontSize: 48, marginBottom: 16 }}>🎉</div>
          <div style={{ fontSize: 20, fontWeight: 700, marginBottom: 8 }}>
            {t('level.result_title')}
          </div>
          <div className="card" style={{ padding: 20, margin: '16px 0' }}>
            <div style={{ fontSize: 14, color: 'var(--tg-hint)', marginBottom: 4 }}>
              {t('progress.level')}
            </div>
            <div style={{ fontSize: 36, fontWeight: 700, color: 'var(--tg-button)', marginBottom: 8 }}>
              {level.toUpperCase()}
            </div>
            <div style={{ fontSize: 14, color: 'var(--tg-hint)' }}>
              {levelLabels[level]}
            </div>
            <div style={{ fontSize: 13, color: 'var(--tg-hint)', marginTop: 12 }}>
              {t('level.score', { correct: correctCount, total })}
            </div>
          </div>

          <button
            className="btn btn-primary"
            style={{ width: '100%', padding: 14 }}
            onClick={() => saveMutation.mutate(level)}
            disabled={saveMutation.isPending}
          >
            {saveMutation.isPending ? t('common.loading') : t('level.save_result')}
          </button>
        </div>
      </div>
    );
  }

  const q = questions[currentIdx];
  const isCorrect = answers.length > 0 && answers[currentIdx - 1] === questions[currentIdx - 1]?.correct;
  const showFeedback = currentIdx > 0 && answers.length > 0;

  return (
    <div className="scroll-area">
      <div style={{ padding: 16 }}>
        <div style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 16,
        }}>
          <span style={{ fontSize: 13, color: 'var(--tg-hint)' }}>
            {t('level.question_of', { current: currentIdx + 1, total: questions.length })}
          </span>
          <span style={{
            fontSize: 11,
            background: 'var(--tg-secondary-bg)',
            padding: '2px 8px',
            borderRadius: 4,
            color: 'var(--tg-hint)',
          }}>
            {q.level?.toUpperCase()}
          </span>
        </div>

        <div className="card" style={{ padding: 16, marginBottom: 16 }}>
          <div style={{ fontSize: 16, fontWeight: 500, marginBottom: 16, lineHeight: 1.5 }}>
            {q.question}
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
            {q.options.map((opt, i) => (
              <button
                key={i}
                className="btn"
                style={{
                  textAlign: 'left',
                  padding: '12px 14px',
                  borderRadius: 10,
                  border: selectedOption === i ? '2px solid var(--tg-button)' : '1px solid var(--tg-border)',
                  background: selectedOption === i ? 'var(--tg-secondary-bg)' : 'var(--tg-bg)',
                  color: 'var(--tg-text)',
                  fontSize: 14,
                  cursor: 'pointer',
                }}
                onClick={() => setSelectedOption(i)}
              >
                <span style={{ fontWeight: 600, marginRight: 8 }}>
                  {String.fromCharCode(65 + i)}.
                </span>
                {opt}
              </button>
            ))}
          </div>
        </div>

        {showFeedback && (
          <div style={{
            padding: 12,
            borderRadius: 10,
            marginBottom: 16,
            background: isCorrect ? '#1b5e2033' : '#5e1b2033',
            fontSize: 13,
            color: isCorrect ? '#4caf50' : 'var(--tg-destructive)',
          }}>
            {questions[currentIdx - 1].correct === answers[currentIdx - 1]
              ? '✅ ' + (i18n.language === 'uz' ? questions[currentIdx - 1].explanation_uz : questions[currentIdx - 1].explanation_ru)
              : '❌ ' + (i18n.language === 'uz' ? questions[currentIdx - 1].explanation_uz : questions[currentIdx - 1].explanation_ru)}
          </div>
        )}

        <button
          className="btn btn-primary"
          style={{ width: '100%', padding: 14 }}
          disabled={selectedOption === null}
          onClick={handleAnswer}
        >
          {currentIdx + 1 >= questions.length ? t('level.finish') : t('level.next')}
        </button>
      </div>
    </div>
  );
}
