[← API Reference](api.md) · [Back to README](../README.md) · [AI Streaming →](streaming.md)

# Конфигурация

Все настройки задаются через переменные окружения или `.env` файл.

## Переменные окружения

| Переменная | Обязательно | По умолчанию | Описание |
|------------|-------------|-------------|----------|
| `BOT_TOKEN` | **Да** | — | Токен Telegram бота от @BotFather |
| `GEMINI_API_KEY` | **Да** | — | API ключ Gemini от Google AI Studio |
| `DATABASE_PATH` | Нет | `lingvo.db` | Путь к файлу SQLite |
| `PORT` | Нет | `8080` | Порт HTTP-сервера |
| `WEBAPP_URL` | Нет | (пусто) | URL Mini App для CORS/deeplink |
| `ADMIN_IDS` | Нет | (пусто) | Telegram ID администраторов (через запятую) |
| `OPENAI_API_KEY` | Нет | — | API ключ для fallback (OpenAI-совместимый) |
| `OPENAI_BASE_URL` | Нет | `https://api.openai.com/v1` | Базовый URL для fallback API (DeepSeek, Groq и т.д.) |
| `OPENAI_MODEL` | Нет | `gpt-4o-mini` | Модель для fallback |
| `AI_QUEUE_ENABLED` | Нет | — | Включить приоритетную очередь AI (premium → free) |
| `DEV_MODE` | Нет | — | Режим разработчика — все фичи открыты, без лимитов (true/false) |
| `TTS_VOICE_UZ` | Нет | `en-US-JennyNeural` | Голос edge-tts (все слова — английские) |
| `TTS_VOICE_RU` | Нет | `en-US-JennyNeural` | Голос edge-tts (все слова — английские) |

## Файл .env

Скопируйте `.env.example` в `.env`:

```bash
cp .env.example .env
```

Файл `.env` загружается автоматически при старте сервера.

## Настройка Telegram Bot

1. Откройте [@BotFather](https://t.me/BotFather)
2. Создайте бота: `/newbot`
3. Получите токен и укажите в `BOT_TOKEN`
4. Настройте Mini App: `/mybots` → ваш бот → Bot Settings → Menu Button
5. Укажите URL вашего фронтенда (например, `http://localhost:5173`)

## Настройка Gemini API

1. Перейдите на [Google AI Studio](https://aistudio.google.com)
2. Создайте API ключ
3. Укажите в `GEMINI_API_KEY`

## Режимы работы бота

Проект использует **long-polling** (не webhook). Сервер сам опрашивает Telegram API.

## См. также

- [Установка и запуск](getting-started.md) — как запустить проект
- [API Reference](api.md) — все эндпоинты
