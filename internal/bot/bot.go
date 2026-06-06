package bot

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func Start(db *sqlx.DB, botToken string, sugar *zap.SugaredLogger) {
	_ = db
	_ = botToken
	sugar.Info("bot started (stub)")
}
