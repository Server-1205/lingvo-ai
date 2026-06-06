package db

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/lingvo-ai/lingvo/internal/models"
)

func GetVocabulary(ctx context.Context, db *sqlx.DB, userID int) ([]models.VocabWord, error) {
	var words []models.VocabWord
	err := db.SelectContext(ctx, &words,
		"SELECT * FROM vocabulary WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	return words, nil
}

func AddVocabulary(ctx context.Context, db *sqlx.DB, userID int, word, translation, example, level string) error {
	if level == "" {
		level = "a1"
	}
	_, err := db.ExecContext(ctx, `
		INSERT INTO vocabulary (user_id, word, translation, example, level)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(user_id, word) DO UPDATE SET
			translation = excluded.translation,
			example = excluded.example,
			level = excluded.level
	`, userID, word, translation, example, level)
	return err
}

func DeleteVocabulary(ctx context.Context, db *sqlx.DB, userID, wordID int) error {
	_, err := db.ExecContext(ctx,
		"DELETE FROM vocabulary WHERE id = ? AND user_id = ?", wordID, userID)
	return err
}
