# Lingvo AI — Session Context

> Загружается автоматически при старте сессии.
> Обновляй этот файл после значимых изменений (новые компоненты, архитектурные решения, переезды API).

---

## Состояние проекта

- **Стадия**: Core — Backend API + middleware работают
- **Backend**: Go-сервер компилируется, 12 эндпоинтов работают (chat, grammar, vocab CRUD, vocab lookup, quiz, level test, level save, progress, subscription, invoice)
- **Middleware**: Auth (HMAC-SHA256 initData), Ratelimit (10/day free, unlimited premium)
- **Frontend**: UI написан (Chat, NavBar, Subscription, Vocabulary, Progress) — React 19, Vite 8, TS 6
- **Бот**: Long-polling работает (/start, /help, /daily, successful_payment)
- **AI**: Gemini-клиент реализован (gemini.go, prompts.go, response.go)
- **БД**: Схема готова (5 таблиц), миграция работает, CRUD готов (users, messages, subscriptions)
- **i18n**: Frontend uz.json + ru.json (31 ключ), Backend uz.json + ru.json (7 ключей)

## Текущий фокус

- Frontend UI (Vocabulary, Progress, Grammar — подключить к реальным API)

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
- `internal/ai/prompts.go` — промпты для всех эндпоинтов
- `internal/api/chat.go` — POST /api/chat
- `internal/api/subscription.go` — GET /api/subscription
- `internal/api/invoice.go` — POST /api/create-invoice (Telegram Stars)
- `internal/bot/bot.go` — long-polling loop
- `internal/bot/handlers.go` — /start, /help, /daily, successful_payment
- `web/src/components/Chat.tsx` — Chat screen
- `web/src/components/NavBar.tsx` — Bottom navigation
- `web/src/components/Subscription.tsx` — Subscription plans
- `web/src/api/client.ts` — HTTP API client
