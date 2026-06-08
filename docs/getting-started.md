[Back to README](../README.md) · [API Reference →](api.md)

# Установка и запуск

## Требования

- Go 1.25+
- Node.js 20+
- pnpm (или npm)
- Telegram Bot Token (от @BotFather)
- Gemini API Key (от [Google AI Studio](https://aistudio.google.com))
- Python 3.10+ и `edge-tts` для озвучки: `pip install edge-tts` (опционально)

## 1. Настройка окружения

```bash
cp .env.example .env
```

Заполните обязательные поля в `.env`:

| Переменная | Описание |
|------------|----------|
| `BOT_TOKEN` | Токен Telegram бота от @BotFather |
| `GEMINI_API_KEY` | API ключ Gemini от Google AI Studio |

Остальные переменные имеют значения по умолчанию — см. [Конфигурация](configuration.md).

## 2. Запуск backend

```bash
go mod tidy
go run ./cmd/server
```

Сервер запустится на `:8080` (или `PORT` из `.env`).

**Логи:** сервер использует `go.uber.org/zap` для структурированного логирования.

## 3. Запуск frontend

В отдельном терминале:

```bash
cd web
pnpm install
pnpm dev
```

Frontend запустится на `http://localhost:5173`.

## Telegram Bot

1. Создайте бота через [@BotFather](https://t.me/BotFather)
2. Получите токен и укажите его в `BOT_TOKEN`
3. Настройте Mini App через BotFather: укажите URL `http://localhost:5173`
4. Запустите сервер — бот начнёт polling автоматически

### Команды бота

| Команда | Описание |
|---------|----------|
| `/start` | Регистрация, приветствие, кнопка Launch |
| `/help` | Список команд |
| `/daily` | Ежедневный урок от AI |
| `/stats` | Статистика пользователя |

## Проверка

```bash
curl http://localhost:8080/api/health
# {"status":"ok"}
```

## См. также

- [API Reference](api.md) — все эндпоинты
- [Конфигурация](configuration.md) — переменные окружения
- [AI Streaming](streaming.md) — SSE-поток
