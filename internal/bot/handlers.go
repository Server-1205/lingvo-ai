package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

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

func handleCommand(bot *tgbotapi.BotAPI, database *sqlx.DB, sugar *zap.SugaredLogger, update tgbotapi.Update) {
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

		msg := tgbotapi.NewMessage(chatID, welcomeText(lang))
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = launchKeyboard(chatID, lang)
		if _, err := bot.Send(msg); err != nil {
			sugar.Errorw("send start message", "error", err)
		}

	case "help":
		msg := tgbotapi.NewMessage(chatID, helpText(lang))
		msg.ParseMode = "Markdown"
		if _, err := bot.Send(msg); err != nil {
			sugar.Errorw("send help message", "error", err)
		}

	case "daily":
		handleDaily(bot, database, sugar, chatID, telegramID, lang)

	case "stats":
		handleStats(bot, database, sugar, chatID, telegramID, lang)

	default:
		msg := tgbotapi.NewMessage(chatID, "Unknown command. /help")
		if _, err := bot.Send(msg); err != nil {
			sugar.Errorw("send unknown command", "error", err)
		}
	}
}

func handleDaily(bot *tgbotapi.BotAPI, database *sqlx.DB, sugar *zap.SugaredLogger, chatID int64, telegramID int64, lang string) {
	user, err := db.GetUserByTelegramID(context.Background(), database, telegramID)
	if err != nil {
		sugar.Errorw("get user for daily", "error", err, "telegram_id", telegramID)
		sendMessage(bot, chatID, "Xatolik yuz berdi. / Ошибка.")
		return
	}

	dueCount, err := db.GetDueWordCount(context.Background(), database, user.ID)
	if err != nil {
		sugar.Errorw("get due count for daily", "error", err, "user_id", user.ID)
		dueCount = 0
	}

	sugar.Debugw("daily check", "user_id", user.ID, "due", dueCount)

	var msgText string
	if dueCount > 0 {
		dueMsgs := map[string]string{
			"uz": fmt.Sprintf("📅 *Bugungi dars*\n\n"+
				"Daraja: *%s*\n"+
				"📚 Takrorlash uchun *%d* ta so'z bor.\n\n"+
				"Quyidagi tugma orqali boshlang:", strings.ToUpper(user.Level), dueCount),
			"ru": fmt.Sprintf("📅 *Урок дня*\n\n"+
				"Уровень: *%s*\n"+
				"📚 Слов для повторения: *%d*\n\n"+
				"Начните прямо сейчас:", strings.ToUpper(user.Level), dueCount),
		}
		msgText = dueMsgs["uz"]
		if m, ok := dueMsgs[lang]; ok {
			msgText = m
		}

		msg := tgbotapi.NewMessage(chatID, msgText)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = reviewButton(lang)
		if _, err := bot.Send(msg); err != nil {
			sugar.Errorw("send daily message with review", "error", err)
		}
	} else {
		emptyMsgs := map[string]string{
			"uz": fmt.Sprintf("📅 *Bugungi dars*\n\n"+
				"Daraja: *%s*\n\n"+
				"✅ Takrorlash uchun so'z yo'q. Yangi so'zlar qo'shing yoki AI bilan suhbatlashing.\n\n"+
				"Mini App ni oching:", strings.ToUpper(user.Level)),
			"ru": fmt.Sprintf("📅 *Урок дня*\n\n"+
				"Уровень: *%s*\n\n"+
				"✅ Нет слов для повторения. Добавьте новые слова или пообщайтесь с AI.\n\n"+
				"Откройте Mini App:", strings.ToUpper(user.Level)),
		}
		msgText = emptyMsgs["uz"]
		if m, ok := emptyMsgs[lang]; ok {
			msgText = m
		}

		msg := tgbotapi.NewMessage(chatID, msgText)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = launchKeyboard(chatID, lang)
		if _, err := bot.Send(msg); err != nil {
			sugar.Errorw("send daily message", "error", err)
		}
	}
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

func launchKeyboard(chatID int64, lang string) interface{} {
	labels := map[string]string{
		"uz": "🚀 Lingvo AI ни ишга тушириш",
		"ru": "🚀 Запустить Lingvo AI",
	}

	label := labels["uz"]
	if l, ok := labels[lang]; ok {
		label = l
	}

	webAppURL := "https://t.me/lingvo_ai_bot/app"
	webAppData := map[string]interface{}{
		"text": label,
		"web_app": map[string]string{
			"url": webAppURL,
		},
	}

	data, _ := json.Marshal([]interface{}{[]interface{}{webAppData}})

	replyMarkup := fmt.Sprintf(`{"inline_keyboard":%s}`, string(data))

	var markup tgbotapi.InlineKeyboardMarkup
	_ = json.Unmarshal([]byte(replyMarkup), &markup)
	return markup
}
