package db

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func setupErrorTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db, err := sqlx.Connect("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS error_history (
		id              INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id         INTEGER NOT NULL,
		original        TEXT NOT NULL,
		corrected       TEXT NOT NULL,
		category        TEXT NOT NULL DEFAULT 'grammar',
		severity        TEXT NOT NULL DEFAULT 'minor',
		rule_violated   TEXT DEFAULT '',
		learning_tip    TEXT DEFAULT '',
		context         TEXT DEFAULT '',
		created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_error_user_date ON error_history(user_id, created_at);
	CREATE INDEX IF NOT EXISTS idx_error_category ON error_history(user_id, category);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatal(err)
	}

	return db
}

func TestSaveError(t *testing.T) {
	db := setupErrorTestDB(t)
	ctx := context.Background()

	err := SaveError(ctx, db, 1, "I go", "I went", "grammar", "major", "Past Simple", "Use 'went' for past", "")
	if err != nil {
		t.Fatalf("SaveError failed: %v", err)
	}

	var count int
	if err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM error_history WHERE user_id = 1"); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 error, got %d", count)
	}
}

func TestSaveErrorMultiple(t *testing.T) {
	db := setupErrorTestDB(t)
	ctx := context.Background()

	errors := []struct {
		userID   int
		original string
		corrected string
		category string
		severity string
	}{
		{1, "I goes", "I go", "grammar", "major"},
		{1, "yesterday I", "I yesterday", "word_order", "minor"},
		{2, "she don't", "she doesn't", "grammar", "critical"},
	}

	for _, e := range errors {
		if err := SaveError(ctx, db, e.userID, e.original, e.corrected, e.category, e.severity, "", "", ""); err != nil {
			t.Fatalf("SaveError failed: %v", err)
		}
	}

	var total int
	db.GetContext(ctx, &total, "SELECT COUNT(*) FROM error_history")
	if total != 3 {
		t.Errorf("expected 3 total, got %d", total)
	}

	var user1Count int
	db.GetContext(ctx, &user1Count, "SELECT COUNT(*) FROM error_history WHERE user_id = 1")
	if user1Count != 2 {
		t.Errorf("expected 2 for user 1, got %d", user1Count)
	}
}

func TestGetErrorStats(t *testing.T) {
	db := setupErrorTestDB(t)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		SaveError(ctx, db, 1, "I goes", "I go", "grammar", "major", "Subject-Verb Agreement", "", "")
	}
	for i := 0; i < 3; i++ {
		SaveError(ctx, db, 1, "yesterday I", "I yesterday", "word_order", "minor", "", "", "")
	}
	for i := 0; i < 2; i++ {
		SaveError(ctx, db, 1, "she don't", "she doesn't", "grammar", "critical", "Subject-Verb Agreement", "", "")
	}

	stats, err := GetErrorStats(ctx, db, 1, 30)
	if err != nil {
		t.Fatalf("GetErrorStats failed: %v", err)
	}

	if stats.TotalErrors != 10 {
		t.Errorf("expected TotalErrors=10, got %d", stats.TotalErrors)
	}
	if stats.CategoryCounts["grammar"] != 7 {
		t.Errorf("expected grammar=7, got %d", stats.CategoryCounts["grammar"])
	}
	if stats.CategoryCounts["word_order"] != 3 {
		t.Errorf("expected word_order=3, got %d", stats.CategoryCounts["word_order"])
	}
	if stats.SeverityCounts["major"] != 5 {
		t.Errorf("expected major=5, got %d", stats.SeverityCounts["major"])
	}
	if len(stats.MostFrequentRules) == 0 {
		t.Error("expected at least 1 frequent rule")
	}
}

func TestGetErrorStatsEmpty(t *testing.T) {
	db := setupErrorTestDB(t)
	ctx := context.Background()

	stats, err := GetErrorStats(ctx, db, 999, 30)
	if err != nil {
		t.Fatalf("GetErrorStats failed: %v", err)
	}

	if stats.TotalErrors != 0 {
		t.Errorf("expected TotalErrors=0, got %d", stats.TotalErrors)
	}
	if len(stats.CategoryCounts) != 0 {
		t.Errorf("expected empty category counts, got %d", len(stats.CategoryCounts))
	}
}

