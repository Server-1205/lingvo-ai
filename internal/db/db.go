package db

import (
	"os"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func InitDB(dbPath string) (*sqlx.DB, error) {
	return sqlx.Connect("sqlite", dbPath)
}

func Migrate(db *sqlx.DB, sugar *zap.SugaredLogger) {
	schema, err := os.ReadFile("internal/db/schema.sql")
	if err != nil {
		sugar.Fatalw("read schema", "error", err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		sugar.Fatalw("exec schema", "error", err)
	}
	sugar.Info("database migrated")
}
