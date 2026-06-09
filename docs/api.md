[← Установка и запуск](getting-started.md) · [Back to README](../README.md) · [Конфигурация →](configuration.md)

# API Reference

Базовый URL: `http://localhost:8080` (или `WEBAPP_URL`)

Все защищённые эндпоинты требуют заголовок `X-Telegram-Init-Data` с initData от Telegram WebApp.

## Аутентификация

```http
X-Telegram-Init-Data: query_id=...&auth_date=...&hash=...
```

Верификация через HMAC-SHA256 в `internal/middleware/auth.go`.

## Эндпоинты

### GET /api/health

Проверка работоспособности.

```json
{"status": "ok"}
```

---

### POST /api/chat

Отправить сообщение AI-репетитору.

**Request:**
```json
{"text": "I go to school yesterday"}
```

**Response (free):**
```json
{
  "reply": "Good! Tell me more about your day.",
  "corrections": [
    {
      "original": "I go",
      "corrected": "I went",
      "explanation_uz": "\"I go\" эмас, \"I went\". \"yesterday\" ўтган замонни кўрсатади.",
      "explanation_ru": "Нужно \"I went\". \"yesterday\" указывает на Past Simple.",
      "type": "grammar"
    }
  ],
  "usage": {
    "daily_used": 4,
    "daily_limit": 10,
    "is_premium": false
  }
}
```

**Response (premium — доп. поля в corrections + premium_analysis):**
```json
{
  "reply": "Good attempt! Let me help you improve.",
  "corrections": [
    {
      "original": "I has",
      "corrected": "I have",
      "explanation_uz": "\"I\" дан кейин \"has\" эмас, \"have\" ишлатилади.",
      "explanation_ru": "После \"I\" используется \"have\", а не \"has\".",
      "type": "grammar",
      "severity": "major",
      "category": "grammar",
      "learning_tip": "\"I\" всегда с \"have\", \"he/she/it\" — с \"has\"",
      "rule_violated": "Subject-Verb Agreement (Present Simple)"
    }
  ],
  "usage": {
    "daily_used": 4,
    "daily_limit": 10,
    "is_premium": true
  },
  "premium_analysis": {
    "overall_grade": "B",
    "strengths": ["good vocabulary", "clear sentence structure"],
    "areas_for_improvement": ["subject-verb agreement"],
    "suggested_topic": "Present Simple vs Present Continuous"
  }
}
```

---

### POST /api/chat/stream

SSE-поток для streaming AI-ответов. Аутентификация и rate-limit как у `/api/chat`.

**Response (SSE):**
```
data: {"type":"token","data":"Hello"}
data: {"type":"token","data":"! I'm"}
data: {"type":"corrections","data":[...]}
data: {"type":"usage","data":{...}}
data: {"type":"done"}
```

Подробнее: [AI Streaming](streaming.md)

---

### POST /api/grammar

Проверка грамматики без диалога.

**Request:**
```json
{"text": "He don't like coffee"}
```

**Response:**
```json
{"corrections": [...]}
```

---

### POST /api/vocab

Добавить слово в словарь.

**Request:**
```json
{"word": "ubiquitous"}
```

**Response:**
```json
{
  "word": "ubiquitous",
  "translation_uz": "барча жойда учрайдиган",
  "translation_ru": "вездесущий",
  "examples": ["Smartphones are ubiquitous these days."],
  "level": "b2"
}
```

---

### GET /api/vocab

Список слов пользователя.

**Query:** `?page=1&per_page=20&due_only=true`

**Response:**
```json
{"words": [...], "total": 42, "due_count": 12}
```

---

### GET /api/vocab/export

Экспорт словаря в CSV. Только для premium-пользователей.

**Response:** `Content-Type: text/csv` с заголовком `Content-Disposition: attachment; filename="vocabulary.csv"`

```csv
word,translation,example,level,review_count,next_review,created_at
hello,salom,Hello world!,a1,3,,2026-06-01
```

---

### POST /api/vocab/review

Повторение слова (SM-2).

**Request:**
```json
{"word_id": 1, "quality": 4}
```

---

### POST /api/quiz

Генерация теста.

**Request:**
```json
{"topic": "past_simple", "count": 5}
```

**Response:**
```json
{
  "questions": [
    {
      "question": "...",
      "options": [...],
      "correct": 1,
      "explanation_uz": "..."
    }
  ]
}
```

---

### POST /api/level

Определение уровня по ответам теста.

---

### GET /api/progress

Текущая статистика пользователя: streak, сообщения, слова, уровень.

---

### GET /api/progress/history

История ежедневной активности для графика прогресса.

**Query:** `?days=7` (по умолчанию 7, макс 30)

**Response:**
```json
{
  "entries": [
    {"date": "2026-06-01", "messages_sent": 5, "words_learned": 3, "quizzes_taken": 1},
    {"date": "2026-06-02", "messages_sent": 8, "words_learned": 2, "quizzes_taken": 0}
  ]
}
```

---

### GET /api/subscription

Статус подписки.

---

### POST /api/create-invoice

Создать счёт в Telegram Stars.

**Request:**
```json
{"plan": "monthly"}
```

**Response:**
```json
{"invoice_link": "tg://pay?invoice=...", "stars": 800}
```

---

### GET /api/tts

Озвучивание текста (Text-to-Speech) через edge-tts.

**Query Parameters:**
- `text` (string, required) — текст для озвучивания (макс. 500 символов)
- `lang` (string, optional, default: `uz`) — язык: `uz` или `ru`

**Response:**
- `200` — аудиофайл в формате MP3 (`Content-Type: audio/mpeg`)

**Errors:**
- `400` — отсутствует `text` или текст длиннее 500 символов
- `401` — unauthorized
- `500` — ошибка синтеза речи (edge-tts не установлен или сбой)

**Headers:**
- `X-Telegram-Init-Data` — обязательный

**Note:** Требует установленного `edge-tts` (Python): `pip install edge-tts`. Максимальная длина текста — 500 символов.

## Коды ошибок

| HTTP | error | Причина |
|------|-------|---------|
| 400 | `invalid_request` | Отсутствует поле |
| 401 | `unauthorized` | Невалидный initData |
| 429 | `daily_limit_exceeded` | Лимит бесплатного тарифа |
| 400 | `text_required` | Не указан text для TTS |
| 400 | `inappropriate_word` | Неприемлемое слово |
| 400 | `text_too_long` | Текст для TTS длиннее 500 символов |
| 500 | `ai_service_unavailable` | Gemini недоступен |
| 500 | `tts_failed` | Ошибка синтеза речи (edge-tts) |
| 500 | `internal_error` | Ошибка БД/сервера |

**429 Response:**
```json
{
  "error": "daily_limit_exceeded",
  "message_uz": "Бугунги лимит тугади. Чексизга обуна бўлинг!",
  "message_ru": "Дневной лимит исчерпан. Оформите подписку!",
  "premium_link": "tg://pay?invoice=abc"
}
```

## См. также

- [Установка и запуск](getting-started.md) — как запустить проект
- [AI Streaming](streaming.md) — SSE-поток
- [Конфигурация](configuration.md) — переменные окружения
