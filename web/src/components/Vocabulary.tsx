import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { initData } from '@telegram-apps/sdk';
import { getVocab, addVocab, deleteVocab, getDueWords, submitReview, getSubscription } from '../api/client';
import type { VocabWord } from '../api/client';
import { ReviewCard } from './ReviewCard';
import { useTTS } from '../hooks/useTTS';
import { ApiError } from '../api/client';

interface VocabularyProps {
  initialTab?: 'my' | 'lookup' | 'review';
}

export function Vocabulary({ initialTab }: VocabularyProps) {
  const { t, i18n } = useTranslation();
  const queryClient = useQueryClient();
  const [tab, setTab] = useState<'my' | 'lookup' | 'review'>(initialTab || 'my');
  const [reviewIndex, setReviewIndex] = useState(0);
  const [lookupWord, setLookupWord] = useState('');
  const [addedWord, setAddedWord] = useState<string | null>(null);
  const [addError, setAddError] = useState<string | null>(null);

  const { data: sub } = useQuery({
    queryKey: ['subscription'],
    queryFn: getSubscription,
    staleTime: 30000,
  });
  const isPremium = sub?.active ?? false;

  const { data: vocabResp, isLoading } = useQuery({
    queryKey: ['vocab'],
    queryFn: getVocab,
  });
  const words = vocabResp?.words ?? [];

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteVocab(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['vocab'] }),
  });

  const addMutation = useMutation({
    mutationFn: () => addVocab({ word: lookupWord, lang: i18n.language }),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['vocab'] });
      const showTranslation = lang === 'ru' && data.translation_ru ? data.translation_ru : data.translation_uz;
      setAddedWord(`${data.word_en || lookupWord} — ${showTranslation}`);
      setLookupWord('');
      setAddError(null);
    },
    onError: (err) => {
      if (err instanceof ApiError && err.status === 400) {
        const data = err.data as Record<string, string>;
        if (data.error === 'inappropriate_word') {
          setAddError('vocab.inappropriate');
        } else if (data.error === 'invalid_language') {
          setAddError('vocab.invalid_language');
        } else {
          setAddError('vocab.add_error');
        }
      } else {
        setAddError('vocab.add_error');
      }
    },
  });

  const { data: dueWords } = useQuery({
    queryKey: ['vocab-review'],
    queryFn: () => getDueWords(20),
    staleTime: 15000,
  });

  const reviewMutation = useMutation({
    mutationFn: ({ wordId, quality }: { wordId: number; quality: number }) => submitReview(wordId, quality),
    onSuccess: () => {
      setReviewIndex((i) => i + 1);
    },
  });

  const { play: playTTS, isPlaying: ttsPlaying } = useTTS();
  const dueList = dueWords ?? [];

  const lang = i18n.language === 'ru' ? 'ru' : 'uz';
  const tr = (w: VocabWord) => lang === 'ru' && w.translation_ru ? w.translation_ru : w.translation;
  const ex = (w: VocabWord) => lang === 'ru' && w.example_ru ? w.example_ru : w.example;

  return (
    <div className="scroll-area">
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: '12px 16px',
      }}>
        <span className="section-title" style={{ padding: 0 }}>{t('vocab.title')}</span>
        {isPremium && (
          <button
            className="btn"
            style={{ fontSize: 12, padding: '4px 10px' }}
            onClick={() => {
              console.debug('[vocab] export clicked');
              const a = document.createElement('a');
              const headers: Record<string, string> = {};
              const raw = initData.raw();
              if (raw) headers['X-Telegram-Init-Data'] = raw;
              fetch('/api/vocab/export', { headers })
                .then(res => {
                  if (!res.ok) throw new Error('export failed');
                  return res.blob();
                })
                .then(blob => {
                  a.href = URL.createObjectURL(blob);
                  a.download = 'vocabulary.csv';
                  a.click();
                  URL.revokeObjectURL(a.href);
                })
                .catch(err => console.debug('[vocab] export error', err));
            }}
          >
            {t('premium.export_vocab')}
          </button>
        )}
      </div>

      <div style={{ display: 'flex', margin: '0 16px 16px', gap: 8 }}>
        <button
          className={`btn ${tab === 'my' ? 'btn-primary' : 'btn-secondary'}`}
          style={{ flex: 1, padding: '10px 0', fontSize: 14, fontWeight: 600 }}
          onClick={() => setTab('my')}
        >
          {t('vocab.my_words')}
        </button>
        <button
          className={`btn ${tab === 'lookup' ? 'btn-primary' : 'btn-secondary'}`}
          style={{ flex: 1, padding: '10px 0', fontSize: 14, fontWeight: 600 }}
          onClick={() => setTab('lookup')}
        >
          {t('vocab.lookup_tab')}
        </button>
        <button
          className={`btn ${tab === 'review' ? 'btn-primary' : 'btn-secondary'}`}
          style={{ flex: 1, padding: '10px 0', fontSize: 14, fontWeight: 600, position: 'relative' }}
          onClick={() => { setTab('review'); setReviewIndex(0); }}
        >
          {t('vocab.review')}
          {dueList.length > 0 && tab !== 'review' && (
            <span style={{
              position: 'absolute', top: -4, right: -4,
              background: 'var(--tg-destructive)', color: '#fff',
              fontSize: 10, padding: '1px 5px', borderRadius: 8, fontWeight: 600,
            }}>
              {dueList.length}
            </span>
          )}
        </button>
      </div>

      {tab === 'my' && (
        <>
          {isLoading && (
            <div className="placeholder">
              <div className="placeholder-icon">⏳</div>
              <div>{t('common.loading')}</div>
            </div>
          )}

          {!isLoading && words.length === 0 && (
            <div className="placeholder">
              <div className="placeholder-icon">📚</div>
              <div>{t('vocab.empty')}</div>
            </div>
          )}

          {words.map((w: VocabWord) => (
            <div key={w.id} className="card" style={{ margin: '10px 16px', padding: 16 }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start' }}>
                <div style={{ flex: 1 }}>
                  <div style={{ fontSize: 18, fontWeight: 700, color: 'var(--tg-text)' }}>{w.word}</div>
                  <div style={{ fontSize: 16, color: 'var(--c-text-secondary)', marginTop: 4 }}>{tr(w)}</div>
                  <div style={{ fontSize: 15, color: 'var(--c-text-hint)', marginTop: 6, fontStyle: 'italic' }}>
                    {ex(w)}
                  </div>
                  <div style={{ display: 'flex', gap: 6, marginTop: 8 }}>
                    {w.level && (
                      <span style={{
                        fontSize: 12,
                        background: 'var(--c-primary-light)',
                        color: 'var(--c-primary-dark)',
                        padding: '2px 8px',
                        borderRadius: 6,
                        fontWeight: 600,
                      }}>
                        {w.level.toUpperCase()}
                      </span>
                    )}
                  </div>
                </div>
                <div style={{ display: 'flex', gap: 6, marginLeft: 12 }}>
                  {!isLoading && (
                    <button
                      className="btn"
                      style={{ fontSize: 16, padding: '6px 10px', lineHeight: 1, borderRadius: 10 }}
                      onClick={() => playTTS(w.word)}
                      disabled={ttsPlaying}
                    >
                      🔊
                    </button>
                  )}
                  <button
                    className="btn"
                    style={{ fontSize: 12, padding: '6px 10px', color: 'var(--tg-destructive)', borderRadius: 10 }}
                    onClick={() => deleteMutation.mutate(w.id)}
                  >
                    {t('vocab.delete')}
                  </button>
                </div>
              </div>
            </div>
          ))}
        </>
      )}

      {tab === 'review' && (
        <>
          {dueList.length === 0 && (
            <div className="placeholder">
              <div className="placeholder-icon">🃏</div>
              <div>{t('vocab.review_empty')}</div>
            </div>
          )}

          {dueList.length > 0 && reviewIndex >= dueList.length && (
            <div className="placeholder">
              <div className="placeholder-icon">🎉</div>
              <div>{t('vocab.review_done')}</div>
              <button
                className="btn btn-primary"
                style={{ marginTop: 12 }}
                onClick={() => {
                  setReviewIndex(0);
                  queryClient.invalidateQueries({ queryKey: ['vocab-review'] });
                  queryClient.invalidateQueries({ queryKey: ['vocab'] });
                }}
              >
                {t('common.retry')}
              </button>
            </div>
          )}

          {dueList.length > 0 && reviewIndex < dueList.length && (
            <div>
              <div style={{ textAlign: 'center', fontSize: 13, color: 'var(--tg-hint)', padding: '8px 0' }}>
                {t('vocab.review_progress', { current: reviewIndex + 1, total: dueList.length })}
              </div>
              <ReviewCard
                word={dueList[reviewIndex]}
                onRate={(quality) => {
                  reviewMutation.mutate({ wordId: dueList[reviewIndex].id, quality });
                }}
              />
            </div>
          )}
        </>
      )}

      {tab === 'lookup' && (
        <div style={{ padding: '0 16px' }}>
          <div style={{ display: 'flex', gap: 8, marginBottom: 12 }}>
            <input
              value={lookupWord}
              onChange={e => { setLookupWord(e.target.value); setAddedWord(null); setAddError(null); }}
              placeholder={t('vocab.lookup_placeholder')}
              style={{
                flex: 1,
                padding: '10px 14px',
                border: '1px solid var(--tg-border)',
                borderRadius: 10,
                background: 'var(--tg-secondary-bg)',
                color: 'var(--tg-text)',
                fontSize: 15,
                outline: 'none',
              }}
            />
            <button
              className="btn btn-primary"
              style={{ padding: '10px 16px' }}
              onClick={() => addMutation.mutate()}
              disabled={!lookupWord.trim() || addMutation.isPending}
            >
              {t('vocab.save_word')}
            </button>
          </div>

          {addMutation.isPending && (
            <div className="placeholder">
              <div className="placeholder-icon">⏳</div>
              <div>{t('common.loading')}</div>
            </div>
          )}

          {addError && !addMutation.isPending && (
            <div className="card" style={{ padding: 16, textAlign: 'center' }}>
              <div style={{ fontSize: 24, marginBottom: 8 }}>🚫</div>
              <div style={{ fontSize: 14, color: '#e53935', fontWeight: 600 }}>
                {t(addError)}
              </div>
            </div>
          )}

          {addedWord && !addMutation.isPending && (
            <div className="card" style={{ padding: 16, textAlign: 'center' }}>
              <div style={{ fontSize: 24, marginBottom: 8 }}>✅</div>
              <div style={{ fontSize: 14, color: '#4caf50', fontWeight: 600 }}>
                {t('vocab.saved')}
              </div>
              <div style={{ fontSize: 13, color: 'var(--tg-hint)', marginTop: 4 }}>
                {addedWord}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
