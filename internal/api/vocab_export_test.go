package api

import (
	"context"
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
	"go.uber.org/zap"
)

func setupExportTest(t *testing.T) (*sqlx.DB, *gin.Engine) {
	t.Helper()

	db, err := sqlx.Connect("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
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
	CREATE TABLE IF NOT EXISTS vocabulary (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL REFERENCES users(id),
		word TEXT NOT NULL,
		translation TEXT NOT NULL,
		example TEXT NOT NULL,
		level TEXT NOT NULL DEFAULT 'a1',
		review_count INTEGER NOT NULL DEFAULT 0,
		ease_factor REAL NOT NULL DEFAULT 2.5,
		interval INTEGER NOT NULL DEFAULT 0,
		last_reviewed_at DATETIME,
		next_review DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.ExecContext(context.Background(), schema)
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		"INSERT INTO users (id, telegram_id, username, lang, level) VALUES (1, 12345, 'testuser', 'uz', 'a1')")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO vocabulary (user_id, word, translation, example, level, review_count, created_at)
		 VALUES (1, 'hello', 'salom', 'Hello world!', 'a1', 3, '2026-06-01')`)
	if err != nil {
		t.Fatalf("failed to insert vocab: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO vocabulary (user_id, word, translation, example, level, review_count, created_at)
		 VALUES (1, 'world', 'dunyo', 'The world is big.', 'a1', 1, '2026-06-02')`)
	if err != nil {
		t.Fatalf("failed to insert vocab: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Set("telegram_id", int64(12345))
		c.Set("lang", "uz")
		c.Set("level", "a1")
		c.Set("is_premium", false)
		c.Set("daily_used", 0)
		c.Set("daily_limit", 10)
	})

	r.GET("/api/vocab/export", vocabExportHandler(db, zap.NewNop().Sugar()))

	return db, r
}

func TestVocabExport_ReturnsCSVWithCorrectHeaders(t *testing.T) {
	_, r := setupExportTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/vocab/export", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/csv" {
		t.Errorf("expected text/csv, got %s", contentType)
	}

	reader := csv.NewReader(w.Body)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse csv: %v", err)
	}

	if len(records) < 2 {
		t.Fatalf("expected at least 2 rows (header + data), got %d", len(records))
	}

	headers := records[0]
	expectedHeaders := []string{"word", "translation", "example", "level", "review_count", "next_review", "created_at"}
	for i, h := range expectedHeaders {
		if headers[i] != h {
			t.Errorf("header[%d] = %s, want %s", i, headers[i], h)
		}
	}
}

func TestVocabExport_ContainsUserWords(t *testing.T) {
	_, r := setupExportTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/vocab/export", nil)
	r.ServeHTTP(w, req)

	reader := csv.NewReader(w.Body)
	records, _ := reader.ReadAll()

	foundHello := false
	foundWorld := false
	for _, row := range records {
		if row[0] == "hello" {
			foundHello = true
			if row[4] != "3" {
				t.Errorf("hello review_count = %s, want 3", row[4])
			}
		}
		if row[0] == "world" {
			foundWorld = true
		}
	}

	if !foundHello {
		t.Error("expected 'hello' in export")
	}
	if !foundWorld {
		t.Error("expected 'world' in export")
	}
}

func TestVocabExport_ContentDisposition(t *testing.T) {
	_, r := setupExportTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/vocab/export", nil)
	r.ServeHTTP(w, req)

	disp := w.Header().Get("Content-Disposition")
	if !strings.Contains(disp, "vocabulary.csv") {
		t.Errorf("expected Content-Disposition with vocabulary.csv, got %s", disp)
	}
}

func TestVocabExport_CacheControl(t *testing.T) {
	_, r := setupExportTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/vocab/export", nil)
	r.ServeHTTP(w, req)

	cc := w.Header().Get("Cache-Control")
	if cc != "no-store" {
		t.Errorf("expected Cache-Control: no-store, got %s", cc)
	}
}

func TestVocabExport_EmptyForOtherUser(t *testing.T) {
	db, _ := setupExportTest(t)
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/vocab/export", nil)

	r2 := gin.New()
	r2.Use(func(c *gin.Context) {
		c.Set("user_id", 999)
	})
	r2.GET("/api/vocab/export", vocabExportHandler(db, zap.NewNop().Sugar()))
	r2.ServeHTTP(w, req)

	reader := csv.NewReader(w.Body)
	records, _ := reader.ReadAll()

	if len(records) != 1 {
		t.Errorf("expected only header row for other user, got %d rows", len(records))
	}
}

func TestVocabExport_NextReviewColumn(t *testing.T) {
	db, err := sqlx.Connect("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			telegram_id INTEGER NOT NULL UNIQUE,
			username TEXT DEFAULT '',
			lang TEXT NOT NULL DEFAULT 'uz',
			level TEXT NOT NULL DEFAULT 'a1',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("failed to create users: %v", err)
	}

	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS vocabulary (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL REFERENCES users(id),
			word TEXT NOT NULL,
			translation TEXT NOT NULL,
			example TEXT NOT NULL,
			level TEXT NOT NULL DEFAULT 'a1',
			review_count INTEGER NOT NULL DEFAULT 0,
			ease_factor REAL NOT NULL DEFAULT 2.5,
			interval INTEGER NOT NULL DEFAULT 0,
			last_reviewed_at DATETIME,
			next_review DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("failed to create vocab: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		"INSERT INTO users (id, telegram_id) VALUES (1, 111)")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO vocabulary (user_id, word, translation, example, level, next_review, created_at)
		 VALUES (1, 'reviewword', 'tarjima', 'example', 'a1', '2026-07-01', '2026-06-01')`)
	if err != nil {
		t.Fatalf("failed to insert vocab with next_review: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
	})
	r.GET("/api/vocab/export", vocabExportHandler(db, zap.NewNop().Sugar()))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/vocab/export", nil)
	r.ServeHTTP(w, req)

	reader := csv.NewReader(w.Body)
	records, _ := reader.ReadAll()

	if len(records) < 2 {
		t.Fatalf("expected at least 2 rows, got %d", len(records))
	}

	if records[1][5] != "2026-07-01" {
		t.Errorf("next_review = %s, want 2026-07-01", records[1][5])
	}
}
