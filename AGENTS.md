# Lingvo AI — AGENTS.md

> AI English Tutor Telegram Mini App. Узбекистан. Go + React + Gemini.

---

## Быстрый старт

```bash
# Backend
go mod init github.com/yourhandle/lingvo-ai
go mod tidy
go run ./cmd/server

# Frontend
cd web && npm install && npm run dev
```

## Команды разработки

| Действие | Команда |
|---|---|
| Запустить backend | `go run ./cmd/server` |
| Собрать backend | `go build -o server ./cmd/server` |
| Запустить frontend (dev) | `cd web && pnpm dev` |
| Собрать frontend | `cd web && npm run build` |
| Установить Go зависимости | `go mod tidy` |
| Установить JS зависимости | `cd web && pnpm install` |

Порядок: **сначала go mod tidy, потом pnpm install, потом go run ./cmd/server**

## Структура проекта

```
sinking/
├── cmd/server/main.go       — точка входа
├── internal/
│   ├── bot/                 — Telegram Bot API
│   ├── api/                 — REST API для Mini App
│   ├── ai/                  — Gemini интеграция + промпты
│   ├── db/                  — SQLite + CRUD
│   └── middleware/          — auth + ratelimit
├── web/src/
│   ├── components/          — React компоненты
│   ├── api/client.ts        — HTTP клиент к backend
│   └── locales/             — UZ + RU переводы
├── SPEC.md                  — полная техническая спецификация
├── ROADMAP.md               — этапы и milestones
└── AGENTS.md                — этот файл
```

## Важные конвенции

### Языки
- UI: узбекский (лат.) + русский
- Определение: `user.lang` из БД → `initData.language_code` → `uz` по умолч.
- i18n: react-i18next (frontend), JSON файлы в internal/i18n/ и web/src/locales/

### База данных
- SQLite через `modernc.org/sqlite` (pure Go, без CGO)
- Схема в `internal/db/schema.sql`, выполняется при старте
- Миграции: если надо изменить схему — изменить schema.sql, добавить версионирование

### Платежи
- Telegram Stars (встроенная валюта Telegram)
- `createInvoiceLink` через Bot API
- Обработка `successful_payment` в bot/handlers
- Комиссия ~32% (Apple/Google), реинвестировать Stars в Telegram Ads (30% бонус)

### AI
- Модель: Gemini 2.0 Flash (температура 0.4)
- Fallback: GPT-4o-mini
- Промпты в `internal/ai/prompts.go`
- Парсинг ответа — JSON (reply + corrections[])

### Аутентификация
- Telegram WebApp initData (HMAC-SHA256)
- Верификация в `internal/middleware/auth.go`
- Передача в заголовке `X-Telegram-Init-Data`

### Rate limit
- Free: 10 сообщений/день (DATE-based, UTC)
- Premium: безлимит (проверка subscriptions.expires_at)
- Проверка в `internal/middleware/ratelimit.go`

## Чеклист перед коммитом
- [ ] `go build ./...` успешно
- [ ] `cd web && npm run build` успешно
- [ ] Нет секретов в коде (только env)
- [ ] `go mod tidy` если меняли зависимости
- [ ] Схема БД актуальна (schema.sql)

## Ссылки
- SPEC.md — полная спецификация API, БД, промптов
- ROADMAP.md — этапы и текущий статус
- docs/streaming.md — Streaming AI responses (SSE)
- Telegram Mini App docs: https://core.telegram.org/bots/webapps
- Gemini API: https://ai.google.dev/gemini-api/docs
