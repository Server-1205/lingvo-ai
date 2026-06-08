# UI Redesign — Современный минимализм

**Branch:** `feature/ui-redesign`
**Created:** 2026-06-07 13:00

## Settings

- **Testing:** no
- **Logging:** verbose
- **Docs:** no
- **Roadmap Linkage:** none

## Цветовая схема

| Роль | Цвет | HEX |
|------|------|-----|
| Primary | Teal | `#00BFA5` |
| Primary Dark | Deep Teal | `#00897B` |
| Accent | Orange | `#FF6D00` |
| Accent Light | Light Orange | `#FF9100` |
| Success | Green | `#00C853` |
| Error | Red | `#E53935` |
| Background | White/Light Gray | `#FFFFFF` / `#F5F5F5` |
| Text Primary | Dark | `#212121` |
| Text Secondary | Gray | `#616161` |
| Card | White | `#FFFFFF` |

## Изменения по файлам

### Task 1: `web/src/index.css` — Дизайн-система

**Что изменить:**
- Переопределить `--tg-button` на Teal (`#00BFA5`)
- Добавить новые CSS-переменные: `--c-primary`, `--c-accent`, `--c-success`, `--c-error`, `--c-bg`, `--c-card`, `--c-text`, `--c-text-secondary`
- `font-size: 16px` → `18px` (базовый)
- `line-height: 1.4` → `1.5`
- Добавить классы для текста: `.text-lg` (20px), `.text-xl` (24px), `.text-sm` (14px)
- `.card`: padding 16→20, border-radius 12→16, box-shadow более выраженный
- `.btn`: padding 10→12px, border-radius 10→14px, font-size 15→16px
- `.btn-primary`: фон Teal, ховер эффект
- `.section-title`: font-size 17→20px
- `.placeholder-icon`: font-size 48→56px
- Добавить `.gradient-bg` для хедера

### Task 2: `web/src/App.tsx` — Хедер и общий лайаут

**Что изменить:**
- Хедер: градиентный фон (Teal → Deep Teal) вместо `var(--tg-bg)`
- Текст хедера: белый, жирный, 18px
- `minHeight: 100svh` оставить, но добавить фоновый градиент
- Увеличить `padding` в хедере

### Task 3: `web/src/components/Chat.tsx` — Чат

**Что изменить:**
- Сообщения пользователя: фон Teal (`#00BFA5`), белый текст, больший border-radius (16→20px)
- Сообщения AI: фон `#F5F5F5` (светло-серый), тёмный текст
- Шрифт сообщений: 15→17px
- Поле ввода: larger padding (10→14px), border-radius (10→14px)
- Кнопка Send: Teal, padding 10→14px
- Placeholder текст: крупнее, контрастнее

### Task 4: `web/src/components/Vocabulary.tsx` — Словарь

**Что изменить:**
- Word card: padding 12→16px, border-radius 12→16px, добавить тень
- Слово: font-size 16→18px
- Перевод: font-size 14→16px, цвет `#616161`
- Пример: font-size 13→15px
- Level badge: padding больше, фон Teal с прозрачностью
- Tab buttons: активный — Teal, неактивный — светло-серый

### Task 5: `web/src/components/ReviewCard.tsx` — Карточка повторения

**Что изменить:**
- Лицевая сторона: слово 24→28px, центрирование, мягкая тень
- Задняя сторона: перевод 15→17px, пример 13→15px
- Кнопки оценок: padding 12→16px, скругление 12→16px
- Цвета кнопок: Again (#E53935), Hard (#FF6D00), Good (#00BFA5), Easy (#2979FF)

### Task 6: `web/src/components/NavBar.tsx` — Нижняя навигация

**Что изменить:**
- Активная вкладка: Teal вместо `var(--tg-button)`
- Неактивные: `#9E9E9E` (темнее, чем hint)
- Иконки: font-size 22→24px
- Лейблы: font-size 10→12px, font-weight 500→600
- Высота навигации: `--nav-height: 64px` → `72px`

### Task 7: `web/src/components/Progress.tsx` — Прогресс

**Что изменить:**
- Stats card: padding 20→24px, более выраженные числа (28→32px)
- Уровень: Teal badge
- Кнопка Level Test: Teal, padding 12→16px
- Отступы между stat-блоками

### Task 8: `web/src/components/Subscription.tsx` — Подписка

**Что изменить:**
- Cards: более выраженные, с тенью и границей
- Акцент на Premium: оранжевый badge
- Кнопка Subscribe: оранжевый градиент
- План Free: серая обводка

### Task 9: `web/src/components/ProgressChart.tsx` — График

**Что изменить:**
- Цвет бара: Teal
- Цвет оси: `#E0E0E0`
- Размер шрифта подписей: 11→13px
- Tooltip: более контрастный фон, скругление

### Task 10: `web/src/components/LevelTest.tsx` — Тест уровня

**Что изменить:**
- Вопрос: font-size 17→20px, жирный
- Опции: padding 12→16px, border-radius 12→16px
- Правильный ответ: зелёный фон
- Неправильный: красный фон
- Прогресс-бар: Teal

---

## Риски

- **Telegram theme override**: цвета могут конфликтовать с тёмной темой Telegram. Проверить оба режима
- **Размер шрифта**: увеличение может сломать лейаут на маленьких экранах. Тестировать на 320px ширине
- **Контраст**: обеспечить WCAG AA (минимум 4.5:1 для текста)
- **Тёмная тема**: все новые цвета должны иметь тёмные варианты

## Commit Plan

1. `style: add design system with teal/orange palette`
2. `style: redesign chat component`
3. `style: redesign vocabulary and review cards`
4. `style: redesign navigation, progress, subscription`
