package db

import (
	"context"
	"strings"

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

func AddVocabulary(ctx context.Context, db *sqlx.DB, userID int, word, translation, example, level, translationRu, exampleRu string) error {
	if level == "" {
		level = "a1"
	}
	word = strings.ToLower(word)
	_, err := db.ExecContext(ctx, `
		INSERT INTO vocabulary (user_id, word, translation, example, level, translation_ru, example_ru)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, word) DO UPDATE SET
			translation = excluded.translation,
			example = excluded.example,
			level = excluded.level,
			translation_ru = excluded.translation_ru,
			example_ru = excluded.example_ru
	`, userID, word, translation, example, level, translationRu, exampleRu)
	return err
}

func DeleteVocabulary(ctx context.Context, db *sqlx.DB, userID, wordID int) error {
	_, err := db.ExecContext(ctx,
		"DELETE FROM vocabulary WHERE id = ? AND user_id = ?", wordID, userID)
	return err
}

func GetWordsMissingRu(ctx context.Context, db *sqlx.DB) ([]models.VocabWord, error) {
	var words []models.VocabWord
	err := db.SelectContext(ctx, &words,
		"SELECT * FROM vocabulary WHERE translation_ru IS NULL OR translation_ru = ''")
	if err != nil {
		return nil, err
	}
	return words, nil
}

func UpdateWordRuTranslation(ctx context.Context, db *sqlx.DB, wordID int, translationRu, exampleRu string) error {
	_, err := db.ExecContext(ctx,
		"UPDATE vocabulary SET translation_ru = ?, example_ru = ? WHERE id = ?",
		translationRu, exampleRu, wordID)
	return err
}

func GetDueWordCount(ctx context.Context, db *sqlx.DB, userID int) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM vocabulary
		WHERE user_id = ?
		  AND (next_review IS NULL OR next_review <= datetime('now'))
	`, userID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetDueWords(ctx context.Context, db *sqlx.DB, userID, limit int) ([]models.VocabWord, error) {
	var words []models.VocabWord
	err := db.SelectContext(ctx, &words, `
		SELECT * FROM vocabulary
		WHERE user_id = ?
		  AND (next_review IS NULL OR next_review <= datetime('now'))
		ORDER BY next_review ASC
		LIMIT ?
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	return words, nil
}

func UpdateReview(ctx context.Context, db *sqlx.DB, wordID, userID int, reviewCount, interval int, easeFactor float64, nextReview string) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE vocabulary
		SET review_count = ?, interval = ?, ease_factor = ?,
		    last_reviewed_at = datetime('now'), next_review = ?
		WHERE id = ? AND user_id = ?
	`, reviewCount, interval, easeFactor, nextReview, wordID, userID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO daily_progress (user_id, date, words_learned)
		VALUES (?, date('now'), 1)
		ON CONFLICT(user_id, date) DO UPDATE SET
			words_learned = words_learned + 1
	`, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
