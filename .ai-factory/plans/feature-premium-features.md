# Premium Features

> Branch: `feature/premium-features`
> Created: 2026-06-07
> Settings: Testing=yes, Logging=verbose, Docs=yes

## Roadmap Linkage

- Milestone: "Premium Features"
- Rationale: Three premium-only features — advanced error analysis, CSV vocabulary export, priority AI queue — that unlock value for paid subscribers.

## Research Context

- Current `Correction` model: `{original, corrected, explanation_uz, explanation_ru, type}` — no severity/category/tips
- Current `BuildChatPrompt` asks for `corrections` array with 1-3 items max — no premium variation
- No export functionality exists anywhere in the codebase
- No queue/priority mechanism exists
- Frontend `GrammarBlock.tsx` renders corrections with type icon, strikethrough original → green corrected, and explanation

---

## Tasks

### Phase 1: Advanced Error Analysis (Premium)

**T1 — Enhance Correction model with premium fields**

Files: `internal/models/types.go`

- Add fields to `Correction` struct:
  ```go
  Severity    string   `json:"severity,omitempty"`    // "critical" | "major" | "minor"
  Category    string   `json:"category,omitempty"`    // "grammar" | "vocabulary" | "spelling" | "word_order" | "punctuation"
  LearningTip string   `json:"learning_tip,omitempty"` // e.g. "Remember: 'yesterday' always triggers Past Simple"
  RuleViolated string  `json:"rule_violated,omitempty"` // e.g. "Subject-Verb Agreement (Present Simple)"
  ```
- Add `PremiumAnalysis` struct (returned only for premium):
  ```go
  type PremiumAnalysis struct {
    OverallGrade    string `json:"overall_grade"`    // "A" | "B" | "C" | "D"
    Strengths       []string `json:"strengths"`
    AreasForImprovement []string `json:"areas_for_improvement"`
    SuggestedTopic  string `json:"suggested_topic"`  // "Next: practice Past Continuous"
  }
  ```
- Add `PremiumAnalysis` field to `AIResponse` and `ChatResponse`

**T2 — Premium chat prompt**

File: `internal/ai/prompts.go`

- Add `BuildPremiumChatPrompt(level, lang, text string) string`:
  - Same base as `BuildChatPrompt` but asks for richer output
  - Requests: `severity`, `category`, `learning_tip`, `rule_violated` per correction
  - Requests `premium_analysis` block with `overall_grade`, `strengths[]`, `areas_for_improvement[]`, `suggested_topic`
  - Allows up to 5 corrections instead of 3
  - Temperature stays 0.4

- Logging: `sugar.Debugw("premium prompt built", "prompt_len", len(prompt))`

**T3 — Modify chat handler for premium**

Files: `internal/api/chat.go`, `internal/api/chat_stream.go`

- Detect `is_premium` from context (`c.Get("is_premium")`)
- If premium → call `BuildPremiumChatPrompt` instead of `BuildChatPrompt`
- Parse enriched response and return `premium_analysis` in response
- For streaming: include premium analysis in final `EventResult`

- Logging:
  - `sugar.Infow("premium chat", "user_id", uid)` when premium path taken
  - `sugar.Debugw("premium analysis parsed", "grade", grade)` on success
  - `sugar.Warnw("premium analysis parse failed", "error", err)` on failure (fall back to standard prompt silently)

**T4 — Create PremiumCorrectionBlock UI component**

File: `web/src/components/PremiumCorrectionBlock.tsx`

- Extends GrammarBlock with:
  - Severity badge (🔴 critical, 🟡 major, 🔵 minor)
  - Rule violation label (`Subject-Verb Agreement`)
  - Learning tip section at the bottom
  - Overall grade pill at the top of the corrections block
  - Strengths / Areas for improvement sections
  - Suggested topic next-step card

- Use `{!isLoading && ...}` pattern for conditional rendering (per skill-context rule)
- Map severity to visual style via Telegram theme vars

**T5 — Integrate PremiumCorrectionBlock into Chat**

File: `web/src/components/Chat.tsx`

- Pass `isPremium` prop to Chat component
- When `msg.role === 'ai' && msg.corrections`: if `isPremium` and `premium_analysis` exists, render `PremiumCorrectionBlock` instead of `GrammarBlock`
- Otherwise render standard `GrammarBlock` (current behavior)

- Logging: `console.debug('[chat] premium corrections rendered')` for premium path

**T6 — i18n keys for premium analysis**

Files: `web/src/locales/uz.json`, `web/src/locales/ru.json`

