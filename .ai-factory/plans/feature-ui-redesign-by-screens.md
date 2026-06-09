# UI Redesign — Lingvo AI (по экранам)

**Branch:** `feature/ui-redesign-by-screens`
**Created:** 2026-06-08
**Source:** Google Stitch MCP — проект `Lingvo AI Tutor` (ID: `13746191202433220969`)
**Design System:** Teal #00685f + Plus Jakarta Sans + Fidelity

## Settings

- **Testing:** yes
- **Logging:** verbose
- **Docs:** yes
- **Roadmap Linkage:** Milestone: "Редизайн", Rationale: "Полный редизайн UI через Google Stitch для подготовки к маркетинговому запуску"

## Roadmap Linkage

Milestone: "Редизайн"
Rationale: "Полный редизайн всех экранов через Google Stitch MCP с заменой React-компонентов и CSS-токенов"

---

## Структура

Каждый экран = отдельная фаза. Фазы 0-2 (фундамент) обязательны перед остальными. Фазы 3-9 независимы и могут выполняться в любом порядке.

```
Фаза 0 (CSS-токены + шрифты)
  └──> Фаза 1 (NavBar)
         └──> Фаза 2 (App.tsx — хедер + лайаут)
                ├──> Фаза 3 (Chat + Grammar + Usage)
                ├──> Фаза 4 (Vocabulary + ReviewCard)
                ├──> Фаза 5 (Progress + Chart)
                ├──> Фаза 6 (Subscription)
                ├──> Фаза 7 (LevelTest)
                ├──> Фаза 8 (Onboarding)
                └──> Фаза 9 (DailyLesson)
                       └──> Фаза 10 (Docs + Build)
```

---

## Фаза 0: CSS-токены и шрифты

### Задача 0.1: Создать CSS-токены из Stitch Design System
- [x] Создан `web/src/styles/tokens.css` с полной палитрой Stitch DS

**Файлы:**
- `web/src/styles/tokens.css` — новый файл с CSS custom properties
- `web/src/index.css` — импорт tokens.css, обновление глобальных стилей

**Что сделать:**
- Создать `web/src/styles/tokens.css` с переменными из Stitch DS:
  - Цвета: `--c-primary: #00685f`, `--c-primary-container: #008378`, палитра surface/on-surface/outline
  - Шрифт: `--font-family: 'Plus Jakarta Sans', sans-serif`
  - Скругления: `--round-sm: 4px`, `--round-md: 8px`, `--round-lg: 12px`, `--round-xl: 16px`
  - Отступы: `--space-unit: 4px`, `--gutter: 16px`, `--card-padding: 20px`, `--nav-height: 64px`
  - Тени: `--shadow-sm`, `--shadow-md` по Stitch
- Обновить `index.css`:
  - Импортировать `tokens.css`
  - Переопределить Telegram theme variables (`--tg-bg`, `--tg-text` и т.д.) на Stitch-токены
  - Обновить классы `.card`, `.btn`, `.scroll-area`, `.section-title` под новый дизайн
  - Добавить классы для Telegram-safe areas
