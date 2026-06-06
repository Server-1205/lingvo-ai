# SM-2 Spaced Repetition для Vocabulary

**Branch:** `feature/spaced-repetition`
**Created:** 2026-06-06 20:15

## Settings

| Key | Value |
|---|---|
| Testing | Yes |
| Logging | Verbose (DEBUG) |
| Docs | Yes — mandatory checkpoint |
| Tests | Yes |

## Roadmap Linkage

**Milestone:** "Этап 2: Vocabulary + Level Test — Spaced Repetition (алгоритм SM-2 упрощённый)"
**Rationale:** SM-2 algorithm is the core of the Vocabulary Builder milestone. Without it, words are learned once and never reviewed, defeating the purpose of a personal dictionary.

## Tasks

### Phase 1: Backend — Schema + SM-2 Algorithm

**Task 1.1 — Add SM-2 columns to vocabulary schema**

- Add migration to schema.sql: `ease_factor REAL DEFAULT 2.5`, `interval INTEGER DEFAULT 0`, `last_reviewed_at DATETIME`
- Update VocabWord model: add `EaseFactor float64`, `Interval int`, `LastReviewedAt *time.Time`
- Update client.ts VocabWord interface with new fields

**Task 1.2 — Create SM-2 algorithm package**

File: `internal/ai/sm2.go`

- Struct `SM2Card` with fields: `Repetitions int`, `Interval int`, `EaseFactor float64`
- Function `NewSM2Card()` — returns default card (0, 0, 2.5)
- Function `ProcessReview(card, quality int) SM2Card` — applies SM-2:
  - quality < 3 (failed): reset repetitions to 0, interval = 1
  - quality >= 3 (passed): interval progression 1→6→(interval×EF)
  - EF update: `EF' = EF + (0.1 - (5-q)*(0.08 + (5-q)*0.02))`, minimum 1.3
- Logging: `[FIX] SM2 review: word=%d quality=%d interval=%d ef=%.2f`

### Phase 2: Backend — DB Functions

**Task 2.1 — GetDueWords**

File: `internal/db/vocabulary.go`

```go
func GetDueWords(ctx context.Context, db *sqlx.DB, userID, limit int) ([]models.VocabWord, error)
```
- Query: `SELECT * FROM vocabulary WHERE user_id = ? AND (next_review IS NULL OR next_review <= datetime('now')) ORDER BY next_review ASC LIMIT ?`
- Logging: `[FIX] GetDueWords user=%d due=%d`

**Task 2.2 — UpdateReview**

File: `internal/db/vocabulary.go`

```go
func UpdateReview(ctx context.Context, db *sqlx.DB, wordID, userID int, reviewCount, interval int, easeFactor float64, nextReview string) error
```
- UPDATE query with all SM-2 fields
- Increment `words_learned` in daily_progress
- Logging: `[FIX] UpdateReview word=%d interval=%d ef=%.2f next=%s`

### Phase 3: Backend — API Handlers

**Task 3.1 — GET /api/vocab/review**

File: `internal/api/vocab.go`

- Handler `vocabReviewHandler` — returns due words list
- Uses `GetDueWords` with limit (default 20)
- Response: `VocabWord[]`
- Logging: `[FIX] Review user=%d due=%d`

**Task 3.2 — POST /api/vocab/review/:id**

File: `internal/api/vocab.go`

- Handler `vocabReviewSubmitHandler`
- Request body: `{ "quality": 0-5 }`
- Fetches word from DB, applies SM-2 via `ProcessReview`
- Saves updated card state via `UpdateReview`
- Response: `{ "next_review": "..." }`
- Validation: quality must be 0-5
- Logging: `[FIX] Review submit word=%d quality=%d next=%s`

**Task 3.3 — Register new routes**

File: `internal/api/router.go`

- Add `protected.GET("/vocab/review", vocabReviewHandler(db))`
- Add `protected.POST("/vocab/review/:id", vocabReviewSubmitHandler(db, sugar))`

### Phase 4: Frontend — API Client

**Task 4.1 — Add review API functions**

File: `web/src/api/client.ts`

- `getDueWords(limit?: number): Promise<VocabWord[]>`
- `submitReview(wordId: number, quality: number): Promise<{ next_review: string }>`
- Add quality constant: `ReviewQuality = { Again: 0, Hard: 2, Good: 4, Easy: 5 }`

