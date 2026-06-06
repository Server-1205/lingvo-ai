# Архитектура: Structured Modules (Technical Layers)

## Обзор

Structured Modules — лёгкая, domain-aware модульная архитектура, организованная по техническим слоям. Каждый модуль (`internal/ai`, `internal/api`, `internal/bot` и т.д.) инкапсулирует функциональную область со своими обработчиками, сервисами и моделями. Этот паттерн выбран как баланс между простотой Layered Architecture и формализмом Explicit Architecture — он уже соответствует текущей структуре проекта и не требует рефакторинга.

## Обоснование решения

- **Тип проекта:** Telegram Mini App (AI English Tutor) — один сервер, одна кодовая база
- **Технологический стек:** Go (Gin) + React (Vite) + SQLite
- **Ключевой фактор:** Проект уже имеет чёткую модульную структуру `internal/{ai,api,bot,db,middleware}`; принудительное навязывание DDD-формализма сейчас добавит сложности без немедленной выгоды
- **Размер команды:** 1-3 разработчика
- **Масштаб:** 500-2000 MAU в первый месяц, один бинарный файл для Railway

## Структура папок

```
sinking/
├── cmd/
│   └── server/
│       └── main.go                # Точка входа, DI, конфигурация
│
├── internal/
│   ├── ai/                        # AI-модуль
│   │   ├── gemini.go              # Клиент Gemini API
│   │   ├── prompts.go             # Промпты для всех эндпоинтов
│   │   ├── sm2.go                 # Алгоритм SM-2
│   │   └── response.go            # Парсинг ответов AI
│   │
│   ├── api/                       # REST API (Gin-обработчики)
│   │   ├── router.go              # Роутер + middleware
│   │   ├── chat.go                # POST /api/chat
│   │   ├── chat_stream.go         # GET /api/chat/stream (SSE)
│   │   ├── grammar.go             # POST /api/grammar
│   │   ├── vocabulary.go          # CRUD словаря
│   │   ├── quiz.go                # POST /api/quiz
│   │   ├── level.go               # POST /api/level
│   │   ├── progress.go            # GET /api/progress
│   │   ├── subscription.go        # GET /api/subscription
│   │   └── invoice.go             # POST /api/create-invoice
│   │
│   ├── bot/                       # Telegram Bot (long-polling)
│   │   ├── bot.go                 # Инициализация, цикл poll
│   │   └── handlers.go            # /start, /help, /daily, /stats, successful_payment
│   │
│   ├── db/                        # Доступ к данным (SQLite)
│   │   ├── schema.sql             # Схема БД
│   │   ├── db.go                  # Init + миграция
│   │   ├── users.go               # CRUD пользователей
│   │   ├── messages.go            # Учёт сообщений
│   │   ├── subscriptions.go       # CRUD подписок
│   │   ├── vocabulary.go          # CRUD словаря
│   │   └── progress.go            # Прогресс пользователя
│   │
│   ├── middleware/                # HTTP middleware (cross-cutting)
│   │   ├── auth.go                # HMAC-SHA256 initData
│   │   └── ratelimit.go           # Rate limiter (10/day free)
│   │
│   ├── i18n/                      # Переводы (backend)
│   │   ├── uz.json
│   │   └── ru.json
│   │
│   └── models/                    # Общие типы и структуры
│       └── types.go
│
├── web/                           # Frontend (React SPA)
│   ├── src/
│   │   ├── components/            # React-компоненты
│   │   ├── api/client.ts          # HTTP-клиент
│   │   ├── hooks/                 # Кастомные хуки
│   │   └── locales/               # i18n (uz, ru)
│   ├── package.json
│   └── vite.config.ts
│
├── docs/                          # Документация
├── .ai-factory/                   # AI Factory конфигурация
└── .opencode/                     # OpenCode настройки
```

## Правила зависимостей

- **Строгий downward flow:** `api → ai/db/middleware`, `bot → ai/db`, `middleware → db`. Внутренние слои (db, models) НИКОГДА не зависят от внешних (api, bot)
- **Нет пропуска слоёв:** api-обработчики не должны вызывать db напрямую, минуя ai (если требуется AI-логика)
- **Изоляция модулей:** Модули внутри `internal/` используют `models` и `i18n` но НЕ импортируют внутренности друг друга напрямую. public API только через определённые экспортируемые функции
- **shared minimal:** `internal/models` и `internal/i18n` — единственные общедоступные модули

