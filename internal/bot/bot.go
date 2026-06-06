package bot

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func Start(database *sqlx.DB, botToken string, sugar *zap.SugaredLogger) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		sugar.Fatalw("bot init", "error", err)
	}

	sugar.Infow("bot started", "username", bot.Self.UserName)

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
			processUpdate(bot, database, sugar, update)
		}
	}
}

func processUpdate(bot *tgbotapi.BotAPI, database *sqlx.DB, sugar *zap.SugaredLogger, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.SuccessfulPayment != nil {
		handlePayment(bot, database, sugar, update)
		return
	}

	if update.Message.IsCommand() {
		handleCommand(bot, database, sugar, update)
		return
	}
}
