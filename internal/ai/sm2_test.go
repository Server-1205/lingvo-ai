package ai

import (
	"math"
	"testing"
	"time"
)

func TestNewSM2Card(t *testing.T) {
	card := NewSM2Card()
	if card.Repetitions != 0 {
		t.Errorf("expected Repetitions=0, got %d", card.Repetitions)
	}
	if card.Interval != 0 {
		t.Errorf("expected Interval=0, got %d", card.Interval)
	}
	if card.EaseFactor != 2.5 {
		t.Errorf("expected EaseFactor=2.5, got %f", card.EaseFactor)
	}
}

func TestProcessReviewFail(t *testing.T) {
	card := NewSM2Card()
	result := ProcessReview(card, 0)

	if result.Repetitions != 0 {
		t.Errorf("fail: expected Repetitions=0, got %d", result.Repetitions)
	}
	if result.Interval != 1 {
		t.Errorf("fail: expected Interval=1, got %d", result.Interval)
	}
	if result.EaseFactor != 2.5 {
		t.Errorf("fail: expected EaseFactor=2.5, got %f", result.EaseFactor)
	}

	// quality < 3 also resets: Hard (2)
	result = ProcessReview(card, 2)
	if result.Repetitions != 0 {
		t.Errorf("hard: expected Repetitions=0, got %d", result.Repetitions)
	}
}

func TestProcessReviewFirstPass(t *testing.T) {
	card := NewSM2Card()
	result := ProcessReview(card, 4)

	if result.Repetitions != 1 {
		t.Errorf("expected Repetitions=1, got %d", result.Repetitions)
	}
	if result.Interval != 1 {
		t.Errorf("expected Interval=1, got %d", result.Interval)
	}
	if result.EaseFactor != 2.5 {
		t.Errorf("expected EF=2.5 for quality=4, got %f", result.EaseFactor)
	}
}

func TestProcessReviewPerfect(t *testing.T) {
	card := NewSM2Card()
	result := ProcessReview(card, 5)

	if result.Repetitions != 1 {
		t.Errorf("expected Repetitions=1, got %d", result.Repetitions)
	}
	if result.Interval != 1 {
		t.Errorf("expected Interval=1, got %d", result.Interval)
	}
	// EF' = 2.5 + (0.1 - 0*(0.08+0*0.02)) = 2.6
	if result.EaseFactor != 2.6 {
		t.Errorf("expected EF=2.6 for quality=5, got %f", result.EaseFactor)
	}
}

func TestProcessReviewSecondPass(t *testing.T) {
	card := SM2Card{
		Repetitions: 1,
		Interval:    1,
		EaseFactor:  2.6,
	}
	result := ProcessReview(card, 4)

	if result.Repetitions != 2 {
		t.Errorf("expected Repetitions=2, got %d", result.Repetitions)
	}
	if result.Interval != 6 {
		t.Errorf("expected Interval=6, got %d", result.Interval)
	}
}

func TestProcessReviewThirdPass(t *testing.T) {
	card := SM2Card{
		Repetitions: 2,
		Interval:    6,
		EaseFactor:  2.6,
	}
	result := ProcessReview(card, 4)

	if result.Repetitions != 3 {
		t.Errorf("expected Repetitions=3, got %d", result.Repetitions)
	}
	expectedInterval := int(math.Round(float64(6) * 2.6))
	if result.Interval != expectedInterval {
		t.Errorf("expected Interval=%d, got %d", expectedInterval, result.Interval)
	}
}

func TestEFFloor(t *testing.T) {
	card := SM2Card{
		Repetitions: 0,
		Interval:    0,
		EaseFactor:  1.3,
	}
	// quality=3 gives the smallest EF increase
	result := ProcessReview(card, 3)
	if result.EaseFactor < 1.3 {
		t.Errorf("expected EF >= 1.3, got %f", result.EaseFactor)
	}
	// quality=0 should keep EF unchanged (fail path)
	result = ProcessReview(card, 0)
	if result.EaseFactor != 1.3 {
		t.Errorf("expected EF=1.3 on fail, got %f", result.EaseFactor)
	}
}

func TestNextReviewDate(t *testing.T) {
	card := NewSM2Card()
	result := ProcessReview(card, 4) // first pass, interval=1

	today := time.Now().Truncate(24 * time.Hour)
	expectedNext := today.AddDate(0, 0, 1)
	if !result.NextReview.Equal(expectedNext) {
		t.Errorf("expected next_review=%v, got %v", expectedNext, result.NextReview)
	}
}

func TestMultipleReviewCycles(t *testing.T) {
	card := NewSM2Card()

	// Day 1: learn (quality=5)
	result := ProcessReview(card, 5)
	if result.Interval != 1 {
		t.Errorf("day1: expected Interval=1, got %d", result.Interval)
	}
	if result.EaseFactor != 2.6 {
		t.Errorf("day1: expected EF=2.6, got %f", result.EaseFactor)
	}

	// Day 2: first review (quality=4)
	card2 := SM2Card{
		Repetitions: result.Repetitions,
		Interval:    result.Interval,
		EaseFactor:  result.EaseFactor,
	}
	result2 := ProcessReview(card2, 4)
	if result2.Repetitions != 2 {
		t.Errorf("day2: expected Repetitions=2, got %d", result2.Repetitions)
	}
	if result2.Interval != 6 {
		t.Errorf("day2: expected Interval=6, got %d", result2.Interval)
	}

	// Day 8: second review (quality=3)
	card3 := SM2Card{
		Repetitions: result2.Repetitions,
		Interval:    result2.Interval,
		EaseFactor:  result2.EaseFactor,
	}
	result3 := ProcessReview(card3, 3)
	if result3.Repetitions != 3 {
		t.Errorf("day8: expected Repetitions=3, got %d", result3.Repetitions)
	}
	expectedInterval := int(math.Round(float64(6) * result2.EaseFactor))
	if result3.Interval != expectedInterval {
		t.Errorf("day8: expected Interval=%d, got %d", expectedInterval, result3.Interval)
	}
}