**Task 4.2 — Add i18n review keys**

Files: `web/src/locales/uz.json`, `web/src/locales/ru.json`

```json
{
  "vocab.review": "Takrorlash / Повторение",
  "vocab.review_due": "{count} ta so'z / {count} слов",
  "vocab.review_again": "Yana / Снова",
  "vocab.review_hard": "Qiyin / Сложно",
  "vocab.review_good": "Yaxshi / Хорошо",
  "vocab.review_easy": "Oson / Легко",
  "vocab.review_progress": "{current}/{total}",
  "vocab.review_done": "Bugungi takrorlash tugadi! / На сегодня всё!",
  "vocab.review_empty": "Takrorlash uchun so'z yo'q / Нет слов для повторения"
}
```

### Phase 5: Frontend — Review UI

**Task 5.1 — Create ReviewCard component**

File: `web/src/components/ReviewCard.tsx`

- Props: `word: VocabWord, onRate: (quality: number) => void`
- States: `flipped: boolean` (show answer side)
- Front: word + level badge
- Back: translation + example + 4 rating buttons (Again/Hard/Good/Easy)
- Smooth flip animation using CSS transform
- Loading state while submitting review

**Task 5.2 — Integrate review into Vocabulary screen**

File: `web/src/components/Vocabulary.tsx`

- Add 3rd tab "Review" with due count badge
- Fetch due words via `getDueWords` on mount
- Show `ReviewCard` + progress `{current}/{total}`
- After all reviewed: show "Done" screen with option to review again
- Auto-invalidate vocab list after review

### Phase 6: Tests

**Task 6.1 — SM-2 algorithm unit tests**

File: `internal/ai/sm2_test.go`

- Test `NewSM2Card` defaults
- Test "quality = 0" (fail) → reset, interval=1
- Test "quality = 5" (perfect) → EF update, interval progression
- Test EF minimum floor (1.3)
- Test multiple review cycles (day 1, day 2, etc.)

**Task 6.2 — Frontend ReviewCard tests**

File: `web/src/components/__tests__/ReviewCard.test.tsx`

- Test renders word on front
- Test flip reveals answer
- Test quality buttons call `onRate`
- Test empty state

## Commit Plan

| # | Tasks | Message |
|---|---|---|
| 1 | 1.1, 1.2 | `feat(sm2): add SM-2 algorithm and schema columns` |
| 2 | 2.1, 2.2 | `feat(db): add GetDueWords and UpdateReview` |
| 3 | 3.1, 3.2, 3.3 | `feat(api): add review endpoints GET/POST /api/vocab/review` |
| 4 | 4.1, 4.2 | `feat(client): add review API and i18n keys` |
| 5 | 5.1, 5.2 | `feat(ui): add ReviewCard component and review tab` |
| 6 | 6.1, 6.2 | `test(sm2): add algorithm and frontend tests` |

## Files to Modify

- `internal/db/schema.sql` — SM-2 columns
- `internal/models/types.go` — VocabWord SM-2 fields
- `internal/ai/sm2.go` — new: SM-2 algorithm
- `internal/ai/sm2_test.go` — new: algorithm tests
- `internal/db/vocabulary.go` — GetDueWords, UpdateReview
- `internal/api/vocab.go` — review handlers
- `internal/api/router.go` — review routes
- `web/src/api/client.ts` — review API functions
- `web/src/locales/uz.json` — review i18n
- `web/src/locales/ru.json` — review i18n
- `web/src/components/ReviewCard.tsx` — new: flashcard component
- `web/src/components/Vocabulary.tsx` — review tab
- `web/src/components/__tests__/ReviewCard.test.tsx` — new: component tests

## Risks & Considerations

- **SM-2 quality scale**: 0-5 must be clearly documented. Use 4-button UI: Again(1)→Hard(2)→Good(4)→Easy(5)
- **Next review date**: Words with NULL next_review are due immediately (new words)
- **Daily limit**: Review is free (no rate limit), unlike chat. Only words_learned is logged
- **Edge case**: User has 0 words → show "No words for review" empty state
- **Edge case**: All words reviewed → show completion screen today, reset interval tomorrow
