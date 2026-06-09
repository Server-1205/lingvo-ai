package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"

	"github.com/lingvo-ai/lingvo/internal/models"
)

func TestIeltsScoresHandler_PremiumGate(t *testing.T) {
	_, router := setupIELTSTest(t, false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/ielts/scores", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for free user, got %d", w.Code)
	}
}

func TestIeltsScoresHandler_Success(t *testing.T) {
	db, router := setupIELTSTest(t, true)

	insertIELTSScore(t, db, 1, "writing_task1", 6.5)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/ielts/scores", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp models.IeltsScoresResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if resp.Total != 1 {
		t.Errorf("expected total 1, got %d", resp.Total)
	}
	if len(resp.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(resp.Entries))
	}
	if resp.Stats == nil {
		t.Fatal("expected non-nil stats")
	}
}

func TestIeltsScoresHandler_WithModuleFilter(t *testing.T) {
	db, router := setupIELTSTest(t, true)

	insertIELTSScore(t, db, 1, "writing_task1", 6.5)
	insertIELTSScore(t, db, 1, "speaking", 7.0)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/ielts/scores?module=speaking", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp models.IeltsScoresResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Total != 1 {
		t.Errorf("expected 1 speaking score, got %d", resp.Total)
	}
}

func TestIeltsWritingHandler_PremiumGate(t *testing.T) {
	_, router := setupIELTSTest(t, false)

	body := `{"type":"task1","user_text":"The chart shows..."}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/ielts/writing", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for free user, got %d", w.Code)
	}
}

func TestIeltsWritingHandler_InvalidRequest(t *testing.T) {
	_, router := setupIELTSTest(t, true)

	body := `{"type":"invalid","user_text":"test"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/ielts/writing", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 (ai unavailable), got %d", w.Code)
	}
}

func TestIeltsSpeakingGenerateHandler_PremiumGate(t *testing.T) {
	_, router := setupIELTSTest(t, false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/ielts/speaking?part=1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for free user, got %d", w.Code)
	}
}

func TestIeltsSpeakingEvaluateHandler_PremiumGate(t *testing.T) {
	_, router := setupIELTSTest(t, false)

	body := `{"part":1,"question":"Tell me about yourself","user_response":"I am a student"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/ielts/speaking", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for free user, got %d", w.Code)
	}
}

func TestIeltsReadingGenerateHandler_PremiumGate(t *testing.T) {
	_, router := setupIELTSTest(t, false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/ielts/reading", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for free user, got %d", w.Code)
	}
}

func TestIeltsReadingEvaluateHandler_PremiumGate(t *testing.T) {
	_, router := setupIELTSTest(t, false)

	body := `{"passage":"test","questions":[{"type":"multiple_choice","question":"Q?","options":["A","B"],"correct":0}],"user_answers":[0]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/ielts/reading/submit", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for free user, got %d", w.Code)
	}
}

func TestIeltsReadingEvaluateHandler_BadRequest(t *testing.T) {
	_, router := setupIELTSTest(t, true)

	body := `{"passage":"test","questions":[{"type":"multiple_choice","question":"Q?","options":["A","B"],"correct":0}],"user_answers":[0,1]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/ielts/reading/submit", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 (ai unavailable), got %d", w.Code)
	}
}

func TestReadingBandScore(t *testing.T) {
	tests := []struct {
		correct, total int
		expected       float64
	}{
		{10, 10, 9.0},
		{9, 10, 8.0},
		{8, 10, 7.0},
		{7, 10, 6.5},
		{6, 10, 6.0},
		{5, 10, 5.5},
		{4, 10, 5.0},
		{3, 10, 4.0},
		{0, 10, 4.0},
		{0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("correct=%d/total=%d", tt.correct, tt.total), func(t *testing.T) {
			got := readingBandScore(tt.correct, tt.total)
			if got != tt.expected {
				t.Errorf("readingBandScore(%d,%d) = %f, want %f", tt.correct, tt.total, got, tt.expected)
			}
		})
	}
}

func TestRoundBandScore(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{6.25, 6.5},
		{6.0, 6.0},
		{6.1, 6.0},
		{6.75, 7.0},
		{7.3, 7.5},
	}

	for _, tt := range tests {
		got := roundBandScore(tt.input)
		if got != tt.expected {
			t.Errorf("roundBandScore(%f) = %f, want %f", tt.input, got, tt.expected)
		}
	}
}

func setupIELTSTest(t *testing.T, isPremium bool) (*sqlx.DB, *gin.Engine) {
	t.Helper()

	gin.SetMode(gin.TestMode)

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
	);`

	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("failed to exec schema: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO users (id, telegram_id, username, lang, level) VALUES (1, 12345, 'testuser', 'uz', 'b1')`)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}

	router := gin.New()
	api := router.Group("/api")

	api.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Set("lang", "uz")
		c.Set("level", "b1")
		c.Set("is_premium", isPremium)
		c.Set("telegram_id", int64(12345))
		c.Next()
	})

	sugar := zap.NewNop().Sugar()

	api.POST("/ielts/writing", ieltsWritingHandler(db, nil, sugar))
	api.GET("/ielts/speaking", ieltsSpeakingGenerateHandler(db, nil, sugar))
	api.POST("/ielts/speaking", ieltsSpeakingEvaluateHandler(db, nil, sugar))
	api.GET("/ielts/reading", ieltsReadingGenerateHandler(db, nil, sugar))
	api.POST("/ielts/reading/submit", ieltsReadingEvaluateHandler(db, nil, sugar))
	api.GET("/ielts/scores", ieltsScoresHandler(db, sugar))

	return db, router
}

func insertIELTSScore(t *testing.T, db *sqlx.DB, userID int, module string, bandScore float64) {
	t.Helper()
	_, err := db.ExecContext(context.Background(),
		`INSERT INTO ielts_scores (user_id, module, band_score, details, prompt, user_response, feedback)
		 VALUES (?, ?, ?, '{}', '', '', '')`,
		userID, module, bandScore)
	if err != nil {
		t.Fatalf("insert ielts score: %v", err)
	}
}
