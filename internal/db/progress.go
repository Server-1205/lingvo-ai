package db

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

func GetDailyProgress(ctx context.Context, db *sqlx.DB, userID int, date string) (map[string]int, error) {
	var row struct {
		MessagesSent int `db:"messages_sent"`
		WordsLearned int `db:"words_learned"`
		QuizzesTaken int `db:"quizzes_taken"`
	}
	err := db.GetContext(ctx, &row,
		`SELECT COALESCE(messages_sent,0) as messages_sent,
		        COALESCE(words_learned,0) as words_learned,
		        COALESCE(quizzes_taken,0) as quizzes_taken
		 FROM daily_progress WHERE user_id = ? AND date = ?`, userID, date)
	if err != nil {
		return map[string]int{"messages_sent": 0, "words_learned": 0, "quizzes_taken": 0}, nil
	}
	return map[string]int{
		"messages_sent": row.MessagesSent,
		"words_learned": row.WordsLearned,
		"quizzes_taken": row.QuizzesTaken,
	}, nil
}

func IncrementProgress(ctx context.Context, db *sqlx.DB, userID int, date, field string) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO daily_progress (user_id, date, `+field+`)
		VALUES (?, ?, 1)
		ON CONFLICT(user_id, date) DO UPDATE SET `+field+" = "+field+" + 1",
		userID, date)
	return err
}

func GetStreakDays(ctx context.Context, db *sqlx.DB, userID int) (int, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT date FROM daily_progress
		WHERE user_id = ? AND messages_sent > 0
		ORDER BY date DESC LIMIT 365`, userID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	streak := 0
	expected := ""
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			return 0, err
		}
		if expected != "" && date != expected {
			break
		}
		streak++
		parsed, _ := time.Parse("2006-01-02", date)
		expected = parsed.AddDate(0, 0, -1).Format("2006-01-02")
	}
	return streak, nil
}
