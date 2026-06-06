# Implementation Plan: /stats Bot Command

Branch: `feature/stats-bot-command`
Created: 2026-06-06

## Settings
- Testing: yes
- Logging: verbose
- Docs: yes

## Roadmap Linkage
Milestone: "Этап 1: Ядро (Дни 1-5) — Telegram Bot"
Rationale: /stats is the last missing bot command from Stage 1 core infrastructure

## Commit Plan
- **Commit 1** (after tasks 1-2): `feat(db): add GetUserStats query and UserStats model`
- **Commit 2** (after tasks 3-4): `feat(bot): add /stats command with help text update`
- **Commit 3** (after tasks 5-6): `test(bot): add tests for GetUserStats and /stats handler`

## Tasks

### Phase 1: Database Query
- [x] **Task 1: Add `GetUserStats` DB function** (depends on: none)

  **Files:**
  - `internal/db/progress.go` — add `GetUserStats`

  **What to do:**
  - Create a `GetUserStats(ctx, db, userID int)` function that returns a `UserStats` struct (defined in models) containing:
    - `Level` (from users table)
    - `TotalMessages` (SELECT COALESCE(SUM(count),0) FROM messages WHERE user_id = ?)
    - `TotalWords` (SELECT COUNT(*) FROM vocabulary WHERE user_id = ?)
    - `WordsDueToday` (SELECT COUNT(*) FROM vocabulary WHERE user_id = ? AND (next_review IS NULL OR next_review <= datetime('now')))
    - `StreakDays` (reuse existing `GetStreakDays`)
    - `IsPremium` (check via existing `GetActiveSubscription`)
    - `AccountCreatedAt` (from users table)
    - `SubscriptionExpiresAt` (from subscriptions, if active)
  - Run multiple queries sequentially, compose into one struct
  - Return nil + nil if user not found (caller handles "no stats" message)

  **Logging:**
  - DEBUG: "db: GetUserStats — user_id=%d"
  - WARN: if any query fails, log user_id + error but continue with zeros for that field

- [x] **Task 2: Add `UserStats` model type** (depends on: none)

  **Files:**
  - `internal/models/types.go` — add UserStats struct

  **Fields:**
  ```go
  type UserStats struct {
      Level              string `json:"level"`
      TotalMessages      int    `json:"total_messages"`
      TotalWords         int    `json:"total_words"`
      WordsDueToday      int    `json:"words_due_today"`
      StreakDays         int    `json:"streak_days"`
      IsPremium          bool   `json:"is_premium"`
      AccountCreatedAt   string `json:"account_created_at"`
      SubscriptionExpiry string `json:"subscription_expiry,omitempty"`
  }
  ```

  No logging needed for this task.

### Phase 2: Bot Handler
- [x] **Task 3: Add `/stats` handler** (depends on: 1, 2)

  **Files:**
  - `internal/bot/handlers.go` — add case "stats", add handleStats function

  **What to do:**
  - Add `case "stats":` to the switch in `handleCommand()`
  - Create `handleStats(bot, database, sugar, chatID, telegramID, lang)` function
  - Call `GetUserStats(ctx, database, user.ID)` (need to get user first via GetUserByTelegramID)
  - Format a stats message with uz/ru translations including:
    - 📊 **Statistika / Статистика** header
    - Level badge
    - Total messages sent
    - Words in vocabulary + due for review today
    - Streak days
    - Premium status (with expiry if active) or "Free plan (10/day)"
    - Account age (days since created_at)
  - If user not found: send "Xatolik / Ошибка"

  **Logging:**
  - DEBUG: "bot: /stats — telegram_id=%d" on entry
  - DEBUG: "bot: /stats — user_id=%d, messages=%d, words=%d, streak=%d" after fetching
  - ERROR: if GetUserByTelegramID or GetUserStats fails
  - INFO: "/stats command completed for user %d"

- [x] **Task 4: Add `/stats` to help text** (depends on: 3)

  **Files:**
  - `internal/bot/handlers.go` — update helpText()

  **What to do:**
  - Add `/stats` to the uz and ru help text in `helpText()`

  **Logging:**
  - None needed

### Phase 3: Tests
- [x] **Task 5: Write tests for DB function** (depends on: 1, 2)

  **Files:**
  - `internal/db/progress_test.go` — new file

  **What to do:**
  - Test `GetUserStats` with:
    - Existing user with messages, vocabulary, subscription
    - User with no data (should return zeros, not fail)
    - Non-existent user (should return nil)
  - Use in-memory SQLite for tests

  **Logging:**
  - DEBUG: "test: GetUserStats — setup fixture"
  - INFO: "test: GetUserStats — passed %d assertions"

- [x] **Task 6: Write tests for bot handler** (depends on: 2, 3)

  **Files:**
  - `internal/bot/handlers_test.go` — new file

  **What to do:**
  - Test that `/stats` command produces a message in expected format
  - Test uz and ru language outputs
  - Test with premium and free user
  - Mock the DB functions or use in-memory SQLite

  **Logging:**
  - DEBUG: "test: handleStats — setup mock user with %s plan"
