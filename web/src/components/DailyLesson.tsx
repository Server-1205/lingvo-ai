import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { initData } from '@telegram-apps/sdk';

interface Exercise {
  question: string;
  answer: string;
  options: string[];
}

interface LessonVocab {
  word: string;
  translation_uz: string;
  translation_ru: string;
}

interface DailyLessonData {
  topic: string;
  explanation_uz: string;
  explanation_ru: string;
  examples: string[];
  exercises: Exercise[];
  vocabulary: LessonVocab[];
}

interface DailyLessonProps {
  onBack: () => void;
}

export function DailyLesson({ onBack }: DailyLessonProps) {
  const { t, i18n } = useTranslation();
  const [lesson, setLesson] = useState<DailyLessonData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedAnswers, setSelectedAnswers] = useState<Record<number, string>>({});
  const [checkedAnswers, setCheckedAnswers] = useState<Record<number, boolean | null>>({});

  const lang = i18n.language === 'ru' ? 'ru' : 'uz';
  const explanation = lang === 'ru' && lesson?.explanation_ru ? lesson.explanation_ru : lesson?.explanation_uz;

  useEffect(() => {
    const fetchLesson = async () => {
      try {
        const raw = initData.raw();
        const headers: Record<string, string> = { 'Content-Type': 'application/json' };
        if (raw) headers['X-Telegram-Init-Data'] = raw;

        const res = await fetch('/api/daily', { method: 'POST', headers, body: '{}' });
        if (!res.ok) throw new Error('failed');

        const data = await res.json() as DailyLessonData;
        setLesson(data);
      } catch {
        setError('error');
      } finally {
        setLoading(false);
      }
    };
    fetchLesson();
  }, []);

  const handleSelect = (exerciseIdx: number, option: string) => {
    setSelectedAnswers(prev => ({ ...prev, [exerciseIdx]: option }));
    setCheckedAnswers(prev => ({ ...prev, [exerciseIdx]: null }));
  };

  const handleCheck = (exerciseIdx: number) => {
    if (!lesson) return;
    const selected = selectedAnswers[exerciseIdx];
    const correct = lesson.exercises[exerciseIdx]?.answer;
    if (!selected) return;
    setCheckedAnswers(prev => ({ ...prev, [exerciseIdx]: selected === correct }));
  };

  if (loading) {
    return (
      <div style={{ display: 'flex', flexDirection: 'column', height: '100%', alignItems: 'center', justifyContent: 'center', gap: 16, padding: 24 }}>
        <div style={{ fontSize: 48 }}>📅</div>
        <div style={{ fontSize: 17, color: 'var(--c-on-surface-variant)' }}>{t('common.loading')}</div>
      </div>
    );
  }

  if (error || !lesson) {
    return (
      <div style={{ display: 'flex', flexDirection: 'column', height: '100%', alignItems: 'center', justifyContent: 'center', gap: 16, padding: 24 }}>
        <div style={{ fontSize: 48 }}>😕</div>
        <div style={{ fontSize: 17, color: 'var(--c-on-surface-variant)', textAlign: 'center' }}>
          {lang === 'ru' ? 'Не удалось загрузить урок. Попробуйте позже.' : 'Darsni yuklashda xatolik. Keyinroq urinib ko\'ring.'}
        </div>
        <button className="btn btn-primary" style={{ marginTop: 12 }} onClick={onBack}>
          {t('common.back')}
        </button>
      </div>
    );
  }

  return (
    <div className="scroll-area">
      <div style={{ padding: '16px', display: 'flex', alignItems: 'center', gap: 8 }}>
        <button
          onClick={onBack}
          style={{ fontSize: 20, padding: '4px 8px', border: 'none', background: 'none', cursor: 'pointer', color: 'var(--c-primary)' }}
        >
          ←
        </button>
        <div style={{ fontSize: 20, fontWeight: 700, flex: 1, color: 'var(--c-on-surface)' }}>{t('daily.title')}</div>
      </div>

      <div style={{
        margin: '0 16px 16px',
        padding: 'var(--card-padding)',
        borderRadius: 'var(--round-lg)',
        border: '1px solid var(--c-outline-variant)',
        background: 'var(--c-surface-container-lowest)',
      }}>
        <div style={{ fontSize: 13, color: 'var(--c-on-surface-variant)', marginBottom: 4 }}>{t('daily.topic')}</div>
        <div style={{ fontSize: 22, fontWeight: 700, color: 'var(--c-primary)', marginBottom: 12 }}>
          {lesson.topic}
        </div>

        <div style={{ fontSize: 12, color: 'var(--c-outline)', marginBottom: 8, textTransform: 'uppercase', letterSpacing: '0.05em', fontWeight: 600 }}>
          {t('daily.explanation')}
        </div>
        <div style={{ fontSize: 16, lineHeight: 1.5, color: 'var(--c-on-surface)', marginBottom: 16 }}>
          {explanation}
        </div>

        {lesson.examples.length > 0 && (
          <>
            <div style={{ fontSize: 12, color: 'var(--c-outline)', marginBottom: 8, textTransform: 'uppercase', letterSpacing: '0.05em', fontWeight: 600 }}>
              {t('daily.examples')}
            </div>
            {lesson.examples.map((ex, i) => (
              <div key={i} style={{
                fontSize: 15,
                fontStyle: 'italic',
                color: 'var(--c-on-surface-variant)',
                padding: '8px 12px',
                marginBottom: 6,
                background: 'var(--c-surface-container)',
                borderRadius: 'var(--round-md)',
              }}>
                • {ex}
              </div>
            ))}
          </>
        )}
      </div>

      {lesson.exercises.length > 0 && (
        <div style={{
          margin: '0 16px 16px',
          padding: 'var(--card-padding)',
          borderRadius: 'var(--round-lg)',
          border: '1px solid var(--c-outline-variant)',
          background: 'var(--c-surface-container-lowest)',
        }}>
          <div style={{ fontSize: 16, fontWeight: 700, marginBottom: 12, color: 'var(--c-on-surface)' }}>{t('daily.exercises')}</div>
          {lesson.exercises.map((ex, i) => {
            const isCorrect = checkedAnswers[i];
            return (
              <div key={i} style={{ marginBottom: 16 }}>
                <div style={{ fontSize: 15, marginBottom: 8, color: 'var(--c-on-surface)' }}>{i + 1}. {ex.question}</div>
                <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap', marginBottom: 6 }}>
                  {ex.options.map(opt => (
                    <button
                      key={opt}
                      onClick={() => handleSelect(i, opt)}
                      style={{
                        padding: '6px 14px',
                        fontSize: 14,
                        border: '2px solid',
                        borderColor: selectedAnswers[i] === opt
                          ? (isCorrect === true ? 'var(--c-tertiary)' : isCorrect === false ? 'var(--c-error)' : 'var(--c-primary)')
                          : 'var(--c-outline-variant)',
                        borderRadius: 'var(--round-full)',
                        background: selectedAnswers[i] === opt
                          ? (isCorrect === true ? 'var(--c-tertiary)' : isCorrect === false ? 'var(--c-error)' : 'var(--c-primary)')
                          : 'var(--c-surface-container)',
                        color: selectedAnswers[i] === opt ? '#fff' : 'var(--c-on-surface)',
                        cursor: 'pointer',
                        fontWeight: selectedAnswers[i] === opt ? 600 : 400,
                      }}
                    >
                      {opt}
                    </button>
                  ))}
                </div>
                {selectedAnswers[i] && isCorrect === null && (
                  <button
                    onClick={() => handleCheck(i)}
                    style={{
                      fontSize: 13,
                      padding: '4px 12px',
                      border: 'none',
                      borderRadius: 'var(--round-full)',
                      background: 'var(--c-primary)',
                      color: '#fff',
                      cursor: 'pointer',
                    }}
                  >
                    {t('daily.check')}
                  </button>
                )}
                {isCorrect === true && (
                  <div style={{ fontSize: 13, color: 'var(--c-tertiary)', fontWeight: 600 }}>✅ {t('daily.correct')}</div>
                )}
                {isCorrect === false && (
                  <div style={{ fontSize: 13, color: 'var(--c-error)', fontWeight: 600 }}>
                    ❌ {t('daily.wrong')} {ex.answer}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      )}

      {lesson.vocabulary.length > 0 && (
        <div style={{
          margin: '0 16px 16px',
          padding: 'var(--card-padding)',
          borderRadius: 'var(--round-lg)',
          border: '1px solid var(--c-outline-variant)',
          background: 'var(--c-surface-container-lowest)',
        }}>
          <div style={{ fontSize: 16, fontWeight: 700, marginBottom: 12, color: 'var(--c-on-surface)' }}>{t('daily.new_words')}</div>
          {lesson.vocabulary.map((v, i) => (
            <div key={i} style={{
              display: 'flex',
              justifyContent: 'space-between',
              padding: '10px 12px',
              marginBottom: 6,
              background: 'var(--c-surface-container)',
              borderRadius: 'var(--round-md)',
            }}>
              <span style={{ fontSize: 16, fontWeight: 600, color: 'var(--c-on-surface)' }}>{v.word}</span>
              <span style={{ fontSize: 15, color: 'var(--c-on-surface-variant)' }}>
                {lang === 'ru' ? v.translation_ru : v.translation_uz}
              </span>
            </div>
          ))}
        </div>
      )}

      <div style={{ padding: '0 16px 24px' }}>
        <button
          className="btn btn-primary"
          style={{ width: '100%', padding: 14, fontSize: 16, fontWeight: 700 }}
          onClick={onBack}
        >
          {t('daily.back_to_chat')}
        </button>
      </div>
    </div>
  );
}
