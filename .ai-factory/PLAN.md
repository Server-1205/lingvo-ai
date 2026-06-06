# Frontend UI — Lingvo AI Mini App

**Created:** 2026-06-06
**Mode:** Fast
**Testing:** Yes
**Logging:** Verbose

---

## Tasks

### Phase 1: Foundation

#### Task 1: i18n + Telegram SDK + QueryClient setup

**Status: ✅ Done**

**Files:**
- `web/src/main.tsx` — инициализация i18next (uz/ru), Telegram SDK init(), QueryClientProvider
- `web/src/i18n.ts` — конфиг i18next с import ресурсов

**Что сделать:**
- Создать `web/src/i18n.ts` с i18next + react-i18next, загрузить uz.json/ru.json
- В `main.tsx`: инициализировать i18n, Telegram init(), обернуть `<App/>` в `<QueryClientProvider>`
- Учесть `verbatimModuleSyntax: true` — `import type` для типов

**Проверка:** `npm run build` без ошибок

#### Task 2: API client

**Status: ✅ Done**

**Files:**
- `web/src/api/client.ts`

**Что сделать:**
- Создать типизированный API клиент: `api<T>(path, body?)` → `Promise<T>`
- Автоматически добавлять `X-Telegram-Init-Data` из initData.raw()
- Базовый URL из `import.meta.env.VITE_API_URL || ''`
- Обработка ошибок: если `!res.ok` — парсить JSON и кидать Error с сообщением
- Экспорт типов: `ChatRequest`, `ChatResponse`, `Correction`, `Usage`, `SubscriptionResponse`, `InvoiceRequest`, `InvoiceResponse`

#### Task 3: useTelegram hook

**Status: ✅ Done**

**Files:**
- `web/src/hooks/useTelegram.ts`

**Что сделать:**
- Хук, возвращающий `{ user, initData, theme, isReady }`
- Читать `initData.state()` и `themeParams.state()` из SDK
- Возвращать `isReady: boolean` — true после монтирования

#### Task 4: CSS + index.html update

**Status: ✅ Done**

**Files:**
- `web/src/index.css` — заменить на Lingvo AI стили
- `web/index.html` — обновить title, lang, meta
- `web/src/App.css` — удалить (перенести нужное в index.css)

**Что сделать:**
- Использовать CSS custom properties с Telegram theme (Telegram SDK vars)
- Mobile-first, Telegram Mini App (390px)
- Стили для: чат-пузырей, навбара, карточек подписки, индикатора лимита
- Темы: светлая/тёмная через Telegram SDK `themeParams`

---

### Phase 2: Core Components

#### Task 5: NavBar component

**Status: ✅ Done**

**Files:**
- `web/src/components/NavBar.tsx`
- `web/src/components/NavBar.css`

**Что сделать:**
- Нижняя навигация с 4 табами: Chat, Vocabulary, Progress, Subscription
- Иконки через SVG/emoji, подпись из i18n
- Active tab выделен цветом `themeParams.buttonColor`
- Проп: `active: string, onTabChange: (tab: string) => void`

#### Task 6: LanguageSwitcher component

**Status: ✅ Done**

**Files:**
- `web/src/components/LanguageSwitcher.tsx`

**Что сделать:**
- Кнопка переключения UZ ↔ RU
- Вызов `i18n.changeLanguage()`
- Отображение текущего языка

#### Task 7: GrammarBlock component

**Status: ✅ Done**

**Files:**
- `web/src/components/GrammarBlock.tsx`
- `web/src/components/GrammarBlock.css`

**Что сделать:**
- Отображение списка исправлений: original → corrected
- Тип ошибки (grammar/vocabulary/spelling)
- Объяснение на языке пользователя (explanation_uz/ru)
- Проп: `corrections: Correction[]`

#### Task 8: UsageIndicator component

**Status: ✅ Done**

