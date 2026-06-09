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
- **AI Fallback**: ✅ DeepSeek fallback работает — оба эндпоинта (/api/chat и /api/chat/stream) успешно фолбекятся при 429 Gemini

## Текущий фокус

- (none — все feature-ветки слиты в master)

## Session Memory MCP

- Настроен `@modelcontextprotocol/server-memory` (knowledge graph) в `.opencode/opencode.json`
- Memory file: `.opencode/memory/memory.jsonl`
- Используется `MEMORY_FILE_PATH` env var для проекта
- Добавлена инструкция в AGENTS.md использовать memory MCP при необходимости

## Архитектурные решения

- База: SQLite (modernc.org/sqlite, без CGO)
- Аутентификация: HMAC-SHA256 initData
- Платежи: Telegram Stars
- Rate limit: Free 10/день, Premium — безлимит
- AI: Gemini 2.0 Flash (температура 0.4), fallback DeepSeek-chat через OpenAI-compatible API

## Важные пути

- `cmd/server/main.go` — точка входа
- `internal/db/schema.sql` — схема БД
- `web/src/locales/` — переводы
- `internal/middleware/auth.go` — HMAC initData verification
- `internal/middleware/ratelimit.go` — daily limit check
- `internal/ai/gemini.go` — Gemini client wrapper (+ DeepSeek fallback)
- `internal/ai/openai.go` — OpenAI-compatible клиент (DeepSeek)
- `internal/ai/sm2.go` — SM-2 algorithm
- `internal/ai/prompts.go` — промпты для всех эндпоинтов
- `internal/api/chat.go` — POST /api/chat
- `internal/api/chat_stream.go` — SSE streaming endpoint
- `internal/api/subscription.go` — GET /api/subscription
- `internal/api/invoice.go` — POST /api/create-invoice (Telegram Stars)
- `internal/bot/bot.go` — long-polling loop
- `internal/bot/handlers.go` — /start, /help, /daily, successful_payment
- `web/src/components/Chat.tsx` — Chat screen (streaming via SSE)
- `web/src/components/NavBar.tsx` — Bottom navigation
- `web/src/components/Subscription.tsx` — Subscription plans
- `web/src/api/client.ts` — HTTP API client (streaming via chatStream)
- `docs/streaming.md` — Streaming AI responses docs
- `.env` — конфиг: Gemini API key, DeepSeek fallback, WebApp URL
- `.ai-factory/patches/` — патчи эволюции (2 файла)

## Сессия 2026-06-07 10:45 — DeepSeek fallback, i18n, bot URL

### Что сделано
1. **ChatStream fallback**: Добавлен вызов `c.fallback.Generate()` при ошибке Gemini stream в `ChatStream()` (gemini.go:234-246). Ответ парсится как `AIResponse`, извлекается только `reply` для токена — JSON больше не отображается в чате
2. **i18n**: Исправлен синтаксис интерполяции `{variable}` → `{{variable}}` в uz.json и ru.json (6 ключей)
3. **Bot URL**: `launchKeyboard()` читает `WEBAPP_URL` из параметра вместо хардкода
4. **Chat padding**: Добавлен `paddingBottom` для предотвращения перекрытия NavBar
5. **Vite proxy**: Настроен proxy `/api` → `:8080` для dev-режима
6. **Telegram SDK init**: Исправлен async `isTMA()`, добавлен `initData.restore()`

### Результаты
- `go build ./...` ✅
- `npm run build` ✅
- `go test ./...` — 13/13 ✅
- Оба эндпоинта (/api/chat и /api/chat/stream) работают через DeepSeek fallback ✅
- Gemini возвращает 429 (quota exceeded), но DeepSeek успешно обрабатывает запросы

## Сессия 2026-06-08 — Google Stitch промпт, Session Memory MCP

### Что сделано
1. **Google Stitch промпт**: Создан единый промпт для всех экранов Lingvo AI (Onboarding, Chat, Vocabulary, Progress, Subscription, LevelTest, DailyLesson, NavBar) для генерации UI-дизайна через Google Stitch
2. **Session Memory MCP**: Настроен `@modelcontextprotocol/server-memory` — knowledge graph-based MCP сервер для сохранения контекста между сессиями
3. **Memory file**: `.opencode/memory/memory.jsonl`
4. **Config**: `MEMORY_FILE_PATH` env var указывает на project-specific memory файл

### Результаты
- Memory MCP добавлен в `.opencode/opencode.json` ✅
- Memory файл инициализирован ✅
- Контекст обновлён ✅
