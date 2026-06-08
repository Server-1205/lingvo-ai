import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ReviewCard } from '../ReviewCard';
import type { VocabWord } from '../../api/client';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      const map: Record<string, string> = {
        'vocab.review_again': 'Снова',
        'vocab.review_hard': 'Сложно',
        'vocab.review_good': 'Хорошо',
        'vocab.review_easy': 'Легко',
      };
      return map[key] || key;
    },
  }),
}));

const mockWord: VocabWord = {
  id: 1,
  word: 'apple',
  translation: 'яблоко',
  translation_ru: 'яблоко',
  example: 'I eat an apple every day.',
  example_ru: 'I eat an apple every day.',
  level: 'a1',
  review_count: 0,
  ease_factor: 2.5,
  interval: 0,
  next_review: null,
  created_at: '2026-06-01T00:00:00Z',
};

describe('ReviewCard', () => {
  it('renders word on front', () => {
    render(<ReviewCard word={mockWord} onRate={() => {}} />);
    expect(screen.getByText('apple')).toBeInTheDocument();
    expect(screen.getByText('A1')).toBeInTheDocument();
  });

  it('reveals answer with rating buttons on click', () => {
    render(<ReviewCard word={mockWord} onRate={() => {}} />);
    fireEvent.click(screen.getByText('apple'));
    expect(screen.getByText('яблоко')).toBeInTheDocument();
    expect(screen.getByText('I eat an apple every day.')).toBeInTheDocument();
    expect(screen.getByText('Снова')).toBeInTheDocument();
    expect(screen.getByText('Сложно')).toBeInTheDocument();
    expect(screen.getByText('Хорошо')).toBeInTheDocument();
    expect(screen.getByText('Легко')).toBeInTheDocument();
  });

  it('calls onRate with correct quality when button clicked', () => {
    const onRate = vi.fn();
    render(<ReviewCard word={mockWord} onRate={onRate} />);
    fireEvent.click(screen.getByText('apple'));
    fireEvent.click(screen.getByText('Хорошо'));
    expect(onRate).toHaveBeenCalledWith(4);
  });
});
