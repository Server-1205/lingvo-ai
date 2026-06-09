# Продвинутый анализ ошибок (Advanced Error Analysis)

**Ветка:** `feature/advanced-error-analysis`
**Дата:** 2026-06-08
**Тип:** Enhancement (Premium Feature)

## Roadmap Linkage

- **Milestone:** "Premium Features"
- **Rationale:** Продвинутый анализ ошибок — ключевая премиум-функция, указанная в ROADMAP.md для Этапа 3

## Settings

| Параметр | Значение |
|----------|----------|
| Тесты | Да |
| Логирование | Verbose (DEBUG) |
| Документация | Нет |

## Краткое описание

Продвинутый анализ ошибок — премиум-функция для глубокого анализа грамматических/лексических ошибок пользователя. Поверх уже существующего premium_analysis (severity, category, learning_tip, rule_violated, OverallGrade) добавляет:

1. **История ошибок** — сохранение каждой correction в БД
2. **Статистика ошибок** — группировка по категориям, severity, датам
3. **Анализ паттернов** — какие ошибки повторяются, прогресс по категориям
4. **Улучшенный промпт** — передача истории ошибок в контекст AI

## Tasks

### Phase 1: Бэкенд — База данных

#### Task 1.1: Создать таблицу error_history в schema.sql

**Файл:** `internal/db/schema.sql`

Добавить новую таблицу:

```sql
CREATE TABLE IF NOT EXISTS error_history (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL REFERENCES users(id),
    message_id      INTEGER REFERENCES messages(id),
    original        TEXT NOT NULL,
    corrected       TEXT NOT NULL,
    category        TEXT NOT NULL DEFAULT 'grammar',
    severity        TEXT NOT NULL DEFAULT 'minor',
    rule_violated   TEXT DEFAULT '',
    learning_tip    TEXT DEFAULT '',
    context         TEXT DEFAULT '',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_error_user_date ON error_history(user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_error_category ON error_history(user_id, category);
```

**Логирование:** INFO при создании таблицы (в `db.go`, уже есть)

#### Task 1.2: Создать error_history CRUD в internal/db/

**Новый файл:** `internal/db/error_history.go`

Функции:
- `SaveError(ctx, db, userID, original, corrected, category, severity, ruleViolated, learningTip, context, messageID) error`
- `GetErrorHistory(ctx, db, userID, limit, offset) ([]ErrorHistoryEntry, error)`
- `GetErrorStats(ctx, db, userID) (*ErrorStats, error)` — группировка: `{category_counts: map[string]int, severity_counts: map[string]int, total: int, most_frequent_rules: []string, category_trend: []CategoryDayEntry}`
- `GetErrorCategoryTrend(ctx, db, userID, days) ([]CategoryDayEntry, error)` — ошибки по категориям за N дней

**Новые типы в `internal/models/types.go`:**
```go
type ErrorHistoryEntry struct {
    ID          int       `db:"id" json:"id"`
    UserID      int       `db:"user_id" json:"user_id"`
    MessageID   *int      `db:"message_id" json:"message_id,omitempty"`
    Original    string    `db:"original" json:"original"`
    Corrected   string    `db:"corrected" json:"corrected"`
    Category    string    `db:"category" json:"category"`
    Severity    string    `db:"severity" json:"severity"`
    RuleViolated string   `db:"rule_violated" json:"rule_violated"`
    LearningTip  string   `db:"learning_tip" json:"learning_tip"`
    Context     string    `db:"context" json:"context"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type ErrorStats struct {
    TotalErrors     int            `json:"total_errors"`
    CategoryCounts  map[string]int `json:"category_counts"`
    SeverityCounts  map[string]int `json:"severity_counts"`
    MostFrequentRules []string     `json:"most_frequent_rules"`
    CategoryTrend   []CategoryDayEntry `json:"category_trend,omitempty"`
}

