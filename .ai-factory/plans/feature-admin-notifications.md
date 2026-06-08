# Уведомления админу о платежах

**Branch:** `feature/admin-notifications`
**Created:** 2026-06-07 14:30

## Settings

- **Testing:** yes — unit tests для handlePayment с adminIDs
- **Logging:** verbose
- **Docs:** no
- **Roadmap Linkage:** Платежи — уведомления админу

## Анализ

**Текущее состояние:**
- `ADMIN_IDS` читается из `.env` и передаётся в `middleware.RateLimitMiddleware`
- `bot.Start()` не принимает `adminIDs` — уведомлений нет
- `handlePayment()` сохраняет подписку и отправляет подтверждение пользователю, но админ не знает о платеже

**Что нужно изменить:**
- `bot.Start()` → добавить параметр `adminIDs []int64`
- `processUpdate()` → передавать `adminIDs` в `handlePayment()`
- `handlePayment()` → после сохранения подписки слать уведомление админам
- `cmd/server/main.go` → передать `adminIDs` в `bot.Start()`

## Задачи

### Task 1: Передать adminIDs в бот

**Файлы:** `internal/bot/bot.go`, `internal/bot/handlers.go`, `internal/bot/payments.go`, `cmd/server/main.go`

**Что сделать:**
- `bot.go` — `Start()` принимает `adminIDs []int64`
- `bot.go` — передать `adminIDs` в `processUpdate()`
- `handlers.go` — `processUpdate()` принимает и передаёт `adminIDs` в `handlePayment()`
- `payments.go` — `handlePayment()` принимает `adminIDs`, после сохранения подписки шлёт уведомление каждому админу
- `main.go` — распарсить `ADMIN_IDS` и передать в `bot.Start()`

**Код уведомления в `handlePayment`:**
```go
adminMsg := fmt.Sprintf("💰 Новый платёж!\n\nПользователь: %d\nПлан: %s\nStars: %d\nДействует до: %s",
    telegramID, plan, starsAmount, expiresAt.Format("2006-01-02"))

for _, adminID := range adminIDs {
    msg := tgbotapi.NewMessage(adminID, adminMsg)
    msg.ParseMode = "Markdown"
    if _, err := bot.Send(msg); err != nil {
        sugar.Errorw("send admin notification", "error", err, "admin_id", adminID)
    }
}
```

**Логирование:**
- `[payment] admin notification sent to N admins`
- `[payment] admin notification failed for admin_id=X`

### Task 2: Тесты

**Файлы:** `internal/bot/payments_test.go`

**Что сделать:**
- Тест: handlePayment отправляет уведомление админу при успешном платеже
- Тест: handlePayment не падает при ошибке отправки админу
- Тест: handlePayment работает без adminIDs (пустой слайс)

## Commit Plan

1. `feat(bot): send admin notifications on successful payment`
