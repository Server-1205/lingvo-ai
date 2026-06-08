# Ежедневный урок (Daily Lesson)

**Branch:** `feature/daily-lesson`
**Created:** 2026-06-07 15:00

## Settings

- **Testing:** yes
- **Logging:** verbose
- **Docs:** yes — добавить в docs/api.md
- **Roadmap Linkage:** Фичи V2 → Ежедневный урок

## Концепция

Бот генерирует персонализированный 5-минутный урок английского на основе уровня пользователя.
Урок приходит по команде `/daily` и по расписанию (раз в день).

## Архитектура

```
User → /daily command (or reminder scheduler)
  → bot handler → AI prompt → lesson JSON
  → send to user as formatted message
  → add new vocabulary words to DB
  → track lesson completion in daily_progress
```

**Модель ответа AI:**
```json
{
  "topic": "Past Simple vs Present Perfect",
  "explanation_uz": "...",
  "explanation_ru": "...",
  "examples": ["example 1", "example 2"],
  "exercises": [
    {"question": "I ___ (see) that movie yesterday.", "answer": "saw", "options": ["saw", "seen", "have seen"]}
  ],
  "vocabulary": [
    {"word": "yesterday", "translation_uz": "kecha", "translation_ru": "vchera"}
  ]
}
```

## Задачи

### Task 1: Создать промпт `BuildDailyLessonPrompt`

**Файлы:** `internal/ai/prompts.go`

- Новый промпт для генерации урока на основе уровня, языка и topic (опционально)
- Возвращает JSON с topic, explanation_{uz,ru}, examples[], exercises[], vocabulary[]
- Инструкция: "Generate a 5-minute English lesson personalized to level {level}"

### Task 2: Создать `internal/ai/daily.go` — модель DailyLesson

**Файлы:** `internal/ai/daily.go`

- Тип `DailyLessonResponse` с полями: Topic, ExplanationUz, ExplanationRu, Examples, Exercises, Vocabulary
- Тип `Exercise` с полями: Question, Answer, Options
- Тип `LessonVocab` с полями: Word, TranslationUz, TranslationRu

### Task 3: Создать `internal/api/daily.go` — API handler

**Файлы:** `internal/api/daily.go`, `internal/api/daily_test.go`

- `POST /api/daily` — генерирует урок для текущего пользователя
- Принимает: `{}` (использует level+lang из контекста)
- Возвращает: `DailyLessonResponse`
- Требует auth middleware
- Вызов `aiClient.Generate()` с `BuildDailyLessonPrompt`

### Task 4: Обновить bot handler `/daily`

**Файлы:** `internal/bot/handlers.go`, `internal/bot/bot.go`

- `handleDaily` — вызывает AI, форматирует результат в читаемое сообщение
- Показывает тему, объяснение, примеры
- Кнопка «Начать упражнения» → открывает Mini App с уроком
- Кнопка «Пропустить»

### Task 5: Обновить reminder scheduler

**Файлы:** `internal/bot/reminder.go`

- Добавить reminder на ежедневный урок (не только для слов)
- Проверять, сделал ли пользователь урок сегодня
- Если нет — слать напоминание с уроком

### Task 6: Интеграция в router.go

**Файлы:** `internal/api/router.go`

- Добавить маршрут `protected.POST("/daily", dailyHandler(db, aiClient, sugar))`

### Task 7: Документация

**Файлы:** `docs/api.md`

- Добавить `POST /api/daily` с параметрами и ответом

### Task 8: Тесты

**Файлы:** `internal/api/daily_test.go`

- `TestDailyHandler_Success` — мок AI, проверить ответ
- `TestDailyHandler_Failure` — AI ошибка → 500
- `TestDailyHandler_NoAuth` — без auth → 401

## Commit Plan

1. `feat(ai): add daily lesson prompt and response model`
2. `feat(api): add POST /api/daily endpoint`
3. `feat(bot): update /daily command with AI lesson`
4. `feat(bot): add daily lesson reminder to scheduler`
5. `docs: add daily lesson API documentation`
