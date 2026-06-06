package db

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sqlx.DB {
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
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL REFERENCES users(id),
		date TEXT NOT NULL,
		count INTEGER NOT NULL DEFAULT 0,
		UNIQUE(user_id, date)
	);
	CREATE TABLE IF NOT EXISTS vocabulary (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL REFERENCES users(id),
		word TEXT NOT NULL,
		translation TEXT NOT NULL,
		example TEXT NOT NULL,
		level TEXT DEFAULT 'a1',
		review_count INTEGER DEFAULT 0,
		ease_factor REAL DEFAULT 2.5,
		interval INTEGER DEFAULT 0,
		last_reviewed_at DATETIME,
		next_review DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, word)
	);
	CREATE TABLE IF NOT EXISTS subscriptions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL UNIQUE REFERENCES users(id),
		plan TEXT NOT NULL CHECK(plan IN ('weekly','monthly')),
		stars_amount INTEGER NOT NULL,
		started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME NOT NULL
	);
	CREATE TABLE IF NOT EXISTS daily_progress (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL REFERENCES users(id),
		date TEXT NOT NULL,
		messages_sent INTEGER DEFAULT 0,
		words_learned INTEGER DEFAULT 0,
		quizzes_taken INTEGER DEFAULT 0,
		UNIQUE(user_id, date)
	);`

	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("failed to exec schema: %v", err)
	}

	return db
}

func TestGetUserStats_ExistingUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	ctx := context.Background()

	_, err := db.ExecContext(ctx,
		`INSERT INTO users (id, telegram_id, username, lang, level) VALUES (1, 12345, 'testuser', 'uz', 'b1')`)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}

	db.MustExecContext(ctx,
		`INSERT INTO messages (user_id, date, count) VALUES (1, '2026-06-01', 5)`)
	db.MustExecContext(ctx,
		`INSERT INTO messages (user_id, date, count) VALUES (1, '2026-06-02', 3)`)

	db.MustExecContext(ctx,
		`INSERT INTO vocabulary (user_id, word, translation, example) VALUES (1, 'hello', 'привет', 'Hello world')`)
	db.MustExecContext(ctx,
		`INSERT INTO vocabulary (user_id, word, translation, example) VALUES (1, 'world', 'мир', 'The world is big')`)

	db.MustExecContext(ctx,
		`INSERT INTO subscriptions (user_id, plan, stars_amount, expires_at)
		 VALUES (1, 'monthly', 800, datetime('now', '+30 days'))`)

	stats, err := GetUserStats(ctx, db, 1)
	if err != nil {
		t.Fatalf("GetUserStats returned error: %v", err)
	}
	if stats == nil {
		t.Fatal("GetUserStats returned nil for existing user")
	}

	if stats.Level != "b1" {
		t.Errorf("expected level b1, got %s", stats.Level)
	}
	if stats.TotalMessages != 8 {
		t.Errorf("expected 8 total messages, got %d", stats.TotalMessages)
	}
	if stats.TotalWords != 2 {
		t.Errorf("expected 2 words, got %d", stats.TotalWords)
	}
	if !stats.IsPremium {
		t.Error("expected isPremium = true")
	}
	if stats.SubscriptionExpiry == "" {
		t.Error("expected non-empty subscription expiry")
	}
	if stats.AccountCreatedAt == "" {
		t.Error("expected non-empty account created_at")
	}
}

func TestGetUserStats_NoData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	ctx := context.Background()

	_, err := db.ExecContext(ctx,
		`INSERT INTO users (id, telegram_id, username, lang) VALUES (2, 99999, 'emptyuser', 'ru')`)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}

	stats, err := GetUserStats(ctx, db, 2)
	if err != nil {
		t.Fatalf("GetUserStats returned error: %v", err)
	}
	if stats == nil {
		t.Fatal("GetUserStats returned nil for existing user")
	}

	if stats.TotalMessages != 0 {
		t.Errorf("expected 0 messages, got %d", stats.TotalMessages)
	}
	if stats.TotalWords != 0 {
		t.Errorf("expected 0 words, got %d", stats.TotalWords)
	}
	if stats.IsPremium {
		t.Error("expected isPremium = false for user without subscription")
	}
}

func TestGetUserStats_NonExistentUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	ctx := context.Background()

	stats, err := GetUserStats(ctx, db, 999)
	if err != nil {
		t.Fatalf("GetUserStats returned error for non-existent user: %v", err)
	}
	if stats != nil {
		t.Fatal("expected nil for non-existent user")
	}
}
