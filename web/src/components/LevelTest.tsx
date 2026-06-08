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
  const { t } = useTranslation();
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
          <div style={{ fontSize: 20, fontWeight: 700, marginBottom: 8, color: 'var(--c-on-surface)' }}>
            {t('level.result_title')}
          </div>
          <div style={{
            padding: 20,
            margin: '16px 0',
            borderRadius: 'var(--round-lg)',
            border: '1px solid var(--c-outline-variant)',
            background: 'var(--c-surface-container-lowest)',
          }}>
            <div style={{ fontSize: 14, color: 'var(--c-on-surface-variant)', marginBottom: 4 }}>
              {t('progress.level')}
            </div>
            <div style={{ fontSize: 36, fontWeight: 700, color: 'var(--c-primary)', marginBottom: 8 }}>
              {level.toUpperCase()}
            </div>
            <div style={{ fontSize: 14, color: 'var(--c-on-surface-variant)' }}>
              {levelLabels[level]}
            </div>
            <div style={{ fontSize: 13, color: 'var(--c-on-surface-variant)', marginTop: 12 }}>
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

  return (
    <div className="scroll-area">
      <div style={{ padding: 16 }}>
        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: 12,
          marginBottom: 16,
        }}>
          <div style={{
            flex: 1,
            height: 6,
            borderRadius: 'var(--round-full)',
            background: 'var(--c-surface-container-highest)',
            overflow: 'hidden',
          }}>
            <div style={{
              width: `${((currentIdx + 1) / questions.length) * 100}%`,
              height: '100%',
              borderRadius: 'var(--round-full)',
              background: 'var(--c-primary)',
              transition: 'width 0.3s',
            }} />
          </div>
          <span style={{ fontSize: 13, color: 'var(--c-on-surface-variant)', whiteSpace: 'nowrap' }}>
            {currentIdx + 1}/{questions.length}
          </span>
        </div>

        <div style={{
          padding: 16,
          marginBottom: 16,
          borderRadius: 'var(--round-lg)',
          border: '1px solid var(--c-outline-variant)',
          background: 'var(--c-surface-container-lowest)',
        }}>
          <div style={{ fontSize: 20, fontWeight: 600, marginBottom: 16, lineHeight: 1.5, color: 'var(--c-on-surface)' }}>
            {q.question}
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
            {q.options.map((opt, i) => (
              <button
                key={i}
                style={{
                  textAlign: 'left',
                  padding: '12px 14px',
                  borderRadius: 'var(--round-md)',
                  border: selectedOption === i ? '2px solid var(--c-primary)' : '1px solid var(--c-outline-variant)',
                  background: selectedOption === i ? 'var(--c-primary-fixed)' : 'var(--c-surface-container)',
                  color: selectedOption === i ? 'var(--c-on-primary-fixed)' : 'var(--c-on-surface)',
                  fontSize: 15,
                  cursor: 'pointer',
                  fontWeight: selectedOption === i ? 600 : 400,
                  transition: 'all 0.15s',
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
