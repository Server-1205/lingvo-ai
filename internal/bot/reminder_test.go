package bot

import (
	"strings"
	"sync"
	"testing"
	"time"
)

func TestReminderText(t *testing.T) {
	uz := reminderText("uz", 5)
	if !strings.Contains(uz, "5") {
		t.Errorf("uz reminder should contain count 5, got: %s", uz)
	}
	if !strings.Contains(uz, "Eslatma") {
		t.Errorf("uz reminder should contain 'Eslatma', got: %s", uz)
	}

	ru := reminderText("ru", 3)
	if !strings.Contains(ru, "3") {
		t.Errorf("ru reminder should contain count 3, got: %s", ru)
	}
	if !strings.Contains(ru, "Напоминание") {
		t.Errorf("ru reminder should contain 'Напоминание', got: %s", ru)
	}

	unknown := reminderText("kk", 2)
	if !strings.Contains(unknown, "2") {
		t.Errorf("unknown lang should fallback to uz, got: %s", unknown)
	}
}

func TestReviewButton(t *testing.T) {
	uzKeyboard := reviewButton("uz")
	if len(uzKeyboard.InlineKeyboard) == 0 || len(uzKeyboard.InlineKeyboard[0]) == 0 {
		t.Fatal("reviewButton should return a keyboard with at least one button")
	}
	btn := uzKeyboard.InlineKeyboard[0][0]
	if !strings.Contains(btn.Text, "Hozir") {
		t.Errorf("uz button should contain 'Hozir', got: %s", btn.Text)
	}
	if btn.URL == nil || !strings.Contains(*btn.URL, "startapp=review") {
		t.Errorf("button URL should contain startapp=review, got: %v", btn.URL)
	}

	ruKeyboard := reviewButton("ru")
	ruBtn := ruKeyboard.InlineKeyboard[0][0]
	if !strings.Contains(ruBtn.Text, "Повторить") {
		t.Errorf("ru button should contain 'Повторить', got: %s", ruBtn.Text)
	}
}

func TestRemindedTodayDedup(t *testing.T) {
	remindedMu = sync.Mutex{}
	remindedToday = map[int64]time.Time{}

	now := time.Now()

	remindedMu.Lock()
	remindedToday[123] = now
	remindedMu.Unlock()

	remindedMu.Lock()
	entry, exists := remindedToday[123]
	remindedMu.Unlock()

	if !exists {
		t.Fatal("expected remindedToday to contain entry for user 123")
	}
	if !entry.Truncate(24 * time.Hour).Equal(now.Truncate(24 * time.Hour)) {
		t.Errorf("entry time should match today, got %v", entry)
	}

	remindedMu.Lock()
	_, exists456 := remindedToday[456]
	remindedMu.Unlock()

	if exists456 {
		t.Error("user 456 should not be in remindedToday")
	}
}

func TestReviewButtonDeepLink(t *testing.T) {
	kbd := reviewButton("uz")
	btn := kbd.InlineKeyboard[0][0]

	expectedPrefix := "https://t.me/lingvo_ai_bot/app"
	if btn.URL == nil || !strings.HasPrefix(*btn.URL, expectedPrefix) {
		t.Errorf("button URL should start with %s, got: %v", expectedPrefix, btn.URL)
	}
}