func TestGetErrorCategoryTrend(t *testing.T) {
	db := setupErrorTestDB(t)
	ctx := context.Background()

	today := time.Now().UTC().Format("2006-01-02")
	yesterday := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02")

	db.MustExecContext(ctx, `INSERT INTO error_history (user_id, original, corrected, category, severity, created_at) VALUES (1, 'I goes', 'I go', 'grammar', 'major', ?)`, yesterday)
	db.MustExecContext(ctx, `INSERT INTO error_history (user_id, original, corrected, category, severity, created_at) VALUES (1, 'she doesn''t', 'she doesn''t', 'grammar', 'critical', ?)`, yesterday)
	db.MustExecContext(ctx, `INSERT INTO error_history (user_id, original, corrected, category, severity, created_at) VALUES (1, 'yesterday I', 'I yesterday', 'word_order', 'minor', ?)`, today)

	trend, err := GetErrorCategoryTrend(ctx, db, 1, 30)
	if err != nil {
		t.Fatalf("GetErrorCategoryTrend failed: %v", err)
	}

	if len(trend) != 2 {
		t.Errorf("expected 2 days in trend, got %d", len(trend))
	}

	for _, entry := range trend {
		if entry.Date == yesterday {
			if entry.Grammar != 2 {
				t.Errorf("expected grammar=2 on %s, got %d", yesterday, entry.Grammar)
			}
			if entry.WordOrder != 0 {
				t.Errorf("expected word_order=0 on %s, got %d", yesterday, entry.WordOrder)
			}
		}
		if entry.Date == today {
			if entry.WordOrder != 1 {
				t.Errorf("expected word_order=1 on %s, got %d", today, entry.WordOrder)
			}
		}
	}
}

func TestGetErrorHistory(t *testing.T) {
	db := setupErrorTestDB(t)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		SaveError(ctx, db, 1, "test", "fixed", "grammar", "minor", "", "", "")
	}

	entries, total, err := GetErrorHistory(ctx, db, 1, 3, 0, "")
	if err != nil {
		t.Fatalf("GetErrorHistory failed: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total=5, got %d", total)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestGetErrorHistoryWithCategory(t *testing.T) {
	db := setupErrorTestDB(t)
	ctx := context.Background()

	SaveError(ctx, db, 1, "I goes", "I go", "grammar", "major", "", "", "")
	SaveError(ctx, db, 1, "yesterday I", "I yesterday", "word_order", "minor", "", "", "")

	entries, total, err := GetErrorHistory(ctx, db, 1, 10, 0, "grammar")
	if err != nil {
		t.Fatalf("GetErrorHistory failed: %v", err)
	}

	if total != 1 {
		t.Errorf("expected total=1, got %d", total)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Category != "grammar" {
		t.Errorf("expected category=grammar, got %s", entries[0].Category)
	}
}

func TestGetErrorHistoryEmpty(t *testing.T) {
	db := setupErrorTestDB(t)
	ctx := context.Background()

	entries, total, err := GetErrorHistory(ctx, db, 999, 10, 0, "")
	if err != nil {
		t.Fatalf("GetErrorHistory failed: %v", err)
	}

	if total != 0 {
		t.Errorf("expected total=0, got %d", total)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestGetRecentErrorContext(t *testing.T) {
	db := setupErrorTestDB(t)
	ctx := context.Background()

	SaveError(ctx, db, 1, "I goes", "I go", "grammar", "major", "Subject-Verb Agreement", "Use 'I go' not 'I goes'", "")
	SaveError(ctx, db, 1, "she don't", "she doesn't", "grammar", "critical", "Subject-Verb Agreement", "Use 'doesn't' with she/he/it", "")

	errors, err := GetRecentErrorContext(ctx, db, 1, 5)
	if err != nil {
		t.Fatalf("GetRecentErrorContext failed: %v", err)
	}

	if len(errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errors))
	}
}

func TestGetRecentErrorSummary(t *testing.T) {
	db := setupErrorTestDB(t)
	ctx := context.Background()

	SaveError(ctx, db, 1, "I goes", "I go", "grammar", "major", "", "", "")

	summary, err := GetRecentErrorSummary(ctx, db, 1, 5)
	if err != nil {
		t.Fatalf("GetRecentErrorSummary failed: %v", err)
	}

	if len(summary) != 1 {
		t.Errorf("expected 1 summary, got %d", len(summary))
	}
}
