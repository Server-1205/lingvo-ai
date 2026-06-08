package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/ai"
	"github.com/lingvo-ai/lingvo/internal/db"
	"github.com/lingvo-ai/lingvo/internal/models"
)

func welcomeText(lang string) string {
	msgs := map[string]string{
		"uz": "👋 Салом! Men Lingvo AI — sizning shaxsiy ingliz tili o'qituvchingiz.\n\n" +
			"Ingliz tilida yozing, men хatolaringizni tuzataman va tushuntiraman.\n\n" +
			"🚀 Quyidagi tugma orqali boshlang:",
		"ru": "👋 Привет! Я Lingvo AI — твой личный учитель английского.\n\n" +
			"Пиши на английском, я буду исправлять ошибки и объяснять.\n\n" +
			"🚀 Нажми кнопку ниже, чтобы начать:",
	}
	if m, ok := msgs[lang]; ok {
		return m
	}
	return msgs["uz"]
}

func helpText(lang string) string {
	msgs := map[string]string{
		"uz": "📚 *Lingvo AI — ёрдам*\n\n" +
			"/start — Бошлаш\n" +
			"/daily — Бугунги дарс\n" +
			"/stats — Статистика\n" +
			"/help — Ёрдам\n\n" +
			"Шунингдек, Mini App орқали сўз бойлигингизни оширинг ва AI билан суҳбатлашинг.",
		"ru": "📚 *Lingvo AI — помощь*\n\n" +
			"/start — Начать\n" +
			"/daily — Урок дня\n" +
			"/stats — Статистика\n" +
			"/help — Помощь\n\n" +
			"Используй Mini App для пополнения словарного запаса и общения с AI.",
	}
	if m, ok := msgs[lang]; ok {
		return m
	}
	return msgs["uz"]
}

func handleCommand(bot *tgbotapi.BotAPI, database *sqlx.DB, webappURL string, sugar *zap.SugaredLogger, update tgbotapi.Update, aiClient *ai.Client) {
	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	chatID := update.Message.Chat.ID
	telegramID := update.Message.From.ID
	username := update.Message.From.UserName
	lang := update.Message.From.LanguageCode
	if lang != "uz" && lang != "ru" {
		lang = "uz"
	}

	cmd := update.Message.Command()
	sugar.Debugw("bot command", "cmd", cmd, "telegram_id", telegramID, "username", username)

	ctx := context.Background()

	switch cmd {
	case "start":
		if err := db.UpsertUser(ctx, database, telegramID, username, lang); err != nil {
			sugar.Errorw("upsert user", "error", err, "telegram_id", telegramID)
			sendMessage(bot, chatID, "Xatolik yuz berdi. / Ошибка.")
			return
		}
		sugar.Infow("new user", "telegram_id", telegramID, "username", username)

		labels := map[string]string{
			"uz": "🚀 Lingvo AI ни ишга тушириш",
			"ru": "🚀 Запустить Lingvo AI",
		}
		btnLabel := labels["uz"]
		if l, ok := labels[lang]; ok {
			btnLabel = l
		}
		sendWebAppMessage(chatID, welcomeText(lang), btnLabel, webappURL, sugar, "Markdown")

	case "help":
		msg := tgbotapi.NewMessage(chatID, helpText(lang))
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			sugar.Errorw("send help message", "error", err)
		}

	case "daily":
		handleDaily(bot, database, webappURL, sugar, chatID, telegramID, lang, aiClient)

	case "stats":
		handleStats(bot, database, sugar, chatID, telegramID, lang)

	default:
		msg := tgbotapi.NewMessage(chatID, "Unknown command. /help")
		if _, err := bot.Send(msg); err != nil {
			sugar.Errorw("send unknown command", "error", err)
		}
	}
}

func handleDaily(bot *tgbotapi.BotAPI, database *sqlx.DB, webappURL string, sugar *zap.SugaredLogger, chatID int64, telegramID int64, lang string, aiClient *ai.Client) {
	dailyMsgs := map[string]string{
		"uz": "📅 *Kunlik dars*\n\nYangi mavzuni o'rganish va mashqlarni bajarish uchun quyidagi tugmani bosing:",
		"ru": "📅 *Урок дня*\n\nНажмите кнопку ниже, чтобы изучить новую тему и выполнить упражнения:",
	}
	msgText := dailyMsgs["uz"]
	if m, ok := dailyMsgs[lang]; ok {
		msgText = m
	}

	dailyURL := webappURL
	if strings.Contains(dailyURL, "?") {
		dailyURL += "&startapp=daily"
	} else {
		dailyURL += "?startapp=daily"
	}

	btnLabel := map[string]string{"uz": "📖 Darsni boshlash", "ru": "📖 Начать урок"}
	lbl := btnLabel["uz"]
	if l, ok := btnLabel[lang]; ok {
		lbl = l
	}

	sendWebAppMessage(chatID, msgText, lbl, dailyURL, sugar, "Markdown")
}

