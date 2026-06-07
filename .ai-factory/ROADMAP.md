# Lingvo AI — Roadmap

> AI English Tutor в Telegram. Узбекистан → Центральная Азия → Мир.

## Milestones

- [x] **Инфраструктура** — Go-модуль, React+Vite шаблон, SQLite, AI Factory, Dockerfile, CI (GitHub Actions: build/lint/tests/security)
- [x] **Ядро: Telegram Bot + Chat** — long-polling, /start, /help, /daily, /stats, initData auth, Gemini 2.0 Flash, SSE streaming, rate limiter (10/day free)
- [x] **Ядро: Mini App Frontend** — Chat, Vocabulary, Progress, Level Test, Subscription, NavBar, LanguageSwitcher
- [x] **Ядро: База данных** — schema + migrations (users, messages, subscriptions, vocabulary, daily_progress)
- [x] **Vocabulary + SM-2** — CRUD словаря, AI-генерация перевода/примеров, Spaced Repetition (SM-2), review API
- [x] **Level Test + Progress** — AI-генерация тестов (A1-C1), сохранение уровня, streak, график прогресса (7/30 дней)
- [x] **Монетизация** — Telegram Stars (300/week, 800/month), createInvoiceLink, successful_payment, subscription middleware
- [x] **AI Fallback** — OpenAI-совместимый клиент (Gemini Flash + Flash-Lite, fallback GPT-4o-mini/DeepSeek и др.)
- [x] **Документация + SPEC** — SPEC.md, ARCHITECTURE.md, AGENTS.md, docs/api.md, docs/configuration.md, docs/streaming.md, docs/getting-started.md
- [ ] **Premium Features** — продвинутый анализ ошибок, экспорт словаря (CSV), приоритетная очередь AI
- [ ] **Referral программа** — уникальные ссылки, +5 сообщений за приглашение, статистика
- [ ] **Маркетинговый запуск** — Telegram канал, размещение в каталогах (findmini.app, tapps.center), офлайн-продвижение
- [ ] **Аналитика + Оптимизация** — retention (D1/D7/D30), conversion rate, ARPU/LTV, A/B тесты, оптимизация AI costs
- [ ] **Telegram Ads** — реинвестиция Stars, таргетинг Узбекистан, A/B тест креативов
- [ ] **Экспансия: языки V2** — казахский (kk), кыргызский (ky), турецкий (tr)
- [ ] **Фичи V2** — ежедневный AI-урок, TTS произношение, групповой режим, IELTS/CEFR подготовка
- [ ] **Платежи V2** — USDT/USDC (GramBase), подписка за TON, автопродление

## Completed

| Milestone | Date |
|-----------|------|
| Инфраструктура | 2026-06-07 |
| Ядро: Telegram Bot + Chat | 2026-06-07 |
| Ядро: Mini App Frontend | 2026-06-07 |
| Ядро: База данных | 2026-06-07 |
| Vocabulary + SM-2 | 2026-06-07 |
| Level Test + Progress | 2026-06-07 |
| Монетизация | 2026-06-07 |
| AI Fallback | 2026-06-07 |
| Документация + SPEC | 2026-06-07 |

## KPI (первые 60 дней)

| Метрика | Цель |
|---------|------|
| MAU (месяц 1) | 500-2,000 |
| MAU (месяц 2) | 3,000-10,000 |
| Conversion rate | 3-5% |
| Paying users (месяц 2) | 100-300 |
| Выручка (месяц 2) | $500-2,500 |
| Retention D7 | >20% |
| Retention D30 | >10% |
| CAC (Telegram Ads) | <$0.50 |
| AI cost per user/mes | <$0.02 |

## Risks

| Риск | Митигация |
|------|-----------|
| Gemini галлюцинирует в грамматике | Промпты с примерами + fallback на OpenAI-совместимую модель |
| Юзеры не платят | Бесплатный tier (10/день) даёт ценность; A/B тест цен |
| Высокие AI costs | Лимит 10/day free, Flash-Lite для quiz/level, fallback для экономии |
| Модерация Telegram | Не нарушать ToS: подписки через Stars разрешены |
| Конкуренты | Скорость: ядро готово, запуск — маркетинг + рефералки |