**Files:**
- `web/src/components/UsageIndicator.tsx`

**Что сделать:**
- Показывает `"Бугун: 4/10"` или `"Чексиз"` для Premium
- Прогресс-бар (used/limit)
- Если лимит исчерпан — красный цвет
- Проп: `usage: Usage`

---

### Phase 3: Screens

#### Task 9: Chat screen

**Status: ✅ Done**

**Files:**
- `web/src/components/Chat.tsx`
- `web/src/components/Chat.css`

**Что сделать:**
- Поле ввода + кнопка отправки
- История сообщений (user ↔ AI) в виде пузырей
- Отображение GrammarBlock под AI-ответами
- UsageIndicator внизу
- Использовать `@tanstack/react-query` для мутации
- При отправке: POST /api/chat { text }
- Отображать loading состояние
- DEBUG-логи в консоль: "chat: sending...", "chat: response received"

#### Task 10: Subscription screen

**Status: ✅ Done**

**Files:**
- `web/src/components/Subscription.tsx`
- `web/src/components/Subscription.css`

**Что сделать:**
- Три карточки: Free (10/день), Weekly (300⭐), Monthly (800⭐)
- Кнопка "Обуна бўлиш" / "Подписаться"
- При клике на платный план: POST /api/create-invoice → открыть ссылку
- GET /api/subscription — показать текущий статус
- Telegram theme-aware стили

#### Task 11: Vocabulary screen

**Status: ✅ Done**

**Files:**
- `web/src/components/Vocabulary.tsx`
- `web/src/components/Vocabulary.css`

**Что сделать:**
- Экран-заглушка с текстом "Ҳали сўзлар йўқ" / "Слов пока нет"
- Кнопка "+ Добавить" (заглушка — API ещё не готов)
- Подготовить структуру для будущего списка слов

#### Task 12: Progress screen

**Status: ✅ Done**

**Files:**
- `web/src/components/Progress.tsx`
- `web/src/components/Progress.css`

**Что сделать:**
- Экран-заглушка с отображением уровня (из user)
- Статистика: сообщения, слова, streak
- Подготовить структуру для будущих графиков

---

### Phase 4: Integration

#### Task 13: App.tsx — root layout + navigation

**Status: ✅ Done**

**Files:**
- `web/src/App.tsx`

**Что сделать:**
- Tab-маршрутизация через useState (Chat / Vocabulary / Progress / Subscription)
- Верхняя панель: LanguageSwitcher + title
- Нижняя панель: NavBar
- Основная область: рендер выбранного экрана
- Telegram theme integration через useTelegram

#### Task 14: Tests

**Status: ✅ Done**

**Files:**
- `web/src/components/__tests__/Chat.test.tsx`
- `web/src/components/__tests__/NavBar.test.tsx`
- `web/src/components/__tests__/UsageIndicator.test.tsx`
- `web/vitest.config.ts`

**Что сделать:**
- Установить vitest + @testing-library/react + jsdom
- Создать vitest.config.ts
- Тест NavBar: рендер 4 табов, клик переключает active
- Тест UsageIndicator: рендер "4/10", рендер "Premium"
- Тест Chat: рендер формы, вызов API при отправке (mock)

---

## Commit Plan

| # | Задачи | Сообщение |
|---|--------|-----------|
| 1 | Tasks 1-4 (Foundation) | `feat(web): add app foundation — i18n, api client, telegram hooks, css` |
| 2 | Tasks 5-8 (Core Components) | `feat(web): add core components — NavBar, LanguageSwitcher, GrammarBlock, UsageIndicator` |
| 3 | Tasks 9-10 (Screens) | `feat(web): add Chat and Subscription screens` |
| 4 | Tasks 11-12 (Placeholder screens) | `feat(web): add Vocabulary and Progress placeholder screens` |
| 5 | Tasks 13-14 (Integration + Tests) | `feat(web): integrate layout, add navigation and tests` |
