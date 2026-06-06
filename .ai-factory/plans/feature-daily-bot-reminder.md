# Daily Bot Reminder — "Repeat Your Words"

**Branch:** `feature/daily-bot-reminder`
**Created:** 2026-06-06

## Settings

| Key | Value |
|---|---|
| Testing | Yes |
| Logging | Verbose (DEBUG) |
| Docs | Yes — mandatory checkpoint |

## Roadmap Linkage

**Milestone:** "Этап 2: Vocabulary + Level Test — Ежедневное напоминание: Repeat 5 words"
**Rationale:** Bot sends a daily reminder to users with due words, driving engagement and returning users to the app. Completes the last unchecked item in Этап 2.

## Tasks

### Phase 1: Backend — DB + Bot

**Task 1.1 — Add GetAllUsers and GetDueWordCount** ✅

Files: `internal/db/users.go`, `internal/db/vocabulary.go`

- `GetAllUsers(ctx, db) ([]models.User, error)` — returns all users (for scheduler iteration)
- `GetDueWordCount(ctx, db, userID int) (int, error)` — `SELECT COUNT(*) FROM vocabulary WHERE user_id = ? AND (next_review IS NULL OR next_review <= datetime('now'))`
- Logging: `[FIX] GetAllUsers count=%d`, `[FIX] GetDueWordCount user=%d count=%d`

**Task 1.2 — Create bot reminder scheduler**

File: `internal/bot/reminder.go` (new)

- Function `StartReminderScheduler(bot, db, sugar)` — goroutine with 1-hour ticker
- On each tick: `GetAllUsers` → for each user `GetDueWordCount` → if > 0, send reminder
- In-memory map `remindedToday map[int64]time.Time` to enforce 1-reminder-per-day per user
- Reminder text (uz/ru):
  - uz: `📚 Eslatma! Sizda {count} ta so'z takrorlash uchun.`
  - ru: `📚 Напоминание! У вас {count} слов для повторения.`
- Inline keyboard button "📖 Review Now" → deep link `https://t.me/lingvo_ai_bot/app?startapp=review`
- Logging: `[FIX] Reminder sent user=%d due=%d`, `[FIX] Reminder tick checked=%d sent=%d`

**Task 1.3 — Enhance /daily command**

File: `internal/bot/handlers.go`

- Modify `handleDaily` to also call `GetDueWordCount` for the user
- If due > 0, show count in the message and add "Review Now" button
- If due == 0, show "No words due" encouragement message
- Message text updated to include due word count
- Logging: `[FIX] Daily check user=%d due=%d`

### Phase 2: Frontend — Deep-link Navigation

**Task 2.1 — Vocabulary accepts initialTab prop**

File: `web/src/components/Vocabulary.tsx`

- Add optional prop `initialTab?: 'my' | 'lookup' | 'review'`
- Use it as initial value for the tab state (with fallback to 'my')
- When `initialTab='review'`, the review tab is auto-fetched and shown

**Task 2.2 — Start param handling in App.tsx**

File: `web/src/App.tsx`

- On mount, read `window.Telegram.WebApp.initDataUnsafe?.start_param`
- If `start_param === 'review'`, set activeTab to `'vocab'` and pass `initialTab="review"` to Vocabulary
- Use `useEffect` with empty deps to read initData once on mount
- Logging: `[FIX] Deep-link start_param=%s → tab=%s`

### Phase 3: Tests

**Task 3.1 — Reminder scheduler tests**

File: `internal/bot/reminder_test.go` (new)

- Test `GetDueWordCount` function
- Test reminder message format (uz/ru variants)
- Test that reminder doesn't send twice in one day (in-memory map check)

**Task 3.2 — Frontend deep-link tests**

Files: `web/src/components/__tests__/App.test.tsx` (new)

- Test that `start_param=review` sets active tab to vocab with initialTab=review
- Test that empty start_param keeps default tab (chat)

## Commit Plan

| # | Tasks | Message |
|---|---|---|
| 1 | 1.1, 1.2, 1.3 | `feat(bot): add daily word reminder scheduler and enhance /daily` |
| 2 | 2.1, 2.2 | `feat(ui): add deep-link support for review tab navigation` |
| 3 | 3.1, 3.2 | `test(bot): add reminder scheduler and deep-link tests` |

## Files to Modify/Create

- `internal/db/users.go` — GetAllUsers
- `internal/db/vocabulary.go` — GetDueWordCount
- `internal/bot/reminder.go` — new: scheduler + reminder logic
- `internal/bot/reminder_test.go` — new: scheduler tests
- `internal/bot/handlers.go` — enhanced /daily
- `web/src/components/Vocabulary.tsx` — initialTab prop
- `web/src/App.tsx` — start_param handling
- `web/src/__tests__/App.test.tsx` — new: deep-link tests

## Risks & Considerations

- **Spam prevention**: In-memory `remindedToday` map resets on bot restart. For production, a DB column `last_reminded_at` would be more persistent
- **Bot token**: The scheduler needs the bot token to create its own API client (or we pass the existing bot instance)
- **Deep link**: `startapp` param only works via `t.me/bot/app?startapp=...` deep links (url-type button). Regular WebApp buttons don't pass `startapp` automatically — need `url` type button with t.me link
- **No words**: If user has 0 words in vocabulary, no reminder is sent (GetDueWordCount returns 0)
- **Review tab state**: Review tab needs `dueWords` query to auto-fetch when navigating via deep link
