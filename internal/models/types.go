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
