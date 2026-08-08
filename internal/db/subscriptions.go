package db

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lingvo-ai/lingvo/internal/models"
)

func GetActiveSubscription(ctx context.Context, db *sqlx.DB, userID int) (*models.Subscription, error) {
	var s models.Subscription
	err := db.GetContext(ctx, &s,
		"SELECT * FROM subscriptions WHERE user_id = ? AND expires_at > datetime('now')",
		userID)
	if err != nil {
		return nil, nil
	}
	return &s, nil
}

func SaveSubscription(ctx context.Context, db *sqlx.DB, userID int, plan string, starsAmount int, expiresAt string) error {
	newExpiry, _ := time.Parse("2006-01-02 15:04:05", expiresAt)
	duration := newExpiry.Sub(time.Now().UTC())

	var existing models.Subscription
	err := db.GetContext(ctx, &existing,
		"SELECT * FROM subscriptions WHERE user_id = ? AND expires_at > datetime('now')",
		userID)

	if err == nil {
		extended := existing.ExpiresAt.Add(duration)
		_, err = db.ExecContext(ctx, `
			UPDATE subscriptions SET plan = ?, stars_amount = stars_amount + ?, expires_at = ? WHERE user_id = ?
		`, plan, starsAmount, extended.Format("2006-01-02 15:04:05"), userID)
	} else {
		_, err = db.ExecContext(ctx, `
			INSERT INTO subscriptions (user_id, plan, stars_amount, expires_at)
			VALUES (?, ?, ?, ?)
		`, userID, plan, starsAmount, expiresAt)
	}
	return err
}

func TogglePremium(ctx context.Context, db *sqlx.DB, userID int) (bool, error) {
	sub, err := GetActiveSubscription(ctx, db, userID)
	if err != nil {
		return false, err
	}

	if sub != nil {
		_, err = db.ExecContext(ctx, "DELETE FROM subscriptions WHERE user_id = ?", userID)
		return false, err
	}

	farFuture := time.Now().AddDate(10, 0, 0).Format(time.RFC3339)
	err = SaveSubscription(ctx, db, userID, "monthly", 0, farFuture)
	return true, err
}
