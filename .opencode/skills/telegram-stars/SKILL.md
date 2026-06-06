---
name: telegram-stars
description: Use when implementing Telegram Stars payment integration — creating invoice links, handling successful_payment updates, managing subscriptions, and processing Star transactions. NOT for general bot commands or chat features.
---

# Telegram Stars Payments

## Creating Invoice Link (bot side)

```go
func createInvoiceLink(bot *tgbotapi.BotAPI, userID int64, plan string) (string, error) {
    prices := map[string]struct {
        title       string
        description string
        amount      int
    }{
        "weekly":  {"Weekly Premium", "7 days unlimited access", 50},
        "monthly": {"Monthly Premium", "30 days unlimited access", 150},
    }

    p, ok := prices[plan]
    if !ok {
        return "", fmt.Errorf("unknown plan: %s", plan)
    }

    params := tgbotapi.CreateInvoiceLinkParams{
        Title:          p.title,
        Description:    p.description,
        Payload:        fmt.Sprintf("%d:%s", userID, plan),
        Currency:       "XTR",
        Prices:         []tgbotapi.Product{{Label: p.title, Amount: int64(p.amount)}},
    }

    return bot.CreateInvoiceLink(params)
}
```

## Handling successful_payment (bot handler)

```go
func handlePayment(update tgbotapi.Update, db *sqlx.DB) {
    payment := update.Message.SuccessfulPayment
    payload := payment.InvoicePayload // "userID:plan"
    parts := strings.Split(payload, ":")
    userID := parts[0]
    plan := parts[1]

    // Parse stars amount
    starsAmount := payment.TotalAmount

    // Calculate expiry
    duration := map[string]string{
        "weekly":  "+7 days",
        "monthly": "+30 days",
    }

    _, err := db.NamedExec(`
        INSERT INTO subscriptions (user_id, plan, stars_amount, expires_at)
        VALUES (:user_id, :plan, :stars, datetime('now', :duration))
        ON CONFLICT(user_id) DO UPDATE SET
            plan = excluded.plan,
            stars_amount = excluded.stars_amount,
            expires_at = excluded.expires_at
    `, map[string]interface{}{
        "user_id":   userID,
        "plan":      plan,
        "stars":     starsAmount,
        "duration":  duration[plan],
    })
}
```

## API Endpoint

```go
// POST /api/create-invoice
func createInvoiceHandler(db *sqlx.DB, bot *tgbotapi.BotAPI) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            Plan string `json:"plan" binding:"required,oneof=weekly monthly"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": "invalid plan"})
            return
        }
        userID := c.GetString("telegram_id")
        link, err := createInvoiceLink(bot, userID, req.Plan)
        if err != nil {
            c.JSON(500, gin.H{"error": "failed to create invoice"})
            return
        }
        c.JSON(200, gin.H{"invoice_link": link})
    }
}
```

## Commission Notes
- Apple/Google take ~32% of Stars
- Reinvest Stars in Telegram Ads (30% bonus)
- Subscription prices must account for commission
