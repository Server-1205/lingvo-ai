package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func GetMessageCount(ctx context.Context, db *sqlx.DB, userID int, date string) (int, error) {
	var count int
	err := db.GetContext(ctx, &count,
		"SELECT COALESCE(count, 0) FROM messages WHERE user_id = ? AND date = ?",
		userID, date)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func IncrementMessageCount(ctx context.Context, db *sqlx.DB, userID int, date string) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO messages (user_id, date, count)
		VALUES (?, ?, 1)
		ON CONFLICT(user_id, date) DO UPDATE SET count = count + 1
	`, userID, date)
	return err
}
