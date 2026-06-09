package models

import "time"

type User struct {
	ID         int       `db:"id" json:"id"`
	TelegramID int64     `db:"telegram_id" json:"telegram_id"`
	Username   string    `db:"username" json:"username"`
	Lang       string    `db:"lang" json:"lang"`
	Level      string    `db:"level" json:"level"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type Correction struct {
	Original      string `json:"original"`
	Corrected     string `json:"corrected"`
	ExplanationUz string `json:"explanation_uz"`
	ExplanationRu string `json:"explanation_ru"`
	Type          string `json:"type"`
	Severity      string `json:"severity,omitempty"`
	Category      string `json:"category,omitempty"`
	LearningTip   string `json:"learning_tip,omitempty"`
	RuleViolated  string `json:"rule_violated,omitempty"`
}
	Corrections     []Correction     `json:"corrections"`
	Usage           Usage            `json:"usage"`
=======
	Severity      string `json:"severity,omitempty"`       // "critical" | "major" | "minor"
	Category      string `json:"category,omitempty"`       // "grammar" | "vocabulary" | "spelling" | "word_order" | "punctuation"
	LearningTip   string `json:"learning_tip,omitempty"`   // e.g. "Remember: 'yesterday' always triggers Past Simple"
	RuleViolated  string `json:"rule_violated,omitempty"`  // e.g. "Subject-Verb Agreement (Present Simple)"
}

type PremiumAnalysis struct {
	OverallGrade        string   `json:"overall_grade"`         // "A" | "B" | "C" | "D"
	Strengths           []string `json:"strengths"`
	AreasForImprovement []string `json:"areas_for_improvement"`
	SuggestedTopic      string   `json:"suggested_topic"`       // "Next: practice Past Continuous"
}

type ChatResponse struct {
	Reply           string          `json:"reply"`
	Corrections     []Correction    `json:"corrections"`
	Usage           Usage           `json:"usage"`
>>>>>>> feature/premium-features
	PremiumAnalysis *PremiumAnalysis `json:"premium_analysis,omitempty"`
}

type Usage struct {
	DailyUsed  int  `json:"daily_used"`
	DailyLimit int  `json:"daily_limit"`
	IsPremium  bool `json:"is_premium"`
}

type ChatRequest struct {
	Text string `json:"text" binding:"required"`
}

type Subscription struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	Plan        string    `db:"plan" json:"plan"`
	StarsAmount int       `db:"stars_amount" json:"stars_amount"`
	StartedAt   time.Time `db:"started_at" json:"started_at"`
	ExpiresAt   time.Time `db:"expires_at" json:"expires_at"`
}

type SubscriptionResponse struct {
	Active    bool   `json:"active"`
	Plan      string `json:"plan,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

type InvoiceRequest struct {
	Plan string `json:"plan" binding:"required"`
}

type InvoiceResponse struct {
	InvoiceLink string `json:"invoice_link"`
	Stars       int    `json:"stars"`
}

type AIResponse struct {
	Reply           string           `json:"reply"`
	Corrections     []Correction     `json:"corrections"`
	PremiumAnalysis *PremiumAnalysis `json:"premium_analysis,omitempty"`
}

type ErrorResponse struct {
	Error      string `json:"error"`
	MessageUz  string `json:"message_uz,omitempty"`
	MessageRu  string `json:"message_ru,omitempty"`
}

type GrammarRequest struct {
	Text string `json:"text" binding:"required"`
}

type GrammarResponse struct {
	Corrections []Correction `json:"corrections"`
}

type VocabLookupRequest struct {
	Word string `json:"word" binding:"required"`
}

type VocabLookupResponse struct {
	TranslationUz string   `json:"translation_uz"`
	TranslationRu string   `json:"translation_ru"`
	Examples      []string `json:"examples"`
	Level         string   `json:"level"`
}

type VocabWord struct {
	ID            int        `db:"id" json:"id"`
	UserID        int        `db:"user_id" json:"user_id"`
	Word          string     `db:"word" json:"word"`
	Translation   string     `db:"translation" json:"translation"`
	Example       string     `db:"example" json:"example"`
	Level         string     `db:"level" json:"level"`
	ReviewCount   int        `db:"review_count" json:"review_count"`
	EaseFactor    float64    `db:"ease_factor" json:"ease_factor"`
	Interval      int        `db:"interval" json:"interval"`
	LastReviewedAt *time.Time `db:"last_reviewed_at" json:"last_reviewed_at,omitempty"`
	NextReview    *time.Time `db:"next_review" json:"next_review,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}

type AddVocabRequest struct {
	Word        string `json:"word" binding:"required"`
	Translation string `json:"translation" binding:"required"`
	Example     string `json:"example" binding:"required"`
	Level       string `json:"level"`
}

type QuizRequest struct {
	Topic string `json:"topic"`
	Count int    `json:"count"`
}

type QuizQuestion struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	Correct       int      `json:"correct"`
	ExplanationUz string   `json:"explanation_uz"`
	ExplanationRu string   `json:"explanation_ru"`
}

type QuizResponse struct {
	Questions []QuizQuestion `json:"questions"`
}

type LevelRequest struct {
	Text string `json:"text"`
}

type LevelQuestion struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	Correct       int      `json:"correct"`
	Level         string   `json:"level"`
	ExplanationUz string   `json:"explanation_uz"`
	ExplanationRu string   `json:"explanation_ru"`
}

type LevelResponse struct {
	Questions []LevelQuestion `json:"questions"`
	Level     string          `json:"level,omitempty"`
}

type UserStats struct {
	Level              string `json:"level"`
	TotalMessages      int    `json:"total_messages"`
	TotalWords         int    `json:"total_words"`
	WordsDueToday      int    `json:"words_due_today"`
	StreakDays         int    `json:"streak_days"`
	IsPremium          bool   `json:"is_premium"`
	AccountCreatedAt   string `json:"account_created_at"`
	SubscriptionExpiry string `json:"subscription_expiry,omitempty"`
}

type ProgressResponse struct {
	MessagesSent int    `json:"messages_sent"`
	WordsLearned int    `json:"words_learned"`
	QuizzesTaken int    `json:"quizzes_taken"`
	StreakDays   int    `json:"streak_days"`
	Level        string `json:"level"`
}

type DailyProgressEntry struct {
	Date         string `json:"date"`
	MessagesSent int    `json:"messages_sent"`
	WordsLearned int    `json:"words_learned"`
	QuizzesTaken int    `json:"quizzes_taken"`
}

type ProgressHistoryResponse struct {
	Entries []DailyProgressEntry `json:"entries"`
}

type VocabListResponse struct {
	Words    []VocabWord `json:"words"`
	Total    int         `json:"total"`
	DueCount int         `json:"due_count"`
}

type ErrorHistoryEntry struct {
	ID            int       `db:"id" json:"id"`
	UserID        int       `db:"user_id" json:"user_id"`
	Original      string    `db:"original" json:"original"`
	Corrected     string    `db:"corrected" json:"corrected"`
	Category      string    `db:"category" json:"category"`
	Severity      string    `db:"severity" json:"severity"`
	RuleViolated  string    `db:"rule_violated" json:"rule_violated"`
	LearningTip   string    `db:"learning_tip" json:"learning_tip"`
	Context       string    `db:"context" json:"context"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type ErrorStats struct {
	TotalErrors       int                       `json:"total_errors"`
	CategoryCounts    map[string]int            `json:"category_counts"`
	SeverityCounts    map[string]int            `json:"severity_counts"`
	MostFrequentRules []string                  `json:"most_frequent_rules"`
	CategoryTrend     []ErrorCategoryDayEntry   `json:"category_trend,omitempty"`
}

type ErrorCategoryDayEntry struct {
	Date        string `json:"date"`
	Grammar     int    `json:"grammar"`
	Vocabulary  int    `json:"vocabulary"`
	Spelling    int    `json:"spelling"`
	WordOrder   int    `json:"word_order"`
	Punctuation int    `json:"punctuation"`
}

type ErrorHistoryResponse struct {
	Entries []ErrorHistoryEntry `json:"entries"`
	Total   int                 `json:"total"`
}

type ErrorStatsRequest struct {
	Days int `form:"days"`
}
