# Правила проекта — Lingvo AI

> Автоопределённые соглашения из анализа кодовой базы. Редактируйте по необходимости.

## Соглашения об именовании

- **Файлы (Go):** snake_case (`chat_stream.go`, `auth.go`)
- **Файлы (React):** PascalCase для компонентов (`Chat.tsx`, `NavBar.tsx`), camelCase для утилит (`client.ts`)
- **Переменные (Go):** camelCase (`userID`, `dbConn`)
- **Функции (Go — экспортируемые):** PascalCase (`NewHandler`, `GetUser`)
- **Функции (Go — неэкспортируемые):** camelCase (`verifyInitData`, `parseMessage`)
- **Переменные (TS):** camelCase (`isLoading`, `userData`)
- **Компоненты (React):** PascalCase (`Chat`, `NavBar`, `Subscription`)
- **Хуки (React):** `use` + PascalCase (`useTelegram`, `useTranslation`)
- **Акронимы (Go):** все заглавные (`HTTP`, `API`, `DB`, `URL`)

## Структура модулей

- **Go backend:** `internal/{ai,api,bot,db,middleware,i18n,models}/`
- **Frontend:** `web/src/{components,api,hooks,locales,types}/`
- **Точка входа:** `cmd/server/main.go`
- **Схема БД:** `internal/db/schema.sql`
- **Переводы:** `web/src/locales/{uz,ru}.json` (frontend), `internal/i18n/{uz,ru}.json` (backend)

## Обработка ошибок

- **Go:** `if err != nil { ... }` — стандартный паттерн Go
- **API:** структурированные JSON-ответы: `{"error": "message"}` с HTTP-кодами
- **HTTP коды:** 400 (invalid_request), 401 (unauthorized), 429 (daily_limit_exceeded), 500 (internal/ai error)
- **Middleware:** auth → ratelimit → handler (всегда верификация initData первой)
- **Cache-Control:** `c.Header("Cache-Control", "no-store")` на всех API-ответах

## Логирование

- **Go:** `go.uber.org/zap`, strykturirovanny (sugar.Infow, sugar.Errorw, sugar.Fatalw)
- **Формат:** всегда key-value пары: `sugar.Infow("message", "user_id", id, "error", err)`
- **Контекст:** включать `user_id`, `error`, `request_id` в каждую запись

## Тестирование

- **Go:** стандартный `testing` пакет, табличные тесты (table-driven tests)
- **Frontend:** Vitest + @testing-library/react, `describe`/`it`/`expect`
- **Проверка:** `go build ./...` и `cd web && npm run build` перед коммитом

## Импорты (Go)

1. Стандартная библиотека
2. Внешние пакеты (алфавитный порядок)
3. Внутренние пакеты (`github.com/lingvo-ai/lingvo/internal/...`)
- Группы разделены пустыми строками

## Импорты (TypeScript)

- Абсолютные импорты для библиотек
- Относительные импорты для внутренних модулей (`../api/client`, `./Component`)
- Named imports предпочтительнее default

## Git

- Conventional Commits: `feat(area): description`, `fix(area): description`, `chore: description`
- Ветки: `feature/<slug>`
- Базовая ветка: `master`
