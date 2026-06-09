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
}

type ChatResponse struct {
	Reply       string       `json:"reply"`
	Corrections []Correction `json:"corrections"`
	Usage       Usage        `json:"usage"`
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
	Reply       string       `json:"reply"`
	Corrections []Correction `json:"corrections"`
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

type IeltsScoreEntry struct {
	ID           int       `db:"id" json:"id"`
	UserID       int       `db:"user_id" json:"user_id"`
	Module       string    `db:"module" json:"module"`
	BandScore    float64   `db:"band_score" json:"band_score"`
	Details      string    `db:"details" json:"details"`
	Prompt       string    `db:"prompt" json:"prompt"`
	UserResponse string    `db:"user_response" json:"user_response"`
	Feedback     string    `db:"feedback" json:"feedback"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type IeltsScoreStats struct {
	WritingTask1Avg float64 `json:"writing_task1_avg"`
	WritingTask2Avg float64 `json:"writing_task2_avg"`
	SpeakingAvg     float64 `json:"speaking_avg"`
	ReadingAvg      float64 `json:"reading_avg"`
	TotalPractices  int     `json:"total_practices"`
}

type IeltsScoresResponse struct {
	Entries []IeltsScoreEntry `json:"entries"`
	Total   int               `json:"total"`
	Stats   *IeltsScoreStats  `json:"stats,omitempty"`
}

type IeltsWritingRequest struct {
	Type            string `json:"type" binding:"required"`
	UserText        string `json:"user_text" binding:"required"`
	TaskDescription string `json:"task_description,omitempty"`
}

type IeltsWritingCriteria struct {
	TaskAchievement  float64 `json:"task_achievement"`
	CoherenceCohesion float64 `json:"coherence_cohesion"`
	LexicalResource  float64 `json:"lexical_resource"`
	GrammaticalRange float64 `json:"grammatical_range"`
}

type IeltsWritingResponse struct {
	BandScore      float64              `json:"band_score"`
	Criteria       IeltsWritingCriteria `json:"criteria"`
	Feedback       string               `json:"feedback"`
	Corrections    []Correction         `json:"corrections"`
	ImprovementTips []string            `json:"improvement_tips"`
}

type IeltsSpeakingRequest struct {
	Part         int    `json:"part" binding:"required"`
	Question     string `json:"question,omitempty"`
	UserResponse string `json:"user_response,omitempty"`
}

type IeltsSpeakingQuestionsResponse struct {
	Part      int      `json:"part"`
	Questions []string `json:"questions"`
	CueCard   string   `json:"cue_card,omitempty"`
}

type IeltsSpeakingCriteria struct {
	FluencyCoherence  float64 `json:"fluency_coherence"`
	LexicalResource   float64 `json:"lexical_resource"`
	GrammaticalRange  float64 `json:"grammatical_range"`
	Pronunciation     float64 `json:"pronunciation"`
}

type IeltsSpeakingResponse struct {
	BandScore      float64               `json:"band_score"`
	Criteria       IeltsSpeakingCriteria `json:"criteria"`
	Feedback       string                `json:"feedback"`
	ImprovementTips []string             `json:"improvement_tips"`
}

type IeltsReadingPassage struct {
	Title     string          `json:"title"`
	Passage   string          `json:"passage"`
	WordCount int             `json:"word_count"`
	Questions []interface{}   `json:"questions"`
}

type IeltsReadingQuestion struct {
	Type     string   `json:"type"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Correct  int      `json:"correct"`
}

type IeltsReadingSubmitRequest struct {
	Passage     string                `json:"passage" binding:"required"`
	Questions   []IeltsReadingQuestion `json:"questions" binding:"required"`
	UserAnswers []int                 `json:"user_answers" binding:"required"`
}

type IeltsReadingResult struct {
	CorrectAnswers  int                     `json:"correct_answers"`
	TotalQuestions  int                     `json:"total_questions"`
	BandScore       float64                 `json:"band_score"`
	Results         []IeltsQuestionResult   `json:"results"`
	Feedback        string                  `json:"feedback"`
}

type IeltsQuestionResult struct {
	QuestionIndex   int    `json:"question_index"`
	UserAnswer      int    `json:"user_answer"`
	CorrectAnswer   int    `json:"correct_answer"`
	IsCorrect       bool   `json:"is_correct"`
	ExplanationUz   string `json:"explanation_uz"`
	ExplanationRu   string `json:"explanation_ru"`
}
