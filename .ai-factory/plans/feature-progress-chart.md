# Implementation Plan: 📊 График прогресса

Branch: feature/progress-chart
Created: 2026-06-07

## Settings
- Testing: yes
- Logging: verbose
- Docs: yes

## Roadmap Linkage
Milestone: "Этап 2: Vocabulary + Level Test (Дни 6-8) — Progress Dashboard"
Rationale: График ежедневной активности — единственный незакрытый пункт в прогресс-дашборде

## Commit Plan
- **Commit 1** (task 1): "feat(api): add GET /api/progress/history endpoint with daily chart data"
- **Commit 2** (tasks 2-3): "feat(web): add recharts progress chart component with daily activity bars"

## Tasks

### Phase 1: Backend — новый эндпоинт для истории

- [x] Task 1: Добавить эндпоинт `GET /api/progress/history?days=7`

  **Backend:**
  - `internal/api/progress.go` — новый handler `progressHistoryHandler`, принимает query `days` (default 7, макс 30)
  - `internal/db/progress.go` — новая функция `GetProgressHistory(ctx, db, userID, days)` — запрос к `daily_progress` за последние N дней
  - `internal/models/types.go` — новый тип `DailyProgressEntry{Date, MessagesSent, WordsLearned, QuizzesTaken}` и `ProgressHistoryResponse{Entries []DailyProgressEntry}`
  - `internal/api/router.go` — добавить маршрут `GET /progress/history`

  **Frontend:**
  - `web/src/api/client.ts` — новый тип `DailyProgressEntry`, `ProgressHistoryResponse`, функция `getProgressHistory(days: number)`
  - Обновить тип `getProgress` если нужно

  **LOGGING REQUIREMENTS:**
  - Логировать вход в handler с user_id, days
  - Логировать количество записей в ответе
  - Логировать ошибки БД
  - `sugar.Debugw("progress history requested", "user_id", uid, "days", days)`
  - `sugar.Debugw("progress history response", "entries", len(entries))`

  **Tests:**
  - `internal/db/progress_test.go` — тест `GetProgressHistory`: вставка тестовых данных, проверка фильтрации по дням
  - `internal/api/progress_test.go` — тест `progressHistoryHandler`: проверка query параметров, формата ответа

### Phase 2: Frontend — компонент графика

- [x] Task 2: Установить библиотеку recharts и создать компонент ProgressChart

  **Files:**
  - `web/package.json` — добавить зависимость `recharts`
  - `web/src/components/ProgressChart.tsx` — новый компонент с bar chart:
    - Принимает проп `data: DailyProgressEntry[]`
    - `<ResponsiveContainer width="100%" height={200}>`
    - `<BarChart data={data}>` с `<Bar dataKey="messages_sent" fill="var(--tg-button)">`
    - `<XAxis dataKey="date">` с форматированием дат (дд.мм)
    - `<Tooltip>` с кастомным содержимым
    - Telegram theme-aware цвета
    - Toggle: 7 дней / 30 дней

  **LOGGING REQUIREMENTS:**
  - `console.debug('[progress-chart] rendering with N data points', data.length)`
  - `console.debug('[progress-chart] date range:', firstDate, '-', lastDate)`

- [x] Task 3: Интегрировать ProgressChart в Progress.tsx

  **Files:**
  - `web/src/components/Progress.tsx`:
    - Добавить `useQuery` для `getProgressHistory(7)`
    - Под cards добавить `<ProgressChart data={history} />`
    - Кнопки переключения "7 дней" / "30 дней"
    - При смене периода — новый запрос через `queryKey`
  - `web/src/locales/uz.json` — добавить `progress.last_7_days`, `progress.last_30_days`
  - `web/src/locales/ru.json` — добавить `progress.last_7_days`, `progress.last_30_days`

  **LOGGING REQUIREMENTS:**
  - `console.debug('[progress] chart mounted, period=7')`
  - `console.debug('[progress] chart updated to period=30')`
  - `console.debug('[progress] history loaded', data.length, 'entries')`

  **Tests:**
  - `web/src/components/__tests__/Progress.test.tsx`:
    - Тест рендера чарта с тестовыми данными
    - Тест переключения 7/30 дней
    - Тест пустого состояния (нет данных)
