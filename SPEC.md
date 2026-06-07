# Lingvo AI — Technical Specification

> AI English Tutor. Telegram Mini App. Uzbekistan → Central Asia.
> Backend: Go. Frontend: React + Vite. AI: Gemini 2.0 Flash (+ Flash-Lite, fallback GPT-4o-mini). DB: SQLite.

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Tech Stack](#2-tech-stack)
3. [Directory Structure](#3-directory-structure)
4. [Database Schema](#4-database-schema)
5. [Backend API](#5-backend-api)
6. [Telegram Bot](#6-telegram-bot)
7. [Mini App (Frontend)](#7-mini-app-frontend)
8. [AI Integration](#8-ai-integration)
9. [Payments & Monetization](#9-payments--monetization)
10. [i18n / Localization](#10-i18n--localization)
11. [Deployment](#11-deployment)
12. [Error Handling & Logging](#12-error-handling--logging)
13. [Security](#13-security)
14. [Testing Strategy](#14-testing-strategy)

---

## 1. Architecture Overview

```
┌──────────────────────────────────────────────────────────────┐
│                       Telegram                                │
│  ┌─────────────────────┐    ┌─────────────────────────────┐  │
│  │ Mini App (WebView)  │    │ Bot (Chat Interface)        │  │
│  │ React SPA           │    │ /start /help /daily /stats  │  │
│  │ tg://webview        │    │ Inline keyboard             │  │
│  └─────────┬───────────┘    └───────────┬─────────────────┘  │
└────────────┼───────────────────────────┼─────────────────────┘
             │                           │
             ▼                           ▼
┌──────────────────────────────────────────────────────────────┐
│                   Go Backend (Gin)                            │
│                                                               │
│  ┌───────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │ Bot Layer      │  │ REST API     │  │ AI Layer          │  │
│  │ long-polling   │  │ /api/*       │  │ gemini.Client     │  │
│  │ commands       │  │ auth MW      │  │ prompts.go        │  │
│  │ payments       │  │ ratelimit MW │  │ streaming         │  │
│  └───────┬───────┘  └──────┬───────┘  └────────┬─────────┘  │
│          │                 │                    │            │
│          └────────┬────────┴────────────────────┘            │
│                   │                                          │
│                   ▼                                          │
│          ┌────────────────┐                                  │
│          │ SQLite          │                                  │
│          │ users, messages │                                  │
│          │ subscriptions   │                                  │
│          │ vocabulary, pro │                                  │
│          └────────────────┘                                  │
└──────────────────────────────────────────────────────────────┘
```

### Data flow: Chat message

```
1. User types in Mini App chat input
2. React sends POST /api/chat { text, tg_init_data }
3. Middleware: auth (verify initData HMAC)
4. Middleware: check subscription / decrement daily limit
5. AI layer: call Gemini with system prompt + user text
6. Parse Gemini response: reply_text + corrections[]
7. Save message count to SQLite
8. Return JSON to Mini App
9. React renders: AI reply + correction blocks
```

### Data flow: Subscription

```
1. User taps "Get Unlimited" in Mini App
2. Backend: POST /api/create-invoice { plan: "weekly"|"monthly" }
3. Backend returns invoice link (Telegram Stars)
4. Frontend opens link in Telegram (tg://)
5. Telegram processes Stars payment (Apple/Google Pay)
6. Bot receives successful_payment via long-polling
7. Bot handler: save subscription to SQLite with expiry
8. Bot sends confirmation message to user
9. Next API call: middleware sees active subscription → no limit
```

---

## 2. Tech Stack

### Backend (Go)

| Package | Purpose |
|---|---|
| github.com/gin-gonic/gin | HTTP router |
| github.com/go-telegram-bot-api/telegram-bot-api/v5 | Telegram Bot API |
| github.com/google/generative-ai-go | Gemini API SDK |
| modernc.org/sqlite | SQLite driver (pure Go, no CGO) |
| github.com/jmoiron/sqlx | SQL query helpers |
| go.uber.org/zap | Structured logging |

### Frontend (React)

| Package | Purpose |
|---|---|
| react, react-dom | UI framework |
| vite | Bundler |
| (plain CSS) | Styling (no framework) |
| @telegram-apps/sdk | Telegram WebApp SDK |
| react-i18next | Internationalization |
| @tanstack/react-query | Server state / API calls |

### AI

| Service | Model | Input cost | Output cost | Use |
|---|---|---|---|---|
| Gemini 2.0 Flash | gemini-2.0-flash | $0.10/1M tok | $0.40/1M tok | Chat, grammar, lessons |
| Gemini 2.0 Flash-Lite | gemini-2.0-flash-lite | $0.075/1M tok | $0.30/1M tok | Level test, quiz |

Fallback: OpenAI-compatible model (default GPT-4o-mini) if Gemini is unavailable. Configurable via OPENAI_API_KEY / OPENAI_BASE_URL / OPENAI_MODEL.

---

## 3. Directory Structure

```
sinking/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point
├── internal/
│   ├── bot/
│   │   ├── bot.go                   # Init, long-polling
│   │   ├── handlers.go              # /start, /help, /daily, /stats
│   │   ├── payments.go              # successful_payment handler
│   │   └── reminder.go              # Daily lesson reminders
│   ├── api/
│   │   ├── router.go                # Gin router + middleware
│   │   ├── chat.go                  # POST /api/chat
│   │   ├── chat_stream.go           # POST /api/chat/stream (SSE)
│   │   ├── grammar.go              # POST /api/grammar (+ vocab lookup)
│   │   ├── vocab.go                 # POST/GET /api/vocab + review + delete
│   │   ├── quiz.go                  # POST /api/quiz
│   │   ├── level.go                 # POST /api/level + /api/level/save
│   │   ├── progress.go             # GET /api/progress + /api/progress/history
│   │   ├── subscription.go         # GET /api/subscription
│   │   └── invoice.go              # POST /api/create-invoice
│   ├── ai/
│   │   ├── gemini.go                # Gemini client wrapper + fallback
│   │   ├── openai.go                # OpenAI-compatible HTTP client
│   │   ├── prompts.go               # System prompts (UZ, RU)
│   │   ├── response.go              # Response parsing
│   │   └── sm2.go                   # SM-2 spaced repetition algorithm
│   ├── db/
│   │   ├── db.go                    # Init SQLite + schema migration
│   │   ├── users.go                 # User CRUD
│   │   ├── messages.go             # Message limit tracking
│   │   ├── subscriptions.go        # Subscription CRUD
│   │   ├── vocabulary.go           # Vocabulary CRUD + SM-2 review
│   │   └── progress.go             # Progress tracking + streak
│   ├── middleware/
│   │   ├── auth.go                  # Telegram initData verify (HMAC-SHA256)
│   │   └── ratelimit.go             # Daily message limit
│   └── models/
│       └── types.go                 # Shared structs
├── web/
│   ├── src/
│   │   ├── main.tsx                 # React entry
│   │   ├── App.tsx                  # Router + layout
│   │   ├── components/
│   │   │   ├── Chat.tsx
│   │   │   ├── GrammarBlock.tsx
│   │   │   ├── Vocabulary.tsx
│   │   │   ├── LevelTest.tsx
│   │   │   ├── Progress.tsx
│   │   │   ├── Subscription.tsx
│   │   │   ├── NavBar.tsx
│   │   │   └── LanguageSwitcher.tsx
│   │   ├── api/
│   │   │   └── client.ts
│   │   ├── hooks/
│   │   │   └── useTelegram.ts
│   │   └── locales/
│   │       ├── uz.json
│   │       └── ru.json
│   ├── index.html
│   ├── vite.config.ts
│   └── package.json
├── docs/
│   ├── api.md                  # API Reference
│   ├── configuration.md       # Environment variables
│   ├── getting-started.md     # Setup guide
│   └── streaming.md           # AI Streaming (SSE)
├── .ai-factory/
│   ├── DESCRIPTION.md         # Project description
│   ├── ARCHITECTURE.md        # Architecture guidelines
│   ├── plans/                 # Implementation plans
│   │   └── sync-spec.md
│   ├── rules/
│   │   └── base.md            # Code conventions
│   └── config.yaml            # AI Factory config
├── .opencode/
│   ├── agents/                # Subagent configs
│   ├── skills/                # Custom skills
│   └── memory/
│       └── context.md         # Session context
├── AGENTS.md
├── ROADMAP.md
├── SPEC.md
├── go.mod
├── .env.example
└── .gitignore
```

---

## 4. Database Schema

```sql
-- internal/db/schema.sql
-- SQLite schema. Switch to PostgreSQL when scaling.

CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id   INTEGER NOT NULL UNIQUE,
    username      TEXT DEFAULT '',
    lang          TEXT NOT NULL DEFAULT 'uz' CHECK(lang IN ('uz','ru')),
    level         TEXT NOT NULL DEFAULT 'a1' CHECK(level IN ('a1','a2','b1','b2','c1')),
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL REFERENCES users(id),
    date          TEXT NOT NULL,                -- '2026-06-06'
    count         INTEGER NOT NULL DEFAULT 0,
    UNIQUE(user_id, date)
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL UNIQUE REFERENCES users(id),
    plan          TEXT NOT NULL CHECK(plan IN ('weekly','monthly')),
    stars_amount  INTEGER NOT NULL,
    started_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at    DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS vocabulary (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL REFERENCES users(id),
    word          TEXT NOT NULL,
    translation   TEXT NOT NULL,
    example       TEXT NOT NULL,
    level         TEXT DEFAULT 'a1',
    review_count  INTEGER DEFAULT 0,
    next_review   DATETIME,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, word)
);

CREATE TABLE IF NOT EXISTS daily_progress (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL REFERENCES users(id),
    date          TEXT NOT NULL,
    messages_sent INTEGER DEFAULT 0,
    words_learned INTEGER DEFAULT 0,
    quizzes_taken INTEGER DEFAULT 0,
    UNIQUE(user_id, date)
);

CREATE INDEX idx_messages_user_date ON messages(user_id, date);
CREATE INDEX idx_vocab_user ON vocabulary(user_id);
CREATE INDEX idx_progress_user_date ON daily_progress(user_id, date);
```

---

## 5. Backend API

Base: `https://lingvo-ai-production.up.railway.app`

All endpoints except health require `X-Telegram-Init-Data` header.

### Endpoints

#### GET /api/health
Returns `{"status": "ok"}`.

#### POST /api/chat
Send message to AI tutor.

Request:
```json
{ "text": "I go to school yesterday" }
```

Response:
```json
{
  "reply": "Good! Tell me more about your day.",
  "corrections": [
    {
      "original": "I go",
      "corrected": "I went",
      "explanation_uz": "\"I go\" эмас, \"I went\". \"yesterday\" ўтган замонни кўрсатади.",
      "explanation_ru": "Нужно \"I went\". \"yesterday\" указывает на Past Simple.",
      "type": "grammar"
    }
  ],
  "usage": {
    "daily_used": 4,
    "daily_limit": 10,
    "is_premium": false
  }
}
```

429 (limit exceeded):
```json
{
  "error": "daily_limit_exceeded",
  "message_uz": "Бугунги лимит тугади. Чексизга обуна бўлинг!",
  "message_ru": "Дневной лимит исчерпан. Оформите подписку!",
  "premium_link": "tg://resolve?domain=lingvo_ai_bot&appname=app"
}
```

#### POST /api/chat/stream
SSE streaming chat. Same request as POST /api/chat. Response is a Server-Sent Events stream.

See [docs/streaming.md](docs/streaming.md).

#### POST /api/grammar
Grammar check without dialogue.

Request: `{ "text": "He don't like coffee" }`

Response: `{ "corrections": [...] }`

#### POST /api/vocab
Add word to dictionary. AI generates translation, example, level automatically.

Request: `{ "word": "ubiquitous" }`

Response:
```json
{
  "word": "ubiquitous",
  "translation_uz": "барча жойда учрайдиган",
  "examples": ["Smartphones are ubiquitous these days."],
  "level": "b2"
}
```

#### POST /api/vocab/lookup
Look up word without saving (can be used for preview).

Request: `{ "word": "ubiquitous" }`

Response: Same as POST /api/vocab.

#### GET /api/vocab
List user's words. Query: `?page=1&per_page=20&due_only=true`

Response: `{ "words": [...], "total": 42, "due_count": 12 }`

#### DELETE /api/vocab/:id
Delete a word from vocabulary.

#### GET /api/vocab/review
Get due words for review (SM-2). Query: `?limit=10`

Response: `[{ "id": 1, "word": "ubiquitous", "translation": "..." }]`

#### POST /api/vocab/review
Submit SM-2 review result. Request: `{ "word_id": 1, "quality": 4 }`

Response: `{ "next_review": "2026-06-08T00:00:00Z", "interval": 6, "ease": 2.5 }`

#### POST /api/quiz
Generate quiz (uses Gemini Flash-Lite). Request: `{ "topic": "past_simple", "count": 5 }`

Response: `{ "questions": [{"question": "...", "options": [...], "correct": 1, "explanation_uz": "..."}] }`

#### POST /api/level
Determine level (uses Gemini Flash-Lite). Request: `{ "answers": [...] }`

Response: `{ "level": "a2", "score": 65, "total": 100 }`

#### POST /api/level/save
Save level test result. Request: `{ "level": "a2" }`

#### GET /api/progress
User stats. Streak, total messages, words, level.

Response: `{ "level": "a2", "streak": 5, "total_messages": 42, "total_words": 30 }`

#### GET /api/progress/history
Daily progress history. Query: `?days=7`

#### GET /api/subscription
Subscription status.

Response: `{ "is_premium": false, "plan": "", "expires_at": "" }`

#### POST /api/create-invoice
Create Stars invoice. Request: `{ "plan": "monthly" }`
Response: `{ "invoice_link": "...", "stars": 800 }`

---

## 6. Telegram Bot

### Bot commands

```
start - 🚀 Запустить Lingvo AI
daily - 📅 Сегодняшний урок
help  - ℹ️ Помощь
```

### Handlers

- `/start` — Create user in DB, send welcome + Launch button (WebApp keyboard)
- `/daily` — Generate daily lesson via Gemini, send message with tasks
- `/help` — Show commands list
- `successful_payment` — Save subscription, confirm to user

### Welcome messages

UZ: "👋 Салом! Men Lingvo AI — sizning shaxsiy ingliz tili o'qituvchingiz."
RU: "👋 Привет! Я Lingvo AI — твой личный учитель английского."

---

## 7. Mini App (Frontend)

### Screens

- **Chat** — Message history, input with streaming, corrections display, limit indicator
  - Лимит не исчерпан: поле ввода + UsageIndicator (прогресс-бар)
  - Лимит исчерпан (free): поле ввода скрывается, показывается панель альтернатив:
    - 📚 Повторение слов (SM-2)
    - 📊 Тест уровня
    - ⭐ Получить безлимит (переход в Subscription)
  - Premium: поле ввода + UsageIndicator (⭐ Unlimited)
- **Vocabulary** — Word list with search, add word, review buttons
- **Progress** — Level, streak, daily chart
- **Subscription** — Plan comparison, subscribe buttons

### Telegram integration

Use `@telegram-apps/sdk` to get initData, expand app, detect theme.

### API client

See `web/src/api/client.ts`. Uses a shared `request<T>()` wrapper (not raw `fetch`) and `chatStream()` for SSE streaming.

```typescript
// All API calls go through request() helper:
async function request<T>(path: string, options?: RequestOptions): Promise<T>

// Streaming:
function chatStream(text: string, onToken: (t: string) => void, onDone: (r: AIResponse) => void, onError: (e: Error) => void): AbortController
```

---

## 8. AI Integration

### Gemini client wrapper (`internal/ai/gemini.go`)

- Creates genai.Client with API key
- Two models: Gemini 2.0 Flash (chat, grammar, vocab) and Gemini 2.0 Flash-Lite (quiz, level test), temperature 0.4
- Fallback: OpenAI-compatible HTTP client (`openai.go`) on Gemini error
- Builds prompts from `prompts.go`
- Parses JSON response into ChatResponse struct
- SM-2 algorithm (`sm2.go`) for spaced repetition scheduling

### Prompts (`internal/ai/prompts.go`)

Each endpoint has a buildXxxPrompt() function that returns a string prompt.

Chat prompt structure:
```
You are an AI English tutor. Level: {level}. Lang: {uz|ru}.
1. Reply naturally in English
2. Return JSON: { reply, corrections: [{original, corrected, explanation_uz, explanation_ru, type}] }
User message: {text}
```

### SM-2 Algorithm (`internal/ai/sm2.go`)

Implements the SM-2 spaced repetition algorithm for vocabulary review:
- Quality rating: 0-5 (0=complete blackout, 5=perfect response)
- Calculates next review date based on ease factor and interval
- Updates ease factor: EF' = EF + (0.1 - (5-q) × (0.08 + (5-q) × 0.02))
- Repetitions reset on quality < 3

### Streaming (`internal/ai/gemini.go`)

- ChatStream method returns `<-chan StreamEvent`
- Events: EventToken (partial text), EventResult (final parsed response), EventError
- Used by POST /api/chat/stream (SSE endpoint)

---

## 9. Payments & Monetization

### Telegram Stars plans

| Plan | Stars | Retail price | Net (mobile) | Net (desktop) |
|---|---|---|---|---|
| Weekly | 300 | ~$3.90 | ~$2.65 | ~$3.78 |
| Monthly | 800 | ~$10.40 | ~$7.07 | ~$10.09 |

### Reinvest loop

Earned Stars → Telegram Ads (30% bonus) → more users → more Stars.

### Free tier

10 AI messages per day. Premium = unlimited.

---

## 10. i18n / Localization

### Backend translations

Inline maps in Go code (no separate files). Strings: greetings, help text, plan names, stats messages.

Language detection: user.lang from DB → initData.language_code → default `uz`.

### Frontend translations

Files: `web/src/locales/uz.json`, `web/src/locales/ru.json`

Strings: all UI text (nav, buttons, labels, placeholders).

### Language detection

1. `user.lang` from DB
2. `initDataUnsafe.user.language_code` from Telegram
3. Default: `uz`

---

## 11. Deployment

### Build

```bash
# Backend
go build -o server ./cmd/server

# Frontend
cd web && pnpm build

# Run
./server  # or: go run ./cmd/server
```

### Environment variables

```
BOT_TOKEN=...
GEMINI_API_KEY=...
DATABASE_PATH=/data/lingvo.db
PORT=8080
WEBAPP_URL=https://lingvo-ai-production.up.railway.app
ADMIN_IDS=123,456
OPENAI_API_KEY=                       # optional, for AI fallback
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-4o-mini
```

### Platform

Railway.app (free tier: $5 credit/month). Go binary + frontend static files.

---

## 12. Error Handling

### HTTP error codes

| Code | Body error | Cause |
|---|---|---|
| 400 | invalid_request | Missing field |
| 400 | invalid_plan | Unknown subscription plan |
| 401 | unauthorized | Bad initData |
| 429 | daily_limit_exceeded | Free limit hit (includes premium_link) |
| 500 | ai_service_unavailable | AI service down |
| 500 | internal_error | DB/unknown |
| 503 | ai_service_unavailable | AI client not configured |

### Logging

Use `go.uber.org/zap` for structured logs. Log telegram IDs, error details, AI fallback usage.

### Graceful shutdown

Handle SIGINT/SIGTERM — currently not implemented.

---

## 13. Security

- Verify Telegram initData HMAC-SHA256 on every request (`internal/middleware/auth.go`)
- Bot token, Gemini key, OpenAI key, DB path — env only, never in code
- No passwords, no sessions (Telegram handles auth)
- Rate limit: 10 msgs/day per free user, unlimited for premium

---

## 14. Testing Strategy

- Go unit tests for: SM-2 algorithm (sm2_test.go), prompt building, hash verification, rate limit logic, chat streaming
- Frontend tests with Vitest: chat rendering, limit display, subscription flow
- Manual checklist: registration, limit, payment, vocab CRUD, progress

---

## Appendix: Cost estimation

Gemini 2.0 Flash per 1M tokens: $0.10 input / $0.40 output.

1 chat request: ~700 tokens → ~$0.00049

100 DAU × 5 req/day × 30 days → ~$7.35/month

1,000 DAU → ~$73.50/month → needs ~15 paying users at $7/mo to break even.
