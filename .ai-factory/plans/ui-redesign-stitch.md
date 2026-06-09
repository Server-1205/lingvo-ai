# UI Redesign — Stitch (Lingvo AI Teal)

**Created:** 2026-06-08
**Source:** Google Stitch MCP — проект `Lingvo AI Tutor` (ID: `13746191202433220969`)
**Design System:** Teal #00685f + Plus Jakarta Sans + Fidelity цветовая схема

## Цель

Полная замена визуала всех React-компонентов на дизайн из Stitch. Сохраняется вся бизнес-логика, API-запросы, i18n, стейт-менеджмент.

## Структура экранов в Stitch

| Экран Stitch | React компонент | Описание |
|---|---|---|
| Onboarding - AI Chat | `Onboarding.tsx` | Приветственный онбординг |
| AI Tutor Chat | `Chat.tsx` | Чат с AI + GrammarBlock + UsageIndicator |
| My Vocabulary | `Vocabulary.tsx` | Список слов + поиск + ReviewCard |
| Learning Progress | `Progress.tsx` + `ProgressChart.tsx` | Статистика и графики |
| Premium Plans | `Subscription.tsx` | Планы подписки |
| Level Test Quiz | `LevelTest.tsx` | Тест на уровень |
| Daily Lesson | `DailyLesson.tsx` | Ежедневный урок |
| NavBar | `NavBar.tsx` | Нижняя навигация |

---

## Фаза 0: Подготовка

### Задача 0.1: Дизайн-система
- Применить `update_design_system` к Stitch проекту (Lingvo AI Teal)
- Убедиться, что все экраны используют единую дизайн-систему

### Задача 0.2: CSS-токены
Создать новый `web/src/styles/tokens.css` с токенами из Stitch:
```css
--c-primary: #00685f;
--c-primary-container: #008378;
--c-primary-fixed: #89f5e7;
--c-secondary: #9d4300;
--c-secondary-container: #fd761a;
--c-tertiary: #006b2d;
--c-tertiary-container: #00873b;
--c-surface: #f5faf8;
--c-surface-dim: #d6dbd9;
--c-surface-bright: #f5faf8;
--c-on-surface: #171d1c;
--c-outline: #6d7a77;
--c-error: #ba1a1a;
--font-family: 'Plus Jakarta Sans', sans-serif;
--round-sm: 4px;
--round-md: 8px;
--round-lg: 12px;
--round-xl: 16px;
--round-full: 9999px;
--space-unit: 4px;
--gutter: 16px;
--card-padding: 20px;
--nav-height: 64px;
```

### Задача 0.3: Google Fonts
Добавить `Plus Jakarta Sans` в `index.html` через Google Fonts CDN

---

## Фаза 1: Базовые компоненты

### Задача 1.1: NavBar (`NavBar.tsx`)
- Фиксированный bottom bar с blur-эффектом
- 4 таба: Chat, Vocabulary, Progress, Subscription
- Иконки stroke-based 24px, активный — primary цвет + dot-индикатор
- Высота 64px, safe-area учёт

### Задача 1.2: App.tsx — хедер и лайаут
- Градиентный хедер (teal) с логотипом и именем пользователя
- Основная область с safe-area
- Onboarding модалка поверх всего

---

## Фаза 2: Чат

### Задача 2.1: Chat (`Chat.tsx`)
- **Сообщения пользователя**: primary teal фон, белый текст, правый угол острый
- **Сообщения AI**: surface-container фон, тёмный текст, левый угол острый
- Пузыри с тенью, max-width 85%
- Поле ввода: pill-shaped, surface-container, кнопка отправки
- GrammarBlock внутри AI-сообщений
- PremiumCorrectionBlock
- UsageIndicator внизу

### Задача 2.2: GrammarBlock (`GrammarBlock.tsx`)
- Badge типа ошибки (grammar/vocabulary/spelling) с 10% opacity фоном
- original → corrected с объяснением

### Задача 2.3: UsageIndicator (`UsageIndicator.tsx`)
- Прогресс-бар 6px высота, rounded-full
- Primary fill, surface-container-highest track
- Текст "Бугун: 4/10" или "Чексиз" для Premium

---

## Фаза 3: Словарь

### Задача 3.1: Vocabulary (`Vocabulary.tsx`)
- Вкладки: My Words / Lookup / Review
- Слова — карточки с 1px border, border-radius 12px
- Слово жирное 18px (word-card-title), перевод 16px (translation-text)
- Level badge — маленькая капсула на карточке
- Кнопки действий: Delete, Edit

### Задача 3.2: ReviewCard (`ReviewCard.tsx`)
- Карточка флип-стиля
- Лицевая сторона: слово, крупно, центрировано
- Обратная: перевод + пример
- 4 кнопки: Again / Hard / Good / Easy (цветные)

---

## Фаза 4: Прогресс

### Задача 4.1: Progress (`Progress.tsx`)
- Заголовок "Learning Progress"
- Стат-карточки: Messages, Words, Streak, Days Active
- Badge уровня (pill, teal)
- Кнопка Level Test

### Задача 4.2: ProgressChart (`ProgressChart.tsx`)
- Бар-чарт без сетки
- primary цвет для баров с 10% opacity треком
- Кнопки 7d / 30d
- label-sm шрифт в hint_color

---

## Фаза 5: Подписка

### Задача 5.1: Subscription (`Subscription.tsx`)
- Две карточки: Free (10msgs/day) и Premium (unlimited)
- Premium — выделена, с оранжевым secondary акцентом
- Кнопка "Buy with Telegram Stars" — primary teal
- Текущий статус подписки

---

## Фаза 6: Тесты и онбординг

### Задача 6.1: LevelTest (`LevelTest.tsx`)
- Прогресс-бар сверху
- Вопрос + 4 варианта ответа
- Правильный/неправильный — зелёный/красный фон
- Результат: CEFR уровень + описание

### Задача 6.2: Onboarding (`Onboarding.tsx`)
- 3-4 шага онбординга
- Pills-индикатор прогресса (8x8 dots, активный 24x8)
- Кнопка "Let's Start"
- Анимация переходов

### Задача 6.3: DailyLesson (`DailyLesson.tsx`)
- Заголовок урока
- Карточки с контентом (слова, фразы)
- Кнопка "Вернуться в чат"

---

## Фаза 7: Сборка и проверка

### Задача 7.1: Сборка
```bash
cd web && npm run build
```

### Задача 7.2: Проверка в туннеле
- Открыть https://acute-visibility-commercial-professional.trycloudflare.com
- Проверить все экраны в Telegram Mini App

### Задача 7.3: Backend
```bash
go build ./...
```

---

## Порядок выполнения

```
Фаза 0 (Токены + Шрифты)
  └─> Фаза 1 (NavBar + App)
       ├─> Фаза 2 (Chat + Grammar + Usage)
       ├─> Фаза 3 (Vocabulary + ReviewCard)
       ├─> Фаза 4 (Progress + Chart)
       ├─> Фаза 5 (Subscription)
       └─> Фаза 6 (LevelTest + Onboarding + DailyLesson)
              └─> Фаза 7 (Build + Test)
```

## Риски

- Stitch HTML может содержать избыточную вёрстку — адаптировать под React
- Telegram theme переменные могут конфликтовать — проверить dark mode
- Plus Jakarta Sans может не загрузиться — fallback на sans-serif
- Safe-area для разных устройств Telegram
