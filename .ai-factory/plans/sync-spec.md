# План: Синхронизация SPEC.md с кодом

> Привести код в соответствие со SPEC.md.
> Приоритет: от критичных багов к косметике.

---

## Фаза 1: Баги и безопасность

### 1.1 Добавить rate limit на `POST /api/grammar` ✅

**Файлы:** `internal/api/router.go`
**Суть:** `grammarHandler` не имеет `rateMw`, хотя использует AI → может быть спам.
**Действие:** добавить `rateMw` в цепочку middleware на строке 41.

**Изменение:**
```go
// Было:
protected.POST("/grammar", grammarHandler(db, aiClient, sugar))
// Стало:
protected.POST("/grammar", rateMw, grammarHandler(db, aiClient, sugar))
```

---

### 1.2 Добавить `premium_link` в 429 response ✅

**Файлы:** `internal/middleware/ratelimit.go`
**Суть:** SPEC описывает `"premium_link"` в ответе 429, код не включает.
**Действие:** добавить `"premium_link": "https://t.me/{bot_username}/app"` в JSON ответа (или динамический invoice link при его наличии).

**Вопрос:** динамический link (через `createInvoiceLink`) или статический deep-link? На данном этапе — статический `https://t.me/{bot_username}/app`, так как SPEC не уточняет формат.

---

## Фаза 2: API-контракты

### 2.1 `POST /api/vocab` — AI-лукап + сохранение ✅

**Файлы:** `internal/api/vocab.go`, `web/src/api/client.ts`, `web/src/components/Vocabulary.tsx`
**SPEC:** request `{"word":"..."}` → AI лукап → сохранение → ответ `{word, translation_uz, translation_ru, examples, level}`
**Текущее:** request `{"word","translation","example","level"}` — все поля от пользователя

**Изменения:**
1. **Backend:** переделать `vocabAddHandler()`:
   - Принимает только `{"word":"..."}`
   - Вызывает AI (через уже существующий `vocabLookupHandler`/`BuildVocabPrompt`)
   - Сохраняет в БД с translation = translation_uz (так как в БД одно поле translation)
   - Возвращает `VocabLookupResponse`
2. **Frontend:** убрать форму ручного ввода, оставить только лукап → авто-сохранение

**Важно:** модель данных vocabulary хранит одно поле `translation` (не translation_uz/ru). Нужно либо:
   - Хранить translation_uz, а translation_ru игнорировать при сохранении
   - Либо расширить схему (но это нарушает SPEC — SPEC не описывает изменения схемы)

Выбираем: сохраняем `translation = translation_uz`, `example = examples[0]`. Это минимальное изменение.

---

### 2.2 `GET /api/vocab` — пагинация + новый формат ответа ✅

**Файлы:** `internal/api/vocab.go`, `internal/db/vocabulary.go`, `web/src/api/client.ts`, `web/src/components/Vocabulary.tsx`
**SPEC:** `?page=1&per_page=20&due_only=true` → `{"words":[], "total":N, "due_count":N}`
**Текущее:** голый массив без пагинации

**Изменения:**
1. **Backend:**
   - Добавить query params: `page` (int, default 1), `per_page` (int, default 20, max 100), `due_only` (bool, default false)
   - Новый response: `{"words": [...], "total": int, "due_count": int}`
   - `due_count` — количество слов, ожидающих повторения (через `GetDueWordCount`)
   - Обернуть в новый тип `VocabListResponse`
2. **Frontend:**
   - Адаптировать `getVocab()` под новый ответ
   - Обновить `Vocabulary.tsx` для работы с пагинацией (page + per_page из стейта)

---

### 2.3 `POST /api/vocab/review` — переезд на SPEC-формат ✅

**Файлы:** `internal/api/vocab.go`, `internal/api/router.go`, `web/src/api/client.ts`, `web/src/components/ReviewCard.tsx`
**SPEC:** единственный эндпоинт `POST /api/vocab/review` с телом `{"word_id": 1, "quality": 4}`
**Текущее:** два эндпоинта — `GET /api/vocab/review` (список due) и `POST /api/vocab/review/:id` (отправка оценки)

