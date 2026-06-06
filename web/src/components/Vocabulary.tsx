import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getVocab, addVocab, deleteVocab, lookupVocab, getDueWords, submitReview } from '../api/client';
import type { VocabWord, VocabLookupResponse } from '../api/client';
import { ReviewCard } from './ReviewCard';

export function Vocabulary() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const [tab, setTab] = useState<'my' | 'lookup' | 'review'>('my');
  const [reviewIndex, setReviewIndex] = useState(0);
  const [lookupWord, setLookupWord] = useState('');
  const [showAdd, setShowAdd] = useState<VocabLookupResponse | null>(null);
  const [addWord, setAddWord] = useState('');
  const [addTranslation, setAddTranslation] = useState('');
  const [addExample, setAddExample] = useState('');
  const [addLevel, setAddLevel] = useState('a1');

  const { data: words, isLoading } = useQuery({
    queryKey: ['vocab'],
    queryFn: getVocab,
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteVocab(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['vocab'] }),
  });

  const addMutation = useMutation({
    mutationFn: () => addVocab({ word: addWord, translation: addTranslation, example: addExample, level: addLevel }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['vocab'] });
      setShowAdd(null);
      setAddWord('');
      setAddTranslation('');
      setAddExample('');
      setAddLevel('a1');
    },
  });

  const lookupMutation = useMutation({
    mutationFn: () => lookupVocab({ word: lookupWord }),
    onSuccess: (data) => {
      setShowAdd(data);
      setAddWord(lookupWord);
      setAddTranslation(data.translation_uz);
      setAddExample(data.examples?.[0] || '');
      setAddLevel(data.level || 'a1');
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

  const list = words ?? [];
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

          {!isLoading && list.length === 0 && (
            <div className="placeholder">
              <div className="placeholder-icon">📚</div>
              <div>{t('vocab.empty')}</div>
            </div>
          )}

          {list.map((w: VocabWord) => (
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
              onChange={e => setLookupWord(e.target.value)}
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
              onClick={() => lookupMutation.mutate()}
              disabled={!lookupWord.trim() || lookupMutation.isPending}
            >
              {t('vocab.lookup_btn')}
            </button>
          </div>

          {lookupMutation.isPending && (
            <div className="placeholder">
              <div className="placeholder-icon">⏳</div>
              <div>{t('common.loading')}</div>
            </div>
          )}

          {showAdd && (
            <div className="card" style={{ padding: 16 }}>
              <div style={{ fontSize: 15, fontWeight: 600, marginBottom: 12 }}>
                {lookupWord}
              </div>
              <div style={{ marginBottom: 8 }}>
                <div style={{ fontSize: 12, color: 'var(--tg-hint)', marginBottom: 2 }}>{t('vocab.translation')}</div>
                <div style={{ fontSize: 14 }}>{showAdd.translation_uz}</div>
              </div>
              {showAdd.examples?.length > 0 && (
                <div style={{ marginBottom: 8 }}>
                  <div style={{ fontSize: 12, color: 'var(--tg-hint)', marginBottom: 2 }}>{t('vocab.example')}</div>
                  <ul style={{ margin: 0, paddingLeft: 16, fontSize: 13 }}>
                    {showAdd.examples.map((ex, i) => (
                      <li key={i} style={{ marginBottom: 4 }}>{ex}</li>
                    ))}
                  </ul>
                </div>
              )}
              <div style={{ marginBottom: 12 }}>
                <div style={{ fontSize: 12, color: 'var(--tg-hint)', marginBottom: 2 }}>{t('vocab.level')}</div>
                <span style={{ fontSize: 12, background: 'var(--tg-secondary-bg)', padding: '2px 6px', borderRadius: 4 }}>
                  {showAdd.level?.toUpperCase()}
                </span>
              </div>

              <div style={{ display: 'flex', flexDirection: 'column', gap: 8, marginBottom: 12 }}>
                <input
                  value={addWord}
                  onChange={e => setAddWord(e.target.value)}
                  placeholder={t('vocab.word_label')}
                  style={{ padding: '8px 12px', border: '1px solid var(--tg-border)', borderRadius: 8, background: 'var(--tg-secondary-bg)', color: 'var(--tg-text)', fontSize: 14, outline: 'none' }}
                />
                <input
                  value={addTranslation}
                  onChange={e => setAddTranslation(e.target.value)}
                  placeholder={t('vocab.translation_label')}
                  style={{ padding: '8px 12px', border: '1px solid var(--tg-border)', borderRadius: 8, background: 'var(--tg-secondary-bg)', color: 'var(--tg-text)', fontSize: 14, outline: 'none' }}
                />
                <input
                  value={addExample}
                  onChange={e => setAddExample(e.target.value)}
                  placeholder={t('vocab.example_label')}
                  style={{ padding: '8px 12px', border: '1px solid var(--tg-border)', borderRadius: 8, background: 'var(--tg-secondary-bg)', color: 'var(--tg-text)', fontSize: 14, outline: 'none' }}
                />
              </div>

              <button
                className="btn btn-primary"
                style={{ width: '100%' }}
                onClick={() => addMutation.mutate()}
                disabled={addMutation.isPending}
              >
                {t('vocab.save_word')}
              </button>

              {addMutation.isSuccess && (
                <div style={{ fontSize: 13, color: '#4caf50', marginTop: 8, textAlign: 'center' }}>
                  {t('vocab.saved')}
                </div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
