package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lingvo-ai/lingvo/internal/models"
)

func SaveIeltsScore(ctx context.Context, db sqlx.ExtContext, userID int, module string, bandScore float64, details, prompt, userResponse, feedback string) error {
	_, err := sqlx.NamedExecContext(ctx, db, `
		INSERT INTO ielts_scores (user_id, module, band_score, details, prompt, user_response, feedback)
		VALUES (:user_id, :module, :band_score, :details, :prompt, :user_response, :feedback)`,
		map[string]interface{}{
			"user_id":       userID,
			"module":        module,
			"band_score":    bandScore,
			"details":       details,
			"prompt":        prompt,
			"user_response": userResponse,
			"feedback":      feedback,
		})
	if err != nil {
		return fmt.Errorf("save ielts score: %w", err)
	}
	return nil
}

func GetIeltsScores(ctx context.Context, db sqlx.ExtContext, userID int, module string, limit, offset int) ([]models.IeltsScoreEntry, int, error) {
	where := "user_id = ?"
	args := []interface{}{userID}
	if module != "" {
		where += " AND module = ?"
		args = append(args, module)
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM ielts_scores WHERE %s", where)
	if err := sqlx.GetContext(ctx, db, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("count ielts scores: %w", err)
	}

	query := fmt.Sprintf("SELECT * FROM ielts_scores WHERE %s ORDER BY created_at DESC LIMIT ? OFFSET ?", where)
	args = append(args, limit, offset)

	var entries []models.IeltsScoreEntry
	if err := sqlx.SelectContext(ctx, db, &entries, query, args...); err != nil {
		return nil, 0, fmt.Errorf("get ielts scores: %w", err)
	}
	if entries == nil {
		entries = []models.IeltsScoreEntry{}
	}

	return entries, total, nil
}

func GetIeltsScoreStats(ctx context.Context, db sqlx.ExtContext, userID int) (*models.IeltsScoreStats, error) {
	stats := &models.IeltsScoreStats{}

	modules := []string{"writing_task1", "writing_task2", "speaking", "reading"}
	for _, m := range modules {
		var avg float64
		err := sqlx.GetContext(ctx, db, &avg,
			"SELECT COALESCE(AVG(band_score), 0) FROM ielts_scores WHERE user_id = ? AND module = ?", userID, m)
		if err != nil {
			return nil, fmt.Errorf("avg ielts %s: %w", m, err)
		}
		switch m {
		case "writing_task1":
			stats.WritingTask1Avg = avg
		case "writing_task2":
			stats.WritingTask2Avg = avg
		case "speaking":
			stats.SpeakingAvg = avg
		case "reading":
			stats.ReadingAvg = avg
		}
	}

	var total int
	if err := sqlx.GetContext(ctx, db, &total,
		"SELECT COUNT(*) FROM ielts_scores WHERE user_id = ?", userID); err != nil {
		return nil, fmt.Errorf("count ielts practices: %w", err)
	}
	stats.TotalPractices = total

	return stats, nil
}
