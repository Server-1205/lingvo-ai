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