**Изменения:**
1. **Backend:**
   - Оставить `GET /api/vocab/review` без изменений (в SPEC нет отдельного "list due words", но он необходим для UX)
   - `POST /api/vocab/review` — новый эндпоинт, тело `{"word_id": int, "quality": int}` → вызывает `ProcessReview`, сохраняет, возвращает `{"next_review": "..."}`
   - `POST /api/vocab/review/:id` — удалить (SPEC не описывает)
   - Добавить в `router.go`: `protected.POST("/vocab/review", ..., vocabReviewSubmitHandler2(...))`

**Важно:** `GET /api/vocab/review` остаётся — SPEC его не описывает, но он критичен для UX (фронтенд должен знать, какие слова повторять).

---

## Фаза 3: Фичи

### 3.1 Gemini 2.0 Flash-Lite для quiz + level test

**Файлы:** `internal/ai/gemini.go`
**SPEC:** quiz и level test должны использовать `gemini-2.0-flash-lite` (дешевле)
**Текущее:** все запросы через `gemini-2.0-flash`

**Изменения:**
1. В `NewClient()` добавить возможность создавать модели с разными именами
2. В `ai.Client` добавить вторую модель `liteModel`
3. В `api/quiz.go` и `api/level.go` передавать контекст для lite-модели

---

### 3.2 Fallback GPT-4o-mini

**Файлы:** `internal/ai/gemini.go`, `cmd/server/main.go`, `docs/configuration.md`
**SPEC:** при недоступности Gemini использовать GPT-4o-mini
**Текущее:** не реализован

**Изменения:**
1. Добавить `OPENAI_API_KEY` в конфигурацию (уже есть в документации)
2. Реализовать Open AI клиент в `internal/ai/openai.go` (или в `gemini.go`)
3. В `Chat()` / `Generate()` при ошибке Gemini делать fallback-запрос к GPT
4. Логировать fallback (zap warn)

---

## Фаза 4: Чистка

### 4.1 Выделить `internal/bot/payments.go`

**Файлы:** `internal/bot/payments.go` (новый), `internal/bot/handlers.go`
**SPEC:** отдельный файл для payment handler
**Текущее:** `handlePayment` в `handlers.go`

**Изменение:** вынести `handlePayment`, `buildSubscriptionConfirmation`, `planTitle`, `sendMessage` в `payments.go`. Оставить в `handlers.go` только команды.

---

### 4.2 Удалить мёртвый код

1. `internal/bot/handlers.go` — удалить `durationMap` (не используется)
2. `internal/i18n/*.json`:
   - `daily_title` — подключить в `handleDaily` вместо хардкоженных строк
   - `subscription_active` / `subscription_expired` — подключить в `handlePayment`
3. `web/src/locales/*.json`:
   - `progress.last_7_days`, `progress.last_30_days` — использовать в `ProgressChart`
   - `vocab.add`, `vocab.review_due` — использовать или удалить

---

### 4.3 Обновить SPEC.md

Синхронизировать SPEC с реальностью:
- Убрать `vocab_get.go`, `payments.go` из директории (или оставить если создадим)
- Убрать Tailwind CSS из стека
- Убрать webhook из документации
- Обновить таблицу моделей AI
- Обновить пример ответа 429

---

### 4.4 Создать недостающие .md файлы

1. `docs/PROMPTS.md` — скопировать/отформатировать содержимое `internal/ai/prompts.go`
2. `.ai-factory/ROADMAP.md` — создать на основе истории коммитов

---

## Порядок выполнения

```
Фаза 1 (баги) ───→ Фаза 2 (API) ───→ Фаза 3 (фичи) ───→ Фаза 4 (чистка)
   1.1, 1.2          2.1, 2.2, 2.3       3.1, 3.2           4.1, 4.2, 4.3, 4.4
```

Каждая фаза сопровождается:
- `go build ./...` — проверка компиляции
- `go test ./...` — проверка тестов
- `cd web && npm run build` — проверка фронтенда
