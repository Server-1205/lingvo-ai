package bot

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/lingvo-ai/lingvo/internal/ai"
)

var botTokenGlobal string

func Start(database *sqlx.DB, botToken, webappURL string, sugar *zap.SugaredLogger, adminIDs []int64, aiClient *ai.Client) {
	botTokenGlobal = botToken
	defer func() {
		if r := recover(); r != nil {
			sugar.Errorw("bot Start panic", "recover", r)
		}
	}()

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		sugar.Errorw("bot init failed", "error", err)
		return
	}

	sugar.Infow("bot started", "username", bot.Self.UserName)

	go StartReminderScheduler(bot, database, sugar, webappURL)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{"message", "callback_query"}

	updates := bot.GetUpdatesChan(u)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ctx.Done():
			sugar.Info("bot shutting down")
			return
		case <-sigCh:
			sugar.Info("bot received signal, shutting down")
			cancel()
			return
		case update := <-updates:
			processUpdate(bot, database, webappURL, sugar, update, adminIDs, aiClient)
		}
	}
}

func processUpdate(bot *tgbotapi.BotAPI, database *sqlx.DB, webappURL string, sugar *zap.SugaredLogger, update tgbotapi.Update, adminIDs []int64, aiClient *ai.Client) {
	defer func() {
		if r := recover(); r != nil {
			sugar.Errorw("bot panic recovered", "recover", r, "update_id", update.UpdateID)
		}
	}()

	if update.Message == nil {
		return
	}

	sugar.Debugw("bot message", "text", update.Message.Text, "from", update.Message.From.ID, "update_id", update.UpdateID)

	if update.Message.SuccessfulPayment != nil {
		handlePayment(bot, database, sugar, update, adminIDs)
		return
	}

	if update.Message.IsCommand() {
		handleCommand(bot, database, webappURL, sugar, update, adminIDs, aiClient)
		return
	}
}
