---
name: go-gin-api
description: Use when creating or modifying Gin HTTP handlers, API endpoints, routes, or middleware for the Lingvo AI Go backend. Covers Gin router setup, handler patterns, JSON responses, request validation, and middleware chaining. NOT for bot handlers or AI integration.
---

# Go Gin API Development

## Router Setup (`internal/api/router.go`)

```go
func SetupRouter(db *sqlx.DB, bot *tgbotapi.BotAPI) *gin.Engine {
    r := gin.Default()
    r.Use(cors.Default())

    api := r.Group("/api")
    api.GET("/health", healthHandler)

    protected := api.Group("")
    protected.Use(authMiddleware(db))
    {
        protected.POST("/chat", chatHandler(db, bot))
    }

    return r
}
```

## Handler Pattern

```go
func chatHandler(db *sqlx.DB, bot *tgbotapi.BotAPI) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.GetString("telegram_id")

        var req ChatRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": "invalid request"})
            return
        }

        // business logic...

        c.JSON(200, ChatResponse{
            Reply:       reply,
            Corrections: corrections,
        })
    }
}
```

## Common Patterns

- Return `c.JSON(code, gin.H{"error": "message"})` for errors
- Extract user from context after auth middleware: `c.GetString("telegram_id")`
- Use `c.ShouldBindJSON(&struct)` for request parsing
- Use `c.Query("param")` for GET query params
- Always validate required fields before processing
- Always set `Cache-Control: no-store` on API responses: `c.Header("Cache-Control", "no-store")` before `c.JSON()`. Telegram Mini Apps may cache responses aggressively.
