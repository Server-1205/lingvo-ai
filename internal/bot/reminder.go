package bot

import (
	"context"
	"fmt"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/db"
)

var (
	remindedMu    sync.Mutex
	remindedToday = map[int64]time.Time{}
)

func reminderText(lang string, count int) string {
	msgs := map[string]string{
		"uz": fmt.Sprintf("📚 Eslatma! Sizda %d ta so'z takrorlash uchun.\n\n"+
			"Kunlik takrorlashni boshlash uchun quyidagi tugmani bosing:", count),
		"ru": fmt.Sprintf("📚 Напоминание! У вас %d слов для повторения.\n\n"+
			"Нажмите кнопку ниже, чтобы начать ежедневное повторение:", count),
	}
	if m, ok := msgs[lang]; ok {
		return m
	}
	return msgs["uz"]
}

func reviewButton(lang string) tgbotapi.InlineKeyboardMarkup {
	labels := map[string]string{
		"uz": "📖 Hozir takrorlash",
		"ru": "📖 Повторить сейчас",
	}
	label := labels["uz"]
	if l, ok := labels[lang]; ok {
		label = l
	}

	deepLink := "https://t.me/lingvo_ai_bot/app?startapp=review"

	btn := tgbotapi.NewInlineKeyboardButtonURL(label, deepLink)
	row := tgbotapi.NewInlineKeyboardRow(btn)
	return tgbotapi.NewInlineKeyboardMarkup(row)
}

func StartReminderScheduler(bot *tgbotapi.BotAPI, database *sqlx.DB, sugar *zap.SugaredLogger) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	sugar.Info("reminder scheduler started")

	for range ticker.C {
		sugar.Debug("reminder tick")

		users, err := db.GetAllUsers(context.Background(), database)
		if err != nil {
			sugar.Errorw("reminder: get all users", "error", err)
			continue
		}

		sentCount := 0
		for _, user := range users {
			remindedMu.Lock()
			lastReminded, already := remindedToday[user.TelegramID]
			remindedMu.Unlock()

			if already && lastReminded.Truncate(24*time.Hour).Equal(time.Now().Truncate(24*time.Hour)) {
				continue
			}

			dueCount, err := db.GetDueWordCount(context.Background(), database, user.ID)
			if err != nil {
				sugar.Errorw("reminder: get due count", "error", err, "user_id", user.ID)
				continue
			}
			if dueCount == 0 {
				continue
			}

			lang := user.Lang
			msg := tgbotapi.NewMessage(user.TelegramID, reminderText(lang, dueCount))
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = reviewButton(lang)

			if _, err := bot.Send(msg); err != nil {
				sugar.Errorw("reminder: send", "error", err, "user_id", user.ID, "telegram_id", user.TelegramID)
				continue
			}

			remindedMu.Lock()
			remindedToday[user.TelegramID] = time.Now()
			remindedMu.Unlock()

			sugar.Debugw("reminder sent", "user_id", user.ID, "telegram_id", user.TelegramID, "due", dueCount)
			sentCount++
		}

		sugar.Infow("reminder tick complete", "users_checked", len(users), "reminders_sent", sentCount)
	}
}