type CategoryDayEntry struct {
    Date     string `json:"date"`
    Grammar  int    `json:"grammar"`
    Vocabulary int  `json:"vocabulary"`
    Spelling int    `json:"spelling"`
    WordOrder int   `json:"word_order"`
    Punctuation int `json:"punctuation"`
}
```

**Логирование:**
- DEBUG: `"error history saved" user_id, category, severity`
- WARN: `"error history save failed" user_id, error`

### Phase 2: Бэкенд — AI и API

#### Task 2.1: Обновить chat.go и chat_stream.go для сохранения ошибок

**Файлы:** `internal/api/chat.go`, `internal/api/chat_stream.go`

После получения ответа от AI:
1. Если есть corrections — сохранить каждую в `error_history` через `db.SaveError()`
2. Передавать `messageID` если есть
3. Сохранять контекст (оригинальное сообщение пользователя, обрезанное до 500 символов)

Для premium-пользователей дополнительно: передавать `context` с полным сообщением пользователя.
Для free — `context: ""`.

**Логирование:**
- DEBUG: `"saved N corrections to error history" user_id`
- WARN: `"failed to save error correction" user_id, original, error`

#### Task 2.2: Создать эндпоинт GET /api/errors/history

**Новый файл:** `internal/api/errors.go`

```go
func errorsHandler(database *sqlx.DB, sugar *zap.SugaredLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // GET /api/errors/history — пагинированная история
        // GET /api/errors/stats — статистика ошибок
    }
}
```

**Эндпоинты:**
- `GET /api/errors/history?page=1&per_page=20&category=grammar` — возвращает `{entries: [...], total: N}`
- `GET /api/errors/stats?days=30` — возвращает `ErrorStats`

**Доступ:** только premium (проверка `is_premium` из контекста).
Для free-юзеров: `403 {"error": "premium_required"}`

**Ответ /stats:**
```json
{
  "total_errors": 127,
  "category_counts": {"grammar": 80, "vocabulary": 25, "spelling": 15, "word_order": 5, "punctuation": 2},
  "severity_counts": {"critical": 12, "major": 45, "minor": 70},
  "most_frequent_rules": ["Subject-Verb Agreement", "Past Simple vs Present Perfect"],
  "category_trend": [
    {"date": "2026-06-01", "grammar": 5, "vocabulary": 2, "spelling": 1, "word_order": 0, "punctuation": 0},
    ...
  ]
}
```

**Логирование:**
- DEBUG: `"error stats requested" user_id`
- WARN: `"error history query failed" user_id, error`

#### Task 2.3: Зарегистрировать эндпоинты в router.go

**Файл:** `internal/api/router.go`

Добавить:
```go
protected.GET("/errors/history", rateMw, errorsHistoryHandler(database, sugar))
protected.GET("/errors/stats", rateMw, errorsStatsHandler(database, sugar))
```

После регистрации authMw.

#### Task 2.4: Улучшить премиум-промпт с контекстом истории ошибок

**Файл:** `internal/ai/prompts.go`

Создать новую функцию:
```go
func BuildPremiumChatPromptWithHistory(level, lang, text string, recentErrors []string) string
```

Где `recentErrors` — строки вида "1. Past Simple (2 times) — 'I go yesterday'", "2. Subject-Verb (3 times) — 'He don't like'"

Промпт включает:
- Все правила обычного premium-промпта
- Дополнительная секция: "The user has been making these errors recently. Pay special attention to them and provide extra detailed explanations when these errors occur."
- Список недавних ошибок пользователя (top 5 по частоте)

В `chat.go` и `chat_stream.go`:
- Для premium: достать топ-5 частых ошибок через `db.GetMostFrequentErrors(ctx, db, userID, 5)`
- Использовать `BuildPremiumChatPromptWithHistory` вместо `BuildPremiumChatPrompt`

**Логирование:**
- DEBUG: `"included N recent errors in prompt context" user_id, errors_count`

### Phase 3: Фронтенд — Error Dashboard

#### Task 3.1: Создать компонент ErrorDashboard

**Новый файл:** `web/src/components/ErrorDashboard.tsx`

Компонент для отображения статистики ошибок. Состоит из:
1. **Total errors counter** — общее количество ошибок (большая цифра)
2. **Category breakdown** — карточки с количеством ошибок по каждой категории (грамматика, лексика, орфография, порядок слов, пунктуация) с цветовой индикацией
3. **Severity distribution** — круговая диаграмма или столбцы critical/major/minor
4. **Most frequent rules** — список самых частых нарушаемых правил
5. **Category trend chart** — BarChart (recharts) с количеством ошибок по дням с разбивкой по категориям

Использует `isLoading &&` pattern для загрузки данных.

**API функции в `client.ts`:**
```typescript
export async function getErrorStats(days?: number): Promise<ErrorStatsResponse>
export async function getErrorHistory(page?: number, perPage?: number, category?: string): Promise<ErrorHistoryResponse>
```

**Логирование:** console.debug для каждого запроса к API.

#### Task 3.2: Обновить App.tsx и NavBar для вкладки Error Analysis

**Файлы:** `web/src/App.tsx`, `web/src/components/NavBar.tsx`

- Добавить новую вкладку `errors` в NavBar (только для premium)
- В `App.tsx`: рендерить `ErrorDashboard` при `activeTab === 'errors'`
- Иконка: 📊 или диаграмма

#### Task 3.3: Добавить i18n ключи для Error Dashboard

**Файлы:** `web/src/locales/uz.json`, `web/src/locales/ru.json`

Добавить ключи:
```json
{
  "errors": {
    "title": "Ошибки / Xatolar",
    "total_errors": "Всего ошибок / Jami xatolar",
    "by_category": "По категориям / Kategoriyalar bo'yicha",
    "by_severity": "По серьёзности / Jiddiyligi bo'yicha",
    "frequent_rules": "Частые правила / Tez-tez uchraydigan qoidalar",
    "trend": "Динамика / Dinamika",
    "premium_only": "Только для Premium / Faqat Premium uchun",
    "no_errors": "Нет ошибок / Xatolar yo'q",
    "grammar": "Грамматика / Grammatika",
    "vocabulary": "Лексика / Leksika",
    "spelling": "Орфография / Imlo",
    "word_order": "Порядок слов / So'z tartibi",
    "punctuation": "Пунктуация / Tinish belgilari",
    "critical": "Критические / Kritik",
    "major": "Важные / Muhim",
    "minor": "Незначительные / Kichik"
  }
}
```

#### Task 3.4: Обновить Progress.tsx с интеграцией ошибок

**Файл:** `web/src/components/Progress.tsx`

Для premium-пользователей добавить блок "Recent Errors Summary" в дашборд прогресса:
- Показывать последние 5 ошибок
- Ссылка на полный анализ (переключение на вкладку errors)

### Phase 4: Тесты

#### Task 4.1: Go unit-тесты для error_history

**Новый файл:** `internal/db/error_history_test.go`

Табличные тесты:
- `TestSaveError` — сохранение одной записи
- `TestGetErrorHistory` — получение истории с пагинацией
- `TestGetErrorStats` — проверка группировки по категориям
- `TestGetErrorCategoryTrend` — проверка тренда за N дней

#### Task 4.2: Go unit-тесты для эндпоинтов

**Новый файл:** `internal/api/errors_test.go`

Тесты:
- `TestErrorsHistoryHandler` — проверка ответа, пагинации, premium gate
- `TestErrorsStatsHandler` — проверка статистики
- `TestErrorsHistoryPermission` — free user получает 403

#### Task 4.3: Vitest тесты для ErrorDashboard

**Новый файл:** `web/src/components/__tests__/ErrorDashboard.test.tsx`

Тесты:
- Рендеринг компонента
- Отображение загрузки (isLoading)
- Отображение данных (total, categories, severity)
- Отображение "premium only" для не-premium
- Отображение "no errors" для пустого списка

## Commit Plan

| # | Коммит | Задачи |
|---|--------|--------|
| 1 | `feat(db): add error_history table and CRUD` | 1.1, 1.2 |
| 2 | `feat(api): save error corrections to history` | 2.1 |
| 3 | `feat(api): add /api/errors/history and /api/errors/stats endpoints` | 2.2, 2.3 |
| 4 | `feat(ai): improve premium prompt with error history context` | 2.4 |
| 5 | `feat(web): add ErrorDashboard component` | 3.1 |
| 6 | `feat(web): integrate error analysis in NavBar and Progress` | 3.2, 3.3, 3.4 |
| 7 | `test: add tests for error history and dashboard` | 4.1, 4.2, 4.3 |
