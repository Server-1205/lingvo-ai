package db

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func TestSaveAndGetIeltsScores(t *testing.T) {
	db := setupIELTSTestDB(t)
	defer db.Close()
	ctx := context.Background()

	insertUser(t, db, ctx, 1, 12345)

	err := SaveIeltsScore(ctx, db, 1, "writing_task1", 6.5, `{"task_achievement":6.5}`, "Describe a chart", "The chart shows...", "Good work")
	if err != nil {
		t.Fatalf("SaveIeltsScore failed: %v", err)
	}

	err = SaveIeltsScore(ctx, db, 1, "speaking", 7.0, `{"fluency_coherence":7.0}`, "Tell me about yourself", "I am a student", "Well done")
	if err != nil {
		t.Fatalf("SaveIeltsScore failed: %v", err)
	}

	entries, total, err := GetIeltsScores(ctx, db, 1, "", 10, 0)
	if err != nil {
		t.Fatalf("GetIeltsScores failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].Module != "speaking" {
		t.Errorf("expected latest entry module 'speaking', got %s", entries[0].Module)
	}

	entries, total, err = GetIeltsScores(ctx, db, 1, "writing_task1", 10, 0)
	if err != nil {
		t.Fatalf("GetIeltsScores failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1 for writing_task1, got %d", total)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry for writing_task1, got %d", len(entries))
	}
}

func TestGetIeltsScores_NoData(t *testing.T) {
	db := setupIELTSTestDB(t)
	defer db.Close()
	ctx := context.Background()

	insertUser(t, db, ctx, 2, 99999)

	entries, total, err := GetIeltsScores(ctx, db, 2, "", 10, 0)
	if err != nil {
		t.Fatalf("GetIeltsScores failed: %v", err)
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
	if entries == nil {
		t.Error("expected non-nil empty slice")
	}
}

func TestGetIeltsScores_Pagination(t *testing.T) {
	db := setupIELTSTestDB(t)
	defer db.Close()
	ctx := context.Background()

	insertUser(t, db, ctx, 3, 11111)

	for i := 0; i < 5; i++ {
		err := SaveIeltsScore(ctx, db, 3, "reading", 7.5, `{}`, "", "", "")
		if err != nil {
			t.Fatalf("SaveIeltsScore failed: %v", err)
		}
	}

	entries, total, err := GetIeltsScores(ctx, db, 3, "", 2, 0)
	if err != nil {
		t.Fatalf("GetIeltsScores failed: %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries (page 1), got %d", len(entries))
	}

	entries, total, err = GetIeltsScores(ctx, db, 3, "", 2, 2)
	if err != nil {
		t.Fatalf("GetIeltsScores failed: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries (page 2), got %d", len(entries))
	}
}

func TestGetIeltsScoreStats(t *testing.T) {
	db := setupIELTSTestDB(t)
	defer db.Close()
	ctx := context.Background()

	insertUser(t, db, ctx, 4, 22222)

	SaveIeltsScore(ctx, db, 4, "writing_task1", 6.0, `{}`, "", "", "")
	SaveIeltsScore(ctx, db, 4, "writing_task1", 7.0, `{}`, "", "", "")
	SaveIeltsScore(ctx, db, 4, "writing_task2", 6.5, `{}`, "", "", "")
	SaveIeltsScore(ctx, db, 4, "speaking", 7.5, `{}`, "", "", "")
	SaveIeltsScore(ctx, db, 4, "reading", 8.0, `{}`, "", "", "")

	stats, err := GetIeltsScoreStats(ctx, db, 4)
	if err != nil {
		t.Fatalf("GetIeltsScoreStats failed: %v", err)
	}

	if stats.WritingTask1Avg != 6.5 {
		t.Errorf("expected WritingTask1Avg 6.5, got %f", stats.WritingTask1Avg)
	}
	if stats.WritingTask2Avg != 6.5 {
		t.Errorf("expected WritingTask2Avg 6.5, got %f", stats.WritingTask2Avg)
	}
	if stats.SpeakingAvg != 7.5 {
		t.Errorf("expected SpeakingAvg 7.5, got %f", stats.SpeakingAvg)
	}
	if stats.ReadingAvg != 8.0 {
		t.Errorf("expected ReadingAvg 8.0, got %f", stats.ReadingAvg)
	}
	if stats.TotalPractices != 5 {
		t.Errorf("expected TotalPractices 5, got %d", stats.TotalPractices)
	}
}

func TestGetIeltsScoreStats_NoData(t *testing.T) {
	db := setupIELTSTestDB(t)
	defer db.Close()
	ctx := context.Background()

	insertUser(t, db, ctx, 5, 33333)

	stats, err := GetIeltsScoreStats(ctx, db, 5)
	if err != nil {
		t.Fatalf("GetIeltsScoreStats failed: %v", err)
	}

	if stats.WritingTask1Avg != 0 {
		t.Errorf("expected 0, got %f", stats.WritingTask1Avg)
	}
	if stats.TotalPractices != 0 {
		t.Errorf("expected 0, got %d", stats.TotalPractices)
	}
}

func setupIELTSTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Connect("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to connect to in-memory db: %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		telegram_id INTEGER NOT NULL UNIQUE,
		username TEXT DEFAULT '',
		lang TEXT NOT NULL DEFAULT 'uz',
		level TEXT NOT NULL DEFAULT 'a1',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS ielts_scores (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL REFERENCES users(id),
		module TEXT NOT NULL CHECK(module IN ('writing_task1','writing_task2','speaking','reading')),
		band_score REAL NOT NULL,
		details TEXT DEFAULT '{}',
		prompt TEXT DEFAULT '',
		user_response TEXT DEFAULT '',
		feedback TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_ielts_user_module ON ielts_scores(user_id, module);
	CREATE INDEX IF NOT EXISTS idx_ielts_user_date ON ielts_scores(user_id, created_at);`

	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("failed to exec schema: %v", err)
	}

	return db
}

func insertUser(t *testing.T, db *sqlx.DB, ctx context.Context, id int, telegramID int64) {
	t.Helper()
	_, err := db.ExecContext(ctx,
		`INSERT INTO users (id, telegram_id, username, lang, level) VALUES (?, ?, 'testuser', 'uz', 'b1')`,
		id, telegramID)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
}
