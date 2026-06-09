package db

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/lingvo-ai/lingvo/internal/models"
)

func UpsertUser(ctx context.Context, db *sqlx.DB, telegramID int64, username, lang string) error {
	if lang != "uz" && lang != "ru" {
		lang = "uz"
	}
	_, err := db.ExecContext(ctx, `
		INSERT INTO users (telegram_id, username, lang)
		VALUES (?, ?, ?)
		ON CONFLICT(telegram_id) DO UPDATE SET
			username = excluded.username,
			lang = excluded.lang,
			updated_at = CURRENT_TIMESTAMP
	`, telegramID, username, lang)
	return err
}

func GetUserByTelegramID(ctx context.Context, db *sqlx.DB, telegramID int64) (*models.User, error) {
	var u models.User
	err := db.GetContext(ctx, &u, "SELECT * FROM users WHERE telegram_id = ?", telegramID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUserLang(ctx context.Context, db *sqlx.DB, userID int, lang string) error {
	if lang != "uz" && lang != "ru" {
		return nil
	}
	_, err := db.ExecContext(ctx,
		"UPDATE users SET lang = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		lang, userID)
	return err
}

func GetAllUsers(ctx context.Context, db *sqlx.DB) ([]models.User, error) {
	var users []models.User
	err := db.SelectContext(ctx, &users, "SELECT * FROM users ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	return users, nil
}
