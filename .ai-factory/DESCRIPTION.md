# Lingvo AI — Описание проекта

> AI English Tutor Telegram Mini App. Узбекистан → Центральная Азия.
> AI-репетитор английского языка в Telegram с оплатой через Telegram Stars.

## Обзор

Lingvo AI — это Telegram Mini App для изучения английского языка с помощью ИИ. Пользователи общаются с AI-репетитором (Gemini 2.0 Flash), получают исправления грамматических ошибок, пополняют словарный запас с интервальным повторением (SM-2), проходят тесты на определение уровня и отслеживают прогресс. Бесплатный тариф: 10 сообщений/день, премиум — безлимит за Telegram Stars.

## Основные возможности

- **AI-чат**: общение с AI-репетитором с исправлением ошибок в реальном времени
- **Streaming**: SSE-поток для мгновенного отображения ответа AI
- **Грамматика**: отдельный эндпоинт для проверки грамматики
- **Словарный запас**: CRUD операции, SM-2 интервальное повторение, поиск слов
- **Тест уровня**: 10 вопросов (A1-C1) для определения уровня
- **Прогресс**: статистика, streak, уровень пользователя
- **Подписка**: Telegram Stars (300/week, 800/month)
- **Бот**: /start, /help, /daily, /stats, successful_payment

## Технологический стек

- **Язык (backend):** Go 1.25
- **Фреймворк:** Gin
- **Язык (frontend):** TypeScript 6
- **Фреймворк (frontend):** React 19, Vite 8
- **База данных:** SQLite (modernc.org/sqlite, без CGO)
- **ORM:** sqlx (без полноценного ORM)
- **AI:** Gemini 2.0 Flash (основной), GPT-4o-mini (fallback)
- **Бот:** go-telegram-bot-api v5 (long-polling)
- **Платежи:** Telegram Stars (createInvoiceLink)
- **Логирование:** zap (uber-go)
- **i18n:** react-i18next (фронтенд), JSON файлы (бекенд + фронтенд)
- **Аутентификация:** Telegram WebApp initData (HMAC-SHA256)

## Архитектурные заметки

- **Разделение слоёв**: backend разделён на bot (Telegram), api (REST), ai (Gemini), db (SQLite), middleware (auth + ratelimit), i18n
- **Middleware**: auth (HMAC верификация) + ratelimit (10/day free, безлимит premium)
- **AI-пайплайн**: промпты → Gemini → парсинг JSON → reply + corrections[]
- **Платежи**: createInvoiceLink → Telegram обрабатывает → successful_payment → сохранение подписки
- **Streaming**: SSE (Server-Sent Events) для AI-ответов
- **Rate limit**: DATE-based (UTC), 10 сообщений/день для free
- **Cache-Control**: no-store на всех API-ответах

## Архитектура

Подробные архитектурные гайдлайны: `.ai-factory/ARCHITECTURE.md`
- **Паттерн:** Structured Modules (Technical Layers)
- **Принципы:** Rich Domain Models, Thin Handlers, Dependency Inversion (lightweight), Composition Root

## Нефункциональные требования

- **Логирование**: zap, структурированные логи (user_id, error, request_id)
- **Обработка ошибок**: структурированные JSON-ответы с HTTP-кодами (400, 401, 429, 500)
- **Безопасность**: HMAC-SHA256 initData, секреты только в env, rate limit 100 req/s global
- **Масштабирование**: SQLite → PostgreSQL при росте, Dockerfile для Railway
- **Тестирование**: Go unit-тесты (табличные) + Vitest (фронтенд)
- **i18n**: узбекский (лат.) + русский язык
