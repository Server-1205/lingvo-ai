package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lingvo-ai/lingvo/internal/models"
)

func SaveError(ctx context.Context, db *sqlx.DB, userID int, original, corrected, category, severity, ruleViolated, learningTip, context string) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO error_history (user_id, original, corrected, category, severity, rule_violated, learning_tip, context)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, original, corrected, category, severity, ruleViolated, learningTip, context)
	return err
}

func GetErrorHistory(ctx context.Context, db *sqlx.DB, userID, limit, offset int, category string) ([]models.ErrorHistoryEntry, int, error) {
	countArgs := []interface{}{userID}
	queryArgs := []interface{}{userID}

	countWhere := "WHERE user_id = ?"
	queryWhere := "WHERE user_id = ?"

	if category != "" {
		countWhere += " AND category = ?"
		queryWhere += " AND category = ?"
		countArgs = append(countArgs, category)
		queryArgs = append(queryArgs, category)
	}

	var total int
	if err := db.GetContext(ctx, &total,
		"SELECT COUNT(*) FROM error_history "+countWhere, countArgs...); err != nil {
		return nil, 0, err
	}

	var entries []models.ErrorHistoryEntry
	if err := db.SelectContext(ctx, &entries,
		"SELECT * FROM error_history "+queryWhere+" ORDER BY created_at DESC LIMIT ? OFFSET ?",
		append(queryArgs, limit, offset)...); err != nil {
		return nil, 0, err
	}

	if entries == nil {
		entries = []models.ErrorHistoryEntry{}
	}

	return entries, total, nil
}

func GetErrorStats(ctx context.Context, db *sqlx.DB, userID int, days int) (*models.ErrorStats, error) {
	stats := &models.ErrorStats{
		CategoryCounts:   make(map[string]int),
		SeverityCounts:   make(map[string]int),
		TotalErrors:      0,
	}

	if days <= 0 {
		days = 30
	}

	dateFilter := time.Now().UTC().AddDate(0, 0, -days).Format("2006-01-02 15:04:05")

	type categoryRow struct {
		Category string `db:"category"`
		Count    int    `db:"cnt"`
	}

	var catRows []categoryRow
	if err := db.SelectContext(ctx, &catRows, `
		SELECT category, COUNT(*) as cnt
		FROM error_history
		WHERE user_id = ? AND created_at >= ?
		GROUP BY category
		ORDER BY cnt DESC`, userID, dateFilter); err != nil {
		return nil, err
	}

	for _, r := range catRows {
		stats.CategoryCounts[r.Category] = r.Count
		stats.TotalErrors += r.Count
	}

	type severityRow struct {
		Severity string `db:"severity"`
		Count    int    `db:"cnt"`
	}

	var sevRows []severityRow
	if err := db.SelectContext(ctx, &sevRows, `
		SELECT severity, COUNT(*) as cnt
		FROM error_history
		WHERE user_id = ? AND created_at >= ?
		GROUP BY severity
		ORDER BY cnt DESC`, userID, dateFilter); err != nil {
		return nil, err
	}

	for _, r := range sevRows {
		stats.SeverityCounts[r.Severity] = r.Count
	}

	type ruleRow struct {
		Rule  string `db:"rule"`
		Count int    `db:"cnt"`
	}

	var ruleRows []ruleRow
	if err := db.SelectContext(ctx, &ruleRows, `
		SELECT rule_violated as rule, COUNT(*) as cnt
		FROM error_history
		WHERE user_id = ? AND created_at >= ? AND rule_violated != ''
		GROUP BY rule_violated
		ORDER BY cnt DESC
		LIMIT 10`, userID, dateFilter); err != nil {
		return nil, err
	}

	for _, r := range ruleRows {
		stats.MostFrequentRules = append(stats.MostFrequentRules, fmt.Sprintf("%s (%dx)", r.Rule, r.Count))
	}

	if stats.MostFrequentRules == nil {
		stats.MostFrequentRules = []string{}
	}

	trend, err := GetErrorCategoryTrend(ctx, db, userID, days)
	if err != nil {
		return stats, nil
	}
	stats.CategoryTrend = trend

	return stats, nil
}

