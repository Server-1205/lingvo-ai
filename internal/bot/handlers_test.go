package bot

import (
	"testing"

	"github.com/lingvo-ai/lingvo/internal/models"
)

func TestFormatStatsMessage_RU(t *testing.T) {
	stats := &models.UserStats{
		Level:              "b1",
		TotalMessages:      42,
		TotalWords:         15,
		WordsDueToday:      3,
		StreakDays:         7,
		IsPremium:          true,
		AccountCreatedAt:   "2026-01-15",
		SubscriptionExpiry: "2026-07-15",
	}

	msg := formatStatsMessage(stats, "ru", 30)

	expectedSubstrings := []string{
		"Статистика",
		"b1",
		"42",
		"15",
		"3",
		"7 дней",
		"30 дней",
		"2026-07-15",
	}
	for _, sub := range expectedSubstrings {
		if !contains(msg, sub) {
			t.Errorf("ru stats message missing: %s\nGot: %s", sub, msg)
		}
	}
}

func TestFormatStatsMessage_UZ(t *testing.T) {
	stats := &models.UserStats{
		Level:              "a2",
		TotalMessages:      10,
		TotalWords:         5,
		WordsDueToday:      2,
		StreakDays:         3,
		IsPremium:          false,
		AccountCreatedAt:   "2026-03-01",
		SubscriptionExpiry: "",
	}

	msg := formatStatsMessage(stats, "uz", 5)

	expectedSubstrings := []string{
		"Statistika",
		"a2",
		"10",
		"5",
		"2",
		"3 kun",
		"5 kun",
		"Yo'q",
	}
	for _, sub := range expectedSubstrings {
		if !contains(msg, sub) {
			t.Errorf("uz stats message missing: %s\nGot: %s", sub, msg)
		}
	}
}

func TestFormatStatsMessage_PremiumFree(t *testing.T) {
	premium := &models.UserStats{
		IsPremium:          true,
		SubscriptionExpiry: "2026-12-31",
	}

	free := &models.UserStats{
		IsPremium: false,
	}

	msgPremium := formatStatsMessage(premium, "uz", 1)
	if !contains(msgPremium, "2026-12-31") {
		t.Errorf("premium message should contain expiry date")
	}

	msgFree := formatStatsMessage(free, "uz", 1)
	if !contains(msgFree, "Yo'q") {
		t.Errorf("free message should indicate no subscription")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
