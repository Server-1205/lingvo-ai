package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/db"
)

var durationMap = map[string]string{
	"weekly":  "+7 days",
	"monthly": "+30 days",
}

var planStars = map[string]int{
	"weekly":  300,
	"monthly": 800,
}

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

	text := fmt.Sprintf("📅 *Bugungi dars / Урок дня*\n\n"+
		"Daraja / Уровень: *%s*\n\n"+
		"Mini App ni oching va AI bilan suhbatlashing yoki lug'at bilan ishlang.\n\n"+
		"Открой Mini App и общайся с AI или работай со словарём.", strings.ToUpper(user.Level))

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = launchKeyboard(chatID, lang)
	if _, err := bot.Send(msg); err != nil {
		sugar.Errorw("send daily message", "error", err)
	}
}

func handlePayment(bot *tgbotapi.BotAPI, database *sqlx.DB, sugar *zap.SugaredLogger, update tgbotapi.Update) {
	if update.Message == nil || update.Message.SuccessfulPayment == nil {
		return
	}

	payment := update.Message.SuccessfulPayment
	payload := payment.InvoicePayload
	starsAmount := payment.TotalAmount
	chatID := update.Message.Chat.ID
	telegramID := update.Message.From.ID

	sugar.Infow("successful payment", "telegram_id", telegramID, "payload", payload, "stars", starsAmount)

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

	premiumStatus := "❌ Yo'q / Нет"
	if stats.IsPremium {
		premiumStatus = fmt.Sprintf("✅ %s", stats.SubscriptionExpiry)
	}

	var msgText string
	if lang == "ru" {
		msgText = fmt.Sprintf(
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
	} else {
		msgText = fmt.Sprintf(
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