Add keys:
- `premium.overall_grade` — "Общая оценка" / "Umumiy baho"
- `premium.strengths` — "Сильные стороны" / "Kuchli tomonlar"
- `premium.areas_to_improve` — "Над чем работать" / "Ustida ishlash"
- `premium.suggested_topic` — "Рекомендуемая тема" / "Tavsiya etilgan mavzu"
- `premium.severity_critical` — "Критическая" / "Tanqidiy"
- `premium.severity_major` — "Важная" / "Muhim"
- `premium.severity_minor` — "Незначительная" / "Kichik"
- `premium.learning_tip` — "Совет" / "Maslahat"
- `premium.export_vocab` — "Экспорт CSV" / "CSV экспорт"
- `premium.export_success` — "Скачивание началось" / "Yuklab olish boshlandi"

**T7 — Tests for advanced analysis**

Files: `internal/ai/prompts_test.go` (new or extend), `internal/api/chat_test.go` (new or extend)

- Test `BuildPremiumChatPrompt` produces valid JSON structure
- Test premium path in chat handler (mock is_premium=true)
- Test premium response parsing with all new fields
- Test fallback to standard prompt when premium parsing fails

---

### Phase 2: CSV Vocabulary Export

**T8 — CSV export endpoint**

File: `internal/api/vocab.go`

- Add `GET /api/vocab/export` handler:
  - Query user's vocabulary from DB (same query as `vocabListHandler` but no pagination)
  - Write CSV with columns: `word, translation, example, level, review_count, next_review, created_at`
  - Headers: `Content-Type: text/csv`, `Content-Disposition: attachment; filename="vocabulary.csv"`, `Cache-Control: no-store`
  - Use `encoding/csv` writer

- Register in router.go: `protected.GET("/vocab/export", vocabExportHandler(db, sugar))`

- Logging:
  - `sugar.Infow("vocab export", "user_id", uid, "count", len(words))`
  - `sugar.Errorw("vocab export failed", "error", err)` on DB error

**T9 — Export button in frontend**

File: `web/src/components/Vocabulary.tsx`

- Add "Export CSV" button in the header area (next to the tab toggle)
- On click: call `GET /api/vocab/export` via fetch with blob response
  - Create a download link and trigger it programmatically
  - Use `URL.createObjectURL` + `<a download>`
- Button shown only for premium users (check `isPremium` from context/prop)

- Button pattern: `{!isLoading && <button>}` per skill-context rule
- Logging: `console.debug('[vocab] export clicked')`

**T10 — Tests for CSV export**

File: `internal/api/vocab_test.go` (new or extend)

- Test CSV response format
- Test Content-Type and Content-Disposition headers
- Test empty vocabulary case

---

### Phase 3: Priority AI Queue

**T11 — Priority queue implementation**

File: `internal/ai/queue.go` (new)

- Create simple channel-based priority queue:
  ```go
  type Priority int
  const (
    PriorityNormal  Priority = 0
    PriorityPremium Priority = 1
  )

  type AIRequest struct {
    Priority  Priority
    Prompt    string
    ResultCh  chan string
    ErrorCh   chan error
    Ctx       context.Context
  }
  ```
- `StartWorker(ctx context.Context, gemini *Client, sugar *zap.SugaredLogger)` — goroutine that reads from priority queue and processes via `gemini.Generate()`
- `Enqueue(req AIRequest)` — inserts into queue (premium requests at front via two internal channels)
- `StopWorker()`

- Logging:
  - `sugar.Debugw("ai queue: premium request enqueued")` for premium
  - `sugar.Debugw("ai queue: normal request enqueued")` for free
  - `sugar.Infow("ai queue worker started")`
  - `sugar.Errorw("ai queue: request failed", "error", err)`

**T12 — Integrate queue with chat handler**

File: `internal/api/chat.go`

- When queue is enabled and request comes in:
  - Premium requests → enqueue with `PriorityPremium`
  - Free requests → enqueue with `PriorityNormal`
  - Wait on `ResultCh` / `ErrorCh`
  - Return response as normal

- Queue integration is behind a feature flag (`AIQueueEnabled` in Client) — disabled by default
- Wire through `cmd/server/main.go` (start worker on server boot if flag is on)

- Logging: `sugar.Debugw("chat via queue", "user_id", uid, "priority", priority)`

**T13 — Priority queue tests**

File: `internal/ai/queue_test.go` (new)

- Test basic enqueue/dequeue with normal priority
- Test premium priority goes before normal
- Test context cancellation
- Test worker starts and stops

---

## Commit Plan

| # | Tasks | Commit Message |
|---|-------|---------------|
| 1 | T1 | `feat(model): add premium correction fields (severity, tip, rule, category)` |
| 2 | T2-T3 | `feat(ai): add premium chat prompt and premium detection in handler` |
| 3 | T4-T6 | `feat(ui): add PremiumCorrectionBlock component with i18n keys` |
| 4 | T7 | `test: premium analysis tests` |
| 5 | T8-T9 | `feat(api): add CSV vocabulary export endpoint and UI button` |
| 6 | T10 | `test: CSV export tests` |
| 7 | T11-T12 | `feat(ai): add priority AI queue with premium/free differentiation` |
| 8 | T13 | `test: priority queue tests` |
| 9 | All | `docs: update docs for premium features` |
