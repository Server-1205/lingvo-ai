# Lingvo AI — AGENTS.md

> AI English Tutor Telegram Mini App. Узбекистан. Go + React + Gemini.
> Карта проекта для AI-агентов. Обновляйте при значительных изменениях структуры.

---

## Обзор проекта

Lingvo AI — AI-репетитор английского языка в Telegram. Пользователи общаются с AI (Gemini 2.0 Flash), получают исправления ошибок, пополняют словарь (SM-2), проходят тесты уровня и отслеживают прогресс. Бесплатно: 10 сообщений/день, премиум: безлимит за Telegram Stars.

Полное описание: `.ai-factory/DESCRIPTION.md`

## Технологический стек

- **Язык (backend):** Go 1.25
- **Фреймворк:** Gin
- **Язык (frontend):** TypeScript 6
- **Фреймворк (frontend):** React 19, Vite 8
- **База данных:** SQLite (modernc.org/sqlite, без CGO)
- **ORM:** sqlx
- **AI:** Gemini 2.0 Flash, fallback GPT-4o-mini
- **Бот:** go-telegram-bot-api v5 (long-polling)
- **Платежи:** Telegram Stars (createInvoiceLink)
- **Логирование:** zap (uber-go)
- **i18n:** react-i18next (frontend), JSON файлы (backend + frontend)
- **Аутентификация:** Telegram WebApp initData (HMAC-SHA256)

## Структура проекта

```
sinking/
├── .ai-factory/               # AI Factory конфигурация
│   ├── DESCRIPTION.md         #   описание проекта
│   ├── config.yaml            #   настройки AI Factory
│   ├── rules/base.md          #   базовые соглашения кода
│   ├── plans/                 #   планы реализации
│   ├── patches/               #   патчи эволюции
│   └── skill-context/         #   переопределения навыков
├── cmd/server/main.go         # точка входа (Go)
├── internal/
│   ├── bot/                   # Telegram Bot (long-polling)
│   ├── api/                   # REST API (Gin)
│   ├── ai/                    # Gemini + промпты + SM-2
│   ├── db/                    # SQLite + CRUD
│   ├── middleware/            # auth + ratelimit
│   ├── i18n/                  # backend переводы (uz, ru)
│   └── models/                # общие типы
├── web/
│   ├── src/
│   │   ├── components/        # React компоненты
│   │   ├── api/client.ts      # HTTP клиент
│   │   ├── hooks/             # React хуки
│   │   └── locales/           # frontend переводы (uz, ru)
│   ├── package.json
│   └── vite.config.ts
├── .opencode/
│   ├── agents/                # subagent конфигурации
│   ├── skills/                # кастомные навыки
│   └── opencode.json          # MCP + настройки
├── docs/                      # документация
├── SPEC.md                    # полная спецификация
├── ROADMAP.md                 # этапы и milestones
└── AGENTS.md                  # этот файл
```

## Ключевые точки входа

| Файл | Назначение |
|------|------------|
| `cmd/server/main.go` | Точка входа сервера: инициализация DB, роутера, бота |
| `internal/api/router.go` | Gin роутер + middleware |
| `internal/bot/bot.go` | Telegram long-polling цикл |
| `internal/ai/gemini.go` | Gemini клиент |
| `internal/db/schema.sql` | Схема SQLite |
| `internal/middleware/auth.go` | HMAC-SHA256 верификация initData |
| `internal/middleware/ratelimit.go` | Rate limit (10/day free) |
| `web/src/main.tsx` | React точка входа |
| `web/src/App.tsx` | Корневой компонент + маршрутизация |
| `web/src/api/client.ts` | HTTP API клиент |

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

## Документация

| Документ | Путь | Описание |
|----------|------|----------|
| README | `README.md` | Landing page проекта |
| README фронтенда | `web/README.md` | Документация фронтенда |
| Установка и запуск | `docs/getting-started.md` | Установка, настройка, первый запуск |
| API Reference | `docs/api.md` | Эндпоинты, форматы запросов/ответов |
| Конфигурация | `docs/configuration.md` | Переменные окружения |
| Streaming | `docs/streaming.md` | Документация AI streaming (SSE) |
| Спецификация | `SPEC.md` | Полная спецификация API, БД, промптов |
| Roadmap | `ROADMAP.md` | Этапы и текущий статус |
| Описание проекта | `.ai-factory/DESCRIPTION.md` | Описание проекта для AI-агентов |
| Архитектура | `.ai-factory/ARCHITECTURE.md` | Архитектурные гайдлайны |

## Файлы AI-контекста

| Файл | Назначение |
|------|------------|
| `AGENTS.md` | Карта проекта для AI-агентов |
| `.ai-factory/DESCRIPTION.md` | Детальное описание проекта |
| `.ai-factory/ARCHITECTURE.md` | Архитектурные решения и паттерны |
| `.opencode/memory/context.md` | Контекст сессии OpenCode |
| `.ai-factory/rules/base.md` | Базовые соглашения проекта |
| `.ai-factory/config.yaml` | Настройки AI Factory |
| `.opencode/opencode.json` | Конфигурация OpenCode (MCP, агенты) |

## Команды разработки

| Действие | Команда |
|---|---|
| Запустить backend | `go run ./cmd/server` |
| Собрать backend | `go build -o server ./cmd/server` |
| Запустить frontend (dev) | `cd web && pnpm dev` |
| Собрать frontend | `cd web && npm run build` |
| Установить Go зависимости | `go mod tidy` |
| Установить JS зависимости | `cd web && pnpm install` |

Порядок: **сначала `go mod tidy`, потом `pnpm install`, потом `go run ./cmd/server`**

## Чеклист перед коммитом

- [ ] `go build ./...` успешно
- [ ] `cd web && npm run build` успешно
- [ ] Нет секретов в коде (только env)
- [ ] `go mod tidy` если меняли зависимости
- [ ] Схема БД актуальна (schema.sql)

## Правила для AI-агентов

- Разделяйте shell-команды на отдельные шаги
  - ❌ Неправильно: `git checkout master && git pull`
  - ✅ Правильно: сначала `git checkout master`, затем `git pull origin master`
- При реализации React компонентов используйте `!isLoading` для кнопок, не `data &&`
- Все API-запросы через shared `request()` из `api/client.ts`, никогда через `fetch()`
- Новые API эндпоинты: добавляйте `Cache-Control: no-store` через `c.Header()`
- Ссылки: `docs/streaming.md` — AI Streaming, SPEC.md — полная спецификация
- Session Memory MCP настроен (`@modelcontextprotocol/server-memory`) — используй `create_entities`, `add_observations`, `search_nodes` для сохранения/извлечения контекста между сессиями. Memory file: `.opencode/memory/memory.jsonl`