- Удалить старые переменные (--c-primary: #00BFA5 и т.д.)

### Задача 0.2: Добавить Plus Jakarta Sans
- [x] Google Fonts link добавлен в `web/index.html`

---

## Фаза 1: NavBar

### Задача 1.1: Полностью заменить NavBar на Stitch-дизайн
- [x] NavBar переписан: stroke-иконки, blur, dot-индикатор

**Файлы:**
- `web/src/components/NavBar.tsx` — полная замена JSX

**Что сделать:**
- Фиксированный bottom bar: `position: fixed`, `bottom: 0`, `backdrop-filter: blur(10px)`
- Высота: 64px + `var(--safe-bottom)`
- 4 таба: Chat, Vocabulary, Progress, Subscription
- Иконки: 24px stroke-based (SVG или Unicode-символы)
- Активный таб: `--c-primary` цвет + dot-индикатор 4px под иконкой
- Неактивные: `--c-outline` или `#9E9E9E`
- Лейблы: `font-size: 12px`, `font-weight: 600`, `letter-spacing: 0.05em`
- Пропсы: `active: Tab`, `onTabChange: (tab: Tab) => void` — без изменений

**Логирование:** `console.debug('[navbar] tab change:', tab)`

---

## Фаза 2: App.tsx — хедер и лайаут

### Задача 2.1: Обновить корневой лайаут
- [x] App.tsx обновлён: новый хедер, Stitch-токены

**Файлы:**
- `web/src/App.tsx` — обновление JSX

**Что сделать:**
- Хедер: градиент `linear-gradient(135deg, var(--c-primary-dark), var(--c-primary))` или сплошной `--c-primary`
- Высота хедера: 56px с safe-top
- Слева: логотип + название "Lingvo AI" (Plus Jakarta Sans, bold, 19px, белый)
- Справа: имя пользователя (из `initData.user()`) в `rgba(255,255,255,0.85)`
- Main: `flex: 1`, `overflow: hidden`, фон `var(--c-surface)`
- Safe-area переменные уже в tokens.css

---

## Фаза 3: Chat + Grammar + Usage

### Задача 3.1: Полностью заменить Chat на Stitch-дизайн

**Файлы:**
- `web/src/components/Chat.tsx` — полная замена JSX

**Что сделать:**
- JSX-структура из Stitch screen "AI Tutor Chat"
- **Сообщения пользователя**: `background: var(--c-primary)`, `color: #fff`, border-radius: верхний-правый острый (4px), остальные 12px
- **Сообщения AI**: `background: var(--c-surface-container)`, `color: var(--c-on-surface)`, border-radius: верхний-левый острый (4px), остальные 12px
- Пузыри: `max-width: 85%`, `box-shadow: var(--shadow-sm)`
- Поле ввода: pill-shaped, `background: var(--c-surface-container)`, `border: 1px solid var(--c-outline-variant)`
- Кнопка отправки: круг, `background: var(--c-primary)`
- Вся логика сохраняется: `chatStream`, `setMessages`, `setInput`, `isStreaming`, `handleVoice`
- GrammarBlock внутри AI-сообщений

### Задача 3.2: Обновить GrammarBlock

**Файлы:**
- `web/src/components/GrammarBlock.tsx` — обновление JSX

**Что сделать:**
- Badge типа ошибки: `background: <primary/error> 10% opacity`, `color: <primary/error> 100%`, pill-shaped
- Исходный текст → исправленный: `text-decoration: line-through` для wrong, жирный для correct
- Объяснение: `font-size: 14px`, `color: var(--c-on-surface-variant)`

### Задача 3.3: Обновить UsageIndicator

**Файлы:**
- `web/src/components/UsageIndicator.tsx` — обновление JSX

**Что сделать:**
- Прогресс-бар: высота 6px, `border-radius: var(--round-full)`
- Track: `background: var(--c-surface-container-highest)`
- Fill: `background: var(--c-primary)`
- Текст: "Бугун: 4/10" или "Чексиз" для Premium, `font-size: 14px`, `font-weight: 500`

**Логирование:** `console.debug('[chat] message sent:', text)`, `console.debug('[chat] stream token received')`

---

## Фаза 4: Vocabulary + ReviewCard

### Задача 4.1: Полностью заменить Vocabulary

**Файлы:**
- `web/src/components/Vocabulary.tsx` — полная замена JSX

**Что сделать:**
- JSX из Stitch screen "My Vocabulary"
- Вкладки: My Words / Lookup / Review — кнопки `pill` стиля
- Активная вкладка: `--c-primary`, неактивная: `--c-surface-container`
- Карточка слова: `border-radius: 12px`, `border: 1px solid var(--c-outline-variant)`, `padding: var(--card-padding)`
- Слово: `font-size: 18px`, `font-weight: 700` (word-card-title)
- Перевод: `font-size: 16px`, `color: var(--c-on-surface-variant)`
- Level badge: капсула, `font-size: 12px`, `background: var(--c-primary-fixed-dim)`
- Пустые состояния: placeholder с иконкой

### Задача 4.2: Обновить ReviewCard

**Файлы:**
- `web/src/components/ReviewCard.tsx` — обновление JSX

**Что сделать:**
- Флип-карточка с анимацией
- Лицевая сторона: слово, центрировано, `font-size: 28px`, `font-weight: 700`
- Обратная: перевод + пример(ы)
- 4 кнопки оценок SM-2: Again (#E53935), Hard (#FF6D00), Good (#00BFA5), Easy (#2979FF)
- Кнопки: pill, padding 12px 24px

**Логирование:** `console.debug('[vocab] tab switched:', tab)`, `console.debug('[vocab] review rating:', rating)`

---

## Фаза 5: Progress + Chart

### Задача 5.1: Полностью заменить Progress

**Файлы:**
- `web/src/components/Progress.tsx` — полная замена JSX

**Что сделать:**
- JSX из Stitch screen "Learning Progress"
- Заголовок: "Learning Progress", `font-size: 24px`, `font-weight: 700`
- Стат-карточки: Messages, Words, Streak, Days Active
- Каждая карточка: число (32px bold) + лейбл (14px, secondary)
- Badge уровня: pill, `--c-primary`
- Кнопка "Level Test": `--c-primary`, primary стиль

### Задача 5.2: Обновить ProgressChart

**Файлы:**
- `web/src/components/ProgressChart.tsx` — обновление JSX

**Что сделать:**
- Recharts: убрать grid lines, `stroke: var(--c-outline-variant)` для оси
- Бары: `fill: var(--c-primary)` с opacity 10% для фона
- Period buttons: 7d / 30d, активная `--c-primary`
- Подписи: `font-size: 12px`, `color: var(--c-hint)`

**Логирование:** `console.debug('[progress] data loaded:', stats)`

---

## Фаза 6: Subscription

### Задача 6.1: Полностью заменить Subscription

**Файлы:**
- `web/src/components/Subscription.tsx` — полная замена JSX

**Что сделать:**
- JSX из Stitch screen "Premium Plans"
- Две карточки: Free и Premium
- Free: серый border, иконка галочки для базовых фич
- Premium: `--c-secondary` (#f97316) акцентный border/badge, "RECOMMENDED"
- Фичи: список с иконками (галочка для доступных, крестик для недоступных в Free)
- Кнопка "Buy with Telegram Stars": `--c-primary`, full width, pill
- Текущий статус подписки сверху

**Логирование:** `console.debug('[subscription] invoice created:', invoiceUrl)`

---

## Фаза 7: LevelTest

### Задача 7.1: Полностью заменить LevelTest

**Файлы:**
- `web/src/components/LevelTest.tsx` — полная замена JSX

**Что сделать:**
- JSX из Stitch screen "Level Test Quiz"
- Прогресс-бар сверху (текущий вопрос / всего)
- Вопрос: `font-size: 20px`, `font-weight: 600`
- 4 варианта ответа: карточки, `border-radius: 12px`
  - Обычное состояние: `--c-surface-container`
  - Выбран: `--c-primary` с opacity
  - Правильный: `--c-tertiary` (зелёный)
  - Неправильный: `--c-error` (красный)
- Результат: CEFR уровень (A1-C2) + описание + кнопка "Начать обучение"

**Логирование:** `console.debug('[leveltest] answer submitted:', questionIndex)`, `console.debug('[leveltest] result:', level)`

---

## Фаза 8: Onboarding

### Задача 8.1: Полностью заменить Onboarding

**Файлы:**
- `web/src/components/Onboarding.tsx` — полная замена JSX

**Что сделать:**
- JSX из Stitch screen "Onboarding - AI Chat"
- 3-4 шага: приветствие, AI-чат, словарь, прогресс
- Pills-индикатор: 8x8px dots, активный 24x8px pill, `--c-primary`
- Заголовок: `font-size: 24px`, `font-weight: 700`
- Описание: `font-size: 16px`, `color: var(--c-on-surface-variant)`
- Кнопка "Let's Start" / "Далее": `--c-primary`, full width, pill
- Анимация переходов между шагами (fade/slide)

**Логирование:** `console.debug('[onboarding] step:', currentStep)`

---

## Фаза 9: DailyLesson

### Задача 9.1: Полностью заменить DailyLesson

**Файлы:**
- `web/src/components/DailyLesson.tsx` — полная замена JSX

**Что сделать:**
- JSX из Stitch screen "Daily Lesson - Ordering Food"
- Заголовок урока: `font-size: 20px`, `font-weight: 700`
- Карточки с контентом: слова, фразы, примеры
- Прогресс выполнения
- Кнопка "Вернуться в чат": `--c-primary`, secondary стиль

**Логирование:** `console.debug('[dailylesson] lesson loaded:', lessonId)`

---

## Фаза 10: Документация и сборка

### Задача 10.1: Обновить документацию

**Файлы:**
- `docs/configuration.md` — если изменились CSS-переменные
- `web/README.md` — если изменилась структура

**Что сделать:**
- Описать новую дизайн-систему (токены, компоненты)
- Обновить информацию о CSS-переменных и стилях

### Задача 10.2: Финальная сборка и проверка

**Команды:**
```bash
cd web && npm run build
cd .. && go build ./...
```

**Что сделать:**
- Проверить `npm run build` без ошибок
- Проверить `go build ./...` без ошибок
- Проверить все экраны через Cloudflare туннель

### Задача 10.3: Тесты

**Файлы:**
- `web/src/components/__tests__/*.test.tsx` — обновить/добавить тесты

**Что сделать:**
- Обновить существующие тесты под новую структуру JSX
- Проверить `vitest run` без ошибок

---

## Commit Plan

| # | Задачи | Сообщение |
|---|--------|-----------|
| 1 | Фаза 0 (CSS-токены + шрифты) | `style(web): add Stitch design system tokens and Plus Jakarta Sans` |
| 2 | Фаза 1-2 (NavBar + App layout) | `style(web): redesign NavBar and App layout` |
| 3 | Фаза 3 (Chat) | `style(web): redesign Chat screen with GrammarBlock and UsageIndicator` |
| 4 | Фаза 4 (Vocabulary + ReviewCard) | `style(web): redesign Vocabulary and ReviewCard` |
| 5 | Фаза 5 (Progress + Chart) | `style(web): redesign Progress and ProgressChart` |
| 6 | Фаза 6 (Subscription) | `style(web): redesign Subscription screen` |
| 7 | Фаза 7 (LevelTest) | `style(web): redesign LevelTest` |
| 8 | Фаза 8 (Onboarding) | `style(web): redesign Onboarding` |
| 9 | Фаза 9 (DailyLesson) | `style(web): redesign DailyLesson` |
| 10 | Фаза 10 (Docs + Build) | `docs(web): update design system docs and verify build` |
