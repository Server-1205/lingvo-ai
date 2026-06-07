import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getVocab, addVocab, deleteVocab, getDueWords, submitReview } from '../api/client';
import type { VocabWord } from '../api/client';
import { ReviewCard } from './ReviewCard';

interface VocabularyProps {
  initialTab?: 'my' | 'lookup' | 'review';
}

export function Vocabulary({ initialTab }: VocabularyProps) {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const [tab, setTab] = useState<'my' | 'lookup' | 'review'>(initialTab || 'my');
  const [reviewIndex, setReviewIndex] = useState(0);
  const [lookupWord, setLookupWord] = useState('');
  const [addedWord, setAddedWord] = useState<string | null>(null);

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
    mutationFn: () => addVocab({ word: lookupWord }),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['vocab'] });
      setAddedWord(`${lookupWord} — ${data.translation_uz}`);
      setLookupWord('');
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

  const dueList = dueWords ?? [];

  return (
    <div className="scroll-area">
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: '12px 16px',
      }}>
        <span className="section-title" style={{ padding: 0 }}>{t('vocab.title')}</span>
      </div>

      <div style={{ display: 'flex', margin: '0 16px 12px', gap: 8 }}>
        <button
          className={`btn ${tab === 'my' ? 'btn-primary' : ''}`}
          style={{ flex: 1, padding: '8px 0', fontSize: 13 }}
          onClick={() => setTab('my')}
        >
          {t('vocab.my_words')}
        </button>
        <button
          className={`btn ${tab === 'lookup' ? 'btn-primary' : ''}`}
          style={{ flex: 1, padding: '8px 0', fontSize: 13 }}
          onClick={() => setTab('lookup')}
        >
          {t('vocab.lookup_tab')}
        </button>
        <button
          className={`btn ${tab === 'review' ? 'btn-primary' : ''}`}
          style={{ flex: 1, padding: '8px 0', fontSize: 13, position: 'relative' }}
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
            <div key={w.id} className="card" style={{ margin: '8px 16px', padding: 12 }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start' }}>
                <div>
                  <div style={{ fontSize: 16, fontWeight: 600 }}>{w.word}</div>
                  <div style={{ fontSize: 14, color: 'var(--tg-hint)', marginTop: 2 }}>{w.translation}</div>
                  <div style={{ fontSize: 13, color: 'var(--tg-hint)', marginTop: 4, fontStyle: 'italic' }}>
                    {w.example}
                  </div>
                  <div style={{ display: 'flex', gap: 6, marginTop: 6 }}>
                    <span style={{ fontSize: 11, background: 'var(--tg-secondary-bg)', padding: '2px 6px', borderRadius: 4 }}>
                      {w.level?.toUpperCase()}
                    </span>
                  </div>
                </div>
                <button
                  className="btn"
                  style={{ fontSize: 12, padding: '4px 8px', color: 'var(--tg-destructive)' }}
                  onClick={() => deleteMutation.mutate(w.id)}
                >
                  {t('vocab.delete')}
                </button>
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
              onChange={e => { setLookupWord(e.target.value); setAddedWord(null); }}
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
