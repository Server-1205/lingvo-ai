package bot

import (
	"context"
	"fmt"
	"strings"
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

func sendReviewWebApp(bot *tgbotapi.BotAPI, chatID int64, lang string, webappURL string, sugar *zap.SugaredLogger) {
	texts := map[string]string{
		"uz": "📚 Eslatma! So'zlaringizni takrorlash vaqti keldi.",
		"ru": "📚 Напоминание! Пора повторить слова.",
	}
	text := texts["uz"]
	if t, ok := texts[lang]; ok {
		text = t
	}

	labels := map[string]string{
		"uz": "📖 Hozir takrorlash",
		"ru": "📖 Повторить сейчас",
	}
	label := labels["uz"]
	if l, ok := labels[lang]; ok {
		label = l
	}

	sendWebAppMessage(chatID, text, label, webappURL+"?startapp=review", sugar, "")
}

func StartReminderScheduler(bot *tgbotapi.BotAPI, database *sqlx.DB, sugar *zap.SugaredLogger, webappURL string) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	sugar.Info("reminder scheduler started")

	baseURL := webappURL
	if !strings.Contains(baseURL, "?") {
		baseURL += "?startapp=review"
	} else {
		baseURL += "&startapp=review"
	}

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
			sendWebAppMessage(
				user.TelegramID,
				reminderText(lang, dueCount),
				map[string]string{"uz": "📖 Hozir takrorlash", "ru": "📖 Повторить сейчас"}[lang],
				baseURL,
				sugar, "Markdown",
			)

			remindedMu.Lock()
			remindedToday[user.TelegramID] = time.Now()
			remindedMu.Unlock()

			sugar.Debugw("reminder sent", "user_id", user.ID, "telegram_id", user.TelegramID, "due", dueCount)
			sentCount++
		}

		sugar.Infow("reminder tick complete", "users_checked", len(users), "reminders_sent", sentCount)
	}
}
