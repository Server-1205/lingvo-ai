# Lingvo AI — Session Context

> Загружается автоматически при старте сессии.
> Обновляй этот файл после значимых изменений (новые компоненты, архитектурные решения, переезды API).

**Пользователь предпочитает русский язык — отвечай на русском.**

---

## Состояние проекта

- **Стадия**: Core — Backend API + middleware работают
- **Backend**: Go-сервер компилируется, 20 эндпоинтов (+6 IELTS: writing, speaking generate/evaluate, reading generate/evaluate, scores)
- **Middleware**: Auth (HMAC-SHA256 initData), Ratelimit (10/day free, unlimited premium)
- **Frontend**: UI написан (Chat, NavBar, Subscription, Vocabulary, Progress, LevelTest, ReviewCard, IeltsDashboard, IeltsWriting, IeltsSpeaking, IeltsReading) — React 19, Vite 8, TS 6
- **Бот**: Long-polling работает (/start, /help, /daily, /stats, successful_payment)
- **AI**: Gemini-клиент + SM-2 алгоритм (sm2.go, prompts.go, response.go), Streaming через SSE (chat_stream.go)
- **БД**: Схема готова (6 таблиц: users, messages, subscriptions, vocabulary, daily_progress, ielts_scores + SM-2 колонки), миграция работает, CRUD готов (users, messages, subscriptions, vocabulary, progress, ielts)
- **i18n**: Frontend uz.json + ru.json (95+ ключей, включая ielts.*), Backend uz.json + ru.json (7 ключей)
- **SM-2 Spaced Repetition**: ✅ Полностью реализован (алгоритм + API + UI + тесты)
- **IELTS Full Prep Bundle**: ✅ Writing (Task 1/2), Speaking (Part 1-3), Reading, Band Score tracking, premium-gated (все 6 эндпоинтов + 4 UI компонента + Go/Vitest тесты)

## Текущий фокус

- (none — IELTS завершён)

## Архитектурные решения

- База: SQLite (modernc.org/sqlite, без CGO)
- Аутентификация: HMAC-SHA256 initData
- Платежи: Telegram Stars
- Rate limit: Free 10/день, Premium — безлимит
- AI: Gemini 2.0 Flash (температура 0.4), fallback GPT-4o-mini

## Важные пути

- `cmd/server/main.go` — точка входа
- `internal/db/schema.sql` — схема БД
- `web/src/locales/` — переводы
- `internal/middleware/auth.go` — HMAC initData verification
- `internal/middleware/ratelimit.go` — daily limit check
- `internal/ai/gemini.go` — Gemini client wrapper
- `internal/ai/sm2.go` — SM-2 algorithm
- `internal/ai/prompts.go` — промпты для всех эндпоинтов
- `internal/api/chat.go` — POST /api/chat
- `internal/api/subscription.go` — GET /api/subscription
- `internal/api/invoice.go` — POST /api/create-invoice (Telegram Stars)
- `internal/bot/bot.go` — long-polling loop
- `internal/bot/handlers.go` — /start, /help, /daily, successful_payment
- `web/src/components/Chat.tsx` — Chat screen (streaming via SSE)
- `web/src/components/NavBar.tsx` — Bottom navigation
- `web/src/components/Subscription.tsx` — Subscription plans
- `web/src/api/client.ts` — HTTP API client (streaming via chatStream)
- `internal/api/chat_stream.go` — SSE streaming endpoint
- `docs/streaming.md` — Streaming AI responses docs
