import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { initData } from '@telegram-apps/sdk';
import { getVocab, addVocab, deleteVocab, getDueWords, submitReview, getSubscription } from '../api/client';
import type { VocabWord } from '../api/client';
import { ReviewCard } from './ReviewCard';
import { useTTS } from '../hooks/useTTS';
import { debug } from '../lib/debug';
import { LoadingDots } from './LoadingDots';
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
                .catch(err => debug('[vocab] export error', err));
            }}
          >
            {t('premium.export_vocab')}
          </button>
        )}
      </div>

      <div style={{ display: 'flex', margin: '0 16px 16px', gap: 8 }}>
        {(['my', 'lookup', 'review'] as const).map((tKey) => (
          <button
            key={tKey}
            className={`btn ${tab === tKey ? 'btn-primary' : 'btn-secondary'}`}
            style={{
              flex: 1,
              padding: '10px 0',
              fontSize: 14,
              fontWeight: 600,
              position: 'relative',
              borderRadius: 'var(--round-full)',
            }}
            onClick={() => { setTab(tKey); if (tKey === 'review') setReviewIndex(0); }}
          >
            {tKey === 'my' ? t('vocab.my_words') : tKey === 'lookup' ? t('vocab.lookup_tab') : t('vocab.review')}
            {tKey === 'review' && dueList.length > 0 && tab !== 'review' && (
              <span style={{
                position: 'absolute', top: -4, right: -4,
                background: 'var(--c-error)', color: '#fff',
                fontSize: 10, padding: '1px 5px', borderRadius: 8, fontWeight: 600,
              }}>
                {dueList.length}
              </span>
            )}
          </button>
        ))}
      </div>

      {tab === 'my' && (
        <>
          {isLoading && (
            <div className="placeholder">
              <LoadingDots size={14} />
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
            <div key={w.id} style={{
              margin: '8px 16px',
              padding: 'var(--card-padding)',
              borderRadius: 'var(--round-lg)',
              border: '1px solid var(--c-outline-variant)',
              background: 'var(--c-surface-container-lowest)',
            }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start' }}>
                <div style={{ flex: 1 }}>
                  <div style={{ fontSize: 18, fontWeight: 700, color: 'var(--c-on-surface)' }}>{w.word}</div>
                  <div style={{ fontSize: 16, color: 'var(--c-on-surface-variant)', marginTop: 4 }}>{tr(w)}</div>
                  {ex(w) && (
                    <div style={{ fontSize: 15, color: 'var(--c-outline)', marginTop: 6, fontStyle: 'italic' }}>
                      {ex(w)}
                    </div>
                  )}
                  <div style={{ display: 'flex', gap: 6, marginTop: 8 }}>
                    {w.level && (
                      <span style={{
                        fontSize: 12,
                        fontWeight: 600,
                        background: 'var(--c-primary-fixed-dim)',
                        color: 'var(--c-on-primary-fixed)',
                        padding: '2px 10px',
                        borderRadius: 'var(--round-full)',
                      }}>
                        {w.level.toUpperCase()}
                      </span>
                    )}
                  </div>
                </div>
                <div style={{ display: 'flex', gap: 6, marginLeft: 12, alignItems: 'center' }}>
                  {!isLoading && (
                    <button
                      className="btn"
                      style={{ fontSize: 16, padding: '6px 10px', lineHeight: 1, borderRadius: 'var(--round-md)' }}
                      onClick={() => playTTS(w.word)}
                      disabled={ttsPlaying}
                    >
                      🔊
                    </button>
                  )}
                  <button
                    className="btn"
                    style={{ fontSize: 12, padding: '6px 10px', color: 'var(--c-error)', borderRadius: 'var(--round-md)' }}
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
              <div style={{ textAlign: 'center', fontSize: 13, color: 'var(--c-on-surface-variant)', padding: '8px 0' }}>
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
            <div style={{
              flex: 1,
              display: 'flex',
              alignItems: 'center',
              padding: '0 16px',
              borderRadius: 'var(--round-full)',
              background: 'var(--c-surface-container)',
              border: '1px solid var(--c-outline-variant)',
            }}>
              <input
                value={lookupWord}
                onChange={e => { setLookupWord(e.target.value); setAddedWord(null); setAddError(null); }}
                placeholder={t('vocab.lookup_placeholder')}
                style={{
                  flex: 1,
                  padding: '12px 0',
                  border: 'none',
                  background: 'transparent',
                  color: 'var(--c-on-surface)',
                  fontSize: 15,
                  outline: 'none',
                }}
              />
            </div>
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
              <LoadingDots size={14} />
              <div>{t('common.loading')}</div>
            </div>
          )}

          {addError && !addMutation.isPending && (
            <div style={{
              padding: 16,
              borderRadius: 'var(--round-lg)',
              border: '1px solid var(--c-error-container)',
              background: 'var(--c-error-container)',
              textAlign: 'center',
            }}>
              <div style={{ fontSize: 14, color: 'var(--c-on-error-container)', fontWeight: 600 }}>
                {t(addError)}
              </div>
            </div>
          )}

          {addedWord && !addMutation.isPending && (
            <div style={{
              padding: 16,
              borderRadius: 'var(--round-lg)',
              border: '1px solid var(--c-tertiary-container)',
              background: 'var(--c-tertiary-container)',
              textAlign: 'center',
            }}>
              <div style={{ fontSize: 14, color: 'var(--c-on-tertiary-container)', fontWeight: 600 }}>
                {t('vocab.saved')}
              </div>
              <div style={{ fontSize: 13, color: 'var(--c-on-surface-variant)', marginTop: 4 }}>
                {addedWord}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