func formatStatsMessage(stats *models.UserStats, lang string, daysActive int) string {
	premiumStatus := "❌ Yo'q / Нет"
	if stats.IsPremium {
		premiumStatus = fmt.Sprintf("✅ %s", stats.SubscriptionExpiry)
	}

	if lang == "ru" {
		return fmt.Sprintf(
			"📊 *Статистика*\n\n"+
				"🎯 Уровень: *%s*\n"+
				"💬 Всего сообщений: *%d*\n"+
				"📚 Слов в словаре: *%d*\n"+
				"📝 На повторении сегодня: *%d*\n"+
				"🔥 Streak: *%d дней*\n"+
				"👤 Аккаунту: *%d дней*\n"+
				"⭐ Подписка: *%s*\n",
			stats.Level, stats.TotalMessages, stats.TotalWords,
			stats.WordsDueToday, stats.StreakDays, daysActive, premiumStatus)
	}

	return fmt.Sprintf(
		"📊 *Statistika*\n\n"+
			"🎯 Daraja: *%s*\n"+
			"💬 Jami xabarlar: *%d*\n"+
			"📚 Lug'atdagi so'zlar: *%d*\n"+
			"📝 Bugungi takrorlash: *%d*\n"+
			"🔥 Streak: *%d kun*\n"+
			"👤 Akkount: *%d kun*\n"+
			"⭐ Obuna: *%s*\n",
		stats.Level, stats.TotalMessages, stats.TotalWords,
		stats.WordsDueToday, stats.StreakDays, daysActive, premiumStatus)
}

func handleStats(bot *tgbotapi.BotAPI, database *sqlx.DB, sugar *zap.SugaredLogger, chatID int64, telegramID int64, lang string) {
	sugar.Debugw("bot: /stats", "telegram_id", telegramID)

	ctx := context.Background()

	user, err := db.GetUserByTelegramID(ctx, database, telegramID)
	if err != nil {
		sugar.Errorw("bot: /stats — get user failed", "telegram_id", telegramID, "error", err)
		sendMessage(bot, chatID, "Xatolik yuz berdi. / Ошибка.")
		return
	}

	stats, err := db.GetUserStats(ctx, database, user.ID)
	if err != nil {
		sugar.Errorw("bot: /stats — get stats failed", "user_id", user.ID, "error", err)
		sendMessage(bot, chatID, "Xatolik yuz berdi. / Ошибка.")
		return
	}

	sugar.Debugw("bot: /stats — fetched",
		"user_id", user.ID,
		"messages", stats.TotalMessages,
		"words", stats.TotalWords,
		"streak", stats.StreakDays,
		"premium", stats.IsPremium)

	daysActive := int(time.Since(user.CreatedAt).Hours() / 24)
	if daysActive < 1 {
		daysActive = 1
	}

	msgText := formatStatsMessage(stats, lang, daysActive)
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ParseMode = "Markdown"
	if _, err := bot.Send(msg); err != nil {
		sugar.Errorw("bot: /stats — send failed", "error", err)
	}

	sugar.Infow("bot: /stats — done", "telegram_id", telegramID)
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		// silent
	}
}

func sendWebAppMessage(chatID int64, text, buttonLabel, webAppURL string, sugar *zap.SugaredLogger, parseMode string) {
	body := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": parseMode,
		"reply_markup": map[string]interface{}{
			"inline_keyboard": [][]map[string]interface{}{
				{{"text": buttonLabel, "web_app": map[string]string{"url": webAppURL}}},
			},
		},
	}

	payload, _ := json.Marshal(body)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botTokenGlobal)

	resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		sugar.Errorw("send webapp message http", "error", err)
		return
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		sugar.Errorw("send webapp message failed", "status", resp.StatusCode, "body", string(raw))
	}
}

func cleanJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	return strings.TrimSpace(raw)
}
