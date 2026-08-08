package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/db"
)

var (
	durationMap = map[string]string{
		"weekly":  "+7 days",
		"monthly": "+30 days",
	}
	processedPayments sync.Map
)

func handlePayment(bot *tgbotapi.BotAPI, database *sqlx.DB, sugar *zap.SugaredLogger, update tgbotapi.Update, adminIDs []int64) {
	if update.Message == nil || update.Message.SuccessfulPayment == nil {
		return
	}

	payment := update.Message.SuccessfulPayment
	payload := payment.InvoicePayload
	chargeID := payment.TelegramPaymentChargeID
	starsAmount := payment.TotalAmount
	chatID := update.Message.Chat.ID
	telegramID := update.Message.From.ID

	if _, loaded := processedPayments.LoadOrStore(chargeID, true); loaded {
		sugar.Warnw("duplicate payment ignored", "charge_id", chargeID, "telegram_id", telegramID)
		return
	}

	sugar.Infow("successful payment", "telegram_id", telegramID, "payload", payload, "stars", starsAmount, "charge_id", chargeID)

	parts := strings.SplitN(payload, "_", 2)
	if len(parts) != 2 {
		sugar.Errorw("invalid payment payload", "payload", payload)
		sendMessage(bot, chatID, "To'lovda xatolik. / Ошибка оплаты.")
		return
	}

	plan := parts[0]
	userIDStr := parts[1]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		sugar.Errorw("invalid user id in payload", "user_id", userIDStr)
		return
	}

	if _, ok := durationMap[plan]; !ok {
		sugar.Errorw("unknown plan", "plan", plan)
		return
	}

	expiresAt := time.Now().UTC().AddDate(0, 0, 7)
	if plan == "monthly" {
		expiresAt = time.Now().UTC().AddDate(0, 1, 0)
	}

	ctx := context.Background()
	if err := db.SaveSubscription(ctx, database, userID, plan, starsAmount, expiresAt.Format("2006-01-02 15:04:05")); err != nil {
		sugar.Errorw("save subscription", "error", err)
		sendMessage(bot, chatID, "Xatolik. / Ошибка.")
		return
	}

	msg := buildSubscriptionConfirmation(chatID, update.Message.From.LanguageCode, plan, expiresAt)
	if _, err := bot.Send(msg); err != nil {
		sugar.Errorw("send payment confirmation", "error", err)
	}

	adminMsg := fmt.Sprintf("💰 *Новый платёж!*\n\nПользователь: `%d`\nПлан: *%s*\nStars: *%d*\nДействует до: `%s`",
		telegramID, plan, starsAmount, expiresAt.Format("2006-01-02"))

	for _, adminID := range adminIDs {
		admMsg := tgbotapi.NewMessage(adminID, adminMsg)
		admMsg.ParseMode = "Markdown"
		if _, err := bot.Send(admMsg); err != nil {
			sugar.Errorw("send admin notification", "error", err, "admin_id", adminID)
		} else {
			sugar.Infow("admin notification sent", "admin_id", adminID, "telegram_id", telegramID)
		}
	}
}

func buildSubscriptionConfirmation(chatID int64, lang, plan string, expiresAt time.Time) tgbotapi.MessageConfig {
	msgs := map[string]string{
		"uz": fmt.Sprintf("✅ *Обуна фаоллаштирилди!*\n\n"+
			"Режа: *%s*\n"+
			"Амал қилади: *%s*\n\n"+
			"Энди сиз чексиз AI хабарлардан фойдалана оласиз! 🎉",
			planTitle(plan, "uz"), expiresAt.Format("2006-01-02")),
		"ru": fmt.Sprintf("✅ *Подписка активирована!*\n\n"+
			"План: *%s*\n"+
			"Действует до: *%s*\n\n"+
			"Теперь у вас безлимитные AI-сообщения! 🎉",
			planTitle(plan, "ru"), expiresAt.Format("2006-01-02")),
	}

	msg := tgbotapi.NewMessage(chatID, msgs[lang])
	msg.ParseMode = "Markdown"
	return msg
}

func planTitle(plan, lang string) string {
	titles := map[string]map[string]string{
		"weekly":  {"uz": "Ҳафталик", "ru": "Недельная"},
		"monthly": {"uz": "Ойлик", "ru": "Месячная"},
	}
	if t, ok := titles[plan][lang]; ok {
		return t
	}
	return plan
}