## Коммуникация слоёв/модулей

- **API → AI:** Обработчики вызывают `ai.Client.GenerateContent()` для запросов к Gemini
- **API → DB:** Обработчики или AI-слой вызывают `db.*Repository` для CRUD
- **Bot → DB:** Обработчики бота вызывают `db` напрямую для пользователей и подписок
- **Bot → AI:** Обработчики бота вызывают `ai.Client` для ежедневных уроков
- **Middleware → DB:** Обработчики middleware (ratelimit, auth) вызывают `db` напрямую
- **Композиция:** `cmd/server/main.go` — единственная точка, где всё связывается (Composition Root)

## Ключевые принципы

1. **Rich Domain Models:** Модели в `internal/models` должны содержать поведение, а не только данные. Например, SM-2 алгоритм — метод на модели карточки, а не внешняя функция
2. **Dependency Inversion (lightweight):** Сервисы получают зависимости через параметры (без DI-фреймворка). Repository interfaces поощряются для подготовки к будущей Explicit Architecture
3. **Domain Awareness:** Core business rules (SM-2, rate limit, подсчёт streak) должны быть в моделях или выделенных сервисных функциях, не в HTTP-обработчиках
4. **Thin Handlers:** Обработчики API/бота парсят ввод, вызывают 1-2 сервисных метода, форматируют ответ. Бизнес-логика — в ai/ или db/
5. **Stateless Services:** Сервисы не хранят состояние запроса. Все данные передаются параметрами

## Примеры кода

### Rich Domain Model — SM-2 Card

```go
// internal/models/types.go
type SM2Card struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    Word      string    `json:"word"`
    Reps      int       `json:"reps"`
    EF        float64   `json:"ef"`
    Interval  int       `json:"interval"`
    NextReview time.Time `json:"next_review"`
}

// CalculateNextReview — бизнес-логика SM-2 на модели
func (c *SM2Card) CalculateNextReview(quality int) {
    if quality < 3 {
        c.Reps = 0
        c.Interval = 1
    } else {
        c.Reps++
        switch c.Reps {
        case 1:
            c.Interval = 1
        case 2:
            c.Interval = 6
        default:
            c.Interval = int(math.Round(float64(c.Interval) * c.EF))
        }
    }
    c.EF = c.EF + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
    if c.EF < 1.3 {
        c.EF = 1.3
    }
    c.NextReview = time.Now().AddDate(0, 0, c.Interval)
}
```

### Thin Handler — API Endpoint

```go
// internal/api/chat.go
func ChatHandler(db *sqlx.DB, aiClient *ai.Client, sugar *zap.SugaredLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Cache-Control", "no-store")

        userID := c.GetInt64("user_id")
        var req ChatRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": "invalid_request"})
            return
        }

        resp, err := aiClient.Chat(c.Request.Context(), userID, req.Text, c.GetString("lang"))
        if err != nil {
            sugar.Errorw("chat failed", "user_id", userID, "error", err)
            c.JSON(500, gin.H{"error": "ai_service_unavailable"})
            return
        }

        c.JSON(200, resp)
    }
}
```

### Composition Root

```go
// cmd/server/main.go — единственное место, где всё связывается
func main() {
    database, _ := sqlx.Connect("sqlite", dbPath)
    db.Migrate(database, sugar)

    aiClient := ai.NewClient(geminiKey, sugar)

    r := gin.Default()
    r.Use(middleware.Auth(database))
    r.Use(middleware.RateLimit(database))

    api.RegisterRoutes(r, database, aiClient, sugar)
    go bot.Start(database, aiClient, botToken, sugar)

    r.Run(":" + port)
}
```

## Анти-паттерны

- ❌ **Anemic Domain Models:** Модели только с геттерами/сеттерами, логика — в обработчиках. Бизнес-правила (SM-2, лимиты) должны быть на моделях
- ❌ **Пропуск слоёв:** Обработчики API, вызывающие db напрямую, минуя AI-слой, когда требуется AI-валидация
- ❌ **God Handlers:** Обработчики, содержащие бизнес-логику, проверку лимитов, вызовы AI и сохранение в БД — всё в одной функции
- ❌ **Cache-Control нарушение:** Не добавлять `Cache-Control: no-store` на API-ответах — клиенты (Telegram WebView) могут кэшировать данные
- ❌ **initData в обработчиках:** Верификация initData должна быть в middleware, а не в каждом обработчике
