package db

import (
	"os"
	"strings"

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

	statements := strings.Split(string(schema), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			// ALTER TABLE ADD COLUMN may fail if column already exists — ignore
			if strings.HasPrefix(stmt, "ALTER TABLE") {
				sugar.Warnw("migration skipped (column may already exist)", "stmt", stmt[:60], "error", err)
			} else {
				sugar.Fatalw("exec schema", "stmt", stmt[:60], "error", err)
			}
		}
	}

	sugar.Info("database migrated")
}