func GetErrorCategoryTrend(ctx context.Context, db *sqlx.DB, userID int, days int) ([]models.ErrorCategoryDayEntry, error) {
	if days <= 0 {
		days = 30
	}

	dateFilter := time.Now().UTC().AddDate(0, 0, -days).Format("2006-01-02")

	type trendRow struct {
		Date        string `db:"date_str"`
		Grammar     int    `db:"grammar"`
		Vocabulary  int    `db:"vocabulary"`
		Spelling    int    `db:"spelling"`
		WordOrder   int    `db:"word_order"`
		Punctuation int    `db:"punctuation"`
	}

	var rows []trendRow
	if err := db.SelectContext(ctx, &rows, `
		SELECT
			DATE(created_at) as date_str,
			SUM(CASE WHEN category = 'grammar' THEN 1 ELSE 0 END) as grammar,
			SUM(CASE WHEN category = 'vocabulary' THEN 1 ELSE 0 END) as vocabulary,
			SUM(CASE WHEN category = 'spelling' THEN 1 ELSE 0 END) as spelling,
			SUM(CASE WHEN category = 'word_order' THEN 1 ELSE 0 END) as word_order,
			SUM(CASE WHEN category = 'punctuation' THEN 1 ELSE 0 END) as punctuation
		FROM error_history
		WHERE user_id = ? AND DATE(created_at) >= ?
		GROUP BY DATE(created_at)
		ORDER BY date_str ASC`, userID, dateFilter); err != nil {
		return nil, err
	}

	entries := make([]models.ErrorCategoryDayEntry, 0, len(rows))
	for _, r := range rows {
		entries = append(entries, models.ErrorCategoryDayEntry{
			Date:        r.Date,
			Grammar:     r.Grammar,
			Vocabulary:  r.Vocabulary,
			Spelling:    r.Spelling,
			WordOrder:   r.WordOrder,
			Punctuation: r.Punctuation,
		})
	}

	return entries, nil
}

func GetMostFrequentErrors(ctx context.Context, db *sqlx.DB, userID int, limit int) ([]string, error) {
	type freqRow struct {
		Rule     string `db:"rule"`
		Category string `db:"category"`
		Count    int    `db:"cnt"`
	}

	var rows []freqRow
	if err := db.SelectContext(ctx, &rows, `
		SELECT rule_violated as rule, category, COUNT(*) as cnt
		FROM error_history
		WHERE user_id = ? AND rule_violated != ''
		GROUP BY rule_violated, category
		ORDER BY cnt DESC
		LIMIT ?`, userID, limit); err != nil {
		return nil, err
	}

	var results []string
	for _, r := range rows {
		results = append(results, fmt.Sprintf("%s (%s, %dx)", r.Rule, r.Category, r.Count))
	}

	return results, nil
}

func GetRecentErrorSummary(ctx context.Context, db *sqlx.DB, userID int, limit int) ([]string, error) {
	type recentRow struct {
		Category string `db:"category"`
		Original string `db:"original"`
		Corrected string `db:"corrected"`
	}

	var rows []recentRow
	if err := db.SelectContext(ctx, &rows, `
		SELECT category, original, corrected FROM error_history
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ?`, userID, limit); err != nil {
		return nil, err
	}

	var result []string
	for _, r := range rows {
		summary := fmt.Sprintf("[%s] \"%s\" → \"%s\"", r.Category, r.Original, r.Corrected)
		if len(summary) > 120 {
			summary = summary[:117] + "..."
		}
		result = append(result, summary)
	}
	return result, nil
}

func GetRecentErrorContext(ctx context.Context, db *sqlx.DB, userID int, limit int) ([]string, error) {
	entries, _, err := GetErrorHistory(ctx, db, userID, limit, 0, "")
	if err != nil {
		return nil, err
	}

	var result []string
	for _, e := range entries {
		result = append(result, fmt.Sprintf("%s (%s, %s) — \"%s\" → \"%s\"",
			e.RuleViolated, e.Category, e.Severity, e.Original, e.Corrected))
	}
	return result, nil
}
